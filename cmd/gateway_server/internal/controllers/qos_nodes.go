package controllers

import (
	"github.com/pokt-network/gateway-server/cmd/gateway_server/internal/common"
	"github.com/pokt-network/gateway-server/cmd/gateway_server/internal/models"
	"github.com/pokt-network/gateway-server/cmd/gateway_server/internal/transform"
	"github.com/pokt-network/gateway-server/internal/session_registry"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// QosNodeController handles requests for staked applications
type QosNodeController struct {
	logger          *zap.Logger
	sessionRegistry session_registry.SessionRegistryService
}

// NewQosNodeController  creates a new instance of QosNodeController.
func NewQosNodeController(sessionRegistry session_registry.SessionRegistryService, logger *zap.Logger) *QosNodeController {
	return &QosNodeController{sessionRegistry: sessionRegistry, logger: logger}
}

// GetAll returns all the qos nodes in the registry and exposes public information about them.
func (c *QosNodeController) GetAll(ctx *fasthttp.RequestCtx) {
	qosNodes := []*models.PublicQosNode{}
	for _, nodes := range c.sessionRegistry.GetNodesMap() {
		for _, node := range nodes.Value() {
			qosNodes = append(qosNodes, transform.ToPublicQosNode(node))
		}
	}
	common.JSONSuccess(ctx, qosNodes, fasthttp.StatusOK)
}
