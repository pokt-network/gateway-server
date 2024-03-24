package chain_configurations_registry

import (
	"context"
	"go.uber.org/zap"
	"pokt_gateway_server/internal/db_query"
	"sync"
	"time"
)

const (
	chainConfigurationUpdateInterval = time.Minute * 1
)

type CachedChainConfigurationRegistry struct {
	dbQuery                 db_query.Querier
	chainConfigurationCache map[string]db_query.GetChainConfigurationsRow // chain id > url
	cacheLock               sync.RWMutex
	logger                  *zap.Logger
}

func NewCachedChainConfigurationRegistry(dbQuery db_query.Querier, logger *zap.Logger) *CachedChainConfigurationRegistry {
	chainConfigurationRegistry := &CachedChainConfigurationRegistry{dbQuery: dbQuery, chainConfigurationCache: map[string]db_query.GetChainConfigurationsRow{}, logger: logger}
	err := chainConfigurationRegistry.updateChainConfigurations()
	if err != nil {
		chainConfigurationRegistry.logger.Sugar().Warnw("Failed to retrieve chain global_config on startup", "err", err)
	}
	chainConfigurationRegistry.startCacheUpdater()
	return chainConfigurationRegistry
}

func (r *CachedChainConfigurationRegistry) GetChainConfiguration(chainId string) (db_query.GetChainConfigurationsRow, bool) {
	r.cacheLock.RLock()
	defer r.cacheLock.RUnlock()
	url, found := r.chainConfigurationCache[chainId]
	return url, found
}

func (r *CachedChainConfigurationRegistry) updateChainConfigurations() error {
	chainConfigurations, err := r.dbQuery.GetChainConfigurations(context.Background())

	if err != nil {
		return err
	}

	chainConfigurationNew := map[string]db_query.GetChainConfigurationsRow{}
	for _, row := range chainConfigurations {
		chainConfigurationNew[row.ChainID.String] = row
	}

	// Update the cache
	r.cacheLock.Lock()
	defer r.cacheLock.Unlock()
	r.chainConfigurationCache = chainConfigurationNew
	return nil
}

// StartCacheUpdater starts a goroutine to periodically update the altruist cache.
func (c *CachedChainConfigurationRegistry) startCacheUpdater() {
	ticker := time.Tick(chainConfigurationUpdateInterval)
	go func() {
		for {
			select {
			case <-ticker:
				// Call the updateChainConfigurations method
				err := c.updateChainConfigurations()
				if err != nil {
					c.logger.Sugar().Warnw("failed to update chain configuration registry", "err", err)
				} else {
					c.logger.Sugar().Infow("successfully updated chain configuration registry", "chainConfigurationLength", len(c.chainConfigurationCache))
				}
			}
		}
	}()
}
