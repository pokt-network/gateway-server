package models

import (
	"github.com/pokt-network/gateway-server/pkg/pokt/pokt_v0/models"
)

type PoktApplicationSigner struct {
	Signer     *models.Ed25519Account
	NetworkApp *models.PoktApplication
	ID         string
}

func NewPoktApplicationSigner(id string, account *models.Ed25519Account) *PoktApplicationSigner {
	return &PoktApplicationSigner{Signer: account, ID: id}
}
