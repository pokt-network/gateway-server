//go:generate ffjson $GOFILE
package models

type Node struct {
	ServiceUrl string `json:"service_url"`
	PublicKey  string `json:"public_key"`
}

type SessionHeader struct {
	SessionHeight uint `json:"session_height"`
}

type GetSessionResponse struct {
	Session *Session `json:"session"`
}

type Session struct {
	Nodes         []*Node        `json:"nodes"`
	SessionHeader *SessionHeader `json:"header"`
}

type GetSessionRequest struct {
	AppPubKey string `json:"app_public_key"`
	Chain     string `json:"chain"`
	Height    uint   `json:"session_height"`
}
