package configmodule_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/cosmos/cosmos-sdk/x/auth"                   // for side effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"                   // for side effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus"              // for side effects
	_ "github.com/cosmos/cosmos-sdk/x/staking"                // for side effects
	_ "github.com/cosmos/interchain-attestation/configmodule" // for side effects

	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	sdklog "cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/configurator"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	configmodulemodulev1 "github.com/cosmos/interchain-attestation/configmodule/api/configmodule/module/v1"
	"github.com/cosmos/interchain-attestation/configmodule/keeper"
	"github.com/cosmos/interchain-attestation/configmodule/types"
)

type suite struct {
	App *runtime.App

	AccountKeeper authkeeper.AccountKeeper
	Keeper        keeper.Keeper
}

func createTestSuite(t *testing.T) suite {
	t.Helper()
	res := suite{}

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
						Config: appconfig.WrapAny(&configmodulemodulev1.Module{}),
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
