//go:generate ffjson $GOFILE
package models

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"github.com/pokt-network/gateway-server/pkg/common"
	"sync"
)

// Ed25519Account represents an account using the Ed25519 cryptographic algorithm.
type Ed25519Account struct {
	privateKeyBytes []byte
	aat             *AAT
	aatOnce         sync.Once
	PrivateKey      string `json:"privateKey"`
	PublicKey       string `json:"publicKey"`
	Address         string `json:"address"`
}

const (
	// CurrentAATVersion represents the current version of the Application Authentication Token (AAT).
	CurrentAATVersion = "0.0.1"

	// privKeyMaxLength represents the required length of the private key.
	privKeyMaxLength = 128
)

var (
	// ErrInvalidPrivateKey is returned when the private key is invalid.
	ErrInvalidPrivateKey = errors.New("invalid private key, requires 128 chars")
)

// NewAccount creates a new Ed25519Account instance.
//
// Parameters:
//   - privateKey: Private key as a string.
//
// Returns:
//   - (*Ed25519Account): New Ed25519Account instance.
//   - (error): Error, if any.
func NewAccount(privateKey string) (*Ed25519Account, error) {
	if len(privateKey) != privKeyMaxLength {
		return nil, ErrInvalidPrivateKey
	}
	publicKey := privateKey[64:]
	addr, err := common.GetAddressFromPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return &Ed25519Account{
		privateKeyBytes: privateKeyBytes,
		PrivateKey:      privateKey,
		PublicKey:       publicKey,
		Address:         addr,
	}, nil
}

// Sign signs a given message using the account's private key.
//
// Parameters:
//   - message: Message to sign.
//
// Returns:
//   - []byte: Signature.
func (a *Ed25519Account) Sign(message []byte) []byte {
	return ed25519.Sign(a.privateKeyBytes, message)
}

// GetAAT retrieves the Application Authentication Token (AAT) associated with the account.
//
// Returns:
//   - (*AAT): AAT for the account.
func (a *Ed25519Account) GetAAT() *AAT {
	a.aatOnce.Do(func() {
		aat := AAT{
			Version:      CurrentAATVersion,
			AppPubKey:    a.PublicKey,
			ClientPubKey: a.PublicKey,
			Signature:    "",
		}
		bytes := common.Sha3_256Hash(aat)
		aat.Signature = hex.EncodeToString(a.Sign(bytes))
		a.aat = &aat
	})
	return a.aat
}
