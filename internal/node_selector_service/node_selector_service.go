package node_selector_service

import (
	"github.com/pokt-network/gateway-server/internal/chain_configurations_registry"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks/solana_data_integrity_check"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks/solana_height_check"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"github.com/pokt-network/gateway-server/internal/session_registry"
	"github.com/pokt-network/gateway-server/pkg/common"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0"
	"go.uber.org/zap"
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
		//evm_height_check.NewEvmHeightCheck(baseCheck, logger.Named("evm_height_checker")),
		//evm_data_integrity_check.NewEvmDataIntegrityCheck(baseCheck, logger.Named("evm_data_integrity_checker")),
		solana_height_check.NewSolanaHeightCheck(baseCheck, logger.Named("solana_height_check")),
		solana_data_integrity_check.NewSolanaDataIntegrityCheck(baseCheck, logger.Named("solana_data_integrity_check")),
		//pokt_height_check.NewPoktHeightCheck(baseCheck, logger.Named("pokt_height_check")),
		//pokt_data_integrity_check.NewPoktDataIntegrityCheck(baseCheck, logger.Named("pokt_data_integrity_check")),
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

	nodes := q.sessionRegistry.GetNodesByChain(chainId)
	if len(nodes) == 0 {
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

// filterBySessionHeightNodes - filter by session height descending. This allows node selector to send relays with
// latest session height which nodes are more likely to serve vs session rollover relays.
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
						for _, nodes := range q.sessionRegistry.GetNodesMap() {
							job.SetNodes(nodes.Value())
							job.Perform()
						}
					}
				}
			}
		}
	}()
}
