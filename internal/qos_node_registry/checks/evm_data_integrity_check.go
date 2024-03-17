package checks

import (
	"encoding/json"
	"pokt_gateway_server/internal/qos_node_registry/models"
	relayer_models "pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"time"
)

const (
	minLastCheckedNodeTime = time.Minute * 1
	timeoutPenalty         = time.Minute * 1
	checkInterval          = time.Second * 5
	blockPayload           = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}`
)

type result struct {
	Hash string `json:"hash"`
}

type evmResponse struct {
	Result result `json:"result"`
}

type EvmDataIntegrityCheck struct {
	*Check
}

type nodeResponse struct {
	node   *models.QosNode
	result result
}

func (c *EvmDataIntegrityCheck) Name() string {
	return "evm_data_integrity_check"
}

func (c *EvmDataIntegrityCheck) Perform() {
	// Initialize a map to store responses and their counts
	nodeResponseCounts := make(map[nodeResponse]int)
	for _, node := range c.nodeList {
		relay, err := c.pocketRelayer.SendRelay(&relayer_models.SendRelayRequest{
			Payload:            &relayer_models.Payload{Data: blockPayload, Method: "POST"},
			Chain:              node.PocketSession.SessionHeader.Chain,
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
			nodeResp.node.SetTimeoutUntil(time.Now().Add(timeoutPenalty), models.InvalidDataTimeout)
		}
	}
	c.nextCheckTime = time.Now().Add(checkInterval)
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
	return c.nextCheckTime.IsZero() || time.Now().After(c.nextCheckTime)
}
