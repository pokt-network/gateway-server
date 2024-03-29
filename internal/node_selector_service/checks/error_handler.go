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
	errPocketInvalidServicerMsg       = "failed to find correct servicer PK"
	errPocketInvalidBlockHeightMsg    = "the block height passed is invalid"
	errPocketRequestTimeoutMsg        = "request timeout"
	errPocketOverServiceMsg           = "the max number of relays serviced for this node is exceeded"
	errPocketMaximumEvidenceSealedMsg = "the evidence is sealed, either max relays reached or claim already submitted"
)

const (
	errHttpSSLExpired    = "tls: failed to verify certificate"
	errHttpNoSuchHostMsg = "no such host"
)

// isKickableSessionErr - determines if a node should be kicked from a session to send relays
func isKickableSessionErr(err error) bool {
	// If evidence is sealed or the node has already overserviced, the node should no longer receive relays.
	if err == relayer_models.ErrPocketEvidenceSealed || err == relayer_models.ErrPocketCoreOverService {
		return true
	}
	// Fallback in the event the error is not parsed correctly due to node operator configurations / custom clients, resort to a simple string check
	// node runner cannot serve with expired ssl
	if err != nil && (strings.Contains(err.Error(), errHttpSSLExpired) || strings.Contains(err.Error(), errPocketOverServiceMsg) || strings.Contains(err.Error(), errPocketMaximumEvidenceSealedMsg) || strings.Contains(err.Error(), errPocketInvalidServicerMsg)) {
		return true
	}
	return false
}

func isTimeoutError(err error) bool {
	// If Invalid block height, pocket  is not caught up to latest session
	if err == relayer_models.ErrPocketCoreInvalidBlockHeight {
		return true
	}

	// Check if pocket error returns 500
	pocketError, ok := err.(relayer_models.PocketRPCError)
	if ok && pocketError.HttpCode >= 500 {
		return true
	}

	// Fallback in the event the error is not parsed correctly due to node operator configurations / custom clients, resort to a simple string check
	return err == fasthttp.ErrTimeout || err == fasthttp.ErrDialTimeout || err == fasthttp.ErrTLSHandshakeTimeout || err != nil && (strings.Contains(err.Error(), errHttpNoSuchHostMsg) || strings.Contains(err.Error(), errPocketRequestTimeoutMsg) || strings.Contains(err.Error(), errPocketInvalidBlockHeightMsg))
}

// DefaultPunishNode generic punisher for whenever a node returns an error independent of a specific check
func DefaultPunishNode(err error, node *models.QosNode, logger *zap.Logger) bool {
	if isKickableSessionErr(err) {
		node.SetTimeoutUntil(time.Now().Add(kickOutSessionPenalty), models.MaximumRelaysTimeout, err)
		return true
	}
	if isTimeoutError(err) {
		node.SetTimeoutUntil(time.Now().Add(timeoutErrorPenalty), models.NodeResponseTimeout, err)
		return true
	}
	logger.Sugar().Warnw("uncategorized error detected from pocket node", "node", node.MorseNode.ServiceUrl, "err", err)
	return false
}
