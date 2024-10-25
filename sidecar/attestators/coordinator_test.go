package attestators

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"

	"github.com/cosmos/interchain-attestation/core/types"
	"github.com/cosmos/interchain-attestation/sidecar/attestators/attestator"
)

const (
	mockChainID        = "mockChainID"
	mockClientID       = "mockClientID"
	mockClientToUpdate = "mockClientToUpdate"
	mockAttestatorID   = "mockAttestatorID"
)

var mockPacketCommits = [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

type MockChainAttestator struct {
	CurrentHeight uint64
	Timestamp     time.Time

	lock sync.Mutex
}

var _ attestator.Attestator = &MockChainAttestator{}

func (m *MockChainAttestator) ChainID() string {
	return mockChainID
}

func (m *MockChainAttestator) CollectAttestation(ctx context.Context) (types.Attestation, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return types.Attestation{
		AttestatorId: []byte(mockAttestatorID),
		AttestedData: types.IBCData{
			ChainId:           mockChainID,
			ClientId:          mockClientID,
			ClientToUpdate:    mockClientToUpdate,
			Height:            clienttypes.NewHeight(1, m.CurrentHeight),
			Timestamp:         m.Timestamp,
			PacketCommitments: mockPacketCommits,
		},
	}, nil
}

func (m *MockChainAttestator) updateHeight(height uint64, timestamp time.Time) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.CurrentHeight = height
	m.Timestamp = timestamp
}

func TestCoordinator_Run(t *testing.T) {
	mockChainAttestator := &MockChainAttestator{}
	mockChainAttestator.CurrentHeight = 1
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(t, err)
	testCoordinator := &coordinator{
		chainAttestators: map[string]attestator.Attestator{
			mockChainID: mockChainAttestator,
		},
		logger:            zap.NewNop(),
		db:                db,
		queryLoopDuration: 50 * time.Millisecond,
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := testCoordinator.Run(ctx)
		require.NoError(t, err)
		wg.Done()
	}()

	timestampAtHeight := make(map[uint64]time.Time)
	for i := 1; i <= 10; i++ {
		height := uint64(i)
		timestamp := time.Now()
		mockChainAttestator.updateHeight(height, timestamp)
		timestampAtHeight[height] = timestamp
		time.Sleep(250 * time.Millisecond)

		latestAttestations, err := testCoordinator.GetLatestAttestations()
		require.NoError(t, err)
		require.Len(t, latestAttestations, 1)
		require.Equal(t, height, latestAttestations[0].AttestedData.Height.RevisionHeight)
		require.Equal(t, mockPacketCommits, latestAttestations[0].AttestedData.PacketCommitments)
		require.Equal(t, mockAttestatorID, string(latestAttestations[0].AttestatorId))
		require.Equal(t, mockChainID, latestAttestations[0].AttestedData.ChainId)
		require.Equal(t, mockClientID, latestAttestations[0].AttestedData.ClientId)
		require.Equal(t, mockClientToUpdate, latestAttestations[0].AttestedData.ClientToUpdate)
		require.Equal(t, timestampAtHeight[height].UnixNano(), latestAttestations[0].AttestedData.Timestamp.UnixNano())

		attestationAtHeight, err := testCoordinator.GetAttestationForHeight(mockChainID, height)
		require.NoError(t, err)
		require.Equal(t, height, attestationAtHeight.AttestedData.Height.RevisionHeight)
		require.Equal(t, latestAttestations[0], attestationAtHeight)

		for j := 1; j <= i; j++ {
			attestationAtHeight, err := testCoordinator.GetAttestationForHeight(mockChainID, uint64(j))
			require.NoError(t, err)
			require.Equal(t, uint64(j), attestationAtHeight.AttestedData.Height.RevisionHeight)
		}
	}

	ctxCancel()
	wg.Wait()
}
