package node_selector_service

import (
	"go.uber.org/zap"
	"pokt_gateway_server/internal/chain_configurations_registry"
	"pokt_gateway_server/internal/node_selector_service/checks"
	"pokt_gateway_server/internal/node_selector_service/models"
	"pokt_gateway_server/internal/session_registry"
	"pokt_gateway_server/pkg/common"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"time"
)

const (
	jobCheckInterval = time.Second
)

type NodeSelectorService interface {
	FindNode(chainId string) (*models.QosNode, bool)
}

type NodeSelectorClient struct {
	sessionRegistry session_registry.SessionRegistryService
	pocketRelayer   pokt_v0.PocketRelayer
	logger          *zap.Logger
	checkJobs       []checks.CheckJob
}

func NewNodeSelectorService(sessionRegistry session_registry.SessionRegistryService, pocketRelayer pokt_v0.PocketRelayer, chainConfiguration chain_configurations_registry.ChainConfigurationsService, logger *zap.Logger) *NodeSelectorClient {

	// base checks will share same node list and pocket relayer
	baseCheck := checks.NewCheck(pocketRelayer, chainConfiguration)

	// enabled checks
	enabledChecks := []checks.CheckJob{
		checks.NewEvmHeightCheck(baseCheck, logger.Named("evm_height_checker")),
		checks.NewEvmDataIntegrityCheck(baseCheck, logger.Named("evm_data_integrity_checker")),
	}
	selectorService := &NodeSelectorClient{
		sessionRegistry: sessionRegistry,
		logger:          logger,
		checkJobs:       enabledChecks,
	}
	selectorService.startJobChecker()
	return selectorService
}

func (q NodeSelectorClient) FindNode(chainId string) (*models.QosNode, bool) {
	var healthyNodes []*models.QosNode
	nodes, found := q.sessionRegistry.GetNodesByChain(chainId)
	if !found {
		return nil, false
	}
	for _, r := range nodes {
		if r.IsHealthy() {
			healthyNodes = append(healthyNodes, r)
		}
	}
	node, ok := common.GetRandomElement(healthyNodes)
	if !ok {
		return nil, false
	}
	return node, true
}

func (q NodeSelectorClient) startJobChecker() {
	ticker := time.Tick(jobCheckInterval)
	go func() {
		for {
			select {
			case <-ticker:
				for _, job := range q.checkJobs {
					if job.ShouldRun() {
						for chain, nodes := range q.sessionRegistry.GetNodesMap() {
							q.logger.Sugar().Infow("running job", "job", job.Name(), "chain", chain)
							job.SetNodes(nodes.Value())
							job.Perform()
						}
					}
				}
			}
		}
	}()
}
