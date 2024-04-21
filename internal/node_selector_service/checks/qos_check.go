package checks

import (
	"github.com/pokt-network/gateway-server/internal/chain_configurations_registry"
	"github.com/pokt-network/gateway-server/internal/chain_network"
	config2 "github.com/pokt-network/gateway-server/internal/global_config"
	qos_models "github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0"
)

const (
	chainMorseMainnetSolanaCustom = "C006"
	chainMorseMainnetPokt         = "0001"
	chainMorseMainnetSolana       = "0006"
)

const (
	chainMorseTestnetPokt   = "0013"
	chainMorseTestnetSolana = "0008"
)

type CheckJob interface {
	Perform()
	Name() string
	ShouldRun() bool
	SetNodes(nodes []*qos_models.QosNode)
}

type Check struct {
	NodeList             []*qos_models.QosNode
	PocketRelayer        pokt_v0.PocketRelayer
	ChainConfiguration   chain_configurations_registry.ChainConfigurationsService
	ChainNetworkProvider config2.ChainNetworkProvider
}

func NewCheck(pocketRelayer pokt_v0.PocketRelayer, chainConfiguration chain_configurations_registry.ChainConfigurationsService, chainNetworkProvider config2.ChainNetworkProvider) *Check {
	return &Check{PocketRelayer: pocketRelayer, ChainConfiguration: chainConfiguration, ChainNetworkProvider: chainNetworkProvider}
}

func (c *Check) IsSolanaChain(node *qos_models.QosNode) bool {
	chainId := node.GetChain()
	if c.ChainNetworkProvider.GetChainNetwork() == chain_network.MorseTestnet {
		return chainId == chainMorseTestnetSolana
	}
	return chainId == chainMorseMainnetSolana || chainId == chainMorseMainnetSolanaCustom
}

func (c *Check) IsPoktChain(node *qos_models.QosNode) bool {
	chainId := node.GetChain()
	if c.ChainNetworkProvider.GetChainNetwork() == chain_network.MorseTestnet {
		return chainId == chainMorseTestnetPokt
	}
	return chainId == chainMorseMainnetPokt
}

func (c *Check) IsEvmChain(node *qos_models.QosNode) bool {
	return !c.IsPoktChain(node) && !c.IsSolanaChain(node)
}
