package chain_configurations_registry

import "pokt_gateway_server/internal/db_query"

type ChainConfigurationsService interface {
	GetChainConfiguration(chainId string) (db_query.GetChainConfigurationsRow, bool)
}
