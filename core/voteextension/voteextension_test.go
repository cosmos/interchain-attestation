package voteextension_test

import (
	fmt "fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"

	storetypes "cosmossdk.io/store/types"

	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	abci "github.com/cometbft/cometbft/abci/types"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"

	"github.com/cosmos/interchain-attestation/core/types"
	"github.com/cosmos/interchain-attestation/core/voteextension"
	"github.com/cosmos/interchain-attestation/core/voteextension/testutil"
)

func TestExtendVote(t *testing.T) {
	testKey := storetypes.NewKVStoreKey("upgrade")
	ctx := sdktestutil.DefaultContext(testKey, storetypes.NewTransientStoreKey("transient_test"))

	randomPort := rand.Intn(65535-49152) + 49152
	addr := fmt.Sprintf("localhost:%d", randomPort)

	mockServer := testutil.NewServer()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := mockServer.Serve(addr)
		require.NoError(t, err)
		wg.Done()
	}()
	time.Sleep(1 * time.Second)

	err := os.Setenv(voteextension.SidecarAddressEnv, addr)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	clientKeeper := testutil.NewMockClientKeeper(ctrl)

	encodingCfg := moduletestutil.MakeTestEncodingConfig()
	appModule := voteextension.NewAppModule(clientKeeper, encodingCfg.Codec)

	mockServer.Response = &types.GetAttestationsResponse{
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
				Signature: []byte("sig"),
			},
		},
	}
	responseExtendVote, err := appModule.ExtendVote(ctx, &abci.RequestExtendVote{})
	require.NoError(t, err)
	require.NotEmpty(t, responseExtendVote.VoteExtension)

	var voteExt voteextension.VoteExtension
	err = encodingCfg.Codec.Unmarshal(responseExtendVote.VoteExtension, &voteExt)
	require.NoError(t, err)
	require.Len(t, voteExt.Attestations, 1)
	require.Equal(t, mockServer.Response.Attestations[0].AttestatorId, voteExt.Attestations[0].AttestatorId)
	require.Equal(t, mockServer.Response.Attestations[0].Signature, voteExt.Attestations[0].Signature)
	require.Equal(t, mockServer.Response.Attestations[0].AttestedData.ChainId, voteExt.Attestations[0].AttestedData.ChainId)
	require.Equal(t, mockServer.Response.Attestations[0].AttestedData.ClientId, voteExt.Attestations[0].AttestedData.ClientId)
	require.Equal(t, mockServer.Response.Attestations[0].AttestedData.ClientToUpdate, voteExt.Attestations[0].AttestedData.ClientToUpdate)
	require.Equal(t, mockServer.Response.Attestations[0].AttestedData.Height, voteExt.Attestations[0].AttestedData.Height)
	require.Equal(t, mockServer.Response.Attestations[0].AttestedData.Timestamp.UnixNano(), voteExt.Attestations[0].AttestedData.Timestamp.UnixNano())

	require.Len(t, voteExt.Attestations[0].AttestedData.PacketCommitments, 2)
	require.Equal(t, mockServer.Response.Attestations[0].AttestedData.PacketCommitments[0], voteExt.Attestations[0].AttestedData.PacketCommitments[0])
	require.Equal(t, mockServer.Response.Attestations[0].AttestedData.PacketCommitments[1], voteExt.Attestations[0].AttestedData.PacketCommitments[1])
}
