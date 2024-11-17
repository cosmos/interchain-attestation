package lightclient_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	suite "github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v9/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v9/modules/core/exported"

	"github.com/cosmos/interchain-attestation/core/lightclient"
	"github.com/cosmos/interchain-attestation/core/types"
)

const (
	mockChainID  = "testchain-1"
	mockClientID = "testclient-1"
)

var (
	initialClientState = lightclient.NewClientState(
		mockChainID,
		sdkmath.NewInt(100),
		clienttypes.Height{},
		clienttypes.NewHeight(1, 42),
	)
	initialConsensusState = lightclient.NewConsensusState(
		time.Now(),
	)
	defaultHeight = clienttypes.NewHeight(1, 42)
)

type AttestationLightClientTestSuite struct {
	suite.Suite

	lightClientModule lightclient.LightClientModule
	trustedUpdateFunc lightclient.TrustedClientUpdateFunc
	storeProvider     clienttypes.StoreProvider

	mockAttestators        []mockAttestator
	mockAttestatorsHandler mockAttestatorsHandler

	ctx     sdk.Context
	testCtx testutil.TestContext
	encCfg  moduletestutil.TestEncodingConfig
}

func TestAttestationLightClientTestSuite(t *testing.T) {
	suite.Run(t, new(AttestationLightClientTestSuite))
}

func (s *AttestationLightClientTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(ibcexported.StoreKey)
	s.storeProvider = clienttypes.NewStoreProvider(key)
	s.testCtx = testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = s.testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	s.encCfg = moduletestutil.MakeTestEncodingConfig(lightclient.AppModuleBasic{})

	s.mockAttestators = generateAttestators(10)
	s.mockAttestatorsHandler = NewMockAttestatorsHandler(s.mockAttestators)

	s.lightClientModule, s.trustedUpdateFunc = lightclient.NewLightClientModule(s.encCfg.Codec, s.storeProvider, s.mockAttestatorsHandler)
}

type mockAttestator struct {
	validatorAddress []byte
}

type mockAttestatorsHandler struct {
	attestators            map[string]mockAttestator
	sufficientAttestations func() (bool, error)
}

var _ lightclient.AttestatorsController = &mockAttestatorsHandler{}

func NewMockAttestatorsHandler(attestators []mockAttestator) mockAttestatorsHandler {
	attestatorsMap := make(map[string]mockAttestator)
	for _, attestator := range attestators {
		attestatorsMap[string(attestator.validatorAddress)] = attestator
	}
	return mockAttestatorsHandler{
		attestators: attestatorsMap,
		sufficientAttestations: func() (bool, error) {
			return true, nil
		},
	}
}

func (m mockAttestatorsHandler) SufficientAttestations(_ context.Context, _ [][]byte) (bool, error) {
	return m.sufficientAttestations()
}

func generateAttestators(n int) []mockAttestator {
	attestators := make([]mockAttestator, n)
	for i := 0; i < n; i++ {
		privKey := secp256k1.GenPrivKey()
		valAddr := sdk.ValAddress(privKey.PubKey().Address())
		attestators[i] = mockAttestator{
			validatorAddress: valAddr,
		}
	}
	return attestators
}

func generateClientMsg(_ codec.BinaryCodec, attestators []mockAttestator, numberOfPacketCommitments int, modifiers ...func(dataToAttestTo *types.IBCData)) *lightclient.AttestationClaim {
	attestations := make([]types.Attestation, len(attestators))
	packetCommitments := generatePacketCommitments(numberOfPacketCommitments)
	timestamp := time.Now()

	for i, attestator := range attestators {
		// Copy so that the test can modify the packet commitments without affecting the other attestations
		packetCommitementsCopy := make([][]byte, len(packetCommitments))
		copy(packetCommitementsCopy, packetCommitments)

		attestationData := types.IBCData{
			ChainId:           mockChainID,
			ClientId:          mockClientID,
			Height:            defaultHeight,
			Timestamp:         timestamp,
			PacketCommitments: packetCommitementsCopy,
		}

		for _, modifier := range modifiers {
			modifier(&attestationData)
		}

		attestations[i] = types.Attestation{
			ValidatorAddress: attestator.validatorAddress,
			AttestedData:     attestationData,
		}
	}
	return &lightclient.AttestationClaim{
		Attestations: attestations,
	}
}

func generatePacketCommitments(n int) [][]byte {
	packetCommitments := make([][]byte, n)
	for i := 0; i < n; i++ {
		packetCommitments[i] = []byte(fmt.Sprintf("packet commitment %d", i))
	}
	return packetCommitments
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

// getConsensusState retrieves the consensus state from the client prefixed store.
// It does no checking if the consensus state exists.
func getConsensusState(store storetypes.KVStore, cdc codec.BinaryCodec, height ibcexported.Height) *lightclient.ConsensusState {
	bz := store.Get(host.ConsensusStateKey(height))
	consensusStateI := clienttypes.MustUnmarshalConsensusState(cdc, bz)
	consensusState, ok := consensusStateI.(*lightclient.ConsensusState)
	if !ok {
		return nil
	}
	return consensusState
}
