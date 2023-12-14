package ttl_cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type TTLCacheClient[K comparable, V any] struct {
	TTLCacheService[K, V]
}

func NewTTLCacheClient[K comparable, V any]() *TTLCacheClient[K, V] {
	return &TTLCacheClient[K, V]{}
}

type TTLCacheService[K comparable, V any] interface {
	Get(key K) *ttlcache.Item[K, V]
	Set(key K, value V, ttl time.Duration) *ttlcache.Item[K, V]
	Start()
}
