package config

import "time"

type EnvironmentStage string

const (
	StageProduction EnvironmentStage = "production"
)

type SecretProvider interface {
	GetPoktApplicationsEncryptionKey() string
	GetAPIKey() string
}

type DBCredentialsProvider interface {
	GetDatabaseConnectionUrl() string
}

type EnvironmentProvider interface {
	GetEnvironmentStage() EnvironmentStage
}

type PoktNodeConfigProvider interface {
	GetPoktRPCFullHost() string
	GetPoktRPCTimeout() time.Duration
}
