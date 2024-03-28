package config

import (
	"pokt_gateway_server/internal/global_config"
)

type GatewayServerProvider interface {
	GetHTTPServerPort() uint
	global_config.DBCredentialsProvider
	global_config.PoktNodeConfigProvider
	global_config.SecretProvider
	global_config.EnvironmentProvider
}
