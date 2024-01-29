package transform

import (
	"pokt_gateway_server/cmd/gateway_server/internal/models"
	internal_model "pokt_gateway_server/internal/pokt_apps_registry/models"
)

func ToPoktApplication(signer *internal_model.PoktApplicationSigner) *models.PoktApplication {
	return &models.PoktApplication{
		ID:        signer.ID,
		MaxRelays: int(signer.MaxRelays),
		Chains:    signer.Chains,
		Address:   signer.PoktApplication.Address,
	}
}
