package relayer

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/apps_registry"
	"pokt_gateway_server/internal/chain_configurations_registry"
	"pokt_gateway_server/internal/global_config"
	"pokt_gateway_server/internal/node_selector_service"
	"pokt_gateway_server/internal/node_selector_service/checks"
	"pokt_gateway_server/internal/session_registry"
	"pokt_gateway_server/pkg/common"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"strconv"
	"time"
)

var (
	counterRelayRequest          *prometheus.CounterVec
	histogramRelayRequestLatency prometheus.Histogram
)

const (
	reasonRelayFailedSessionErr = "relay_session_failure"
	reasonRelayFailedPocketErr  = "relay_pocket_error"
)

var (
	errAltruistNotFound = errors.New("altruist url not found")
	errSelectNodeFail   = errors.New("node selector can't find node")
)

func init() {
	counterRelayRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "relay_counter",
			Help: "Request to send an actual relay and if it succeeded",
		},
		[]string{"success", "altruist", "reason"},
	)
	histogramRelayRequestLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "relay_latency",
			Help: "percentile on the request to send a relay",
		},
	)
	prometheus.MustRegister(counterRelayRequest, histogramRelayRequestLatency)
}

type Relayer struct {
	globalConfigProvider       global_config.GlobalConfigProvider
	pocketClient               pokt_v0.PocketService
	chainConfigurationRegistry chain_configurations_registry.ChainConfigurationsService
	sessionRegistry            session_registry.SessionRegistryService
	nodeSelector               node_selector_service.NodeSelectorService
	applicationRegistry        apps_registry.AppsRegistryService
	httpRequester              httpRequester
	logger                     *zap.Logger
}

func NewRelayer(pocketService pokt_v0.PocketService, sessionRegistry session_registry.SessionRegistryService, applicationRegistry apps_registry.AppsRegistryService, nodeSelector node_selector_service.NodeSelectorService, altruistRegistry chain_configurations_registry.ChainConfigurationsService, globalConfigProvider global_config.GlobalConfigProvider, logger *zap.Logger) *Relayer {
	return &Relayer{
		pocketClient:               pocketService,
		sessionRegistry:            sessionRegistry,
		logger:                     logger,
		chainConfigurationRegistry: altruistRegistry,
		applicationRegistry:        applicationRegistry,
		nodeSelector:               nodeSelector,
		httpRequester:              fastHttpRequester{},
		globalConfigProvider:       globalConfigProvider,
	}
}

func (r *Relayer) SendRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {

	startTime := time.Now()
	// Measure end to end latency for send relay
	defer func() {
		histogramRelayRequestLatency.Observe(float64(time.Since(startTime)))
	}()

	rsp, err := r.sendNodeSelectorRelay(req)

	// Node selector relay was successful
	if err == nil {
		counterRelayRequest.WithLabelValues("true", "false", "").Inc()
		return rsp, nil
	}

	counterRelayRequest.WithLabelValues("false", "true", reasonRelayFailedPocketErr).Inc()

	r.logger.Sugar().Errorw("failed to send to pokt", "poktErr", err)
	altruistRsp, altruistErr := r.altruistRelay(req)
	if altruistErr != nil {
		r.logger.Sugar().Errorw("failed to send to altruist", "altruistError", altruistErr)
		// Prefer to return the network error vs altruist error if both fails.
		return nil, err
	}
	return altruistRsp, nil
}

func (r *Relayer) sendNodeSelectorRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {
	// find a node to send too first.
	node, ok := r.nodeSelector.FindNode(req.Chain)
	if !ok {
		return nil, errSelectNodeFail
	}
	req.Signer = node.MorseSigner
	req.Session = node.MorseSession
	req.SelectedNodePubKey = node.GetPublicKey()
	if err := req.Validate(); err != nil {
		return nil, err
	}

	start := time.Now()
	rsp, err := r.pocketClient.SendRelay(req)
	node.GetLatencyTracker().RecordMeasurement(float64(time.Now().Sub(start).Milliseconds()))
	// Node returned an error, potentially penalize the node operator dependent on error
	if err != nil {
		checks.DefaultPunishNode(err, node, r.logger)
	}
	return rsp, err
}

func (r *Relayer) sendRandomNodeRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {
	// Healthy node could not be found, attempting to use random node
	applications, ok := r.applicationRegistry.GetApplicationsByChainId(req.Chain)
	if !ok {
		return nil, fmt.Errorf("no app found for chain id %s", req.Chain)
	}

	// Get a random app stake from the available list.
	appStake, ok := common.GetRandomElement(applications)
	if !ok {
		return nil, fmt.Errorf("random app stake cannot be found")
	}

	sessionResp, err := r.sessionRegistry.GetSession(&models.GetSessionRequest{
		AppPubKey: appStake.Signer.PublicKey,
		Chain:     req.Chain,
	})

	if err != nil {
		counterRelayRequest.WithLabelValues("false", "true", reasonRelayFailedSessionErr).Inc()
		return nil, err
	}

	randomNode, ok := common.GetRandomElement(sessionResp.Nodes)
	if !ok {
		return nil, errors.New("random node in session cannot be found")
	}

	requestTimeout := r.getPocketRequestTimeout(req.Chain)
	// populate request with session metadata
	req.Session = randomNode.MorseSession
	req.Signer = appStake.Signer
	req.Timeout = &requestTimeout
	req.SelectedNodePubKey = randomNode.GetPublicKey()
	rsp, err := r.pocketClient.SendRelay(req)

	// record if relay was successful
	counterRelayRequest.WithLabelValues(strconv.FormatBool(err == nil), "false", "").Inc()

	return rsp, err
}

func (r *Relayer) altruistRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {

	chainConfig, ok := r.chainConfigurationRegistry.GetChainConfiguration(req.Chain)

	if !ok {
		return nil, errAltruistNotFound
	}

	// Send to altruist
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	requestTimeout := r.getAltruistRequestTimeout(req.Chain)
	request.SetRequestURI(chainConfig.AltruistUrl.String)

	if req.Payload.Method == "POST" {
		request.SetBody([]byte(req.Payload.Data))
	}

	err := r.httpRequester.DoTimeout(request, response, requestTimeout)

	success := err == nil
	counterRelayRequest.WithLabelValues(strconv.FormatBool(success), "true", "").Inc()

	if !success {
		return nil, err
	}

	str := string(response.Body())
	return &models.SendRelayResponse{Response: str}, nil
}

func (r *Relayer) getAltruistRequestTimeout(chainId string) time.Duration {
	chainConfig, ok := r.chainConfigurationRegistry.GetChainConfiguration(chainId)
	if !ok {
		return r.globalConfigProvider.GetAltruistRequestTimeout()
	}
	configTime, err := time.ParseDuration(chainConfig.AltruistRequestTimeoutDuration.String)
	if err != nil {
		return r.globalConfigProvider.GetAltruistRequestTimeout()
	}
	return configTime
}

func (r *Relayer) getPocketRequestTimeout(chainId string) time.Duration {
	chainConfig, ok := r.chainConfigurationRegistry.GetChainConfiguration(chainId)
	if !ok {
		return r.globalConfigProvider.GetPoktRPCRequestTimeout()
	}
	configTime, err := time.ParseDuration(chainConfig.AltruistRequestTimeoutDuration.String)
	if err != nil {
		return r.globalConfigProvider.GetPoktRPCRequestTimeout()
	}
	return configTime
}
