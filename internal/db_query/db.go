package db_query

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	root "github.com/pokt-network/gateway-server"
	"github.com/pokt-network/gateway-server/internal/global_config"
	"go.uber.org/zap"
)

// InitDB - runs DB migrations and provides a code-generated query interface
func InitDB(logger *zap.Logger, config global_config.DBCredentialsProvider, maxConnections uint) (Querier, *pgxpool.Pool, error) {

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
	pool, err := pgxpool.Connect(context.Background(), fmt.Sprintf("%s&pool_max_conns=%d", config.GetDatabaseConnectionUrl(), maxConnections))
	if err != nil {
		return nil, pool, err
	}
	return NewQuerier(pool), nil, nil

}
