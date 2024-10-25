package voteextension_test

import (
	fmt "fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	abci "github.com/cometbft/cometbft/abci/types"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"

	lightclient "github.com/cosmos/interchain-attestation/core/lightclient"
	"github.com/cosmos/interchain-attestation/core/types"
	"github.com/cosmos/interchain-attestation/core/voteextension"
	"github.com/cosmos/interchain-attestation/core/voteextension/testutil"
)

type VoteExtensionTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	encodingCfg    moduletestutil.TestEncodingConfig
	mockServer     *testutil.Server
	appModule      voteextension.AppModule
	mockUpdateFunc lightclient.TrustedClientUpdateFunc
}

func (s *VoteExtensionTestSuite) SetupSuite() {
	s.encodingCfg = moduletestutil.MakeTestEncodingConfig()

	randomPort := rand.Intn(65535-49152) + 49152
	addr := fmt.Sprintf("localhost:%d", randomPort)

	s.mockServer = testutil.NewServer()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := s.mockServer.Serve(addr)
		require.NoError(s.T(), err)
		wg.Done()
	}()
	time.Sleep(1 * time.Second)

	err := os.Setenv(voteextension.SidecarAddressEnv, addr)
	require.NoError(s.T(), err)

	s.mockServer.Response = &types.GetAttestationsResponse{
		Attestations: []types.Attestation{
			{
				AttestatorId: []byte("mock-attestor-id"),
				AttestedData: types.IBCData{
					ChainId:        "mock-chain-id",
					ClientId:       "mock-client-id",
					ClientToUpdate: "mock-client-to-update",
					Height:         clienttypes.NewHeight(1, 1),
					Timestamp:      time.Now(),
					PacketCommitments: [][]byte{
						[]byte("pckt1"),
						[]byte("pckt2"),
					},
				},
			},
		},
	}

	s.mockUpdateFunc = nilUpdateFunc // Default to no updates, change in test if you need another
	s.appModule = voteextension.NewAppModule(func(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
		return s.mockUpdateFunc(ctx, clientID, clientMsg)
	}, s.encodingCfg.Codec)

	testKey := storetypes.NewKVStoreKey("upgrade")
	s.ctx = sdktestutil.DefaultContext(testKey, storetypes.NewTransientStoreKey("transient_test")).WithLogger(log.NewLogger(os.Stdout))
}

func TestVoteExtensionTestSuite(t *testing.T) {
	suite.Run(t, new(VoteExtensionTestSuite))
}

var nilUpdateFunc = func(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	return nil
}

// panicUpdateFunc can be used for tests that want to ensure the update function is never called
var panicUpdateFunc = func(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	panic("should-not-happen")
}

func (s *VoteExtensionTestSuite) TestExtendVote() {
	responseExtendVote, err := s.appModule.ExtendVote(s.ctx, &abci.RequestExtendVote{})
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), responseExtendVote.VoteExtension)

	var voteExt voteextension.VoteExtension
	err = s.encodingCfg.Codec.Unmarshal(responseExtendVote.VoteExtension, &voteExt)
	require.NoError(s.T(), err)
	require.Len(s.T(), voteExt.Attestations, 1)
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestatorId, voteExt.Attestations[0].AttestatorId)
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestedData.ChainId, voteExt.Attestations[0].AttestedData.ChainId)
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestedData.ClientId, voteExt.Attestations[0].AttestedData.ClientId)
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestedData.ClientToUpdate, voteExt.Attestations[0].AttestedData.ClientToUpdate)
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestedData.Height, voteExt.Attestations[0].AttestedData.Height)
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestedData.Timestamp.UnixNano(), voteExt.Attestations[0].AttestedData.Timestamp.UnixNano())

	require.Len(s.T(), voteExt.Attestations[0].AttestedData.PacketCommitments, 2)
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestedData.PacketCommitments[0], voteExt.Attestations[0].AttestedData.PacketCommitments[0])
	require.Equal(s.T(), s.mockServer.Response.Attestations[0].AttestedData.PacketCommitments[1], voteExt.Attestations[0].AttestedData.PacketCommitments[1])
}

func (s *VoteExtensionTestSuite) TestPreBlocker() {
	// TODO: Add a mocked light client to test with
	tests := []struct {
		name           string
		clientUpdates  *voteextension.ClientUpdates
		mockUpdateFunc lightclient.TrustedClientUpdateFunc
	}{
		{
			"success: no attestation tx",
			nil,
			panicUpdateFunc,
		},
		{
			"success: empty attestation list in tx",
			&voteextension.ClientUpdates{},
			panicUpdateFunc,
		},
		{
			"success: panic on update should not panic the app",
			&voteextension.ClientUpdates{
				ClientUpdates: []voteextension.ClientUpdate{
					{
						ClientToUpdate: "non-existent",
						AttestationClaim: lightclient.AttestationClaim{
							Attestations: []types.Attestation{
								{
									AttestatorId: []byte("whatever"),
									AttestedData: types.IBCData{
										ChainId:           "whateverchain",
										ClientId:          "whateverclient",
										ClientToUpdate:    "non-existent",
										Height:            clienttypes.Height{},
										Timestamp:         time.Now(),
										PacketCommitments: [][]byte{},
									},
								},
							},
						},
					},
				},
			},
			panicUpdateFunc,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := &abci.RequestFinalizeBlock{
				Txs: [][]byte{},
			}
			if tt.clientUpdates != nil {
				txBz, err := s.encodingCfg.Codec.Marshal(tt.clientUpdates)
				s.Require().NoError(err)
				req.Txs = append(req.Txs, txBz)
			}
			s.mockUpdateFunc = tt.mockUpdateFunc
			err := s.appModule.PreBlocker(s.ctx, req, 0)
			require.NoError(s.T(), err)
		})
	}
}
