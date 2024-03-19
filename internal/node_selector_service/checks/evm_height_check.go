package checks

import (
	"encoding/json"
	"pokt_gateway_server/internal/node_selector_service/models"
	relayer_models "pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"sync"
	"time"
)

const (
	evmHeightCheckInterval = time.Minute * 1
	heightJsonPayload      = `{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`
	defaultHeightThreshold = 100
)

type evmHeightResponse struct {
	Height uint64 `json:"result"`
}

type EvmHeightCheck struct {
	*Check
	nextCheckTime time.Time
}

func NewEvmHeightCheck(check *Check) *EvmHeightCheck {
	return &EvmHeightCheck{Check: check, nextCheckTime: time.Time{}}
}

type nodeRelayResponse struct {
	Node    *models.QosNode
	Relay   *relayer_models.SendRelayResponse
	Success bool
}

func (c *EvmHeightCheck) Name() string {
	return "evm_height_check"
}

func (c *EvmHeightCheck) Perform() {

	var highestHeight uint64
	var wg sync.WaitGroup

	// Define a channel to receive relay responses
	relayResponses := make(chan *nodeRelayResponse)

	// Define a function to handle sending relay requests concurrently
	sendRelayAsync := func(node *models.QosNode) {
		defer wg.Done()
		relay, err := c.pocketRelayer.SendRelay(&relayer_models.SendRelayRequest{
			Payload:            &relayer_models.Payload{Data: heightJsonPayload, Method: "POST"},
			Chain:              node.PocketSession.SessionHeader.Chain,
			SelectedNodePubKey: node.MorseNode.PublicKey,
		})
		relayResponses <- &nodeRelayResponse{
			Node:    node,
			Relay:   relay,
			Success: err == nil,
		}
	}

	// Start a goroutine for each node to send relay requests concurrently
	for _, node := range c.nodeList {
		wg.Add(1)
		go sendRelayAsync(node)
	}

	wg.Wait()
	close(relayResponses)

	// Process relay responses
	for resp := range relayResponses {
		if resp.Success {
			var evmHeightResp evmHeightResponse
			err := json.Unmarshal([]byte(resp.Relay.Response), &evmHeightResp)
			if err != nil {
				continue
			}
			resp.Node.SetLastHeightCheckTime(time.Now())
			resp.Node.SetLastKnownHeight(evmHeightResp.Height)
			// We track the session's highest height to make a majority decision
			reportedHeight := evmHeightResp.Height
			if reportedHeight > highestHeight {
				highestHeight = reportedHeight
			}
		}
	}

	// Compare each node's reported height against the highest reported height
	for _, node := range c.nodeList {
		heightDifference := int64(highestHeight) - int64(node.GetLastKnownHeight())
		// Penalize nodes whose reported height is significantly lower than the highest reported height
		if heightDifference > defaultHeightThreshold {
			node.SetSynced(false)
			node.SetTimeoutUntil(time.Now().Add(timeoutPenalty), models.OutOfSyncTimeout)
		} else {
			node.SetSynced(true)
		}
	}
	c.nextCheckTime = time.Now().Add(evmHeightCheckInterval)
}

func (c *EvmHeightCheck) SetNodes(nodes []*models.QosNode) {
	c.nodeList = nodes
}

func (c *EvmHeightCheck) ShouldRun() bool {
	return time.Now().After(c.nextCheckTime)
}

func (c *EvmHeightCheck) getEligibleNodes() []*models.QosNode {
	// Filter nodes based on last checked time
	var eligibleNodes []*models.QosNode
	for _, node := range c.nodeList {
		if node.GetLastHeightCheckTime().IsZero() || time.Since(node.GetLastHeightCheckTime()) >= minLastCheckedNodeTime {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	return eligibleNodes
}
