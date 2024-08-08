package server_test

import (
	"context"
	"fmt"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/chainattestor"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/server"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type mockCoordinator struct{}

type mockChainAttestor struct{}

var _ attestors.Coordinator = &mockCoordinator{}
var _ chainattestor.ChainAttestor = &mockChainAttestor{}

func (m mockCoordinator) GetChainProver(_ string) chainattestor.ChainAttestor {
	return &mockChainAttestor{}
}

func (m mockCoordinator) Run(_ context.Context) error {
	panic("should not be called in this test")
}

func (m mockCoordinator) GetLatestAttestation(chainID string) (types.Attestation, error) {
	return types.Attestation{
		AttestatorId: []byte("mockAttestatorID"),
		AttestedData: types.IBCData {
			Height:            clienttypes.NewHeight(1, 42),
			Timestamp:         time.Now(),
			PacketCommitments: [][]byte{{0x01}, {0x02}, {0x03}},
		},
		Signature: []byte("mockSignature"),
	}, nil
}

func (m mockCoordinator) GetAttestationForHeight(chainID string, height uint64) (types.Attestation, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockChainAttestor) ChainID() string {
	return "mockChainID"
}

func (m mockChainAttestor) CollectAttestation(ctx context.Context) (types.Attestation, error) {
	panic("should not be called in this test")
}

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

	claimClient := types.NewSidecarClient(client)
	claim, err := claimClient.GetAttestation(context.Background(), &types.AttestationRequest{
		ChainId: "mockChainID",
	})
	require.NoError(t, err)
	require.NotNil(t, claim)

	s.Stop()

	wg.Wait()
}
