package checks

import (
	"pokt_gateway_server/internal/node_selector_service/models"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"time"
)

type CheckJob interface {
	Perform()
	Name() string
	ShouldRun() bool
}

type Check struct {
	nextCheckTime time.Time
	nodeList      []*models.QosNode
	pocketRelayer pokt_v0.PocketRelayer
}

func NewCheck(nodeList []*models.QosNode, pocketRelayer pokt_v0.PocketRelayer) *Check {
	return &Check{nodeList: nodeList, pocketRelayer: pocketRelayer, nextCheckTime: time.Time{}}
}
