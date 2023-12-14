package pokt_v0

import (
	"os-gateway/pkg/pokt/pokt_v0/models"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type Cache struct{}

func NewCache(ttl time.Duration) *ttlcache.Cache[string, *models.GetSessionResponse] {
	return ttlcache.New[string, *models.GetSessionResponse](
		ttlcache.WithTTL[string, *models.GetSessionResponse](ttl),
	)
}
