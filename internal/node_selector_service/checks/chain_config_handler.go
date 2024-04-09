package checks

import "github.com/pokt-network/gateway-server/internal/chain_configurations_registry"

// GetBlockHeightTolerance - helper function to retrieve block height tolerance across checks
func GetBlockHeightTolerance(chainConfiguration chain_configurations_registry.ChainConfigurationsService, chainId string, defaultValue int) int {
	chainConfig, ok := chainConfiguration.GetChainConfiguration(chainId)
	if !ok {
		return defaultValue
	}
	return int(*chainConfig.HeightCheckBlockTolerance)
}

// GetDataIntegrityHeightLookback - helper function ro retrieve data integrity lookback across checks
func GetDataIntegrityHeightLookback(chainConfiguration chain_configurations_registry.ChainConfigurationsService, chainId string, defaultValue int) int {
	chainConfig, ok := chainConfiguration.GetChainConfiguration(chainId)
	if !ok {
		return defaultValue
	}
	return int(*chainConfig.DataIntegrityCheckLookbackHeight)
}
