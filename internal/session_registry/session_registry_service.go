package session_registry

import "pokt_gateway_server/pkg/pokt/pokt_v0/models"

type SessionRegistryService interface {
	GetSession(req *models.GetSessionRequest) (*models.GetSessionResponse, error)
}
