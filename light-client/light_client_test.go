package lightclient_test

import (
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/gjermundgaraba/pessimistic-validation/lightclient"
	testifysuite "github.com/stretchr/testify/suite"
	"testing"
	"time"
)

var (
	initialClientState = lightclient.NewClientState(
		"testchain-1",
		sdkmath.NewInt(100),
		clienttypes.Height{},
		clienttypes.NewHeight(1, 42),
	)
	initialConsensusState = lightclient.NewConsensusState(
		time.Now(),
		[][]byte{},
	)
)

type PessimisticLightClientTestSuite struct {
	testifysuite.Suite

	lightClientModule lightclient.LightClientModule
	storeProvider    ibcexported.ClientStoreProvider

	ctx        sdk.Context
	encCfg moduletestutil.TestEncodingConfig
}

func TestPessimisticLightClientTestSuite(t *testing.T) {
	testifysuite.Run(t, new(PessimisticLightClientTestSuite))
}

func (s *PessimisticLightClientTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(ibcexported.StoreKey)
	s.storeProvider = clienttypes.NewStoreProvider(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	s.encCfg = moduletestutil.MakeTestEncodingConfig(lightclient.AppModuleBasic{})

	s.lightClientModule = lightclient.NewLightClientModule(s.encCfg.Codec)
	s.lightClientModule.RegisterStoreProvider(s.storeProvider)
}

func createClientID(n int) string {
	return fmt.Sprintf("%s-%d", lightclient.ModuleName, n)
}

// getClientState retrieves the client state from the store using the provided KVStore and codec.
// it does no checking if the client store or client state exists.
func getClientState(store storetypes.KVStore, cdc codec.BinaryCodec) *lightclient.ClientState {
	bz := store.Get(host.ClientStateKey())
	clientStateI := clienttypes.MustUnmarshalClientState(cdc, bz)
	clientState, ok := clientStateI.(*lightclient.ClientState)
	if !ok {
		return nil
	}
	return clientState
}