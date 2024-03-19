package models

import (
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"time"
)

type TimeoutReason string

const (
	OutOfSyncTimeout   TimeoutReason = "out_of_sync_timeout"
	InvalidDataTimeout TimeoutReason = "invalid_data_timeout"
)

// QosNode a FAT model to store the QoS information of a specific node in a session.
type QosNode struct {
	MorseNode              *models.Node
	PocketSession          *models.Session
	Signer                 *models.Ed25519Account
	p90Latency             float64
	timeoutUntil           time.Time
	timeoutReason          TimeoutReason
	lastDataIntegrityCheck time.Time
	latestKnownHeight      uint64
	synced                 bool
	lastHeightCheckTime    time.Time
}

func (n QosNode) IsHealthy() bool {
	return !n.isInTimeout() && n.IsSynced()
}

func (n QosNode) IsSynced() bool {
	return n.synced
}

func (n QosNode) SetSynced(synced bool) {
	n.synced = synced
}

func (n QosNode) isInTimeout() bool {
	return !n.timeoutUntil.IsZero() && time.Now().Before(n.timeoutUntil)
}

func (n QosNode) GetLastHeightCheckTime() time.Time {
	return n.lastHeightCheckTime
}

func (n QosNode) SetTimeoutUntil(time time.Time, reason TimeoutReason) {
	n.timeoutReason = reason
	n.timeoutUntil = time
}

func (n QosNode) SetLastKnownHeight(lastKnownHeight uint64) {
	n.latestKnownHeight = lastKnownHeight
}

func (n QosNode) SetLastHeightCheckTime(time time.Time) {
	n.lastHeightCheckTime = time
}

func (n QosNode) GetLastKnownHeight() uint64 {
	return n.latestKnownHeight
}
