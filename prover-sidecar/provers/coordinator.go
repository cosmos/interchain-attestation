package provers

import (
	"context"
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/config"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/provers/chainprover"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/provers/cosmos"
	"time"
)

const (
	defaultMinQueryLoopDuration      = 1 * time.Second
)

type Coordinator interface {
	GetChainProver(chainID string) chainprover.ChainProver
	Run(ctx context.Context) error
}

type coordinator struct {
	logger *zap.Logger

	chainProvers map[string]chainprover.ChainProver
}

func NewCoordinator(logger *zap.Logger, sidecarConfig config.Config) (Coordinator, error) {
	chainProvers := make(map[string]chainprover.ChainProver)
	for _, cosmosConfig := range sidecarConfig.CosmosChains {
		prover, err := cosmos.NewCosmosProver(logger, sidecarConfig.AttestatorID, cosmosConfig.ChainID, cosmosConfig.RPC, cosmosConfig.ClientID, func(msg []byte) ([]byte, error) {
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

func (c *coordinator) GetChainProver(chainID string) chainprover.ChainProver {
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
			err := c.chainProverLoop(runCtx, chainProver)
			runCtxCancel() // Signal the other chain processors to exit.
			return err
		})
	}

	err := eg.Wait()
	runCtxCancel()
	return err
}

func (c *coordinator) chainProverLoop(ctx context.Context, chainProver chainprover.ChainProver) error {
	ticker := time.NewTicker(defaultMinQueryLoopDuration) // TODO: Make this configurable per chain
	defer ticker.Stop()

	for {
		// TODO: Add retry logic
		c.logger.Info("Collecting proofs", zap.String("chain_id", chainProver.ChainID()))
		if err := chainProver.CollectProofs(ctx); err != nil {
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
