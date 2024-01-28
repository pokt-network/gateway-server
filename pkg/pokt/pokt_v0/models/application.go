//go:generate ffjson $GOFILE
package models

import (
	"encoding/json"
	"strconv"
)

type PoktApplicationStatus uint

const (
	StatusJailed    PoktApplicationStatus = 0
	StatusUnstaking PoktApplicationStatus = 1
	StatusStaked    PoktApplicationStatus = 2
)

// MaxRelays is a custom type for handling the MaxRelays field.
type MaxRelays int

// UnmarshalJSON customizes the JSON unmarshalling for MaxRelays.
func (mr *MaxRelays) UnmarshalJSON(data []byte) error {
	var maxRelaysString string
	if err := json.Unmarshal(data, &maxRelaysString); err != nil {
		return err
	}

	// Parse the string into an integer
	maxRelays, err := strconv.Atoi(maxRelaysString)
	if err != nil {
		return err
	}

	// Set the MaxRelays field with the parsed integer
	*mr = MaxRelays(maxRelays)

	return nil
}

type GetApplicationResponse struct {
	Result []*PoktApplication `json:"result"`
}

type PoktApplication struct {
	Address   string                `json:"address"`
	Chains    []string              `json:"chains"`
	PublicKey string                `json:"public_key"`
	Status    PoktApplicationStatus `json:"status"`
	MaxRelays MaxRelays             `json:"max_relays"`
}
