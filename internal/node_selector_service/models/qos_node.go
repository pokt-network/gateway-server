package models

import (
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"sync"
	"time"
)

type TimeoutReason string

const (
	maxErrorStr int = 100
)

const (
	chainSolanaCustom = "C006"
	chainSolana       = "0006"
)
const (
	OutOfSyncTimeout     TimeoutReason = "out_of_sync_timeout"
	DataIntegrityTimeout TimeoutReason = "invalid_data_timeout"
	MaximumRelaysTimeout TimeoutReason = "maximum_relays_timeout"
	NodeResponseTimeout  TimeoutReason = "node_response_timeout"
)

type LatencyTracker struct {
	measurements []float64
	lock         sync.RWMutex
}

func (l *LatencyTracker) GetP90Latency() float64 {
	return -1.0
}

func (l *LatencyTracker) RecordMeasurement(time float64) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.measurements = append(l.measurements, time)
}

// QosNode a FAT model to store the QoS information of a specific node in a session.
type QosNode struct {
	MorseNode                  *models.Node
	PocketSession              *models.Session
	AppSigner                  *models.Ed25519Account
	LatencyTracker             *LatencyTracker
	timeoutUntil               time.Time
	timeoutReason              TimeoutReason
	lastDataIntegrityCheckTime time.Time
	latestKnownHeight          uint64
	synced                     bool
	lastKnownError             error
	lastHeightCheckTime        time.Time
}

func NewQosNode(morseNode *models.Node, pocketSession *models.Session, appSigner *models.Ed25519Account) *QosNode {
	return &QosNode{MorseNode: morseNode, PocketSession: pocketSession, AppSigner: appSigner, LatencyTracker: &LatencyTracker{measurements: []float64{}}}
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

func (n *QosNode) SetTimeoutUntil(time time.Time, reason TimeoutReason, attachedErr error) {
	n.timeoutReason = reason
	n.timeoutUntil = time
	n.lastKnownError = attachedErr
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

func (n *QosNode) IsSolanaChain() bool {
	chainId := n.GetChain()
	return chainId == chainSolana || chainId == chainSolanaCustom
}

func (n *QosNode) IsEvmChain() bool {
	return !n.IsSolanaChain()
}

func (n *QosNode) GetTimeoutReason() TimeoutReason {
	return n.timeoutReason
}

func (n *QosNode) GetLastKnownErrorStr() string {
	if n.lastKnownError == nil {
		return ""
	}
	errStr := n.lastKnownError.Error()
	if len(errStr) > maxErrorStr {
		return errStr[:maxErrorStr]
	}
	return errStr
}

func (n *QosNode) GetTimeoutUntil() time.Time {
	return n.timeoutUntil
}

func (n *QosNode) GetLatencyTracker() *LatencyTracker {
	return n.LatencyTracker
}
