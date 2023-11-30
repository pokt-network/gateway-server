package models

import "errors"

var (
	ErrMissingFullNodes          = errors.New("require full node host")
	ErrSessionHasZeroNodes       = errors.New("session missing valid nodes")
	ErrNodeNotFound              = errors.New("node not found")
	ErrMalformedSendRelayRequest = errors.New("malformed send relay request")
)
