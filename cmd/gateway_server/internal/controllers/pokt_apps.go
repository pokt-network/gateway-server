package controllers

import (
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/cmd/gateway_server/internal/models"
	"pokt_gateway_server/cmd/gateway_server/internal/transform"
	"pokt_gateway_server/internal/pokt_apps_registry"
	"pokt_gateway_server/pkg/pokt/pokt_v0"
)

// RelayController handles relay requests for a specific chain.
type PoktAppsController struct {
	logger      *zap.Logger
	poktClient  pokt_v0.PocketService
	appRegistry pokt_apps_registry.AppsRegistryService
}

// NewRelayController creates a new instance of RelayController.
func NewPoktAppsController(appRegistry pokt_apps_registry.AppsRegistryService, logger *zap.Logger) *PoktAppsController {
	return &PoktAppsController{appRegistry: appRegistry, logger: logger}
}

// GetAllPoktApps is the path for relay requests.
const GetAllPoktApps = "/poktapps"

// HandleRelay handles incoming relay requests.
func (c *PoktAppsController) GetAll(ctx *fasthttp.RequestCtx) {
	applications := c.appRegistry.GetApplications()
	appsPublic := []*models.PoktApplication{}
	for _, app := range applications {
		appsPublic = append(appsPublic, transform.ToPoktApplication(app))
	}
	result, err := ffjson.Marshal(appsPublic)
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	// Send a successful response back to the client.
	ctx.Response.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetBody(result)
}
