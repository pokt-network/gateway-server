//go:generate ffjson $GOFILE
package models

import (
	"pokt_gateway_server/pkg/common"
)

type AAT struct {
	Version      string `json:"version"`
	AppPubKey    string `json:"app_pub_key"`
	ClientPubKey string `json:"client_pub_key"`
	Signature    string `json:"signature"`
}

func (a AAT) Hash() string {
	return common.Sha3_256HashHex(a)
}
