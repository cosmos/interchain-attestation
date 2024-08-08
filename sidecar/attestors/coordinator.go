package attestors

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dgraph-io/badger/v4"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/chainattestor"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/cosmos"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"time"
)

const (
	defaultMinQueryLoopDuration = 1 * time.Second
)

// TODO: Document
type Coordinator interface {
	Run(ctx context.Context) error
	GetLatestAttestation(chainID string) (types.Attestation, error)
	GetAttestationForHeight(chainID string, height uint64) (types.Attestation, error)
}

type coordinator struct {
	logger *zap.Logger
	db     *badger.DB

	chainProvers map[string]chainattestor.ChainAttestor
	queryLoopDuration time.Duration
}

var _ Coordinator = &coordinator{}

func NewCoordinator(logger *zap.Logger, db *badger.DB, sidecarConfig config.Config) (Coordinator, error) {
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
		db:           db,
		chainProvers: chainProvers,
		queryLoopDuration: defaultMinQueryLoopDuration,
	}, nil
}

func (c *coordinator) GetLatestAttestation(chainID string) (types.Attestation, error) {
	var bz []byte
	if err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(latestKey(chainID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			bz = val
			return nil
		})
	}); err != nil {
		return types.Attestation{}, err
	}

	var attestation types.Attestation
	if err := proto.Unmarshal(bz, &attestation); err != nil {
		return types.Attestation{}, err
	}

	return attestation, nil
}

func (c *coordinator) GetAttestationForHeight(chainID string, height uint64) (types.Attestation, error) {
	var attestation types.Attestation
	if err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(heightKey(chainID, height))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return attestation.Unmarshal(val)
		})
	}); err != nil {
		return attestation, err
	}

	return attestation, nil
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
	ticker := time.NewTicker(c.queryLoopDuration) // TODO: Make this configurable per chain
	defer ticker.Stop()

	// TODO: Refactor all database stuff into a separate file (and probably the coordinator stuff into its own package)
	for {
		c.collectOnce(ctx, chainProver)

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			ticker.Reset(c.queryLoopDuration)
		}
	}
}

func (c *coordinator) collectOnce(ctx context.Context, chainProver chainattestor.ChainAttestor) {
	c.logger.Info("Collecting claims", zap.String("chain_id", chainProver.ChainID()))
	attestation, err := chainProver.CollectAttestation(ctx)
	if err != nil {
		c.logger.Error("Failed to collect claims", zap.String("chain_id", chainProver.ChainID()), zap.Error(err))
		return
	}
	if err := c.db.Update(func(txn *badger.Txn) error {
		aBz, err := attestation.Marshal()
		if err != nil {
			return err
		}
		height := attestation.AttestedData.Height.RevisionHeight
		if err := txn.Set(heightKey(chainProver.ChainID(), height), aBz); err != nil {
			return err
		}
		if err := txn.Set(latestKey(chainProver.ChainID()), aBz); err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.logger.Error("Failed to store attestation", zap.String("chain_id", chainProver.ChainID()), zap.Error(err))
		return
	}
}

func heightKey(chainID string, height uint64) []byte {
	return []byte(fmt.Sprintf("%s/attestations/%d", chainID, height))
}

func latestKey(chainID string) []byte {
	return []byte(fmt.Sprintf("%s/attestations/latest", chainID))
}