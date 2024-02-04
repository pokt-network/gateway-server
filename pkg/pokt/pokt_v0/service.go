package pokt_v0

import (
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
)

type PocketRelayer interface {
	SendRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error)
}

type PocketDispatcher interface {
	GetSession(req *models.GetSessionRequest) (*models.GetSessionResponse, error)
}
type PocketService interface {
	PocketRelayer
	PocketDispatcher
	GetLatestBlockHeight() (*models.GetLatestBlockHeightResponse, error)
	GetLatestStakedApplications() ([]*models.PoktApplication, error)
}
