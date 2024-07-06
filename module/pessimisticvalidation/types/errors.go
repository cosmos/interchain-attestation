package types

import "cosmossdk.io/errors"

// pessimisticvalidation module sentinel errors
var ErrUnauthorized = errors.Register(ModuleName, 1, "unauthorized")
