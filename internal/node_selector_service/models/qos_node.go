package models

import (
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"time"
)

type TimeoutReason string

const (
	OutOfSyncTimeout     TimeoutReason = "out_of_sync_timeout"
	DataIntegrityTimeout TimeoutReason = "invalid_data_timeout"
	MaximumRelaysTimeout TimeoutReason = "maximum_relays_timeout"
	NodeResponseTimeout  TimeoutReason = "node_response_timeout"
)

// QosNode a FAT model to store the QoS information of a specific node in a session.
type QosNode struct {
	MorseNode                  *models.Node
	PocketSession              *models.Session
	AppSigner                  *models.Ed25519Account
	p90Latency                 float64
	timeoutUntil               time.Time
	timeoutReason              TimeoutReason
	lastDataIntegrityCheckTime time.Time
	latestKnownHeight          uint64
	synced                     bool
	lastHeightCheckTime        time.Time
}

func (n *QosNode) IsHealthy() bool {
	return !n.isInTimeout() && n.IsSynced()
}

func (n *QosNode) IsSynced() bool {
	return n.synced
}

func (n *QosNode) SetSynced(synced bool) {
	n.synced = synced
}

func (n *QosNode) isInTimeout() bool {
	return !n.timeoutUntil.IsZero() && time.Now().Before(n.timeoutUntil)
}

func (n *QosNode) GetLastHeightCheckTime() time.Time {
	return n.lastHeightCheckTime
}

func (n *QosNode) SetTimeoutUntil(time time.Time, reason TimeoutReason) {
	n.timeoutReason = reason
	n.timeoutUntil = time
}

func (n *QosNode) SetLastKnownHeight(lastKnownHeight uint64) {
	n.latestKnownHeight = lastKnownHeight
}

func (n *QosNode) SetLastHeightCheckTime(time time.Time) {
	n.lastHeightCheckTime = time
}

func (n *QosNode) GetLastKnownHeight() uint64 {
	return n.latestKnownHeight
}

func (n *QosNode) GetChain() string {
	return n.PocketSession.SessionHeader.Chain
}

func (n *QosNode) GetPublicKey() string {
	return n.MorseNode.PublicKey
}

func (n *QosNode) GetAppStakeSigner() *models.Ed25519Account {
	return n.AppSigner
}

func (n *QosNode) GetLastDataIntegrityCheckTime() time.Time {
	return n.lastDataIntegrityCheckTime
}
func (n *QosNode) SetLastDataIntegrityCheckTime(lastDataIntegrityCheckTime time.Time) {
	n.lastDataIntegrityCheckTime = lastDataIntegrityCheckTime
}
