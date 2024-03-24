package relayer

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/altruist_registry"
	"pokt_gateway_server/internal/apps_registry"
	"pokt_gateway_server/internal/node_selector_service"
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
	pocketClient        pokt_v0.PocketService
	altruistRegistry    altruist_registry.AltruistRegistryService
	altruistTimeout     time.Duration
	sessionRegistry     session_registry.SessionRegistryService
	nodeSelector        node_selector_service.NodeSelectorService
	applicationRegistry apps_registry.AppsRegistryService
	logger              *zap.Logger
}

func NewRelayer(pocketService pokt_v0.PocketService, sessionRegistry session_registry.SessionRegistryService, applicationRegistry apps_registry.AppsRegistryService, nodeSelector node_selector_service.NodeSelectorService, altruistRegistry altruist_registry.AltruistRegistryService, altruistTimeout time.Duration, logger *zap.Logger) *Relayer {
	return &Relayer{
		pocketClient:        pocketService,
		sessionRegistry:     sessionRegistry,
		altruistTimeout:     altruistTimeout,
		logger:              logger,
		altruistRegistry:    altruistRegistry,
		applicationRegistry: applicationRegistry,
		nodeSelector:        nodeSelector,
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
		return nil, errors.New("node selector can't find node")
	}
	req.Signer = node.AppSigner
	req.Session = node.PocketSession
	req.SelectedNodePubKey = node.GetPublicKey()
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return r.pocketClient.SendRelay(req)
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

	// populate request with session metadata
	req.Session = randomNode.PocketSession
	req.Signer = appStake.Signer
	req.SelectedNodePubKey = randomNode.GetPublicKey()
	rsp, err := r.pocketClient.SendRelay(req)

	// record if relay was successful
	counterRelayRequest.WithLabelValues(strconv.FormatBool(err == nil), "false", "").Inc()

	return rsp, err
}

func (r *Relayer) altruistRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error) {

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

	success := err == nil
	counterRelayRequest.WithLabelValues(strconv.FormatBool(success), "true", "").Inc()

	if !success {
		return nil, err
	}

	str := string(response.Body())
	return &models.SendRelayResponse{Response: str}, nil
}
