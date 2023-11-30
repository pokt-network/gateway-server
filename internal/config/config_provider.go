package config

type EnvironmentStage string

const (
	StageProduction EnvironmentStage = "production"
)

type EnvironmentProvider interface {
	GetEnvironmentStage() EnvironmentStage
}
