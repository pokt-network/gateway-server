package qos_node_registry

import (
	"go.uber.org/zap"
	"pokt_gateway_server/internal/qos_node_registry/checks"
	"pokt_gateway_server/internal/qos_node_registry/models"
	"pokt_gateway_server/internal/session_registry"
	"pokt_gateway_server/pkg/common"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	"time"
)

const (
	jobCheckInterval = time.Second
)

type NodeSelectorService struct {
	sessionRegistry session_registry.SessionRegistryService
	pocketRelayer   pokt_v0.PocketRelayer
	logger          *zap.Logger
}

func NewNodeSelectorService(sessionRegistry session_registry.SessionRegistryService, pocketRelayer pokt_v0.PocketRelayer, logger *zap.Logger) *NodeSelectorService {
	selectorService := &NodeSelectorService{
		sessionRegistry: sessionRegistry,
		logger:          logger,
	}
	selectorService.startJobChecker()
	return selectorService
}

func (q NodeSelectorService) getEnabledJobs() []checks.CheckJob {
	nodes := q.sessionRegistry.GetNodes()
	baseCheck := checks.NewCheck(nodes, q.pocketRelayer)
	return []checks.CheckJob{

		&checks.EvmHeightCheck{
			Check: baseCheck,
		},

		&checks.EvmDataIntegrityCheck{
			Check: baseCheck,
		},
	}
}

func (q NodeSelectorService) FindNode(chainId string) (*models.QosNode, bool) {
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
	return common.GetRandomElement(healthyNodes), true
}

func (q NodeSelectorService) startJobChecker() {
	ticker := time.Tick(jobCheckInterval)
	go func() {
		for {
			select {
			case <-ticker:
				for _, job := range q.getEnabledJobs() {
					if job.ShouldRun() {
						q.logger.Sugar().Infow("running job", "job", job.Name())
						job.Perform()
					}
				}
			}
		}
	}()
}
