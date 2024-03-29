package main

import (
	"fmt"
	"github.com/fasthttp/router"
	fasthttpprometheus "github.com/flf2ko/fasthttp-prometheus"
	"github.com/jellydator/ttlcache/v3"
	"github.com/valyala/fasthttp"
	"pokt_gateway_server/cmd/gateway_server/internal/config"
	"pokt_gateway_server/cmd/gateway_server/internal/controllers"
	"pokt_gateway_server/cmd/gateway_server/internal/middleware"
	"pokt_gateway_server/internal/apps_registry"
	"pokt_gateway_server/internal/chain_configurations_registry"
	"pokt_gateway_server/internal/db_query"
	"pokt_gateway_server/internal/logging"
	"pokt_gateway_server/internal/node_selector_service"
	qos_models "pokt_gateway_server/internal/node_selector_service/models"
	"pokt_gateway_server/internal/relayer"
	"pokt_gateway_server/internal/session_registry"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
)

const (
	// Maximum amount of DB connections opened at a time. This should not have to be modified
	// as most of our database queries are periodic and not ran concurrently.
	maxDbConns = 50
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

	querier, pool, err := db_query.InitDB(logger, gatewayConfigProvider, maxDbConns)
	if err != nil {
		logger.Sugar().Fatal(err)
		return
	}

	// Close connection to pool afterward
	defer pool.Close()

	// Initialize a POKT client using the configured POKT RPC host and timeout
	client, err := pokt_v0.NewBasicClient(gatewayConfigProvider.GetPoktRPCFullHost(), gatewayConfigProvider.GetPoktRPCRequestTimeout())
	if err != nil {
		// If POKT client initialization fails, log the error and exit
		logger.Sugar().Fatal(err)
		return
	}

	// Initialize a TTL cache for session caching
	sessionCache := ttlcache.New[string, *session_registry.Session](
		ttlcache.WithTTL[string, *session_registry.Session](gatewayConfigProvider.GetSessionCacheTTL()),
	)

	nodeCache := ttlcache.New[qos_models.SessionChainKey, []*qos_models.QosNode](
		ttlcache.WithTTL[qos_models.SessionChainKey, []*qos_models.QosNode](gatewayConfigProvider.GetSessionCacheTTL()),
	)

	poktApplicationRegistry := apps_registry.NewCachedAppsRegistry(client, querier, gatewayConfigProvider, logger.Named("pokt_application_registry"))
	chainConfigurationRegistry := chain_configurations_registry.NewCachedChainConfigurationRegistry(querier, logger.Named("chain_configurations_registry"))
	sessionRegistry := session_registry.NewCachedSessionRegistryService(client, poktApplicationRegistry, sessionCache, nodeCache, logger.Named("session_registry"))
	nodeSelectorService := node_selector_service.NewNodeSelectorService(sessionRegistry, client, chainConfigurationRegistry, logger.Named("node_selector"))

	relayer := relayer.NewRelayer(client, sessionRegistry, poktApplicationRegistry, nodeSelectorService, chainConfigurationRegistry, gatewayConfigProvider, logger.Named("relayer"))

	// Define routers
	r := router.New()

	// Create a relay controller with the necessary dependencies (logger, registry, cached relayer)
	relayController := controllers.NewRelayController(relayer, logger.Named("relay_controller"))

	relayRouter := r.Group("/relay")
	relayRouter.POST("/{catchAll:*}", relayController.HandleRelay)

	poktAppsController := controllers.NewPoktAppsController(poktApplicationRegistry, querier, gatewayConfigProvider, logger.Named("pokt_apps_controller"))
	poktAppsRouter := r.Group("/poktapps")

	poktAppsRouter.GET("/", middleware.XAPIKeyAuth(poktAppsController.GetAll, gatewayConfigProvider))
	poktAppsRouter.POST("/", middleware.XAPIKeyAuth(poktAppsController.AddApplication, gatewayConfigProvider))
	poktAppsRouter.DELETE("/{app_id}", middleware.XAPIKeyAuth(poktAppsController.DeleteApplication, gatewayConfigProvider))

	// Create qos controller for debugging purposes
	qosNodeController := controllers.NewQosNodeController(sessionRegistry, logger.Named("qos_node_controller"))
	qosNodeRouter := r.Group("/qosnodes")
	qosNodeRouter.GET("/", middleware.XAPIKeyAuth(qosNodeController.GetAll, gatewayConfigProvider))

	// Add Middleware for Generic E2E Prom Tracking
	p := fasthttpprometheus.NewPrometheus("fasthttp")
	fastpHandler := p.WrapHandler(r)

	logger.Info("Gateway Server Started")
	// Start the fasthttp server and listen on the configured server port
	if err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", gatewayConfigProvider.GetHTTPServerPort()), fastpHandler); err != nil {
		// If an error occurs during server startup, log the error and exit
		logger.Sugar().Fatalw("Error in ListenAndServe", "err", err)
	}
}
