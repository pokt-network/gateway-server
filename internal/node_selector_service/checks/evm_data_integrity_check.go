package checks

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/node_selector_service/models"
	"pokt_gateway_server/pkg/common"
	relayer_models "pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"strconv"
	"time"
)

const (
	// penalty whenever a pocket node doesn't match other node providers responses
	dataIntegrityTimePenalty = time.Minute * 15

	// how often to check a node's data integrity
	dataIntegrityNodeCheckInterval = time.Minute * 10

	// how often the job should run
	dataIntegrityCheckInterval = time.Second * 1

	// the lookback we will use to determine which block number to do a data integrity against (latestBlockHeight - lookBack)
	dataIntegrityHeightLookback = 25

	//json rpc payload to send a data integrity check
	blockPayloadFmt = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", false],"id":1}`
)

type blockByNumberResult struct {
	Hash string `json:"hash"`
}

type evmResponse struct {
	Result blockByNumberResult `json:"result"`
}

type EvmDataIntegrityCheck struct {
	*Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

func NewEvmDataIntegrityCheck(check *Check, logger *zap.Logger) *EvmDataIntegrityCheck {
	return &EvmDataIntegrityCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

type nodeResponsePair struct {
	node   *models.QosNode
	result blockByNumberResult
}

func (c *EvmDataIntegrityCheck) Name() string {
	return "evm_data_integrity_check"
}

func (c *EvmDataIntegrityCheck) SetNodes(nodes []*models.QosNode) {
	c.nodeList = nodes
}

func (c *EvmDataIntegrityCheck) Perform() {

	// Find a node that has been reported as healthy
	healthyNode := c.findHealthyNode()

	// Node that is synced cannot be found, so we cannot run data integrity checks since we need a trusted source
	if healthyNode == nil {
		c.logger.Sugar().Warnw("cannot find node for data integrity checks for", "chain", c.nodeList[0].GetChain())
		return
	}

	heightSource := healthyNode.GetLastKnownHeight()

	// Init map that stores the amount of nodes that return the same response
	nodeResponseCounts := make(map[string]int)
	var nodeResponsePairs []*nodeResponsePair
	for _, node := range c.getEligibleNodes() {

		relay, err := c.pocketRelayer.SendRelay(&relayer_models.SendRelayRequest{
			Signer:             node.GetSigner(),
			Payload:            &relayer_models.Payload{Data: getBlockByNumberPayload(heightSource - dataIntegrityHeightLookback), Method: "POST"},
			Chain:              node.GetChain(),
			SelectedNodePubKey: node.GetPublicKey(),
			Session:            node.PocketSession,
		})

		if err != nil {
			defaultPunishNode(err, node, c.logger)
			continue
		}

		var resp evmResponse
		err = json.Unmarshal([]byte(relay.Response), &resp)
		if err != nil {
			defaultPunishNode(err, node, c.logger)
			continue
		}
		node.SetLastDataIntegrityCheckTime(time.Now())
		nodeResponsePairs = append(nodeResponsePairs, &nodeResponsePair{
			node:   node,
			result: resp.Result,
		})
		nodeResponseCounts[resp.Result.Hash]++
	}

	majorityBlockHash := findMajorityBlockHash(nodeResponseCounts)

	// Penalize other node operators with a timeout if they don't have same block hash.
	for _, nodeResp := range nodeResponsePairs {
		if nodeResp.result.Hash != majorityBlockHash {
			c.logger.Sugar().Errorw("punishing node for failed data integrity check", "node", nodeResp.node.MorseNode.ServiceUrl, "nodeBlockHash", nodeResp.result.Hash, "trustedSourceBlockHash", majorityBlockHash)
			nodeResp.node.SetTimeoutUntil(time.Now().Add(dataIntegrityTimePenalty), models.DataIntegrityTimeout)
		} else {
			fmt.Println("good")
		}
	}

	c.nextCheckTime = time.Now().Add(dataIntegrityCheckInterval)
}

// findMajorityBlockHash finds the hash with the highest response count
func findMajorityBlockHash(responseCounts map[string]int) string {
	var highestResponseHash string
	var highestResponseCount int
	for rsp, count := range responseCounts {
		if count > highestResponseCount {
			highestResponseHash = rsp
			highestResponseCount = count
		}
	}
	return highestResponseHash
}

func (c *EvmDataIntegrityCheck) ShouldRun() bool {
	return c.nextCheckTime.IsZero() || time.Now().After(c.nextCheckTime)
}

// findHealthyNode - returns a healthy node that is synced so we can use it as a source of truth for data integrity checks
func (c *EvmDataIntegrityCheck) findHealthyNode() *models.QosNode {
	var healthyNodes []*models.QosNode
	for _, node := range c.nodeList {
		if node.IsHealthy() {
			healthyNodes = append(healthyNodes, node)
		}
	}
	healthyNode := common.GetRandomElement(healthyNodes)
	return healthyNode
}

func (c *EvmDataIntegrityCheck) getEligibleNodes() []*models.QosNode {
	// Filter nodes based on last checked time
	var eligibleNodes []*models.QosNode
	for _, node := range c.nodeList {
		if (node.GetLastDataIntegrityCheckTime().IsZero() || time.Since(node.GetLastDataIntegrityCheckTime()) >= dataIntegrityNodeCheckInterval) && node.IsHealthy() {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	return eligibleNodes
}

func getBlockByNumberPayload(blockNumber uint64) string {
	return fmt.Sprintf(blockPayloadFmt, "0x"+strconv.FormatInt(int64(blockNumber), 16))
}
