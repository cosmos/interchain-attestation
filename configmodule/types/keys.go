package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name
	ModuleName = "attestationconfig"

	// StoreKey is the store key string for configmodule
	StoreKey = ModuleName
)

var (
	// ParamsKey is the prefix for configmodule parameters
	ParamsKey = collections.NewPrefix(0)
	AttestatorsKey = collections.NewPrefix(1)
)
