//go:generate ffjson $GOFILE
package models

type GetLatestBlockHeightResponse struct {
	Height uint `json:"height"`
}
