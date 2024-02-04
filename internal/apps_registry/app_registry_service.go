package apps_registry

import "pokt_gateway_server/internal/apps_registry/models"

type AppsRegistryService interface {
	GetApplications() []*models.PoktApplicationSigner
	GetApplicationsByChainId(chainId string) ([]*models.PoktApplicationSigner, bool)
}
