package config

import (
	"os-gateway/internal/config"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"time"
)

type GatewayServerProvider interface {
	GetPoktRPCFullHost() string
	GetHTTPServerPort() uint
	GetPoktRPCTimeout() time.Duration
	GetAppStakes() []*models.Ed25519Account
	config.EnvironmentProvider
}
