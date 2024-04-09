package transform

import (
	"github.com/pokt-network/gateway-server/cmd/gateway_server/internal/models"
	internal_model "github.com/pokt-network/gateway-server/internal/apps_registry/models"
)

func ToPoktApplication(app *internal_model.PoktApplicationSigner) *models.PublicPoktApplication {
	return &models.PublicPoktApplication{
		ID:        app.ID,
		MaxRelays: int(app.NetworkApp.MaxRelays),
		Chains:    app.NetworkApp.Chains,
		Address:   app.NetworkApp.Address,
	}
}
