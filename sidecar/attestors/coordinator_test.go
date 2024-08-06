package attestors

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/chainattestor"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"sync"
	"testing"
	"time"
)

const mockChainID = "mockChainID"

type MockChainAttestor struct {
	MockCollectFunc    func(ctx context.Context) error
	CollectCalledCount int
}

var _ chainattestor.ChainAttestor = &MockChainAttestor{}

func (m *MockChainAttestor) ChainID() string {
	return mockChainID
}

func (m *MockChainAttestor) CollectAttestations(ctx context.Context) error {
	m.CollectCalledCount++
	if m.MockCollectFunc != nil {
		return m.MockCollectFunc(ctx)
	}
	return nil
}

func (m *MockChainAttestor) GetLatestAttestation() *types.Attestation {
	return nil
}

func TestCoordinator_Run(t *testing.T) {
	mockChainProvider := &MockChainAttestor{}
	coordinator := &coordinator{
		chainProvers: map[string]chainattestor.ChainAttestor{
			mockChainID: mockChainProvider,
		},
		logger: zap.NewNop(),
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := coordinator.Run(ctx)
		require.NoError(t, err)
		wg.Done()
	}()

	time.Sleep(3 * time.Second)
	ctxCancel()
	wg.Wait()

	require.Equal(t, 3, mockChainProvider.CollectCalledCount)
}


