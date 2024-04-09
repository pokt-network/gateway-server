package models

import "time"

type PublicQosNode struct {
	NodePublicKey   string    `json:"node_public_key"`
	ServiceUrl      string    `json:"service_url"`
	Chain           string    `json:"chain"`
	SessionHeight   uint      `json:"session_height"`
	AppPublicKey    string    `json:"app_public_key"`
	TimeoutUntil    time.Time `json:"timeout_until"`
	TimeoutReason   string    `json:"timeout_reason"`
	LastKnownErr    string    `json:"last_known_err"`
	IsHealthy       bool      `json:"is_healthy"`
	IsSynced        bool      `json:"is_synced"`
	LastKnownHeight uint64    `json:"last_known_height"`
	P90Latency      float64   `json:"p90_latency"`
}
