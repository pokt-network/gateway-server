//go:generate ffjson $GOFILE
package models

import (
	"github.com/pquerna/ffjson/ffjson"
	"pokt_gateway_server/pkg/common"
	"time"
)

type SendRelayRequest struct {
	Payload            *Payload
	Signer             *Ed25519Account
	Chain              string
	SelectedNodePubKey string
	Session            *Session
	Timeout            *time.Duration
}

func (req SendRelayRequest) Validate() error {
	if req.Payload == nil || req.Signer == nil {
		return ErrMalformedSendRelayRequest
	}
	return nil
}

type SendRelayResponse struct {
	Response string `json:"response"`
}

// ffjson: skip
type Payload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers,omitempty"`
}

// "payload" - A structure used for custom json marshalling/unmarshalling
// ffjson: skip
type payload struct {
	Data    string            `json:"data"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
}

// "MarshalJSON" - Overrides json marshalling
func (p Payload) MarshalJSON() ([]byte, error) {
	pay := payload{
		Data:    p.Data,
		Method:  p.Method,
		Path:    p.Path,
		Headers: p.Headers,
	}
	return ffjson.Marshal(pay)
}

type Relay struct {
	Payload    *Payload    `json:"payload"`
	Metadata   *RelayMeta  `json:"meta"`
	RelayProof *RelayProof `json:"proof"`
}

type RelayMeta struct {
	BlockHeight uint `json:"block_height"`
}

// RelayProof represents proof of a relay
type RelayProof struct {
	Entropy            uint64 `json:"entropy"`
	SessionBlockHeight uint   `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	AAT                *AAT   `json:"aat"`
	Signature          string `json:"signature"`
	RequestHash        string `json:"request_hash"`
}

// RequestHashPayload struct holding data needed to create a request hash
// ffjson: skip
type RequestHashPayload struct {
	Payload  *Payload   `json:"payload"`
	Metadata *RelayMeta `json:"meta"`
}

func (a *RequestHashPayload) Hash() string {
	return common.Sha3_256HashHex(a)
}

// ffjson: skip
type RelayProofHashPayload struct {
	Entropy            uint64 `json:"entropy"`
	SessionBlockHeight uint   `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Signature          string `json:"signature"`
	UnsignedAAT        string `json:"token"`
	RequestHash        string `json:"request_hash"`
}
