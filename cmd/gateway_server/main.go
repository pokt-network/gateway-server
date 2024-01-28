package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fasthttp/router"
	fasthttpprometheus "github.com/flf2ko/fasthttp-prometheus"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jellydator/ttlcache/v3"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	root "pokt_gateway_server"
	"pokt_gateway_server/cmd/gateway_server/internal/config"
	"pokt_gateway_server/cmd/gateway_server/internal/controllers"
	"pokt_gateway_server/internal/db_query"
	"pokt_gateway_server/internal/logging"
	"pokt_gateway_server/internal/pokt_applications_registry"
	"pokt_gateway_server/internal/pokt_client_decorators"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
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

	querier, pool, err := initDB(logger, gatewayConfigProvider)
	if err != nil {
		logger.Sugar().Fatal(err)
		return
	}

	// Close connection to pool afterward
	defer pool.Close()

	// Initialize a POKT client using the configured POKT RPC host and timeout
	client, err := pokt_v0.NewBasicClient(gatewayConfigProvider.GetPoktRPCFullHost(), gatewayConfigProvider.GetPoktRPCTimeout())
	if err != nil {
		// If POKT client initialization fails, log the error and exit
		logger.Sugar().Fatal(err)
		return
	}

	// Initialize a TTL cache for session caching
	sessionCache := ttlcache.New[string, *models.GetSessionResponse](
		ttlcache.WithTTL[string, *models.GetSessionResponse](gatewayConfigProvider.GetSessionCacheTTL()), //@todo: make this configurable via env ?
	)

	poktApplicationRegistry := pokt_applications_registry.NewCachedRegistry(client, querier, gatewayConfigProvider, logger.Named("pokt_application_registry"))

	// Create a relay controller with the necessary dependencies (logger, registry, cached relayer)
	relayController := controllers.NewRelayController(pokt_client_decorators.NewCachedClient(client, sessionCache), poktApplicationRegistry, logger)

	// Define routers
	r := router.New()
	r.POST(controllers.RelayHandlerPath, relayController.HandleRelay)

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

func initDB(logger *zap.Logger, config config.GatewayServerProvider) (db_query.Querier, *pgxpool.Pool, error) {

	// initialize database
	sqldb, err := sql.Open("postgres", config.GetDatabaseConnectionUrl())
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to init db")
	}

	postgresDriver, err := postgres.WithInstance(sqldb, &postgres.Config{})
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to init postgres driver")
	}

	source, err := iofs.New(root.Migrations, "db_migrations")

	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to create migration fs")
	}

	// Automatic Migrations
	m, err := migrate.NewWithInstance("iofs", source, "postgres", postgresDriver)

	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to migrate")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Sugar().Warn("Migration warning", "err", err)
		return nil, nil, err
	}

	// DB only needs to be open for migration
	err = postgresDriver.Close()
	if err != nil {
		return nil, nil, err
	}
	err = sqldb.Close()
	if err != nil {
		return nil, nil, err
	}

	// open up connection pool for actual sql queries
	pool, err := pgxpool.Connect(context.Background(), config.GetDatabaseConnectionUrl())
	if err != nil {
		return nil, pool, err
	}
	return db_query.NewQuerier(pool), nil, nil

}
