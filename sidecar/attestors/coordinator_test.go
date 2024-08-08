package attestors

import (
	"context"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/dgraph-io/badger/v4"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/chainattestor"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sync"
	"testing"
	"time"
)

const (
	mockChainID = "mockChainID"
	mockAttestorID = "mockAttestorID"
	mockAttestorSignature = "mockAttestorSignature"
)

var mockPacketCommits = [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

type MockChainAttestor struct {
	CurrentHeight      uint64
	Timestamp time.Time

	lock sync.Mutex
}

var _ chainattestor.ChainAttestor = &MockChainAttestor{}

func (m *MockChainAttestor) ChainID() string {
	return mockChainID
}

func (m *MockChainAttestor) CollectAttestation(ctx context.Context) (types.Attestation, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return types.Attestation{
		AttestatorId: []byte(mockAttestorID),
		AttestedData: types.IBCData{
			Height:            clienttypes.NewHeight(1, m.CurrentHeight),
			Timestamp:         m.Timestamp,
			PacketCommitments: mockPacketCommits,
		},
		Signature: []byte(mockAttestorSignature),
	}, nil
}

func (m *MockChainAttestor) updateHeight(height uint64, timestamp time.Time) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.CurrentHeight = height
	m.Timestamp = timestamp
}

func TestCoordinator_Run(t *testing.T) {
	mockChainAttestor := &MockChainAttestor{}
	mockChainAttestor.CurrentHeight = 1
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(t, err)
	testCoordinator := &coordinator{
		chainProvers: map[string]chainattestor.ChainAttestor{
			mockChainID: mockChainAttestor,
		},
		logger: zap.NewNop(),
		db:     db,
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
		mockChainAttestor.updateHeight(height, timestamp)
		timestampAtHeight[height] = timestamp
		time.Sleep(250 * time.Millisecond)

		latestAttestation, err := testCoordinator.GetLatestAttestation(mockChainID)
		require.NoError(t, err)
		require.Equal(t, height, latestAttestation.AttestedData.Height.RevisionHeight)

		attestationAtHeight, err := testCoordinator.GetAttestationForHeight(mockChainID, height)
		require.NoError(t, err)
		require.Equal(t, height, attestationAtHeight.AttestedData.Height.RevisionHeight)
		require.Equal(t, latestAttestation, attestationAtHeight)

		for j := 1; j <= i; j++ {
			attestationAtHeight, err := testCoordinator.GetAttestationForHeight(mockChainID, uint64(j))
			require.NoError(t, err)
			require.Equal(t, uint64(j), attestationAtHeight.AttestedData.Height.RevisionHeight)
		}
	}

	ctxCancel()
	wg.Wait()
}
