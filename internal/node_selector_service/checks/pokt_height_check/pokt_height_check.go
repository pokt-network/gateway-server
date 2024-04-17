package pokt_height_check

import (
	"encoding/json"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"go.uber.org/zap"
	"time"
)

const (

	// interval to run the pokt height check job
	poktHeightCheckInterval = time.Second * 1

	// jsonrpc payload to pokt evm height
	heightJsonPayload = ``
)

type poktHeightResponse struct {
	Height uint64 `json:"height"`
}

type PoktHeightCheck struct {
	*checks.Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

func NewPoktHeightCheck(check *checks.Check, logger *zap.Logger) *PoktHeightCheck {
	return &PoktHeightCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

func (c *PoktHeightCheck) Name() string {
	return "pokt_height_check"
}

func (c *PoktHeightCheck) Perform() {

	// Session is not meant for EVM
	if len(c.NodeList) == 0 || !c.NodeList[0].IsPoktChain() {
		return
	}
	checks.PerformDefaultHeightCheck(c.Check, heightJsonPayload, "/v1/query/height", c.getHeightFromNodeResponse, c.logger)
	c.nextCheckTime = time.Now().Add(poktHeightCheckInterval)
}

func (c *PoktHeightCheck) SetNodes(nodes []*models.QosNode) {
	c.NodeList = nodes
}

func (c *PoktHeightCheck) ShouldRun() bool {
	return time.Now().After(c.nextCheckTime)
}

func (c *PoktHeightCheck) getHeightFromNodeResponse(response string) (uint64, error) {
	var poktRsp poktHeightResponse
	err := json.Unmarshal([]byte(response), &poktRsp)
	if err != nil {
		return 0, err
	}
	return poktRsp.Height, nil
}
