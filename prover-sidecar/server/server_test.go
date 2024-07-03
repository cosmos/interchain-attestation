package server_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"math/rand"
	"proversidecar/server"
	"sync"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	s := server.NewServer(zap.NewNop(), nil)
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

	s.Stop()

	wg.Wait()
}
