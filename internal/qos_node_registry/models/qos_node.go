package models

import (
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"time"
)

type QosNode struct {
	MorseNode         *models.Node
	Signer            *models.Node
	sessionHeight     uint64
	chainId           uint64
	p90Latency        float64
	timeoutUntil      time.Time
	timeoutReason     string
	latestKnownHeight uint64
	synced            bool
}

func NewQosNode(morseNode *models.Node) *QosNode {
	return &QosNode{MorseNode: morseNode}
}

func (n QosNode) IsHealthy() bool {
	return n.isInSync() && !n.isInTimeout()
}

func (n QosNode) isInTimeout() bool {
	return !n.timeoutUntil.IsZero() && time.Now().Before(n.timeoutUntil)
}

func (n QosNode) isInSync() bool {
	return n.synced
}

func (n QosNode) SetTimeoutUntil(time time.Time) {
	n.timeoutUntil = time
}

func (n QosNode) SetLastKnownHeight(lastKnownHeight uint64) {
	n.latestKnownHeight = lastKnownHeight
}

func (n QosNode) GetLastKnownHeight() uint64 {
	return n.latestKnownHeight
}
