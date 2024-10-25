package lightclient

import (
	errorsmod "cosmossdk.io/errors"
)

// Attestation Light Client sentinel errors
var (
	ErrInvalidChainID            = errorsmod.Register(ModuleName, 2, "invalid chain-id")
	ErrInvalidRequiredTokenPower = errorsmod.Register(ModuleName, 3, "invalid required token power")
	ErrInvalidHeaderHeight       = errorsmod.Register(ModuleName, 4, "invalid header height")
	ErrInvalidClientMsg          = errorsmod.Register(ModuleName, 5, "invalid client message")
	ErrPacketCommitmentNotFound  = errorsmod.Register(ModuleName, 6, "packet commitment not found")
	ErrInvalidUpdateMethod       = errorsmod.Register(ModuleName, 7, "invalid update method, can only be done through code")
)
