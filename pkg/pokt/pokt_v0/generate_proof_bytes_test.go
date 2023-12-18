package pokt_v0

import (
	"github.com/stretchr/testify/assert"
	"os-gateway/pkg/pokt/pokt_v0/models"
	"testing"
)

func Test_generateRelayProof(t *testing.T) {
	type relayProofArgs struct {
		entropy        uint64
		chainId        string
		sessionHeight  uint
		servicerPubKey string
		relayMetadata  *models.RelayMeta
		reqPayload     *models.Payload
		account        *models.Ed25519Account
	}
	tests := []struct {
		name string
		args relayProofArgs
		want *models.RelayProof
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, generateRelayProof(tt.args.entropy, tt.args.chainId, tt.args.sessionHeight, tt.args.servicerPubKey, tt.args.relayMetadata, tt.args.reqPayload, tt.args.account), tt.want)
		})
	}
}
