package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os-gateway/internal/config"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"strconv"
	"strings"
	"time"
)

// Environment variable names
const (
	poktRPCFullHostEnv  = "POKT_RPC_FULL_HOST"
	httpServerPortEnv   = "HTTP_SERVER_PORT"
	poktRPCTimeoutEnv   = "POKT_RPC_TIMEOUT"
	environmentStageEnv = "ENVIRONMENT_STAGE"
	// TODO: Use an encrypted key file or move over to a encrypted DB. (@Blade)
	appStakesPrivateKeyEnv = "APPSTAKE_PRIVATE_KEYS"
)

// DotEnvConfigProvider implements the GatewayServerProvider interface.
type DotEnvConfigProvider struct {
	poktRPCFullHost  string
	httpServerPort   uint
	poktRPCTimeout   time.Duration
	environmentStage config.EnvironmentStage
	appStakes        []*models.Ed25519Account
}

// GetPoktRPCFullHost returns the PoktRPCFullHost value.
func (c *DotEnvConfigProvider) GetPoktRPCFullHost() string {
	return c.poktRPCFullHost
}

// GetHTTPServerPort returns the HTTPServerPort value.
func (c *DotEnvConfigProvider) GetHTTPServerPort() uint {
	return c.httpServerPort
}

// GetPoktRPCTimeout returns the PoktRPCTimeout value.
func (c *DotEnvConfigProvider) GetPoktRPCTimeout() time.Duration {
	return c.poktRPCTimeout
}

// GetEnvironmentStage returns the EnvironmentStage value.
func (c DotEnvConfigProvider) GetEnvironmentStage() config.EnvironmentStage {
	return c.environmentStage
}

// GetAppStakes returns the app stakes for sending a relay
func (c *DotEnvConfigProvider) GetAppStakes() []*models.Ed25519Account {
	return c.appStakes
}

// NewDotEnvConfigProvider creates a new instance of DotEnvConfigProvider.
func NewDotEnvConfigProvider() *DotEnvConfigProvider {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("Error loading .env file: %s", err))
	}

	poktRPCFullHost, err := getEnvVar(poktRPCFullHostEnv)
	if err != nil {
		panic(fmt.Sprintf("Error getting %s: %s", poktRPCFullHostEnv, err))
	}

	httpServerPortStr, err := getEnvVar(httpServerPortEnv)
	if err != nil {
		panic(fmt.Sprintf("Error getting %s: %s", httpServerPortEnv, err))
	}

	poktRPCTimeoutStr, err := getEnvVar(poktRPCTimeoutEnv)
	if err != nil {
		panic(fmt.Sprintf("Error getting %s: %s", poktRPCTimeoutEnv, err))
	}

	environmentStage, err := getEnvVar(environmentStageEnv)
	if err != nil {
		panic(fmt.Sprintf("Error getting %s: %s", environmentStageEnv, err))
	}

	httpServerPort, err := strconv.ParseUint(httpServerPortStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Error parsing %s: %s", httpServerPortEnv, err))
	}

	poktRPCTimeout, err := time.ParseDuration(poktRPCTimeoutStr)
	if err != nil {
		panic(fmt.Sprintf("Error parsing %s: %s", poktRPCTimeoutEnv, err))
	}

	return &DotEnvConfigProvider{
		poktRPCFullHost:  poktRPCFullHost,
		httpServerPort:   uint(httpServerPort),
		poktRPCTimeout:   poktRPCTimeout,
		environmentStage: config.EnvironmentStage(environmentStage),
		appStakes:        getAppStakesFromEnv(),
	}
}

// getEnvVar retrieves the value of the environment variable with error handling.
func getEnvVar(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("%s not set", name)
	}
	return value, nil
}

func getAppStakesFromEnv() []*models.Ed25519Account {
	privateKeys, err := getEnvVar(appStakesPrivateKeyEnv)
	if err != nil {
		panic(fmt.Sprintf("Error parsing %s", appStakesPrivateKeyEnv))
	}
	var appStakePrivateKeys []*models.Ed25519Account
	for _, key := range strings.Split(privateKeys, ",") {
		appStake, err := models.NewAccount(key)
		if err != nil {
			panic("Failed to parse appstake key")
		}
		appStakePrivateKeys = append(appStakePrivateKeys, appStake)
	}

	if len(appStakePrivateKeys) == 0 {
		panic("app stakes were not provided or unable to parse successfully")
	}

	return appStakePrivateKeys
}
