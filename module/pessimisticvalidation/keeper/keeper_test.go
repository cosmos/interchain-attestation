package keeper_test

import (
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/keeper"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

var govAcct = authtypes.NewModuleAddress(govtypes.ModuleName)

type KeeperTestSuite struct {
	suite.Suite

	cdc         codec.Codec
	ctx         sdk.Context
	keeper      keeper.Keeper
	queryClient types.QueryClient
	msgSrvr     types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	k := keeper.NewKeeper(storeService, encCfg.Codec, govAcct.String())
	err := k.Params.Set(ctx, types.DefaultParams())
	suite.Require().NoError(err)

	msr := baseapp.NewMsgServiceRouter()
	msr.SetInterfaceRegistry(encCfg.InterfaceRegistry)
	msgSrvr	:= keeper.NewMsgServerImpl(k)
	types.RegisterMsgServer(msr, msgSrvr)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	queryClient := types.NewQueryClient(queryHelper)
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServer(k))

	suite.cdc = encCfg.Codec
	suite.ctx = ctx
	suite.keeper = k
	suite.queryClient = queryClient
	suite.msgSrvr = msgSrvr
}
