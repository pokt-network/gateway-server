package logging

import (
	"go.uber.org/zap"
	"os-gateway/internal/config"
)

func NewLogger(provider config.EnvironmentProvider) (*zap.Logger, error) {
	if provider.GetEnvironmentStage() == config.StageProduction {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
