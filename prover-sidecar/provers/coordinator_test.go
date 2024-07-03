package provers

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"proversidecar/provers/chainprover"
	"sync"
	"testing"
	"time"
)

const mockChainID = "mockChainID"

type MockChainProver struct {
	MockCollectProofsFunc func(ctx context.Context) error
	CollectProofsCalledCount int
}

var _ chainprover.ChainProver = &MockChainProver{}

func (m *MockChainProver) ChainID() string {
	return mockChainID
}

func (m *MockChainProver) CollectProofs(ctx context.Context) error {
	m.CollectProofsCalledCount++
	if m.MockCollectProofsFunc != nil {
		return m.MockCollectProofsFunc(ctx)
	}
	return nil
}

func (m *MockChainProver) GetProof() []byte {
	return nil
}

func TestCoordinator_Run(t *testing.T) {
	mockChainProvider := &MockChainProver{}
	coordinator := &Coordinator{
		chainProvers: map[string]chainprover.ChainProver{
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

	require.Equal(t, 3, mockChainProvider.CollectProofsCalledCount)
}


