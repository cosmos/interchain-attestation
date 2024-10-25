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
	mockChainID           = "mockChainID"
	mockClientID          = "mockClientID"
	mockChainAttestatorID = "mockChainAttestatorID"
)

type mockCoordinator struct{}

type mockChainAttestator struct{}

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

func (m mockCoordinator) GetLatestAttestations() ([]types.Attestation, error) {
	return []types.Attestation{
		{
			AttestatorId: []byte(mockChainAttestatorID),
			AttestedData: types.IBCData{
				ChainId:           mockChainID,
				ClientId:          mockClientID,
				Height:            clienttypes.NewHeight(1, 42),
				Timestamp:         time.Now(),
				PacketCommitments: [][]byte{{0x01}, {0x02}, {0x03}},
			},
		},
	}, nil
}

func (m mockCoordinator) GetAttestationForHeight(chainID string, height uint64) (types.Attestation, error) {
	// TODO implement me
	panic("implement me")
}

func (m mockChainAttestator) ChainID() string {
	return mockChainID
}

func (m mockChainAttestator) CollectAttestation(ctx context.Context) (types.Attestation, error) {
	panic("should not be called in this test")
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
	resp, err := sidecarClient.GetAttestations(context.Background(), &types.GetAttestationsRequest{})

	require.NoError(t, err)
	require.Len(t, resp.Attestations, 1)
	require.Equal(t, []byte(mockChainAttestatorID), resp.Attestations[0].AttestatorId)
	require.Equal(t, mockChainID, resp.Attestations[0].AttestedData.ChainId)
	require.Equal(t, mockClientID, resp.Attestations[0].AttestedData.ClientId)

	s.Stop()

	wg.Wait()
}
