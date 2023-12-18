package pokt_v0

import (
	"github.com/stretchr/testify/assert"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"testing"
)

func Test_generateRelayProof(t *testing.T) {

	account, err := models.NewAccount("3fe64039816c44e8872e4ef981725b968422e3d49e95a1eb800707591df30fe374039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849")
	assert.Equal(t, err, nil)

	chainId := "0001"
	sessionHeight := uint(1)
	servicerPubKey := "0x"
	relayMetadata := &models.RelayMeta{BlockHeight: sessionHeight}
	requestPayload := &models.Payload{
		Data:    "randomJsonPayload",
		Method:  "post",
		Path:    "",
		Headers: nil,
	}
	entropy := uint64(1)
	assert.Equal(t, generateRelayProof(entropy, chainId, sessionHeight, servicerPubKey, relayMetadata, requestPayload, account), &models.RelayProof{
		Entropy:            1,
		SessionBlockHeight: 1,
		ServicerPubKey:     "0x",
		Blockchain:         "0001",
		AAT: &models.AAT{
			Version:      "0.0.1",
			AppPubKey:    "74039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849",
			ClientPubKey: "74039dbe881dd2744e2e0c469cc2241e1e45f14af6975dd89079d22938377849",
			Signature:    "f233ca857b4ada2ca4996e0da8c1761cfbc855edf282fc5a753d4631785946d6c2b08c781c84abbca2dc929de50008729079124e5c5c16921a81139279020a05",
		},
		Signature:   "befcc42130fb9e46fb9874acfb5bd8a9f783db60f86d1b1eb61cdba23fdb7e9e17544cb99afb480c9e1308532e07cdf6f4e2da27790f47dae30133725191b309",
		RequestHash: "c5b64f9a7901ed8c3341f7440913a5ddd7b694dc7b4daeb234a47a9c42b653bb",
	})

}
