package solana_height_check

import (
	"encoding/json"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/checks"
	"github.com/pokt-network/gateway-server/internal/node_selector_service/models"
	"go.uber.org/zap"
	"time"
)

const (

	// interval to run the solana height check
	solanaHeightCheckInterval = time.Second * 1

	// jsonrpc payload to retrieve solana height
	heightJsonPayload = `{"jsonrpc":"2.0","method":"getSlot","params": [],"id":1}`
)

type solanaHeightResponse struct {
	Result uint64 `json:"result"`
}

type SolanaHeightCheck struct {
	*checks.Check
	nextCheckTime time.Time
	logger        *zap.Logger
}

func NewSolanaHeightCheck(check *checks.Check, logger *zap.Logger) *SolanaHeightCheck {
	return &SolanaHeightCheck{Check: check, nextCheckTime: time.Time{}, logger: logger}
}

func (c *SolanaHeightCheck) Name() string {
	return "solana_height_check"
}

func (c *SolanaHeightCheck) Perform() {

	// Session is not meant for Solana
	if len(c.NodeList) == 0 || !c.IsSolanaChain(c.NodeList[0]) {
		return
	}
	checks.PerformDefaultHeightCheck(c.Check, heightJsonPayload, "", c.getHeightFromNodeResponse, c.logger)
	c.nextCheckTime = time.Now().Add(solanaHeightCheckInterval)
}

func (c *SolanaHeightCheck) SetNodes(nodes []*models.QosNode) {
	c.NodeList = nodes
}

func (c *SolanaHeightCheck) ShouldRun() bool {
	return time.Now().After(c.nextCheckTime)
}

func (c *SolanaHeightCheck) getHeightFromNodeResponse(response string) (uint64, error) {
	var solanaRsp solanaHeightResponse
	err := json.Unmarshal([]byte(response), &solanaRsp)
	if err != nil {
		return 0, err
	}
	return solanaRsp.Result, nil
}
