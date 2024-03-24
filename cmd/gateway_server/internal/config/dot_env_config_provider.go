package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"pokt_gateway_server/internal/global_config"
	"strconv"
	"time"
)

// Environment variable names
const (
	poktRPCFullHostEnv               = "POKT_RPC_FULL_HOST"
	httpServerPortEnv                = "HTTP_SERVER_PORT"
	poktRPCTimeoutEnv                = "POKT_RPC_TIMEOUT"
	altruistRequestTimeoutEnv        = "ALTRUIST_REQUEST_TIMEOUT"
	dbConnectionUrlEnv               = "DB_CONNECTION_URL"
	sessionCacheTTLEnv               = "SESSION_CACHE_TTL"
	environmentStageEnv              = "ENVIRONMENT_STAGE"
	poktApplicationsEncryptionKeyEnv = "POKT_APPLICATIONS_ENCRYPTION_KEY"
	apiKey                           = "API_KEY"
)

// DotEnvGlobalConfigProvider implements the GatewayServerProvider interface.
type DotEnvGlobalConfigProvider struct {
	poktRPCFullHost               string
	httpServerPort                uint
	poktRPCRequestTimeout         time.Duration
	sessionCacheTTL               time.Duration
	environmentStage              global_config.EnvironmentStage
	poktApplicationsEncryptionKey string
	databaseConnectionUrl         string
	apiKey                        string
	altruistRequestTimeout        time.Duration
}

func (c DotEnvGlobalConfigProvider) GetAPIKey() string {
	return c.apiKey
}

// GetPoktRPCFullHost returns the PoktRPCFullHost value.
func (c DotEnvGlobalConfigProvider) GetPoktRPCFullHost() string {
	return c.poktRPCFullHost
}

// GetHTTPServerPort returns the HTTPServerPort value.
func (c DotEnvGlobalConfigProvider) GetHTTPServerPort() uint {
	return c.httpServerPort
}

// GetPoktRPCTimeout returns the PoktRPCTimeout value.
func (c DotEnvGlobalConfigProvider) GetPoktRPCRequestTimeout() time.Duration {
	return c.poktRPCRequestTimeout
}

// GetSessionCacheTTL returns the time value for session to expire in cache.
func (c DotEnvGlobalConfigProvider) GetSessionCacheTTL() time.Duration {
	return c.poktRPCRequestTimeout
}

// GetEnvironmentStage returns the EnvironmentStage value.
func (c DotEnvGlobalConfigProvider) GetEnvironmentStage() global_config.EnvironmentStage {
	return c.environmentStage
}

// GetPoktApplicationsEncryptionKey: Key used to decrypt pokt applications private key.
func (c DotEnvGlobalConfigProvider) GetPoktApplicationsEncryptionKey() string {
	return c.poktApplicationsEncryptionKey
}

// GetDatabaseConnectionUrl returns the PoktRPCFullHost value.
func (c DotEnvGlobalConfigProvider) GetDatabaseConnectionUrl() string {
	return c.databaseConnectionUrl
}

// GetDatabaseConnectionUrl returns the PoktRPCFullHost value.
func (c DotEnvGlobalConfigProvider) GetAltruistRequestTimeout() time.Duration {
	return c.altruistRequestTimeout
}

// NewDotEnvConfigProvider creates a new instance of DotEnvGlobalConfigProvider.
func NewDotEnvConfigProvider() *DotEnvGlobalConfigProvider {
	_ = godotenv.Load()

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

	altruistRequestTimeoutDuration, err := time.ParseDuration(getEnvVar(altruistRequestTimeoutEnv))
	if err != nil {
		panic(fmt.Sprintf("Error parsing %s: %s", sessionCacheTTLDuration, err))
	}
	return &DotEnvGlobalConfigProvider{
		poktRPCFullHost:               getEnvVar(poktRPCFullHostEnv),
		httpServerPort:                uint(httpServerPort),
		poktRPCRequestTimeout:         poktRPCTimeout,
		sessionCacheTTL:               sessionCacheTTLDuration,
		databaseConnectionUrl:         getEnvVar(dbConnectionUrlEnv),
		environmentStage:              global_config.EnvironmentStage(getEnvVar(environmentStageEnv)),
		poktApplicationsEncryptionKey: getEnvVar(poktApplicationsEncryptionKeyEnv),
		apiKey:                        getEnvVar(apiKey),
		altruistRequestTimeout:        altruistRequestTimeoutDuration,
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
