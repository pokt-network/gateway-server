package controllers

import (
	"context"
	"github.com/jackc/pgtype"
	"github.com/pokt-network/gateway-server/cmd/gateway_server/internal/common"
	"github.com/pokt-network/gateway-server/cmd/gateway_server/internal/models"
	"github.com/pokt-network/gateway-server/cmd/gateway_server/internal/transform"
	"github.com/pokt-network/gateway-server/internal/apps_registry"
	"github.com/pokt-network/gateway-server/internal/db_query"
	"github.com/pokt-network/gateway-server/internal/global_config"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0"
	pokt_models "github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type addApplicationBody struct {
	PrivateKey string `json:"private_key"`
}

// PoktAppsController handles requests for staked applications
type PoktAppsController struct {
	logger         *zap.Logger
	query          db_query.Querier
	poktClient     pokt_v0.PocketService
	appRegistry    apps_registry.AppsRegistryService
	secretProvider global_config.SecretProvider
}

// NewPoktAppsController creates a new instance of PoktAppsController.
func NewPoktAppsController(appRegistry apps_registry.AppsRegistryService, query db_query.Querier, secretProvider global_config.SecretProvider, logger *zap.Logger) *PoktAppsController {
	return &PoktAppsController{appRegistry: appRegistry, query: query, secretProvider: secretProvider, logger: logger}
}

// GetAll returns all the apps in the registry
func (c *PoktAppsController) GetAll(ctx *fasthttp.RequestCtx) {
	applications := c.appRegistry.GetApplications()
	appsPublic := []*models.PublicPoktApplication{}
	for _, app := range applications {
		appsPublic = append(appsPublic, transform.ToPoktApplication(app))
	}
	common.JSONSuccess(ctx, appsPublic, fasthttp.StatusOK)
}

// AddApplication - enables users to add an application programmatically.
// Not recommended since it requires transmitting creds over wire and opens up to MITM (if not encrypted, or user error).
func (c *PoktAppsController) AddApplication(ctx *fasthttp.RequestCtx) {
	var body addApplicationBody
	err := ffjson.Unmarshal(ctx.PostBody(), &body)
	if err != nil {
		common.JSONError(ctx, "Faiiled to unmarshal req", fasthttp.StatusInternalServerError, err)
		return
	}

	account, err := pokt_models.NewAccount(body.PrivateKey)
	if err != nil {
		common.JSONError(ctx, "Faiiled to convert to ed25519 account", fasthttp.StatusBadRequest, err)
		return
	}
	_, err = c.query.InsertPoktApplications(context.Background(), account.PrivateKey, c.secretProvider.GetPoktApplicationsEncryptionKey())
	if err != nil {
		common.JSONError(ctx, "Something went wrong", fasthttp.StatusInternalServerError, err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
}

// DeleteApplication - enables users to delete an application programmatically.
// Not recommended since it requires transmitting creds over wire and opens up to MITM (if not encrypted, or user error).
func (c *PoktAppsController) DeleteApplication(ctx *fasthttp.RequestCtx) {
	applicationId := ctx.UserValue("app_id")
	uuid := pgtype.UUID{}
	uuid.Set(applicationId)
	_, err := c.query.DeletePoktApplication(context.Background(), uuid)
	if err != nil {
		common.JSONError(ctx, "Something went wrong", fasthttp.StatusInternalServerError, err)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
