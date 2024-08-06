package types

import "cosmossdk.io/errors"

// configmodule module sentinel errors
var ErrUnauthorized = errors.Register(ModuleName, 1, "unauthorized")
