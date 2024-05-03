package pessimist_test

import (
	"testing"

	keepertest "hub/testutil/keeper"
	"hub/testutil/nullify"
	pessimist "hub/x/pessimist/module"
	"hub/x/pessimist/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PessimistKeeper(t)
	pessimist.InitGenesis(ctx, k, genesisState)
	got := pessimist.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
