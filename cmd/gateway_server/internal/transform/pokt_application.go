package transform

import (
	"pokt_gateway_server/cmd/gateway_server/internal/models"
	internal_model "pokt_gateway_server/internal/apps_registry/models"
)

func ToPoktApplication(app *internal_model.PoktApplicationSigner) *models.PoktApplication {
	return &models.PoktApplication{
		ID:        app.ID,
		MaxRelays: int(app.NetworkApp.MaxRelays),
		Chains:    app.NetworkApp.Chains,
		Address:   app.NetworkApp.Address,
	}
}
