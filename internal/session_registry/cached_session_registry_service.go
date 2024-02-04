package session_registry

import (
	"errors"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/apps_registry"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"pokt_gateway_server/pkg/ttl_cache"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	counterSessionRequest          *prometheus.CounterVec
	histogramSessionRequestLatency *prometheus.HistogramVec
	ErrRecentlyFailed              = errors.New("dispatch recently failed, returning early")
)

const (
	blocksPerSession                      = 4
	sessionPrimerInterval                 = time.Second * 15
	reasonSessionSuccessCached            = "session_cached"
	reasonSessionSuccessColdHit           = "session_cold_hit"
	reasonSessionFailedBackoff            = "session_failed_backoff"
	reasonSessionFailedUnderlyingProvider = "session_failed_from_client"
	backoffThreshold                      = time.Second * 2
	maxConcurrentDispatch                 = 50
)

func init() {
	counterSessionRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cached_client_session_request_counter",
			Help: "Request to get a session and if it succeeded",
		},
		[]string{"success", "reason"},
	)

	histogramSessionRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cached_client_session_request_latency",
			Help: "percentile on the request to get a session",
		},
		[]string{"cached"},
	)

	prometheus.MustRegister(counterSessionRequest, histogramSessionRequestLatency)
}

type CachedSessionRegistryService struct {
	poktClient             pokt_v0.PocketService
	appRegistry            apps_registry.AppsRegistryService
	lastFailure            time.Time
	concurrentDispatchPool chan struct{}
	logger                 *zap.Logger
	lastPrimedHeight       uint64
	// Consist of sessions for a given app stake+chain+height. Cache exists to prevent round trip request
	sessionCache ttl_cache.TTLCacheService[string, *models.GetSessionResponse]
	// Cache that contains all of the nodes for all sessions for a specific chain+height.  Cache exists to prevent round trip request
	sessionNodesCache ttl_cache.TTLCacheService[string, []*models.Node] // sessionHeight -> nodes
}

func NewCachedSessionRegistryService(poktClient pokt_v0.PocketService, appRegistry apps_registry.AppsRegistryService, sessionCache ttl_cache.TTLCacheService[string, *models.GetSessionResponse], nodeCache ttl_cache.TTLCacheService[string, []*models.Node], logger *zap.Logger) *CachedSessionRegistryService {
	cachedRegistry := &CachedSessionRegistryService{poktClient: poktClient, appRegistry: appRegistry, sessionCache: sessionCache, lastFailure: time.Time{}, concurrentDispatchPool: make(chan struct{}, maxConcurrentDispatch), sessionNodesCache: nodeCache, logger: logger}
	go sessionCache.Start()
	go nodeCache.Start()
	cachedRegistry.startSessionUpdater()
	return cachedRegistry
}

func (c CachedSessionRegistryService) GetSession(req *models.GetSessionRequest) (*models.GetSessionResponse, error) {
	cacheKey := getSessionCacheKey(req)
	cachedSession := c.sessionCache.Get(cacheKey)
	isCached := cachedSession != nil && cachedSession.Value() != nil
	startTime := time.Now()
	// Measure end to end latency for send relay
	defer func() {
		histogramSessionRequestLatency.WithLabelValues(strconv.FormatBool(isCached)).Observe(float64(time.Since(startTime)))
	}()

	if isCached {
		counterSessionRequest.WithLabelValues("true", reasonSessionSuccessCached).Inc()
		return cachedSession.Value(), nil
	}

	// Backoff check
	if c.shouldBackoffDispatchFailure() {
		counterSessionRequest.WithLabelValues("false", reasonSessionFailedBackoff).Inc()
		return nil, ErrRecentlyFailed
	}

	// Limits the number of concurrent calls going out to a node
	// to prevent overloading the node during session rollover
	c.concurrentDispatchPool <- struct{}{}
	defer func() {
		<-c.concurrentDispatchPool
	}()

	// Call underlying provider
	response, err := c.poktClient.GetSession(req)
	if err != nil {
		counterSessionRequest.WithLabelValues("false", reasonSessionFailedUnderlyingProvider).Inc()
		c.lastFailure = time.Now()
		return nil, err
	}

	counterSessionRequest.WithLabelValues("true", reasonSessionSuccessColdHit).Inc()
	// Update session cache
	c.sessionCache.Set(cacheKey, response, ttlcache.DefaultTTL)

	// Update node cache
	nodeCacheKey := getChainNodesCacheKey(req)
	nodes := c.sessionNodesCache.Get(nodeCacheKey)
	if nodes == nil || nodes.Value() == nil {
		// No values in session and chain cache
		c.sessionNodesCache.Set(cacheKey, response.Session.Nodes, ttlcache.DefaultTTL)
	} else {
		// Values already exist in session and chain cache, so append value.
		c.sessionNodesCache.Set(cacheKey, append(nodes.Value(), response.Session.Nodes...), ttlcache.DefaultTTL)
	}

	c.lastFailure = time.Time{} // Reset last failure since it succeeded
	return response, nil
}

