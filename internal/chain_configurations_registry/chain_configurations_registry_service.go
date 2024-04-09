package chain_configurations_registry

import "github.com/pokt-network/gateway-server/internal/db_query"

type ChainConfigurationsService interface {
	GetChainConfiguration(chainId string) (db_query.GetChainConfigurationsRow, bool)
}
