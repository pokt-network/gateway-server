package node_selector_service

import (
	"go.uber.org/zap"
	"pokt_gateway_server/internal/chain_configurations_registry"
	"pokt_gateway_server/internal/node_selector_service/checks"
	"pokt_gateway_server/internal/node_selector_service/checks/evm_data_integrity_check"
	"pokt_gateway_server/internal/node_selector_service/checks/evm_height_check"
	"pokt_gateway_server/internal/node_selector_service/models"
	"pokt_gateway_server/internal/session_registry"
	"pokt_gateway_server/pkg/common"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"sort"
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
		evm_height_check.NewEvmHeightCheck(baseCheck, logger.Named("evm_height_checker")),
		evm_data_integrity_check.NewEvmDataIntegrityCheck(baseCheck, logger.Named("evm_data_integrity_checker")),
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

	nodes, ok := q.sessionRegistry.GetNodesByChain(chainId)
	if !ok {
		return nil, false
	}

	// Filter nodes by health
	healthyNodes := filterByHealthyNodes(nodes)

	// Find a node that's closer to session height
	sortedSessionHeights, nodeMap := filterBySessionHeightNodes(healthyNodes)
	for _, sessionHeight := range sortedSessionHeights {
		node, ok := common.GetRandomElement(nodeMap[sessionHeight])
		if ok {
			return node, true
		}
	}
	return nil, false
}

func filterBySessionHeightNodes(nodes []*models.QosNode) ([]uint, map[uint][]*models.QosNode) {
	nodesBySessionHeight := map[uint][]*models.QosNode{}

	// Create map to retrieve nodes by session height
	for _, r := range nodes {
		sessionHeight := r.MorseSession.SessionHeader.SessionHeight
		nodesBySessionHeight[sessionHeight] = append(nodesBySessionHeight[sessionHeight], r)
	}

	// Create slice to hold sorted session heights
	var sortedSessionHeights []uint
	for sessionHeight := range nodesBySessionHeight {
		sortedSessionHeights = append(sortedSessionHeights, sessionHeight)
	}

	// Sort the slice of session heights by descending order
	sort.Slice(sortedSessionHeights, func(i, j int) bool {
		return sortedSessionHeights[i] > sortedSessionHeights[j]
	})

	return sortedSessionHeights, nodesBySessionHeight
}
func filterByHealthyNodes(nodes []*models.QosNode) []*models.QosNode {
	var healthyNodes []*models.QosNode

	for _, r := range nodes {
		if r.IsHealthy() {
			healthyNodes = append(healthyNodes, r)
		}
	}
	return healthyNodes
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
