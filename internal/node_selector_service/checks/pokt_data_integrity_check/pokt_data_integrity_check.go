package pokt_data_integrity_check

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"go.uber.org/zap"
	"time"
)

const poktBlockTxEndpoint = "/v1/query/blocktxs"

const (
	// how often the job should run
	dataIntegrityCheckInterval = time.Second * 1

	//json rpc payload to send a data integrity check
	blockPayloadFmt = `{"height": %d}`
)

type blockTxResponse struct {
	TotalTxs int `json:"total_txs"`
}

type PoktDataIntegrityCheck struct {
	*checks.Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

// getBlockIdentifierFromNodeResponse: We use total txs as the block identifier because retrieving block hash from POKT RPC can lead up to
// 8MB+ payloads per node operator response, whereas blocktxs is only ~110kb
func (c *PoktDataIntegrityCheck) getBlockIdentifierFromNodeResponse(response string) (string, error) {
	var blockTxRsp blockTxResponse
	err := json.Unmarshal([]byte(response), &blockTxRsp)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", blockTxRsp.TotalTxs), nil
}

func NewPoktDataIntegrityCheck(check *checks.Check, logger *zap.Logger) *PoktDataIntegrityCheck {
	return &PoktDataIntegrityCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

func (c *PoktDataIntegrityCheck) Name() string {
	return "pokt_data_integrity_check"
}

func (c *PoktDataIntegrityCheck) SetNodes(nodes []*models.QosNode) {
	c.NodeList = nodes
}

func (c *PoktDataIntegrityCheck) Perform() {

	// Session is not meant for POKT
	if len(c.NodeList) == 0 || !c.IsPoktChain(c.NodeList[0]) {
		return
	}
	checks.PerformDataIntegrityCheck(c.Check, getBlockByNumberPayload, poktBlockTxEndpoint, c.getBlockIdentifierFromNodeResponse, c.logger)
	c.nextCheckTime = time.Now().Add(dataIntegrityCheckInterval)
}

func (c *PoktDataIntegrityCheck) ShouldRun() bool {
	return c.nextCheckTime.IsZero() || time.Now().After(c.nextCheckTime)
}

func getBlockByNumberPayload(blockNumber uint64) string {
	return fmt.Sprintf(blockPayloadFmt, blockNumber)
}
