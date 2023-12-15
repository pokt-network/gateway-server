package pokt_client_decorators

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"os-gateway/pkg/pokt/pokt_v0"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"os-gateway/pkg/ttl_cache"
	"strconv"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

const backoffThreshold = time.Second * 5
const maxConcurrentDispatch = 50

var ErrRecentlyFailed = errors.New("dispatch recently failed, returning early")

var (
	counterSessionRequest          *prometheus.CounterVec
	counterRelayRequest            *prometheus.CounterVec
	histogramSessionRequestLatency *prometheus.HistogramVec
	histogramRelayRequestLatency   prometheus.Histogram
)

const (
	reasonSessionFailedBackoff            = "session_backoff"
	reasonSessionFailedUnderlyingProvider = "failed_from_client"
	reasonRelayFailedSessionErr           = "session_failure"
	reasonRelayFailedUnderlyingProvider   = "relay_failure"
)

func init() {
	counterSessionRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cached_client_session_request_counter",
			Help: "Request to get a session and if it succeeded",
		},
		[]string{"success", "reason"},
	)
	counterRelayRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cached_client_relay_counter",
			Help: "Request to send an actual relay and if it succeeded",
		},
		[]string{},
	)
	histogramSessionRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cached_client_session_request_latency",
			Help: "percentile on the request to get a session",
		},
		[]string{"cached"},
	)
	histogramRelayRequestLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "cached_client_relay_request_latency",
			Help: "percentile on the request to send a relay",
		},
	)
	prometheus.MustRegister(counterRelayRequest, counterSessionRequest, histogramRelayRequestLatency, histogramSessionRequestLatency)
}

type CachedClient struct {
	pokt_v0.PocketService
	lastFailure            time.Time
	concurrentDispatchPool chan struct{}
	sessionCache           ttl_cache.TTLCacheService[string, *models.GetSessionResponse]
}

func NewCachedClient(pocketService pokt_v0.PocketService, sessionCache ttl_cache.TTLCacheService[string, *models.GetSessionResponse]) *CachedClient {

	return &CachedClient{
		PocketService:          pocketService,
		lastFailure:            time.Time{},
		sessionCache:           sessionCache,
		concurrentDispatchPool: make(chan struct{}, maxConcurrentDispatch),
	}

}

func (c *CachedClient) GetSession(req *models.GetSessionRequest) (*models.GetSessionResponse, error) {
	cacheKey := getCacheKey(req)
	cachedSession := c.sessionCache.Get(cacheKey)

	isCached := cachedSession != nil && cachedSession.Value() != nil
	startTime := time.Now()
	// Measure end to end latency for send relay
	defer func() {
		histogramSessionRequestLatency.WithLabelValues(strconv.FormatBool(isCached)).Observe(float64(time.Since(startTime)))
	}()

	if isCached {
		return cachedSession.Value(), nil
	}

	// Backoff check
	if c.shouldBackoff() {
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
	response, err := c.PocketService.GetSession(req)
	if err != nil {
		counterSessionRequest.WithLabelValues("false", reasonSessionFailedUnderlyingProvider).Inc()
		c.lastFailure = time.Now()
		return nil, err
	}

	counterSessionRequest.WithLabelValues("true").Inc()
	c.sessionCache.Set(cacheKey, response, ttlcache.DefaultTTL)
	c.lastFailure = time.Time{} // Reset last failure since it succeeded
	return response, nil
}

func (r *CachedClient) SendRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}

	startTime := time.Now()
	// Measure end to end latency for send relay
	defer func() {
		histogramRelayRequestLatency.Observe(float64(time.Since(startTime)))
	}()

	session, err := pokt_v0.GetSessionFromRequest(r, req)

	if err != nil {
		counterRelayRequest.WithLabelValues("false", reasonRelayFailedSessionErr).Inc()
		return nil, err
	}

	req.Session = session
	rsp, err := r.PocketService.SendRelay(req)

	// Emit metrics on success
	if err != nil {
		counterRelayRequest.WithLabelValues("false", reasonRelayFailedUnderlyingProvider).Inc()
	} else {
		counterRelayRequest.WithLabelValues("true").Inc()
	}

	return rsp, err
}

func (c *CachedClient) shouldBackoff() bool {
	return !c.lastFailure.IsZero() && time.Since(c.lastFailure) < backoffThreshold
}

func getCacheKey(req *models.GetSessionRequest) string {
	return req.AppPubKey + "-" + req.Chain
}
