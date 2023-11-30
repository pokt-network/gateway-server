package controllers

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"os-gateway/cmd/gateway_server/internal/common"
	slice_common "os-gateway/pkg/common"
	"os-gateway/pkg/pokt/pokt_v0"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"strings"
)

// RelayController handles relay requests for a specific chain.
type RelayController struct {
	logger     *zap.Logger
	poktClient pokt_v0.PocketService
	appStakes  []*models.Ed25519Account
}

// NewRelayController creates a new instance of RelayController.
func NewRelayController(poktClient pokt_v0.PocketService, appStakes []*models.Ed25519Account, logger *zap.Logger) *RelayController {
	return &RelayController{poktClient: poktClient, appStakes: appStakes, logger: logger}
}

// RelayHandlerPath is the path for relay requests.
const RelayHandlerPath = "/relay/{catchAll:*}"

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

	// Get a random app stake from the available list.
	appStake := slice_common.GetRandomElement(c.appStakes)
	if appStake == nil {
		common.JSONError(ctx, "App stake not provided", fasthttp.StatusInternalServerError)
		return
	}

	// Send relay request to the Pocket Network.
	relayRsp, err := c.poktClient.SendRelay(&models.SendRelayRequest{
		Payload: &models.Payload{
			Data:   string(ctx.PostBody()),
			Method: string(ctx.Method()),
			Path:   path,
		},
		Signer: appStake,
		Chain:  chainID,
	})

	if err != nil {
		c.logger.Error("Error relaying", zap.Error(err))
		common.JSONError(ctx, "Something went wrong", fasthttp.StatusInternalServerError)
		return
	}

	// Send a successful response back to the client.
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBodyString(relayRsp.Response)
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
