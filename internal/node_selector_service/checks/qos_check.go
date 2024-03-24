package checks

import (
	"pokt_gateway_server/internal/chain_configurations_registry"
	qos_models "pokt_gateway_server/internal/node_selector_service/models"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
)

type CheckJob interface {
	Perform()
	Name() string
	ShouldRun() bool
	SetNodes(nodes []*qos_models.QosNode)
}

type Check struct {
	NodeList           []*qos_models.QosNode
	PocketRelayer      pokt_v0.PocketRelayer
	ChainConfiguration chain_configurations_registry.ChainConfigurationsService
}

func NewCheck(pocketRelayer pokt_v0.PocketRelayer, chainConfiguration chain_configurations_registry.ChainConfigurationsService) *Check {
	return &Check{PocketRelayer: pocketRelayer}
}
