package types

import "cosmossdk.io/errors"

// configmodule module sentinel errors
var (
	ErrUnauthorized            = errors.Register(ModuleName, 1, "unauthorized")
	ErrAttestatorAlreadyExists = errors.Register(ModuleName, 2, "attestator already exists")
	ErrInvalidAttestator       = errors.Register(ModuleName, 3, "invalid attestator")
)
