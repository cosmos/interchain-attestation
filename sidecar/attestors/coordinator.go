package attestors

import (
	"context"
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/chainattestor"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/cosmos"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"time"
)

const (
	defaultMinQueryLoopDuration      = 1 * time.Second
)

type Coordinator interface {
	GetChainProver(chainID string) chainattestor.ChainAttestor
	Run(ctx context.Context) error
}

type coordinator struct {
	logger *zap.Logger

	chainProvers map[string]chainattestor.ChainAttestor
}

func NewCoordinator(logger *zap.Logger, sidecarConfig config.Config) (Coordinator, error) {
	chainProvers := make(map[string]chainattestor.ChainAttestor)
	for _, cosmosConfig := range sidecarConfig.CosmosChains {
		prover, err := cosmos.NewCosmosAttestor(logger, sidecarConfig.AttestatorID, cosmosConfig.ChainID, cosmosConfig.RPC, cosmosConfig.ClientID, func(msg []byte) ([]byte, error) {
			signerPrivKeyBase64, err := os.ReadFile(sidecarConfig.SigningPrivateKeyPath)
			if err != nil {
				return nil, err
			}
			signerPrivKeyBz, err := base64.StdEncoding.DecodeString(string(signerPrivKeyBase64))
			if err != nil {
				return nil, err
			}

			signerPrivKey := secp256k1.PrivKey{
				Key: signerPrivKeyBz,
			}

			return signerPrivKey.Sign(msg)
		})
		if err != nil {
			return nil, err
		}
		chainProvers[cosmosConfig.ChainID] = prover
	}

	return &coordinator{
		logger:       logger,
		chainProvers: chainProvers,
	}, nil
}

func (c *coordinator) GetChainProver(chainID string) chainattestor.ChainAttestor {
	return c.chainProvers[chainID]
}

func (c *coordinator) Run(ctx context.Context) error {
	c.logger.Debug("Coordinator.Run")

	var eg errgroup.Group
	runCtx, runCtxCancel := context.WithCancel(ctx)
	for _, chainProver := range c.chainProvers {
		c.logger.Info("Starting chain prover loop", zap.String("chain_id", chainProver.ChainID()))

		chainProver := chainProver
		eg.Go(func() error {
			err := c.collectionLoop(runCtx, chainProver)
			runCtxCancel() // Signal the other chain processors to exit.
			return err
		})
	}

	err := eg.Wait()
	runCtxCancel()
	return err
}

func (c *coordinator) collectionLoop(ctx context.Context, chainProver chainattestor.ChainAttestor) error {
	ticker := time.NewTicker(defaultMinQueryLoopDuration) // TODO: Make this configurable per chain
	defer ticker.Stop()

	for {
		// TODO: Add retry/error handling logic
		c.logger.Info("Collecting claims", zap.String("chain_id", chainProver.ChainID()))
		if err := chainProver.CollectAttestations(ctx); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			ticker.Reset(defaultMinQueryLoopDuration)
		}
	}
}