func (c CachedSessionRegistryService) startSessionUpdater() {
	ticker := time.Tick(sessionPrimerInterval)
	go func() {
		for {
			select {
			case <-ticker:
				err := c.primeSessions()
				if err != nil {
					c.logger.Sugar().Error(err)
				}
			}
		}
	}()
}

// shouldPrimeSession: Track the latest time we primed a session and only prime if there's a new session
func (c CachedSessionRegistryService) shouldPrimeSessions(latestHeight uint64) bool {
	isSessionBlock := latestHeight%blocksPerSession == 1
	isNewSessionBlock := latestHeight > c.lastPrimedHeight
	return c.lastPrimedHeight == 0 || isSessionBlock && isNewSessionBlock
}

// primeSession: used as a background service to optimistically grab sessions
// before relays are handled to prevent unneeded round trips.
func (c CachedSessionRegistryService) primeSessions() error {

	resp, err := c.poktClient.GetLatestBlockHeight()

	if err != nil {
		return err
	}

	if !c.shouldPrimeSessions(resp.Height) {
		return nil
	}

	sessionCount := atomic.Int32{}

	wg := sync.WaitGroup{}
	for _, app := range c.appRegistry.GetApplications() {
		for _, chain := range app.NetworkApp.Chains {
			wg.Add(1)
			app := app
			chain := chain
			go func() {
				defer wg.Done()
				// Goroutine unbounded
				req := &models.GetSessionRequest{
					AppPubKey: app.NetworkApp.PublicKey,
					Chain:     chain,
					Height:    uint(resp.Height),
				}
				_, err = c.GetSession(req)
				if err != nil {
					c.logger.Sugar().Warnw("primeSessions: failed to prime session", "req", req, "err", err)
				} else {
					sessionCount.Add(1)
				}
			}()
		}
	}
	wg.Wait()
	totalSessionsPrimed := sessionCount.Load()
	// As long as we prime at least one session,
	// we consider it as a success and will wait until
	// next session block height to continue priming
	if totalSessionsPrimed > 0 {
		c.logger.Sugar().Infow("primeSessions: successfully primed sessions", "sessionsPrimed", totalSessionsPrimed)
		c.lastPrimedHeight = resp.Height
	}
	return nil
}

// shouldBackOffDispatchFailure: whenever pokt nodes receive too many dispatches at once, it results in overloaded pokt nodes
// and subsequent dispatch failures.
// TODO: Optimization: Allow for % backoff/retry mechanism instead of constant backoff threshhold.
func (c CachedSessionRegistryService) shouldBackoffDispatchFailure() bool {
	return !c.lastFailure.IsZero() && time.Since(c.lastFailure) < backoffThreshold
}

// getSessionCacheKey - used to keep track of a session for a specific app stake, height, and chain.
func getSessionCacheKey(req *models.GetSessionRequest) string {
	return fmt.Sprintf("%s-%s-%d", req.AppPubKey, req.Chain, req.Height)
}

// getSessionNodeCacheKey - used to keep track of all nodes for a specific chain and height.
func getChainNodesCacheKey(req *models.GetSessionRequest) string {
	return fmt.Sprintf("%s-%d", req.Chain, req.Height)
}
