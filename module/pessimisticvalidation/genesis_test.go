package pessimisticvalidation_test

import (
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenesis(t *testing.T) {
	testCases := []struct {
		name string
		genesis types.GenesisState
		// expectedError bool // TODO: Add later when this can fail
	} {
		{
			name: "default",
			genesis: *types.DefaultGenesisState(),
		},
		{
			name: "custom",
			genesis: types.GenesisState{
				Params: &types.Params{
					 MinimumPower: 42,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suite := createTestSuite(t)
			app := suite.App
			ctx := app.BaseApp.NewContext(false)

			pessimisticvalidation.InitGenesis(ctx, suite.Keeper, tc.genesis)

			exportedGenesis := pessimisticvalidation.ExportGenesis(ctx, suite.Keeper)
			require.Equal(t, tc.genesis, *exportedGenesis)
		})
	}
}
