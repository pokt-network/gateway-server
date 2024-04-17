package evm_data_integrity_check

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"go.uber.org/zap"
	"strconv"
	"time"
)

const (
	// how often the job should run
	dataIntegrityCheckInterval = time.Second * 1

	//json rpc payload to send a data integrity check
	blockPayloadFmt = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", false],"id":1}`
)

type blockByNumberResponse struct {
	Result struct {
		Hash string `json:"hash"`
	} `json:"result"`
}

type EvmDataIntegrityCheck struct {
	*checks.Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

func (c *EvmDataIntegrityCheck) getBlockHashFromNodeResponse(response string) (string, error) {
	var evmRsp blockByNumberResponse
	err := json.Unmarshal([]byte(response), &evmRsp)
	if err != nil {
		return "", err
	}
	return evmRsp.Result.Hash, nil
}

func NewEvmDataIntegrityCheck(check *checks.Check, logger *zap.Logger) *EvmDataIntegrityCheck {
	return &EvmDataIntegrityCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

func (c *EvmDataIntegrityCheck) Name() string {
	return "evm_data_integrity_check"
}

func (c *EvmDataIntegrityCheck) SetNodes(nodes []*models.QosNode) {
	c.NodeList = nodes
}

func (c *EvmDataIntegrityCheck) Perform() {

	// Session is not meant for EVM
	if len(c.NodeList) == 0 || !c.NodeList[0].IsEvmChain() {
		return
	}
	checks.PerformDataIntegrityCheck(c.Check, getBlockByNumberPayload, "", c.getBlockHashFromNodeResponse, c.logger)
	c.nextCheckTime = time.Now().Add(dataIntegrityCheckInterval)
}

func (c *EvmDataIntegrityCheck) ShouldRun() bool {
	return c.nextCheckTime.IsZero() || time.Now().After(c.nextCheckTime)
}

func getBlockByNumberPayload(blockNumber uint64) string {
	return fmt.Sprintf(blockPayloadFmt, "0x"+strconv.FormatInt(int64(blockNumber), 16))
}
