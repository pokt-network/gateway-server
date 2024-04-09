package apps_registry

import "github.com/pokt-network/gateway-server/internal/apps_registry/models"

type AppsRegistryService interface {
	GetApplications() []*models.PoktApplicationSigner
	GetApplicationsByChainId(chainId string) ([]*models.PoktApplicationSigner, bool)
	GetApplicationByPublicKey(publicKey string) (*models.PoktApplicationSigner, bool)
}
