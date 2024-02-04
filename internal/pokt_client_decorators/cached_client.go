package pokt_client_decorators

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/altruist_registry"
	"pokt_gateway_server/internal/session_registry"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"time"
)

var (
	counterRelayRequest          *prometheus.CounterVec
	histogramRelayRequestLatency prometheus.Histogram
)

const (
	reasonRelayFailedSessionErr         = "relay_session_failure"
	reasonRelayFailedUnderlyingProvider = "relay_provider_failure"
)

func init() {
	counterRelayRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cached_client_relay_counter",
			Help: "Request to send an actual relay and if it succeeded",
		},
		[]string{"success", "reason"},
	)
	histogramRelayRequestLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "cached_client_relay_request_latency",
			Help: "percentile on the request to send a relay",
		},
	)
	prometheus.MustRegister(counterRelayRequest, histogramRelayRequestLatency)
}

type CachedClient struct {
	pokt_v0.PocketService
	altruistRegistry altruist_registry.AltruistRegistryService
	sessionRegistry  session_registry.SessionRegistryService
	altruistTimeout  time.Duration
	logger           *zap.Logger
}

func NewCachedClient(pocketService pokt_v0.PocketService, sessionRegistry session_registry.SessionRegistryService, altruistRegistry altruist_registry.AltruistRegistryService, altruistTimeout time.Duration, logger *zap.Logger) *CachedClient {
	return &CachedClient{
		PocketService:    pocketService,
		sessionRegistry:  sessionRegistry,
		altruistTimeout:  altruistTimeout,
		logger:           logger,
		altruistRegistry: altruistRegistry,
	}
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

	session, err := pokt_v0.GetSessionFromRequest(r.sessionRegistry, req)

	if err != nil {
		counterRelayRequest.WithLabelValues("false", reasonRelayFailedSessionErr).Inc()
		return nil, err
	}

	req.Session = session
	rsp, err := r.PocketService.SendRelay(req)

	// If request fails, send to altruist.
	if err != nil {
		r.logger.Sugar().Errorw("failed to send to pokt", "poktErr", err)
		counterRelayRequest.WithLabelValues("false", reasonRelayFailedUnderlyingProvider).Inc()
		altruistRsp, altruistErr := r.altruistRelay(req)

		if altruistErr != nil {
			r.logger.Sugar().Errorw("failed to send to altruist", "altruistError", altruistErr)
			// Prefer to return the network error vs altruist error if both fails.
			return nil, err
		}
		return altruistRsp, nil
	}

	counterRelayRequest.WithLabelValues("true", "").Inc()
	return rsp, nil
}

func (r *CachedClient) altruistRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {

	url, ok := r.altruistRegistry.GetAltruistURL(req.Chain)

	if !ok {
		return nil, errors.New("altruist url not found")
	}

	// Send to altruist
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	request.SetRequestURI(url)

	if req.Payload.Method == "POST" {
		request.SetBody([]byte(req.Payload.Data))
	}

	err := fasthttp.DoTimeout(request, response, r.altruistTimeout)
	if err != nil {
		return nil, err
	}

	str := string(response.Body())
	return &models.SendRelayResponse{Response: str}, nil
}
