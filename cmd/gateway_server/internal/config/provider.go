package config

import (
	"pokt_gateway_server/internal/config"
)

type GatewayServerProvider interface {
	GetHTTPServerPort() uint
	config.DBCredentialsProvider
	config.PoktNodeConfigProvider
	config.SecretProvider
	config.EnvironmentProvider
}
