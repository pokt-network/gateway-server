package apps_registry

import (
	"context"
	"errors"
	"github.com/pokt-network/gateway-server/internal/apps_registry/models"
	"github.com/pokt-network/gateway-server/internal/db_query"
	"github.com/pokt-network/gateway-server/internal/global_config"
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0"
	pokt "github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
	"go.uber.org/zap"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	applicationUpdateInterval = time.Second * 15
)

// CachedAppsRegistry is a caching layer for storing and retrieving from internal DB and POKT's blockchain state
type CachedAppsRegistry struct {
	pocketClient        pokt_v0.PocketService
	dbQuery             db_query.Querier
	logger              *zap.Logger
	secretProvider      global_config.SecretProvider
	applications        []*models.PoktApplicationSigner
	applicationChainMap map[string][]*models.PoktApplicationSigner
	lockCache           sync.RWMutex
}

// NewCachedAppsRegistry creates a new instance of CachedAppsRegistry.
func NewCachedAppsRegistry(pocketClient pokt_v0.PocketService, dbQuery db_query.Querier, secretProvider global_config.SecretProvider, logger *zap.Logger) *CachedAppsRegistry {
	cachedRegistry := CachedAppsRegistry{pocketClient: pocketClient, dbQuery: dbQuery, logger: logger, secretProvider: secretProvider}
	err := cachedRegistry.updateApplicationCache()
	if err != nil {
		cachedRegistry.logger.Sugar().Warnw("failed to retrieve applications on init", "err", err)
	}
	cachedRegistry.startCacheUpdater()
	return &cachedRegistry
}

// GetApplications returns all the cached Pocket applications.
func (c *CachedAppsRegistry) GetApplications() []*models.PoktApplicationSigner {
	c.lockCache.RLock()
	defer c.lockCache.RUnlock()
	return c.applications
}

// GetApplicationsByChainId returns Pocket applications filtered by a specific chain ID.
func (c *CachedAppsRegistry) GetApplicationsByChainId(chainId string) ([]*models.PoktApplicationSigner, bool) {
	c.lockCache.RLock()
	defer c.lockCache.RUnlock()
	app, ok := c.applicationChainMap[chainId]
	return app, ok
}

func (c *CachedAppsRegistry) GetApplicationByPublicKey(pubKey string) (*models.PoktApplicationSigner, bool) {
	c.lockCache.RLock()
	defer c.lockCache.RUnlock()
	for _, i := range c.applications {
		if i.Signer.PublicKey == pubKey {
			return i, true
		}
	}
	return nil, false
}

// updateApplicationCache refreshes the cache with the latest Pocket applications and their associated information.
func (c *CachedAppsRegistry) updateApplicationCache() error {
	// Retrieve Pocket applications from the database
	storedAppsPK, err := c.dbQuery.GetPoktApplications(context.Background(), c.secretProvider.GetPoktApplicationsEncryptionKey())
	if err != nil {
		return err
	}

	// Convert the encrypted private keys to PoktApplicationSigner
	poktApplicationSigners := []*models.PoktApplicationSigner{}
	for _, app := range storedAppsPK {
		account, err := pokt.NewAccount(app.DecryptedPrivateKey)
		if err != nil {
			// Log a warning if there is an error converting the private key
			c.logger.Sugar().Warnw("failed to update application", "err", err, "app", app.ID)
			continue
		}
		id, _ := app.ID.Value()
		poktApplicationSigners = append(poktApplicationSigners, models.NewPoktApplicationSigner(id.(string), account))
	}

	// Retrieve the latest staked applications from the Pocket service
	networkStakedApps, err := c.pocketClient.GetLatestStakedApplications()
	if err != nil {
		return err
	}

	storedNetworkAppFound := 0
	for _, networkApp := range networkStakedApps {
		for _, storedAccount := range poktApplicationSigners {
			// Check if the account address matches, and create a PoktApplicationSigner if there's a match
			if strings.EqualFold(networkApp.Address, storedAccount.Signer.Address) {
				storedAccount.NetworkApp = networkApp
				storedNetworkAppFound++
			}
		}
	}

	// Check if all provided keys have an associated network app
	if storedNetworkAppFound != len(poktApplicationSigners) {
		return errors.New("failed to identify all applications private key from network state. POKT RPC potentially failed")
	}

	// Create a map to organize PoktApplicationSigners by chain ID
	applicationChainMap := make(map[string][]*models.PoktApplicationSigner)

	// Iterate through each PoktApplicationSigner and associate it with the corresponding chain IDs
	for _, signer := range poktApplicationSigners {
		for _, chainID := range signer.NetworkApp.Chains {
			// Append the PoktApplicationSigner to the chain ID entry in the map
			applicationChainMap[chainID] = append(applicationChainMap[chainID], signer)
		}
	}

	// No changes needed so will not replace. We do this to also prevent regenerating AAT's
	if arePoktApplicationSignersEqual(c.applications, poktApplicationSigners) {
		return nil
	}
	// Acquire a write lock and update the cache with the newly retrieved information
	c.lockCache.Lock()
	defer c.lockCache.Unlock()
	c.applications = poktApplicationSigners
	c.applicationChainMap = applicationChainMap
	return nil
}

// StartCacheUpdater starts a goroutine to periodically update the application cache.
func (c *CachedAppsRegistry) startCacheUpdater() {
	ticker := time.Tick(applicationUpdateInterval)
	go func() {
		for {
			select {
			case <-ticker:
				// Call the updateApplicationCache method
				err := c.updateApplicationCache()
				if err != nil {
					c.logger.Sugar().Warnw("failed to update application cache", "err", err)
				} else {
					c.logger.Sugar().Infow("successfully updated application registry", "applicationsLength", len(c.applications))
				}
			}
		}
	}()
}

// arePoktApplicationSignersEqual checks if two slices of PoktApplicationSigner are equal.
func arePoktApplicationSignersEqual(slice1, slice2 []*models.PoktApplicationSigner) bool {
	// Check if the lengths are different
	if len(slice1) != len(slice2) {
		return false
	}
	// Create sorted copies of the slices without modifying the original slices to prevent any latency delays with main relay function
	sortedSlice1 := make([]*models.PoktApplicationSigner, len(slice1))
	copy(sortedSlice1, slice1)
	sort.Slice(sortedSlice1, func(i, j int) bool {
		return sortedSlice1[i].NetworkApp.Address < sortedSlice1[j].NetworkApp.Address
	})

	sortedSlice2 := make([]*models.PoktApplicationSigner, len(slice2))
	copy(sortedSlice2, slice2)
	sort.Slice(sortedSlice2, func(i, j int) bool {
		return sortedSlice2[i].NetworkApp.Address < sortedSlice2[j].NetworkApp.Address
	})

	// Now that slices are sorted, check if address keys are same.
	for i := range slice1 {
		// Check if any field is different
		networkApp1 := sortedSlice1[i].NetworkApp
		networkApp2 := sortedSlice2[i].NetworkApp
		if !strings.EqualFold(networkApp1.Address, networkApp2.Address) {
			return false
		}

		if !strings.EqualFold(strings.Join(networkApp1.Chains, ","), strings.Join(networkApp2.Chains, ",")) {
			return false
		}

	}

	// Slices are equal
	return true
}
