package checks

import (
	"encoding/json"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/stat"
	"math"
	"pokt_gateway_server/internal/node_selector_service/models"
	relayer_models "pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// zScore to remove outliers for determing highest height
	zScoreHeightThreshold = 3

	// interval to check a node height again
	checkNodeHeightInterval = time.Minute * 5

	// interval to run the evm height check
	evmHeightCheckInterval = time.Second * 1

	// jsonrpc payload to retrieve evm height
	heightJsonPayload = `{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`

	// height allowance
	defaultHeightThreshold = 100
)

type evmHeightResponse struct {
	Height uint64 `json:"blockByNumberResult"`
}

func (r *evmHeightResponse) UnmarshalJSON(data []byte) error {

	type evmHeightResponseStr struct {
		Result string `json:"blockByNumberResult"`
	}

	// Unmarshal the JSON into the custom type
	var hr evmHeightResponseStr
	if err := ffjson.Unmarshal(data, &hr); err != nil {
		return err
	}

	// Remove the "0x" prefix if present
	heightStr := strings.TrimPrefix(hr.Result, "0x")

	// Parse the hexadecimal string to uint64
	height, err := strconv.ParseUint(heightStr, 16, 64)
	if err != nil {
		return fmt.Errorf("failed to parse height: %v", err)
	}
	// Assign the parsed height to the struct field
	r.Height = height
	return nil
}

type EvmHeightCheck struct {
	*Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

func NewEvmHeightCheck(check *Check, logger *zap.Logger) *EvmHeightCheck {
	return &EvmHeightCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

type nodeRelayResponse struct {
	Node  *models.QosNode
	Relay *relayer_models.SendRelayResponse
	Error error
}

func (c *EvmHeightCheck) Name() string {
	return "evm_height_check"
}

func (c *EvmHeightCheck) Perform() {

	var wg sync.WaitGroup

	// Define a channel to receive relay responses
	relayResponses := make(chan *nodeRelayResponse, len(c.nodeList))

	// Define a function to handle sending relay requests concurrently
	sendRelayAsync := func(node *models.QosNode) {
		defer wg.Done()
		relay, err := c.pocketRelayer.SendRelay(&relayer_models.SendRelayRequest{
			Signer:             node.GetSigner(),
			Payload:            &relayer_models.Payload{Data: heightJsonPayload, Method: "POST"},
			Chain:              node.GetChain(),
			SelectedNodePubKey: node.GetPublicKey(),
			Session:            node.PocketSession,
		})
		relayResponses <- &nodeRelayResponse{
			Node:  node,
			Relay: relay,
			Error: err,
		}
	}

	// Start a goroutine for each node to send relay requests concurrently
	for _, node := range c.nodeList {
		wg.Add(1)
		go sendRelayAsync(node)
	}

	wg.Wait()
	close(relayResponses)

	var nodesResponded []*models.QosNode

	// Process relay responses
	for resp := range relayResponses {

		err := resp.Error
		if err != nil {
			defaultPunishNode(err, resp.Node, c.logger)
			continue
		}

		var evmHeightResp evmHeightResponse
		err = json.Unmarshal([]byte(resp.Relay.Response), &evmHeightResp)

		if err != nil {
			// Treat a invalid response as a timeout error
			defaultPunishNode(fasthttp.ErrTimeout, resp.Node, c.logger)
			continue
		}

		resp.Node.SetLastHeightCheckTime(time.Now())
		resp.Node.SetLastKnownHeight(evmHeightResp.Height)
		nodesResponded = append(nodesResponded, resp.Node)
	}

	highestNodeHeight := getHighestNodeHeight(nodesResponded)
	// Compare each node's reported height against the highest reported height
	for _, node := range nodesResponded {
		heightDifference := int64(highestNodeHeight) - int64(node.GetLastKnownHeight())
		// Penalize nodes whose reported height is significantly lower than the highest reported height
		if heightDifference > defaultHeightThreshold {
			c.logger.Sugar().Infow("node is out of sync", "node", node.MorseNode.ServiceUrl, "heightDifference", heightDifference, "nodeSyncedHeight", node.GetLastKnownHeight(), "highestNodeHeight", highestNodeHeight, "chain", node.GetChain())
			// Punish Node specifically due to timeout.
			node.SetSynced(false)
			node.SetTimeoutUntil(time.Now().Add(dataIntegrityTimePenalty), models.OutOfSyncTimeout)
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
		if node.GetLastHeightCheckTime().IsZero() || time.Since(node.GetLastHeightCheckTime()) >= checkNodeHeightInterval {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	return eligibleNodes
}

// getHighestHeight returns the highest height reported from a slice of nodes
// by using z-score threshhold to prevent any misconfigured or malicious node
func getHighestNodeHeight(nodes []*models.QosNode) uint64 {

	var nodeHeights []float64
	for _, node := range nodes {
		nodeHeights = append(nodeHeights, float64(node.GetLastKnownHeight()))
	}

	// Calculate mean and standard deviation
	meanValue := stat.Mean(nodeHeights, nil)
	stdDevValue := stat.StdDev(nodeHeights, nil)

	var highestNodeHeight float64
	for _, nodeHeight := range nodeHeights {

		zScore := stat.StdScore(nodeHeight, meanValue, stdDevValue)

		// height is an outlier according to zscore threshold
		if math.Abs(zScore) > zScoreHeightThreshold {
			continue
		}
		// Height is higher than last recorded height
		if nodeHeight > highestNodeHeight {
			highestNodeHeight = nodeHeight
		}
	}
	return uint64(highestNodeHeight)
}
