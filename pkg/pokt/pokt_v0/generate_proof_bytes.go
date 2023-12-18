package pokt_v0

import (
	"encoding/hex"
	"math/rand"
	"os-gateway/pkg/common"
	"os-gateway/pkg/pokt/pokt_v0/models"
)

// generateRelayProof generates a relay proof.
// Parameters:
//   - chainId: Blockchain ID.
//   - sessionHeight: Session block height.
//   - servicerPubKey: Servicer public key.
//   - requestMetadata: Request metadata.
//   - account: Ed25519 account used for signing.
//
// Returns:
//   - models.RelayProof: Generated relay proof.
func generateRelayProof(chainId string, sessionHeight uint, servicerPubKey string, relayMetadata *models.RelayMeta, reqPayload *models.Payload, account *models.Ed25519Account) *models.RelayProof {
	entropy := uint64(rand.Int63())
	aat := account.GetAAT()

	requestMetadata := models.RequestHashPayload{
		Metadata: relayMetadata,
		Payload:  reqPayload,
	}

	requestHash := requestMetadata.Hash()

	unsignedAAT := &models.AAT{
		Version:      aat.Version,
		AppPubKey:    aat.AppPubKey,
		ClientPubKey: aat.ClientPubKey,
		Signature:    "",
	}

	proofObj := &models.RelayProofHashPayload{
		RequestHash:        requestHash,
		Entropy:            entropy,
		SessionBlockHeight: sessionHeight,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         chainId,
		Signature:          "",
		UnsignedAAT:        unsignedAAT.Hash(),
	}

	hashedPayload := common.Sha3_256Hash(proofObj)
	hashSignature := hex.EncodeToString(account.Sign(hashedPayload))
	return &models.RelayProof{
		RequestHash:        requestHash,
		Entropy:            entropy,
		SessionBlockHeight: sessionHeight,
		ServicerPubKey:     servicerPubKey,
		Blockchain:         chainId,
		AAT:                aat,
		Signature:          hashSignature,
	}
}
