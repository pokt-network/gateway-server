package session_registry

import (
	qos_models "pokt_gateway_server/internal/qos_node_registry/models"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
)

type Session struct {
	IsValid bool
	Nodes   []*qos_models.QosNode
}

type SessionRegistryService interface {
	GetSession(req *models.GetSessionRequest) (*Session, error)
	GetNodes() []*qos_models.QosNode
	GetNodesByChain(chainId string) ([]*qos_models.QosNode, bool)
}
