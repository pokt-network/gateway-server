package checks

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/node_selector_service/models"
	relayer_models "pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"strings"
	"time"
)

// default timeout penalty whenever a node doesn't respond
const timeoutErrorPenalty = time.Second * 15

// 24 hours is analogous to indefinite
const kickOutSessionPenalty = time.Hour * 24

const (
	errOverServiceMsg           = "the max number of relays serviced for this node is exceeded"
	errMaximumEvidenceSealedMsg = "the evidence is sealed, either max relays reached or claim already submitted"
)

// isMaximumRelaysServicedErr - determines if a node should be kicked from a session to send relays
func isMaximumRelaysServicedErr(err error) bool {
	// If evidence is sealed or the node has already overserviced, the node should no longer receive relays.
	if err == relayer_models.ErrPocketEvidenceSealed || err == relayer_models.ErrPocketCoreOverService {
		return true
	}
	// Fallback in the event the error is not parsed correctly due to node operator configurations / custom clients, resort to a simple string check
	pocketError, ok := err.(relayer_models.PocketRPCError)
	if ok {
		return strings.Contains(pocketError.Message, errOverServiceMsg) || strings.Contains(pocketError.Message, errMaximumEvidenceSealedMsg)
	}
	return false
}

func isTimeoutError(err error) bool {
	pocketError, ok := err.(relayer_models.PocketRPCError)
	if ok {
		return pocketError.HttpCode >= 500 || strings.Contains(pocketError.Message, "request timeout")
	}
	return err == fasthttp.ErrTimeout || err == fasthttp.ErrDialTimeout || err == fasthttp.ErrTLSHandshakeTimeout
}

// DefaultPunishNode generic punisher for whenever a node returns an error independent of a specific check
func DefaultPunishNode(err error, node *models.QosNode, logger *zap.Logger) bool {
	if isMaximumRelaysServicedErr(err) {
		node.SetTimeoutUntil(time.Now().Add(kickOutSessionPenalty), models.MaximumRelaysTimeout)
		return true
	}
	if isTimeoutError(err) {
		node.SetTimeoutUntil(time.Now().Add(timeoutErrorPenalty), models.NodeResponseTimeout)
		return true
	}
	logger.Sugar().Errorw("unknown error for punishing node", "node", node.MorseNode.ServiceUrl)
	return false
}
