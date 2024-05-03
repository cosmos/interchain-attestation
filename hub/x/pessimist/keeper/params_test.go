package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "hub/testutil/keeper"
	"hub/x/pessimist/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.PessimistKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
