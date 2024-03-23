package checks

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/node_selector_service/models"
	relayer_models "pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"strings"
	"time"
)

const timeoutErrorPenalty = time.Second * 15

// isMaximumRelaysServicedErr - determines if a node should be kicked from a session to send relays
func isMaximumRelaysServicedErr(err error) bool {
	// If evidence is sealed or the node has already overserviced, the node should no longer receive relays.
	return err == relayer_models.ErrPocketEvidenceSealed || err == relayer_models.ErrPocketCoreOverService
}

func isTimeoutError(err error) bool {
	pocketError, ok := err.(relayer_models.PocketRPCError)
	if ok {
		return pocketError.HttpCode >= 500 || strings.Contains(pocketError.Message, "request timeout")
	}
	return err == fasthttp.ErrTimeout || err == fasthttp.ErrDialTimeout || err == fasthttp.ErrTLSHandshakeTimeout
}

// defaultPunishNode: generic punisher for whenever a node returns an error independent of a specific check
func defaultPunishNode(err error, node *models.QosNode, logger *zap.Logger) bool {
	if isMaximumRelaysServicedErr(err) {
		// 24 hours is analogous to indefinite
		node.SetTimeoutUntil(time.Now().Add(time.Hour*24), models.MaximumRelaysTimeout)
		return true
	}
	if isTimeoutError(err) {
		node.SetTimeoutUntil(time.Now().Add(timeoutErrorPenalty), models.NodeResponseTimeout)
		return true
	}
	logger.Sugar().Errorw("unknown error for punishing node", "node", node.MorseNode.ServiceUrl)
	return false
}
