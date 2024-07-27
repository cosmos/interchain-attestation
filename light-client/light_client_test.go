package lightclient_test

import (
	"context"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdkcrypto "github.com/cosmos/cosmos-sdk/crypto/types"
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
	storeProvider     ibcexported.ClientStoreProvider

	ctx    sdk.Context
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

	attestators := generateAttestators(10)
	attestatorsHandler := NewMockAttestatorsHandler(attestators)

	s.lightClientModule = lightclient.NewLightClientModule(s.encCfg.Codec, attestatorsHandler)
	s.lightClientModule.RegisterStoreProvider(s.storeProvider)
}

type mockAttestator struct {
	id         []byte
	publicKey  sdkcrypto.PubKey
	privateKey sdkcrypto.PrivKey
}

type mockAttestatorsHandler struct {
	attestators map[string]mockAttestator
	sufficientAttestations func() (bool, error)
}

var _ lightclient.AttestatorsHandler = &mockAttestatorsHandler{}

func NewMockAttestatorsHandler(attestators []mockAttestator) mockAttestatorsHandler {
	attestatorsMap := make(map[string]mockAttestator)
	for _, attestator := range attestators {
		attestatorsMap[string(attestator.id)] = attestator
	}
	return mockAttestatorsHandler{
		attestators: attestatorsMap,
		sufficientAttestations: func() (bool, error) {
			return true, nil
		},
	}
}

func (m mockAttestatorsHandler) GetPublicKey(ctx context.Context, attestatorId []byte) (sdkcrypto.PubKey, error) {
	attestator, ok := m.attestators[string(attestatorId)]
	if !ok {
		return nil, fmt.Errorf("attestator not found")
	}
	return attestator.publicKey, nil
}

func (m mockAttestatorsHandler) SufficientAttestations(_ context.Context, _ [][]byte) (bool, error) {
	return m.sufficientAttestations()
}

func (m mockAttestatorsHandler) getMockAttestator(attestatorId []byte) mockAttestator {
	attestator, ok := m.attestators[string(attestatorId)]
	if !ok {
		panic("attestator not found")
	}
	return attestator
}

func (m mockAttestatorsHandler) reSignClaim(claim *lightclient.SignedPacketCommitmentsClaim) {
	attestator := m.getMockAttestator(claim.AttestatorId)
	signableBytes := lightclient.GetSignableBytes(claim.PacketCommitmentsClaim.PacketCommitments)
	signature, err := attestator.privateKey.Sign(signableBytes)
	if err != nil {
		panic(err)
	}
	claim.Signature = signature
}

func generateAttestators(n int) []mockAttestator {
	attestators := make([]mockAttestator, n)
	for i := 0; i < n; i++ {
		privKey := secp256k1.GenPrivKey()
		valAddr := sdk.ValAddress(privKey.PubKey().Address())
		attestators[i] = mockAttestator{
			id:         valAddr,
			publicKey:  privKey.PubKey(),
			privateKey: privKey,
		}
	}
	return attestators
}

func generateClientMsg(attestators []mockAttestator, numberOfPacketCommitments int) *lightclient.PessimisticClaims {
	claims := make([]lightclient.SignedPacketCommitmentsClaim, len(attestators))
	packetCommitments := generatePacketCommitments(numberOfPacketCommitments)
	signableBytes := lightclient.GetSignableBytes(packetCommitments)

	for i, attestator := range attestators {
		signature, err := attestator.privateKey.Sign(signableBytes)
		if err != nil {
			panic(err)
		}

		// Copy so that the test can modify the packet commitments without affecting the other claims
		packetCommitementsCopy := make([][]byte, len(packetCommitments))
		copy(packetCommitementsCopy, packetCommitments)
		claims[i] = lightclient.SignedPacketCommitmentsClaim{
			AttestatorId: attestator.id,
			PacketCommitmentsClaim: lightclient.PacketCommitmentsClaim{
				PacketCommitments: packetCommitementsCopy,
				Height: clienttypes.Height{
					RevisionNumber: 1,
					RevisionHeight: 42,
				},
			},
			Signature: signature,
		}
	}
	return &lightclient.PessimisticClaims{
		Claims: claims,
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

