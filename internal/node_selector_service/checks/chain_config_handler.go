package checks

import "pokt_gateway_server/internal/chain_configurations_registry"

// getBlockHeightTolerance - helper function to retrieve block height tolerance across checks
func getBlockHeightTolerance(chainConfiguration chain_configurations_registry.ChainConfigurationsService, chainId string, defaultValue int) int {
	chainConfig, ok := chainConfiguration.GetChainConfiguration(chainId)
	if !ok {
		return defaultValue
	}
	return int(*chainConfig.HeightCheckBlockTolerance)
}

// getDataIntegrityHeightLookback - helper function ro retrieve data integrity lookback across checks
func getDataIntegrityHeightLookback(chainConfiguration chain_configurations_registry.ChainConfigurationsService, chainId string, defaultValue int) int {
	chainConfig, ok := chainConfiguration.GetChainConfiguration(chainId)
	if !ok {
		return defaultValue
	}
	return int(*chainConfig.DataIntegrityCheckLookbackHeight)
}
