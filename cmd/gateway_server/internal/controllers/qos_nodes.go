package controllers

import (
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"pokt_gateway_server/cmd/gateway_server/internal/common"
	"pokt_gateway_server/cmd/gateway_server/internal/models"
	"pokt_gateway_server/cmd/gateway_server/internal/transform"
	"pokt_gateway_server/internal/session_registry"
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
