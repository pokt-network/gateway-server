package checks

import (
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
	nodeList      []*qos_models.QosNode
	pocketRelayer pokt_v0.PocketRelayer
}

func NewCheck(pocketRelayer pokt_v0.PocketRelayer) *Check {
	return &Check{pocketRelayer: pocketRelayer}
}
