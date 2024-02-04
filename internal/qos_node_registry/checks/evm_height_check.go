package checks

import (
	"encoding/json"
	"math"
	"pokt_gateway_server/internal/qos_node_registry/models"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
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
	nodesToCheck []*models.QosNode
	relayer      pokt_v0.PocketRelayer
	chainId      string
	lastChecked  time.Time
	lock         sync.Mutex
}

func (c *EvmHeightCheck) PerformJob() {
	c.lock.Lock()
	defer c.lock.Unlock()

	var highestHeight uint64

	// Gather heights from all nodes and keep track of the highest reported height
	for _, node := range c.nodesToCheck {
		relay, err := c.relayer.SendRelay(&relayer_models.SendRelayRequest{
			Payload:            &relayer_models.Payload{Data: heightJsonPayload, Method: "POST"},
			Chain:              c.chainId,
			SelectedNodePubKey: node.MorseNode.PublicKey,
		})
		if err != nil {
			continue
		}

		var evmHeightResp evmHeightResponse
		err = json.Unmarshal([]byte(relay.Response), &evmHeightResp)
		if err != nil {
			continue
		}

		node.SetLastKnownHeight(evmHeightResp.Height)
		reportedHeight := evmHeightResp.Height
		if reportedHeight > highestHeight {
			highestHeight = reportedHeight
		}
	}

	// Compare each node's reported height against the highest reported height
	for _, node := range c.nodesToCheck {

		heightDifference := int64(highestHeight) - int64(node.GetLastKnownHeight())
		// Penalize nodes whose reported height is significantly lower than the highest reported height
		if math.Abs(float64(heightDifference)) > defaultHeightThreshold {
			node.SetTimeoutUntil(time.Now().Add(timeoutPenalty))
		}
	}

	c.lastChecked = time.Now()
}

func (c *EvmHeightCheck) ShouldRun() bool {
	return time.Now().Sub(c.lastChecked) > evmHeightCheckInterval
}
