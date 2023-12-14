package pokt_v0

import (
	"os-gateway/pkg/pokt/pokt_v0/models"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type PocketService interface {
	GetSession(req *models.GetSessionRequest) (*models.GetSessionResponse, error)
	SendRelay(req *models.SendRelayRequest) (*models.SendRelayResponse, error)
	GetLatestBlockHeight() (*models.GetLatestBlockHeightResponse, error)
}

type CacheService interface {
	NewCache(time.Duration) *ttlcache.Cache[string, *models.GetSessionResponse]
}
