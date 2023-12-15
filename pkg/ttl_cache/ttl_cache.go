package ttl_cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type TTLCacheService[K comparable, V any] interface {
	Get(key K, opts ...ttlcache.Option[K, V]) *ttlcache.Item[K, V]
	Set(key K, value V, ttl time.Duration) *ttlcache.Item[K, V]
}
