package pokt_apps_registry

import "pokt_gateway_server/internal/pokt_apps_registry/models"

type AppsRegistryService interface {
	GetApplications() []*models.PoktApplicationSigner
	GetApplicationsByChainId(chainId string) ([]*models.PoktApplicationSigner, bool)
}
