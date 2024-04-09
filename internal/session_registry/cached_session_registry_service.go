package session_registry

import (
	"errors"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"github.com/pokt-network/gateway-server/internal/apps_registry"
	qos_models "github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
	"github.com/pokt-network/gateway-server/pkg/ttl_cache"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
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
	sessionPrimerInterval                 = time.Second * 5
	ttlCacheCleanerInterval               = time.Second * 15
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
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 30, 40, 50, 60},
			Name:    "cached_client_session_request_latency",
			Help:    "percentile on the request to get a session",
		},
		[]string{"cached"},
	)

	prometheus.MustRegister(counterSessionRequest, histogramSessionRequestLatency)
}

type CachedSessionRegistryService struct {
	poktClient              pokt_v0.PocketService
	appRegistry             apps_registry.AppsRegistryService
	lastFailure             time.Time
	concurrentDispatchPool  chan struct{}
	logger                  *zap.Logger
	lastPrimedSessionHeight uint
	// Lock used to synchronize inserting sessions and append sessions nodes.
	sessionCacheLock sync.RWMutex
	// Consist of sessions for a given app stake+chain+height. Cache exists to prevent round trip request
	sessionCache ttl_cache.TTLCacheService[string, *Session]
	// Cache that contains all nodes by chain (chainId -> Nodes)
	chainNodes ttl_cache.TTLCacheService[qos_models.SessionChainKey, []*qos_models.QosNode] // sessionHeight -> nodes
}

func NewCachedSessionRegistryService(poktClient pokt_v0.PocketService, appRegistry apps_registry.AppsRegistryService, sessionCache ttl_cache.TTLCacheService[string, *Session], nodeCache ttl_cache.TTLCacheService[qos_models.SessionChainKey, []*qos_models.QosNode], logger *zap.Logger) *CachedSessionRegistryService {
	cachedRegistry := &CachedSessionRegistryService{poktClient: poktClient, appRegistry: appRegistry, sessionCache: sessionCache, lastFailure: time.Time{}, concurrentDispatchPool: make(chan struct{}, maxConcurrentDispatch), chainNodes: nodeCache, logger: logger}
	go sessionCache.Start()
	go nodeCache.Start()
	cachedRegistry.startTTLCacheCleaner()
	cachedRegistry.startSessionUpdater()
	return cachedRegistry
}

func (c *CachedSessionRegistryService) startTTLCacheCleaner() {
	ticker := time.Tick(ttlCacheCleanerInterval)
	go func() {
		for {
			select {
			case <-ticker:
				c.sessionCache.DeleteExpired()
				c.chainNodes.DeleteExpired()
			}
		}
	}()
}

func (c *CachedSessionRegistryService) GetNodesByChain(chainId string) []*qos_models.QosNode {
	items := c.GetNodesMap()
	nodes := []*qos_models.QosNode{}
	for sessionKey, item := range items {
		if sessionKey.Chain == chainId {
			nodes = append(nodes, item.Value()...)
		}
	}
	return nodes
}

func (c *CachedSessionRegistryService) GetNodesMap() map[qos_models.SessionChainKey]*ttlcache.Item[qos_models.SessionChainKey, []*qos_models.QosNode] {
	c.sessionCacheLock.RLock()
	defer c.sessionCacheLock.RUnlock()
	return c.chainNodes.Items()
}

