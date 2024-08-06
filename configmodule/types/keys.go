package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "configmodule"

	// StoreKey is the store key string for configmodule
	StoreKey = ModuleName

	ClientType = "42-pessimist"
)

var (
	// ParamsKey is the prefix for configmodule parameters
	ParamsKey = collections.NewPrefix(0)
)
