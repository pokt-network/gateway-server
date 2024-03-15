package checks

import (
	"pokt_gateway_server/internal/qos_node_registry/models"
	"time"
)

type QosJob interface {
	PerformJob()
	ShouldRun() bool
}

type Check struct {
	LastChecked time.Time
	NodeList    []*models.QosNode
	ChainId     string
}
