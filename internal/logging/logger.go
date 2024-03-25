package logging

import (
	"go.uber.org/zap"
	"pokt_gateway_server/internal/global_config"
)

func NewLogger(provider global_config.EnvironmentProvider) (*zap.Logger, error) {
	if provider.GetEnvironmentStage() == global_config.StageProduction {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
