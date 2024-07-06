package pessimisticvalidation_test

import (
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	sdklog "cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/configurator"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	_ "github.com/cosmos/cosmos-sdk/x/auth" // for side effects
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/bank"      // for side effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // for side effects
	_ "github.com/cosmos/cosmos-sdk/x/staking"   // for side effects
	pessimisticvalidationmodulev1 "github.com/gjermundgaraba/pessimistic-validation/api/pessimisticvalidation/module/v1"
	_ "github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation" // for side effects
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/keeper"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type suite struct {
	App    *runtime.App

	AccountKeeper      authkeeper.AccountKeeper
	Keeper keeper.Keeper
}

func createTestSuite(t *testing.T) suite {
	res := suite{}

	//config := simtestutil.DefaultStartUpConfig()

	app, err := simtestutil.SetupWithConfiguration(
		depinject.Configs(
			configurator.NewAppConfig(
				configurator.AuthModule(),
				configurator.StakingModule(),
				configurator.ConsensusModule(),
				configurator.BankModule(),
				func(config *configurator.Config) {
					config.ModuleConfigs[types.ModuleName] = &appv1alpha1.ModuleConfig{
						Name:   types.ModuleName,
						Config: appconfig.WrapAny(&pessimisticvalidationmodulev1.Module{}),
					}
					config.InitGenesisOrder = append(config.InitGenesisOrder, types.ModuleName)
				},
			),
			depinject.Supply(sdklog.NewNopLogger()),
		),
		simtestutil.DefaultStartUpConfig(),
		&res.AccountKeeper, &res.Keeper,
	)
	require.NoError(t, err)

	res.App = app

	return res
}
