package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"pokt_gateway_server/internal/config"
	"strconv"
	"time"
)

// Environment variable names
const (
	poktRPCFullHostEnv               = "POKT_RPC_FULL_HOST"
	httpServerPortEnv                = "HTTP_SERVER_PORT"
	poktRPCTimeoutEnv                = "POKT_RPC_TIMEOUT"
	dbConnectionUrlEnv               = "DB_CONNECTION_URL"
	sessionCacheTTLEnv               = "SESSION_CACHE_TTL"
	environmentStageEnv              = "ENVIRONMENT_STAGE"
	poktApplicationsEncryptionKeyEnv = "POKT_APPLICATIONS_ENCRYPTION_KEY"
)

// DotEnvConfigProvider implements the GatewayServerProvider interface.
type DotEnvConfigProvider struct {
	poktRPCFullHost               string
	httpServerPort                uint
	poktRPCTimeout                time.Duration
	sessionCacheTTL               time.Duration
	environmentStage              config.EnvironmentStage
	poktApplicationsEncryptionKey string
	databaseConnectionUrl         string
}

// GetPoktRPCFullHost returns the PoktRPCFullHost value.
func (c DotEnvConfigProvider) GetPoktRPCFullHost() string {
	return c.poktRPCFullHost
}

// GetHTTPServerPort returns the HTTPServerPort value.
func (c DotEnvConfigProvider) GetHTTPServerPort() uint {
	return c.httpServerPort
}

// GetPoktRPCTimeout returns the PoktRPCTimeout value.
func (c DotEnvConfigProvider) GetPoktRPCTimeout() time.Duration {
	return c.poktRPCTimeout
}

// GetSessionCacheTTL returns the time value for session to expire in cache.
func (c DotEnvConfigProvider) GetSessionCacheTTL() time.Duration {
	return c.poktRPCTimeout
}

// GetEnvironmentStage returns the EnvironmentStage value.
func (c DotEnvConfigProvider) GetEnvironmentStage() config.EnvironmentStage {
	return c.environmentStage
}

// GetPoktApplicationsEncryptionKey: Key used to decrypt pokt applications private key.
func (c DotEnvConfigProvider) GetPoktApplicationsEncryptionKey() string {
	return c.poktApplicationsEncryptionKey
}

// GetDatabaseConnectionUrl returns the PoktRPCFullHost value.
func (c DotEnvConfigProvider) GetDatabaseConnectionUrl() string {
	return c.databaseConnectionUrl
}

// NewDotEnvConfigProvider creates a new instance of DotEnvConfigProvider.
func NewDotEnvConfigProvider() *DotEnvConfigProvider {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %s", err))
	}

	poktRPCTimeout, err := time.ParseDuration(getEnvVar(poktRPCTimeoutEnv))
	if err != nil {
		panic(fmt.Sprintf("Error parsing %s: %s", poktRPCTimeoutEnv, err))
	}

	httpServerPort, err := strconv.ParseUint(getEnvVar(httpServerPortEnv), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Error parsing %s: %s", httpServerPortEnv, err))
	}

	sessionCacheTTLDuration, err := time.ParseDuration(getEnvVar(sessionCacheTTLEnv))
	if err != nil {
		panic(fmt.Sprintf("Error parsing %s: %s", sessionCacheTTLDuration, err))
	}

	return &DotEnvConfigProvider{
		poktRPCFullHost:               getEnvVar(poktRPCFullHostEnv),
		httpServerPort:                uint(httpServerPort),
		poktRPCTimeout:                poktRPCTimeout,
		sessionCacheTTL:               sessionCacheTTLDuration,
		databaseConnectionUrl:         getEnvVar(dbConnectionUrlEnv),
		environmentStage:              config.EnvironmentStage(getEnvVar(environmentStageEnv)),
		poktApplicationsEncryptionKey: getEnvVar(poktApplicationsEncryptionKeyEnv),
	}
}

// getEnvVar retrieves the value of the environment variable with error handling.
func getEnvVar(name string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		panic(fmt.Errorf("%s not set", name))
	}
	return value
}
