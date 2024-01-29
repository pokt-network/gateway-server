package models

type PoktApplication struct {
	ID        string   `json:"id"`
	MaxRelays int      `json:"max_relays"`
	Chains    []string `json:"chain"`
	Address   string   `json:"address"`
}
