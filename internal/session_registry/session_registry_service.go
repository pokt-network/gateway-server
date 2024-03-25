package session_registry

import (
	"github.com/jellydator/ttlcache/v3"
	qos_models "pokt_gateway_server/internal/node_selector_service/models"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
)

type Session struct {
	IsValid bool
	Nodes   []*qos_models.QosNode
}

type SessionRegistryService interface {
	GetSession(req *models.GetSessionRequest) (*Session, error)
	GetNodesMap() map[string]*ttlcache.Item[string, []*qos_models.QosNode]
	GetNodesByChain(chainId string) ([]*qos_models.QosNode, bool)
}
