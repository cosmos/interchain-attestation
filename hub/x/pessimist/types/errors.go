package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/pessimist module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")

	ErrInvalidCommitteeProposal = sdkerrors.Register(ModuleName, 1, "invalid committee proposal")
	ErrValidatorNotBonded       = sdkerrors.Register(ModuleName, 2, "validator is not bonded")
	ErrObjectiveNotFound        = sdkerrors.Register(ModuleName, 3, "objective not found")
	ErrClientNotActive          = sdkerrors.Register(ModuleName, 4, "client is not active")
	ErrNotSupported             = sdkerrors.Register(ModuleName, 5, "not supported")
)
