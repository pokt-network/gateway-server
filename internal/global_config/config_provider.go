package global_config

import (
	"github.com/pokt-network/gateway-server/internal/chain_network"
	"time"
)

type EnvironmentStage string

const (
	StageProduction EnvironmentStage = "production"
)

type GlobalConfigProvider interface {
	SecretProvider
	DBCredentialsProvider
	EnvironmentProvider
	PoktNodeConfigProvider
	AltruistConfigProvider
	PromMetricsProvider
	ChainNetworkProvider
}

type PromMetricsProvider interface {
	ShouldEmitServiceUrlPromMetrics() bool
}

type SecretProvider interface {
	GetPoktApplicationsEncryptionKey() string
	GetAPIKey() string
}

type DBCredentialsProvider interface {
	GetDatabaseConnectionUrl() string
}

type EnvironmentProvider interface {
	GetEnvironmentStage() EnvironmentStage
}

type PoktNodeConfigProvider interface {
	GetPoktRPCFullHost() string
	GetPoktRPCRequestTimeout() time.Duration
}

type AltruistConfigProvider interface {
	GetAltruistRequestTimeout() time.Duration
}

type ChainNetworkProvider interface {
	GetChainNetwork() chain_network.ChainNetwork
}
