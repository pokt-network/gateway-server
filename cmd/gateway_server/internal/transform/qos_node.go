package transform

import (
	"math"
	"pokt_gateway_server/cmd/gateway_server/internal/models"
	internal_model "pokt_gateway_server/internal/node_selector_service/models"
)

func ToPublicQosNode(node *internal_model.QosNode) *models.PublicQosNode {
	latency := node.LatencyTracker.GetP90Latency()
	if math.IsNaN(latency) {
		latency = 0.0
	}
	return &models.PublicQosNode{
		NodePublicKey:   node.MorseNode.PublicKey,
		ServiceUrl:      node.MorseNode.ServiceUrl,
		Chain:           node.GetChain(),
		SessionHeight:   node.MorseSession.SessionHeader.SessionHeight,
		AppPublicKey:    node.MorseSigner.PublicKey,
		TimeoutReason:   string(node.GetTimeoutReason()),
		LastKnownErr:    node.GetLastKnownErrorStr(),
		IsHeathy:        node.IsHealthy(),
		IsSynced:        node.IsSynced(),
		LastKnownHeight: node.GetLastKnownHeight(),
		TimeoutUntil:    node.GetTimeoutUntil(),
		P90Latency:      latency,
	}
}
