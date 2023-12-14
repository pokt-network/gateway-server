package main

import (
	"fmt"
	"log"
	"os-gateway/cmd/gateway_server/internal/config"
	"os-gateway/cmd/gateway_server/internal/controllers"
	"os-gateway/internal/logging"
	"os-gateway/internal/pokt_client_decorators"
	"os-gateway/pkg/pokt/pokt_v0"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"os-gateway/pkg/ttl_cache"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func main() {
	// Initialize configuration provider from environment variables
	gatewayConfigProvider := config.NewDotEnvConfigProvider()

	// Initialize logger using the configured settings
	logger, err := logging.NewLogger(gatewayConfigProvider)
	if err != nil {
		// If logger initialization fails, panic with the error
		panic(err)
	}

	// Initialize a POKT client using the configured POKT RPC host and timeout
	client, err := pokt_v0.NewBasicClient(gatewayConfigProvider.GetPoktRPCFullHost(), gatewayConfigProvider.GetPoktRPCTimeout())
	if err != nil {
		// If POKT client initialization fails, log the error and exit
		logger.Sugar().Fatal(err)
		return
	}

	// Initialize a TTL cache for session caching
	sessionCache := ttl_cache.NewTTLCacheClient[string, *models.GetSessionResponse]() // Initialize the cache client

	sessionCache.Start() // Start the cache client

	// Create a relay controller with a cached POKT client and the loggeri
	relayController := controllers.NewRelayController(pokt_client_decorators.NewCachedClient(client, sessionCache), gatewayConfigProvider.GetAppStakes(), logger)

	// Define routers
	r := router.New()
	r.POST(controllers.RelayHandlerPath, relayController.HandleRelay)

	logger.Info("Gateway Server Started")
	// Start the fasthttp server and listen on the configured server port
	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", gatewayConfigProvider.GetHTTPServerPort()), r.Handler); err != nil {
		// If an error occurs during server startup, log the error and exit
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
