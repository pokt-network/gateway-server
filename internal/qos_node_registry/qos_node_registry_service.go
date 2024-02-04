package qos_node_registry

import (
	"pokt_gateway_server/internal/qos_node_registry/models"
	"pokt_gateway_server/internal/session_registry"
)

type QosNodeRegistryService struct {
	sessionRegistry session_registry.SessionRegistryService
	qosNodes        []*models.QosNode
}

func (q QosNodeRegistryService) runChecks() {

}