func (c *CachedSessionRegistryService) GetSession(req *models.GetSessionRequest) (*Session, error) {
	sessionCacheKey := getSessionCacheKey(req)
	cachedSession := c.sessionCache.Get(sessionCacheKey)
	isCached := cachedSession != nil && cachedSession.Value() != nil
	startTime := time.Now()
	// Measure end to end latency for send relay
	defer func() {
		histogramSessionRequestLatency.WithLabelValues(strconv.FormatBool(isCached)).Observe(time.Since(startTime).Seconds())
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

	appSigner, found := c.appRegistry.GetApplicationByPublicKey(req.AppPubKey)

	if !found {
		return nil, errors.New("cannot find signer from session")
	}

	wrappedNodes := []*qos_models.QosNode{}
	for _, a := range response.Session.Nodes {
		wrappedNodes = append(wrappedNodes, qos_models.NewQosNode(a, response.Session, appSigner.Signer))
	}

	// session with metadata
	wrappedSession := &Session{IsValid: true, Nodes: wrappedNodes, PocketSession: response.Session}

	counterSessionRequest.WithLabelValues("true", reasonSessionSuccessColdHit).Inc()

	c.sessionCacheLock.Lock()
	defer c.sessionCacheLock.Unlock()
	// Update session cache
	c.sessionCache.Set(sessionCacheKey, wrappedSession, ttlcache.DefaultTTL)

	// Update node cache
	chainNodeCacheKey := qos_models.SessionChainKey{Chain: req.Chain, SessionHeight: wrappedSession.PocketSession.SessionHeader.SessionHeight}
	if !c.chainNodes.Has(chainNodeCacheKey) {
		// No values in session and chain cache
		c.chainNodes.Set(chainNodeCacheKey, wrappedNodes, ttlcache.DefaultTTL)
	} else {
		// Values already exist in session and chain cache, so append value.
		item := c.chainNodes.Get(chainNodeCacheKey)
		if item != nil {
			chainNodesAppended := append(item.Value(), wrappedNodes...)
			c.chainNodes.Set(chainNodeCacheKey, chainNodesAppended, ttlcache.DefaultTTL)
		}
	}

	c.lastFailure = time.Time{} // Reset last failure since it succeeded
	return wrappedSession, nil
}

func (c *CachedSessionRegistryService) startSessionUpdater() {
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
func (c *CachedSessionRegistryService) shouldPrimeSessions(latestSessionHeight uint) bool {
	isNewSessionBlock := latestSessionHeight > c.lastPrimedSessionHeight
	return c.lastPrimedSessionHeight == 0 || isNewSessionBlock
}

// primeSession: used as a background service to optimistically grab sessions
// before relays are handled to prevent unneeded round trips.
func (c *CachedSessionRegistryService) primeSessions() error {

	resp, err := c.poktClient.GetLatestBlockHeight()

	if err != nil {
		return err
	}

	latestBlockHeight := resp.Height
	latestSessionHeight := getLatestSessionHeight(latestBlockHeight)
	shouldPrimeSessions := c.shouldPrimeSessions(latestSessionHeight)
	c.logger.Sugar().Infow("priming sessions async", "currentBlockHeight", resp.Height, "latestSessionHeight", latestSessionHeight, "lastPrimedSessionHeight", c.lastPrimedSessionHeight, "shouldPrimeSessions", shouldPrimeSessions)

	if !shouldPrimeSessions {
		return nil
	}

	errCount := atomic.Int32{}
	successCount := atomic.Int32{}
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
					AppPubKey:     app.NetworkApp.PublicKey,
					Chain:         chain,
					SessionHeight: latestSessionHeight,
				}
				_, err := c.GetSession(req)
				if err != nil {
					errCount.Add(1)
					c.logger.Sugar().Warnw("primeSessions: failed to prime session", "req", req, "err", err, "latestSessionHeight", latestSessionHeight)
				} else {
					successCount.Add(1)
				}
			}()
		}
	}
	wg.Wait()
	successes := successCount.Load()
	errs := errCount.Load()
	if errs == 0 && successes > 0 {
		c.logger.Sugar().Infow("primeSessions: successfully primed sessions", "successCount", successes, "errorCount", errs)
		c.lastPrimedSessionHeight = latestSessionHeight
	}
	return nil
}

func getLatestSessionHeight(nodeHeight uint) uint {
	// if block height / blocks per session remainder is zero, just subtract blocks per session and add 1
	if nodeHeight%blocksPerSession == 0 {
		return nodeHeight - blocksPerSession + 1
	} else {
		// calculate the latest session block height by diving the current block height by the blocksPerSession
		return (nodeHeight/blocksPerSession)*blocksPerSession + 1
	}
}

// shouldBackOffDispatchFailure: whenever pokt nodes receive too many dispatches at once, it results in overloaded pokt nodes
// and subsequent dispatch failures.
// TODO: Optimization: Allow for % backoff/retry mechanism instead of constant backoff threshhold.
func (c *CachedSessionRegistryService) shouldBackoffDispatchFailure() bool {
	return !c.lastFailure.IsZero() && time.Since(c.lastFailure) < backoffThreshold
}

// getSessionCacheKey - used to keep track of a session for a specific app stake, height, and chain.
func getSessionCacheKey(req *models.GetSessionRequest) string {
	return fmt.Sprintf("%s-%s-%d", req.AppPubKey, req.Chain, req.SessionHeight)
}
