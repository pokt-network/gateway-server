package models

import "pokt_gateway_server/pkg/pokt/pokt_v0/models"

type PoktApplicationSigner struct {
	ID string `json:"id"`
	*models.Ed25519Account
	*models.PoktApplication
}

func NewPoktApplicationSigner(account *models.Ed25519Account, application *models.PoktApplication) *PoktApplicationSigner {
	return &PoktApplicationSigner{Ed25519Account: account, PoktApplication: application}
}
