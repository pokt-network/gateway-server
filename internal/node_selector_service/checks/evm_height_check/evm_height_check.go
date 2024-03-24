package evm_height_check

import (
	"encoding/json"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/stat"
	"math"
	"pokt_gateway_server/internal/node_selector_service/checks"
	"pokt_gateway_server/internal/node_selector_service/models"
	"strconv"
	"strings"
	"time"
)

const (
	// zScore to remove outliers for determing highest height
	zScoreHeightThreshold = 3

	// interval to check a node height again
	checkNodeHeightInterval = time.Minute * 5

	// interval to run the evm height check
	evmHeightCheckInterval = time.Second * 1

	// penalty for being out of sync
	evmHeightCheckPenalty = time.Minute * 5

	// jsonrpc payload to retrieve evm height
	heightJsonPayload = `{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`

	// default height allowance
	defaultHeightTolerance int = 100
)

type evmHeightResponse struct {
	Height uint64 `json:"blockByNumberResponse"`
}

func (r *evmHeightResponse) UnmarshalJSON(data []byte) error {

	type evmHeightResponseStr struct {
		Result string `json:"result"`
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
	*checks.Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

func NewEvmHeightCheck(check *checks.Check, logger *zap.Logger) *EvmHeightCheck {
	return &EvmHeightCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

func (c *EvmHeightCheck) Name() string {
	return "evm_height_check"
}

func (c *EvmHeightCheck) Perform() {

	// Send request to all nodes
	relayResponses := checks.SendRelaysAsync(c.PocketRelayer, c.NodeList, heightJsonPayload, "POST")

	var nodesResponded []*models.QosNode
	// Process relay responses
	for resp := range relayResponses {

		err := resp.Error
		if err != nil {
			checks.DefaultPunishNode(err, resp.Node, c.logger)
			continue
		}

		var evmHeightResp evmHeightResponse
		err = json.Unmarshal([]byte(resp.Relay.Response), &evmHeightResp)

		if err != nil {
			c.logger.Sugar().Warnw("failed to unmarshal response", "err", err)
			// Treat a invalid response as a timeout error
			checks.DefaultPunishNode(fasthttp.ErrTimeout, resp.Node, c.logger)
			continue
		}

		resp.Node.SetLastHeightCheckTime(time.Now())
		resp.Node.SetLastKnownHeight(evmHeightResp.Height)
		nodesResponded = append(nodesResponded, resp.Node)
	}

	highestNodeHeight := getHighestNodeHeight(nodesResponded)
	// Compare each node's reported height against the highest reported height
	for _, node := range nodesResponded {
		heightDifference := int(highestNodeHeight - node.GetLastKnownHeight())
		// Penalize nodes whose reported height is significantly lower than the highest reported height
		if heightDifference > checks.GetBlockHeightTolerance(c.ChainConfiguration, node.GetChain(), defaultHeightTolerance) {
			c.logger.Sugar().Infow("node is out of sync", "node", node.MorseNode.ServiceUrl, "heightDifference", heightDifference, "nodeSyncedHeight", node.GetLastKnownHeight(), "highestNodeHeight", highestNodeHeight, "chain", node.GetChain())
			// Punish Node specifically due to timeout.
			node.SetSynced(false)
			node.SetTimeoutUntil(time.Now().Add(evmHeightCheckPenalty), models.OutOfSyncTimeout)
		} else {
			node.SetSynced(true)
		}
	}
	c.nextCheckTime = time.Now().Add(evmHeightCheckInterval)
}

func (c *EvmHeightCheck) SetNodes(nodes []*models.QosNode) {
	c.NodeList = nodes
}

func (c *EvmHeightCheck) ShouldRun() bool {
	return time.Now().After(c.nextCheckTime)
}

func (c *EvmHeightCheck) getEligibleNodes() []*models.QosNode {
	// Filter nodes based on last checked time
	var eligibleNodes []*models.QosNode
	for _, node := range c.NodeList {
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

		// height is an outlier according to zScore threshold
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
