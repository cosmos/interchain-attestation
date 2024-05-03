package keeper

import (
	"hub/x/pessimist/types"
)

var _ types.QueryServer = Keeper{}
