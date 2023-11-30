//go:generate ffjson $GOFILE
package common

type EvmJsonRpcPayload struct {
	Id     string `json:"id"`
	Method string `json:"method"`
}
