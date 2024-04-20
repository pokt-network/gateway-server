package solana_data_integrity_check

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"go.uber.org/zap"
	"time"
)

const (
	// how often the job should run
	dataIntegrityCheckInterval = time.Second * 1

	//json rpc payload to send a data integrity check
	// we use signatures for transaction detail to prevent large payloads and we don't need anything but block hash
	blockPayloadFmt = `{"jsonrpc":"2.0","method":"getBlock","params":[%d, {"encoding": "jsonParsed", "maxSupportedTransactionVersion":0, "transactionDetails":"signatures"}],"id":1}`
)

type blockByNumberResponse struct {
	Result struct {
		Hash string `json:"blockhash"`
	} `json:"result"`
}

type SolanaDataIntegrityCheck struct {
	*checks.Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

func (c *SolanaDataIntegrityCheck) getBlockIdentifierFromNodeResponse(response string) (string, error) {
	var blockRsp blockByNumberResponse
	err := json.Unmarshal([]byte(response), &blockRsp)
	if err != nil {
		return "", err
	}
	return blockRsp.Result.Hash, nil
}

func NewSolanaDataIntegrityCheck(check *checks.Check, logger *zap.Logger) *SolanaDataIntegrityCheck {
	return &SolanaDataIntegrityCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

func (c *SolanaDataIntegrityCheck) Name() string {
	return "solana_data_integrity_check"
}

func (c *SolanaDataIntegrityCheck) SetNodes(nodes []*models.QosNode) {
	c.NodeList = nodes
}

func (c *SolanaDataIntegrityCheck) Perform() {

	// Session is not meant for Solana
	if len(c.NodeList) == 0 || !c.NodeList[0].IsSolanaChain() {
		return
	}
	checks.PerformDataIntegrityCheck(c.Check, getBlockByNumberPayload, "", c.getBlockIdentifierFromNodeResponse, c.logger)
	c.nextCheckTime = time.Now().Add(dataIntegrityCheckInterval)
}

func (c *SolanaDataIntegrityCheck) ShouldRun() bool {
	return c.nextCheckTime.IsZero() || time.Now().After(c.nextCheckTime)
}

func getBlockByNumberPayload(blockNumber uint64) string {
	return fmt.Sprintf(blockPayloadFmt, blockNumber)
}
