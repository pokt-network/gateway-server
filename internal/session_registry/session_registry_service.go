package session_registry

import (
	"github.com/jellydator/ttlcache/v3"
	qos_models "github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
)

type Session struct {
	IsValid       bool
	PocketSession *models.Session
	Nodes         []*qos_models.QosNode
}

type SessionRegistryService interface {
	GetSession(req *models.GetSessionRequest) (*Session, error)
	GetNodesMap() map[qos_models.SessionChainKey]*ttlcache.Item[qos_models.SessionChainKey, []*qos_models.QosNode]
	GetNodesByChain(chainId string) []*qos_models.QosNode
}
