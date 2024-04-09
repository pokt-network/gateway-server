package config

import (
	"github.com/pokt-network/gateway-server/internal/global_config"
)

type GatewayServerProvider interface {
	GetHTTPServerPort() uint
	global_config.DBCredentialsProvider
	global_config.PoktNodeConfigProvider
	global_config.SecretProvider
	global_config.EnvironmentProvider
}
