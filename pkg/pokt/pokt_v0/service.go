package pokt_v0

import (
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
)

type PocketService interface {
	GetSession(req *models.GetSessionRequest) (*models.GetSessionResponse, error)
	SendRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error)
	GetLatestBlockHeight() (*models.GetLatestBlockHeightResponse, error)
	GetLatestStakedApplications() ([]*models.PoktApplication, error)
}
