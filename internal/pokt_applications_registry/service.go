package pokt_applications_registry

import "pokt_gateway_server/internal/pokt_applications_registry/models"

type Service interface {
	GetApplications() []*models.PoktApplicationSigner
	GetApplicationsByChainId(chainId string) ([]*models.PoktApplicationSigner, bool)
}
