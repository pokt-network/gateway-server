package models

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name              string
		privateKey        string
		expectedPublicKey string
		expectedAddress   string
		err               error
	}{
		{
			name:              "BadPrivateKey",
			privateKey:        "badKey",
			expectedPublicKey: "",
			expectedAddress:   "",
			err:               ErrInvalidPrivateKey,
		},
		{
			name:              "Success",
			privateKey:        "3fe64039816c44e8872e4ef981725b968422e3d49e95a1eb800707591df30fe374039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849",
			expectedPublicKey: "74039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849",
			expectedAddress:   "d873127df524d172276e4da8193a2e2b19ef825f",
			err:               nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := NewAccount(tt.privateKey)
			assert.Equal(t, err, tt.err)
			if err == nil {
				assert.Equal(t, acc.PublicKey, tt.expectedPublicKey)
				assert.Equal(t, acc.Address, tt.expectedAddress)
			}
		})
	}
}

func TestEd25519Account_GetAAT(t *testing.T) {
	a, err := NewAccount("3fe64039816c44e8872e4ef981725b968422e3d49e95a1eb800707591df30fe374039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849")
	assert.Equal(t, err, nil)
	assert.Equal(t, &AAT{
		Version:      "0.0.1",
		AppPubKey:    "74039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849",
		ClientPubKey: "74039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849",
		Signature:    "f233ca857b4ada2ca4996e0da8c1761cfbc855edf282fc5a753d4631785946d6c2b08c781c84abbca2dc929de50008729079124e5c5c16921a81139279020a05",
	}, a.GetAAT())
}

func TestEd25519Account_Sign(t *testing.T) {
	a, err := NewAccount("3fe64039816c44e8872e4ef981725b968422e3d49e95a1eb800707591df30fe374039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849")
	assert.Equal(t, err, nil)
	assert.Equal(t, hex.EncodeToString(a.Sign([]byte("TestMessage"))), "6cf23f8aa00793ef6aec4d3c408f5be249f01ddc96778f3ea03ef8fcdd301e09ce175fbcb97778b222de57469857d99ef97ad978dc49992f70a108aafd3d3001")
}
