package keeper_test

import (
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/keeper"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/testutil"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

const testValidatorAddress = "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8"

var govAcct = authtypes.NewModuleAddress(govtypes.ModuleName)

type KeeperTestSuite struct {
	suite.Suite

	cdc         codec.Codec
	ctx         sdk.Context
	keeper      keeper.Keeper
	queryClient types.QueryClient
	msgSrvr     types.MsgServer

	mockValidator stakingtypes.Validator
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := sdktestutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()
	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	ctrl := gomock.NewController(suite.T())
	stakingKeeper := testutil.NewMockStakingKeeper(ctrl)

	consPubKey := ed25519.GenPrivKey().PubKey()
	suite.Require().NotNil(consPubKey)
	mockValidator, err := stakingtypes.NewValidator(testValidatorAddress, consPubKey, stakingtypes.Description{})
	suite.Require().NoError(err)
	suite.mockValidator = mockValidator

	validatorAddressCodec := address.NewBech32Codec("cosmosvaloper")
	valAddr, err := validatorAddressCodec.StringToBytes(testValidatorAddress)
	suite.Require().NoError(err)
	stakingKeeper.EXPECT().GetValidator(ctx, valAddr).Return(mockValidator, nil).AnyTimes()

	k := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		validatorAddressCodec,
		govAcct.String(),
		stakingKeeper,
	)
	err = k.Params.Set(ctx, types.DefaultParams())
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	queryClient := types.NewQueryClient(queryHelper)
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServer(k))

	msr := baseapp.NewMsgServiceRouter()
	msr.SetInterfaceRegistry(encCfg.InterfaceRegistry)
	msgSrvr := keeper.NewMsgServerImpl(k)
	types.RegisterMsgServer(msr, msgSrvr)

	suite.cdc = encCfg.Codec
	suite.ctx = ctx
	suite.keeper = k
	suite.queryClient = queryClient
	suite.msgSrvr = msgSrvr
}

// TODO: Support creating multiple attestators. Currently calling this multiple times will fail
func (suite *KeeperTestSuite) registerAttestator() []byte {
	pubKey := secp256k1.GenPrivKey().PubKey()
	suite.Require().NotNil(pubKey)
	pubKeyAny, err := codectypes.NewAnyWithValue(pubKey)
	suite.Require().NoError(err)

	attestatorID := []byte("attestator-id")

	_, err = suite.msgSrvr.RegisterAttestator(suite.ctx, &types.MsgRegisterAttestator{
		ValidatorAddress:     testValidatorAddress,
		AttestatorId:         attestatorID,
		AttestationPublicKey: pubKeyAny,
	})
	suite.Require().NoError(err)

	return attestatorID
}
