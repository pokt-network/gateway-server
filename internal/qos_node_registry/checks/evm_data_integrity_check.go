package checks

import (
	"encoding/json"
	"pokt_gateway_server/internal/qos_node_registry/models"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
	relayer_models "pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"sync"
	"time"
)

const (
	timeoutPenalty = time.Minute * 10
	checkInterval  = time.Minute * 1
	blockPayload   = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}`
)

type result struct {
	Hash string `json:"hash"`
}

type evmResponse struct {
	Result result `json:"result"`
}

type EvmDataIntegrityCheck struct {
	nodesToCheck []*models.QosNode
	relayer      pokt_v0.PocketRelayer
	chainId      string
	lastChecked  time.Time
	lock         sync.Mutex // Add a mutex for concurrent access
}

type nodeResponse struct {
	node   *models.QosNode
	result result
}

func (c *EvmDataIntegrityCheck) PerformJob() {

	c.lock.Lock()
	defer c.lock.Unlock()

	// Initialize a map to store responses and their counts
	nodeResponseCounts := make(map[nodeResponse]int)

	for _, node := range c.nodesToCheck {
		relay, err := c.relayer.SendRelay(&relayer_models.SendRelayRequest{
			Payload:            &relayer_models.Payload{Data: blockPayload, Method: "POST"},
			Chain:              c.chainId,
			SelectedNodePubKey: node.MorseNode.PublicKey,
		})
		if err != nil {
			continue
		}

		var resp evmResponse
		err = json.Unmarshal([]byte(relay.Response), &resp)
		if err != nil {
			continue
		}

		nodeResponseCounts[nodeResponse{
			node:   node,
			result: resp.Result,
		}]++
	}

	highestResponseHash := findMajorityResponse(nodeResponseCounts)

	// Penalize other node operators with a timeout
	for nodeResp := range nodeResponseCounts {
		if nodeResp.result.Hash != highestResponseHash {
			nodeResp.node.SetTimeoutUntil(time.Now().Add(timeoutPenalty))
		}
	}
	c.lastChecked = time.Now()
}

// findMajorityResponse finds the hash with the highest response count
func findMajorityResponse(responseCounts map[nodeResponse]int) string {
	var highestResponseHash string
	var highestResponseCount int
	for rsp, count := range responseCounts {
		if count > highestResponseCount {
			highestResponseHash = rsp.result.Hash
			highestResponseCount = count
		}
	}
	return highestResponseHash
}

func (c *EvmDataIntegrityCheck) ShouldRun() bool {
	return time.Now().Sub(c.lastChecked) > checkInterval
}
