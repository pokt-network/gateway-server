package evm_height_check

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"github.com/pquerna/ffjson/ffjson"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

const (

	// interval to run the evm height check
	evmHeightCheckInterval = time.Second * 1

	// jsonrpc payload to retrieve evm height
	heightJsonPayload = `{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`
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

	// Session is not meant for EVM
	if len(c.NodeList) == 0 || !c.IsEvmChain(c.NodeList[0]) {
		return
	}
	checks.PerformDefaultHeightCheck(c.Check, heightJsonPayload, "", c.getHeightFromNodeResponse, c.logger)
	c.nextCheckTime = time.Now().Add(evmHeightCheckInterval)
}

func (c *EvmHeightCheck) SetNodes(nodes []*models.QosNode) {
	c.NodeList = nodes
}

func (c *EvmHeightCheck) ShouldRun() bool {
	return time.Now().After(c.nextCheckTime)
}

func (c *EvmHeightCheck) getHeightFromNodeResponse(response string) (uint64, error) {
	var evmRsp evmHeightResponse
	err := json.Unmarshal([]byte(response), &evmRsp)
	if err != nil {
		return 0, err
	}
	return evmRsp.Height, nil
}
