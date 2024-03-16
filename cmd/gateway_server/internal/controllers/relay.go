package controllers

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/cmd/gateway_server/internal/common"
	"pokt_gateway_server/internal/altruist_registry"
	"pokt_gateway_server/internal/apps_registry"
	"pokt_gateway_server/internal/qos_node_registry"
	"pokt_gateway_server/internal/relayer"
	"pokt_gateway_server/internal/session_registry"
	slice_common "pokt_gateway_server/pkg/common"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
	"strings"
)

var ErrRelayChannelClosed = errors.New("concurrent relay channel closed")

// RelayController handles relay requests for a specific chain.
type RelayController struct {
	logger           *zap.Logger
	relayer          *relayer.Relayer
	appRegistry      apps_registry.AppsRegistryService
	altruistRegistry altruist_registry.AltruistRegistryService
	sessionRegistry  session_registry.SessionRegistryService
	nodeSelector     *qos_node_registry.NodeSelectorService
}

// NewRelayController creates a new instance of RelayController.
func NewRelayController(relayer *relayer.Relayer, appRegistry apps_registry.AppsRegistryService, sessionRegistry session_registry.SessionRegistryService, altruistRegistry altruist_registry.AltruistRegistryService, nodeSelector *qos_node_registry.NodeSelectorService, logger *zap.Logger) *RelayController {
	return &RelayController{relayer: relayer, appRegistry: appRegistry, sessionRegistry: sessionRegistry, nodeSelector: nodeSelector, altruistRegistry: altruistRegistry, logger: logger}
}

// chainIdLength represents the expected length of chain IDs.
const chainIdLength = 4

// HandleRelay handles incoming relay requests.
func (c *RelayController) HandleRelay(ctx *fasthttp.RequestCtx) {

	chainID, path := getPathSegmented(ctx.Path())

	// Check if the chain ID is empty or has an incorrect length.
	if chainID == "" || len(chainID) != chainIdLength {
		common.JSONError(ctx, "Incorrect chain id", fasthttp.StatusBadRequest)
		return
	}

	// Use a healthy node determined by node selector
	node, ok := c.nodeSelector.FindNode(chainID)

	if ok {
		req := &models.SendRelayRequest{
			Payload: &models.Payload{
				Path: path,
			},
			Signer:             node.Signer,
			Chain:              chainID,
			SelectedNodePubKey: node.MorseNode.PublicKey,
		}
		relay, err := c.relayer.SendRelay(&relayer.RelayRequest{
			PocketRequest: req,
			UseAltruist:   false,
		})
		if err != nil {
			c.logger.Error("Error relaying", zap.Error(err))
			common.JSONError(ctx, "Something went wrong", fasthttp.StatusInternalServerError)
			return
		}

		// Send a successful response back to the client.
		ctx.Response.SetStatusCode(fasthttp.StatusOK)
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Response.SetBodyString(relay.Response)
		return
	}

	// Healthy node could not be found, attempting to use random node
	applications, ok := c.appRegistry.GetApplicationsByChainId(chainID)
	if !ok || len(applications) == 0 {
		common.JSONError(ctx, fmt.Sprintf("%s chainId not supported with existing application registry", chainID), fasthttp.StatusBadRequest)
		return
	}

	// Get a random app stake from the available list.
	appStake := slice_common.GetRandomElement(applications)
	if appStake == nil {
		common.JSONError(ctx, "App stake not provided", fasthttp.StatusInternalServerError)
		return
	}

	sessionResp, err := c.sessionRegistry.GetSession(&models.GetSessionRequest{
		AppPubKey: appStake.Signer.PublicKey,
		Chain:     chainID,
	})

	if err != nil {
		c.logger.Error("Error dispatching session", zap.Error(err))
		common.JSONError(ctx, "Something went wrong", fasthttp.StatusInternalServerError)
		return
	}

	randomNode := slice_common.GetRandomElement(sessionResp.Nodes)

	if randomNode == nil {
		c.logger.Error("Error finding a node from session", zap.Error(err))
		common.JSONError(ctx, "Something went wrong", fasthttp.StatusInternalServerError)
		return
	}

	req := &models.SendRelayRequest{
		Payload: &models.Payload{
			Data:   string(ctx.PostBody()),
			Method: string(ctx.Method()),
			Path:   path,
		},
		Signer:             appStake.Signer,
		Chain:              chainID,
		SelectedNodePubKey: randomNode.MorseNode.PublicKey,
	}
	relay, err := c.relayer.SendRelay(&relayer.RelayRequest{
		PocketRequest: req,
		UseAltruist:   true,
	})
	if err != nil {
		c.logger.Error("Error relaying", zap.Error(err))
		common.JSONError(ctx, "Something went wrong", fasthttp.StatusInternalServerError)
		return
	}
	// Send a successful response back to the client.
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(relay.Response)
}

func getPathSegmented(path []byte) (chain, otherParts string) {
	paths := strings.Split(string(path), "/")

	if len(paths) >= 3 {
		chain = paths[2]
	}

	if len(paths) > 3 {
		otherParts = "/" + strings.Join(paths[3:], "/")
	}

	return chain, otherParts
}
