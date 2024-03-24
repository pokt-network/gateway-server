package global_config

import "time"

type EnvironmentStage string

const (
	StageProduction EnvironmentStage = "production"
)

type GlobalConfigProvider interface {
	SecretProvider
	DBCredentialsProvider
	EnvironmentProvider
	PoktNodeConfigProvider
	AltruistConfigProvider
}

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
	GetPoktRPCRequestTimeout() time.Duration
}

type AltruistConfigProvider interface {
	GetAltruistRequestTimeout() time.Duration
}
