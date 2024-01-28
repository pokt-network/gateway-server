package pokt_v0

import (
	"encoding/hex"
	"pokt_gateway_server/pkg/common"
	"pokt_gateway_server/pkg/pokt/pokt_v0/models"
)

// generateRelayProof generates a relay proof.
// Parameters:
//   - entropy - random generated number to signify unique proof
//   - chainId: Blockchain ID.
//   - sessionHeight: Session block height.
//   - servicerPubKey: Servicer public key.
//   - requestMetadata: Request metadata.
//   - account: Ed25519 account used for signing.
//
// Returns:
//   - models.RelayProof: Generated relay proof.
func generateRelayProof(entropy uint64, chainId string, sessionHeight uint, servicerPubKey string, relayMetadata *models.RelayMeta, reqPayload *models.Payload, account *models.Ed25519Account) *models.RelayProof {
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
