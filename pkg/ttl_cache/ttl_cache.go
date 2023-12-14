package ttl_cache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type TTLCacheService[K comparable, V any] interface {
	New(opts ...ttlcache.Option[K, V]) *ttlcache.Cache[K, V]
	// New[K comparable, V any](opts ttlcache.Option[K, V]) *ttlcache.Cache[K, V]
	Delete(key K)
	DeleteAll()
	DeleteExpired()
	Get(key K, opts ttlcache.Option[K, V]) *ttlcache.Item[K, V]
	GetAndDelete(key K, opts ttlcache.Option[K, V]) (*ttlcache.Item[K, V], bool)
	GetOrSet(key K, value V, opts ttlcache.Option[K, V]) (*ttlcache.Item[K, V], bool)
	Has(key K) bool
	Items() map[K]*ttlcache.Item[K, V]
	Keys() []K
	Len() int
	Metrics() ttlcache.Metrics
	OnEviction(fn func(context.Context, ttlcache.EvictionReason, *ttlcache.Item[K, V])) func()
	OnInsertion(fn func(context.Context, *ttlcache.Item[K, V])) func()
	Range(fn func(item *ttlcache.Item[K, V]) bool)
	Set(key K, value V, ttl time.Duration) *ttlcache.Item[K, V]
	Start()
	Stop()
	Touch(key K)
}
