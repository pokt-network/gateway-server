package altruist_registry

import (
	"context"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/db_query"
	"sync"
	"time"
)

const (
	altruistUpdateInterval = time.Minute * 1
)

type CachedAltruistRegistryService struct {
	dbQuery   db_query.Querier
	altruists map[string]string // chain id > url
	lockCache sync.RWMutex
	logger    *zap.Logger
}

func NewCachedAltruistRegistryService(dbQuery db_query.Querier, logger *zap.Logger) *CachedAltruistRegistryService {
	altruistRegistry := &CachedAltruistRegistryService{dbQuery: dbQuery, altruists: map[string]string{}, logger: logger}
	err := altruistRegistry.updateAltruistUrls()
	if err != nil {
		altruistRegistry.logger.Sugar().Warnw("Failed to retrieve altruist urls on startup", "err", err)
	}
	altruistRegistry.startCacheUpdater()
	return altruistRegistry
}

func (r *CachedAltruistRegistryService) GetAltruistURL(chainId string) (string, bool) {
	r.lockCache.RLock()
	defer r.lockCache.RUnlock()
	url, found := r.altruists[chainId]
	return url, found
}

func (r *CachedAltruistRegistryService) updateAltruistUrls() error {
	altruists, err := r.dbQuery.GetAltruists(context.Background())

	if err != nil {
		return err
	}

	altruistNew := map[string]string{}
	for _, row := range altruists {
		altruistNew[row.ChainID.String] = row.Url.String
	}

	// Update the cache
	r.lockCache.Lock()
	defer r.lockCache.Unlock()
	r.altruists = altruistNew
	return nil
}

// StartCacheUpdater starts a goroutine to periodically update the altruist cache.
func (c *CachedAltruistRegistryService) startCacheUpdater() {
	ticker := time.Tick(altruistUpdateInterval)
	go func() {
		for {
			select {
			case <-ticker:
				// Call the updateApplicationCache method
				err := c.updateAltruistUrls()
				if err != nil {
					c.logger.Sugar().Warnw("failed to update altruist registry", "err", err)
				} else {
					c.logger.Sugar().Infow("successfully updated altruist registry", "altruistLength", len(c.altruists))
				}
			}
		}
	}()
}
