//go:generate ffjson $GOFILE
package models

import "math/rand"

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

func FindNodeFromPubKey(nodes []*Node, pubKey string) *Node {
	for _, n := range nodes {
		if n.PublicKey == pubKey {
			return n
		}
	}
	return nil
}

func GetRandomNode(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	randomIndex := rand.Intn(len(nodes))

	return nodes[randomIndex]
}
