//go:generate ffjson $GOFILE
package models

type GetLatestBlockHeightResponse struct {
	Height uint64 `json:"height"`
}
