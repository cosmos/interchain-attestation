package server_test

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"

	"github.com/cosmos/interchain-attestation/core/types"
	"github.com/cosmos/interchain-attestation/sidecar/attestators"
	"github.com/cosmos/interchain-attestation/sidecar/attestators/attestator"
	"github.com/cosmos/interchain-attestation/sidecar/server"
)

const (
	mockChainID  = "mockChainID"
	mockClientID = "mockClientID"
)

func (m mockChainAttestator) ChainID() string {
	return mockChainID
}

func (m mockChainAttestator) CollectIBCData(ctx context.Context) (types.IBCData, error) {
	panic("should not be called in this test")
}

type mockChainAttestator struct{}
type mockCoordinator struct{}

var (
	_ attestators.Coordinator = &mockCoordinator{}
	_ attestator.Attestator   = &mockChainAttestator{}
)

func (m mockCoordinator) GetChainProver(_ string) attestator.Attestator {
	return &mockChainAttestator{}
}

func (m mockCoordinator) Run(_ context.Context) error {
	panic("should not be called in this test")
}

func (m mockCoordinator) GetLatestIBCData() ([]types.IBCData, error) {
	return []types.IBCData{
		{
			ChainId:           mockChainID,
			ClientId:          mockClientID,
			Height:            clienttypes.NewHeight(1, 42),
			Timestamp:         time.Now(),
			PacketCommitments: [][]byte{{0x01}, {0x02}, {0x03}},
		},
	}, nil
}

func (m mockCoordinator) GetIBCDataForHeight(chainID string, height uint64) (types.IBCData, error) {
	// TODO: implement me
	panic("implement me")
}

// TestServe is mostly just a smoke test that the server can start and serve requests. Everything is mocked except the server itself.
func TestServe(t *testing.T) {
	s := server.NewServer(zap.NewNop(), mockCoordinator{})
	randomPort := rand.Intn(65535-49152) + 49152
	addr := fmt.Sprintf("localhost:%d", randomPort)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := s.Serve(addr)
		require.NoError(t, err)
		wg.Done()
	}()

	time.Sleep(1 * time.Second)

	client, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	sidecarClient := types.NewSidecarClient(client)
	resp, err := sidecarClient.GetIBCData(context.Background(), &types.GetIBCDataRequest{})

	require.NoError(t, err)
	require.Len(t, resp.IbcData, 1)
	require.Equal(t, mockChainID, resp.IbcData[0].ChainId)
	require.Equal(t, mockClientID, resp.IbcData[0].ClientId)

	s.Stop()

	wg.Wait()
}
