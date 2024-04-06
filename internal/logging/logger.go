package logging

import (
	"github.com/pokt-network/gateway-server/internal/global_config"
	"go.uber.org/zap"
)

func NewLogger(provider global_config.EnvironmentProvider) (*zap.Logger, error) {
	if provider.GetEnvironmentStage() == global_config.StageProduction {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
