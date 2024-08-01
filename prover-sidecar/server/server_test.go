package server_test

import (
	"context"
	"fmt"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/provers"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/provers/chainprover"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/server"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/types"
	"sync"
	"testing"
	"time"
)

type mockCoordinator struct{}
type mockChainProver struct{}

var _ provers.Coordinator = &mockCoordinator{}
var _ chainprover.ChainProver = &mockChainProver{}

func (m mockCoordinator) GetChainProver(chainID string) chainprover.ChainProver {
	return &mockChainProver{}
}

func (m mockCoordinator) Run(ctx context.Context) error {
	panic("should not be called in this test")
}

func (m mockChainProver) ChainID() string {
	return "mockChainID"
}

func (m mockChainProver) CollectProofs(ctx context.Context) error {
	panic("should not be called in this test")
}

func (m mockChainProver) GetProof() *types.SignedPacketCommitmentsClaim {
	return &types.SignedPacketCommitmentsClaim{
		AttestatorId: []byte("mockAttestatorID"),
		PacketCommitmentsClaim: types.PacketCommitmentsClaim{
			Height:            clienttypes.NewHeight(1, 42),
			Timestamp:         time.Now(),
			PacketCommitments: [][]byte{{0x01}, {0x02}, {0x03}},
		},
		Signature: []byte("mockSignature"),
	}
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

	proofClient := server.NewProofClient(client)
	proof, err := proofClient.GetProof(context.Background(), &server.ProofRequest{
		ChainId: "mockChainID",
	})
	require.NoError(t, err)
	require.NotNil(t, proof)

	s.Stop()

	wg.Wait()
}
