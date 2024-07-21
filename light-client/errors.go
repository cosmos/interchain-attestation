package lightclient

import (
	errorsmod "cosmossdk.io/errors"
)

// Pessimistic Light Client sentinel errors
var (
	ErrInvalidChainID          = errorsmod.Register(ModuleName, 2, "invalid chain-id")
	ErrInvalidRequiredTokenPower  = errorsmod.Register(ModuleName, 3, "invalid required token power")
	ErrInvalidHeaderHeight     = errorsmod.Register(ModuleName, 4, "invalid header height")
)
