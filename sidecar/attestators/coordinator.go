package attestators

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cosmos/gogoproto/proto"
	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/cosmos/interchain-attestation/core/types"
	"github.com/cosmos/interchain-attestation/sidecar/attestators/attestator"
	"github.com/cosmos/interchain-attestation/sidecar/attestators/cosmos"
	"github.com/cosmos/interchain-attestation/sidecar/config"
)

const (
	defaultMinQueryLoopDuration = 1 * time.Second
)

// TODO: Document
type Coordinator interface {
	Run(ctx context.Context) error
	GetLatestAttestations() ([]types.Attestation, error)
	GetAttestationForHeight(chainID string, height uint64) (types.Attestation, error)
}

type coordinator struct {
	logger *zap.Logger
	db     *badger.DB

	chainAttestators  map[string]attestator.Attestator
	queryLoopDuration time.Duration
}

var _ Coordinator = &coordinator{}

func NewCoordinator(logger *zap.Logger, db *badger.DB, sidecarConfig config.Config) (Coordinator, error) {
	chainProvers := make(map[string]attestator.Attestator)
	for _, cosmosConfig := range sidecarConfig.CosmosChains {
		if !cosmosConfig.Attestation {
			logger.Debug("Skipping chain", zap.String("chain_id", cosmosConfig.ChainID), zap.String("reason", "attestation disabled"))
			continue
		}

		att, err := cosmos.NewCosmosAttestator(
			logger,
			sidecarConfig.AttestatorID,
			cosmosConfig,
		)
		if err != nil {
			return nil, err
		}
		chainProvers[cosmosConfig.ChainID] = att
	}

	return &coordinator{
		logger:            logger,
		db:                db,
		chainAttestators:  chainProvers,
		queryLoopDuration: defaultMinQueryLoopDuration,
	}, nil
}

func (c *coordinator) GetLatestAttestations() ([]types.Attestation, error) {
	var wg sync.WaitGroup
	attestationChan := make(chan types.Attestation, len(c.chainAttestators))
	errChan := make(chan error, len(c.chainAttestators))

	for _, chainAttestator := range c.chainAttestators {
		wg.Add(1)
		go func(chainAttestator attestator.Attestator) {
			defer wg.Done()

			var bz []byte
			if err := c.db.View(func(txn *badger.Txn) error {
				item, err := txn.Get(latestKey(chainAttestator.ChainID()))
				if err != nil {
					return err
				}
				return item.Value(func(val []byte) error {
					bz = val
					return nil
				})
			}); err != nil {
				errChan <- err
				return
			}

			var attestation types.Attestation
			if err := proto.Unmarshal(bz, &attestation); err != nil {
				errChan <- err
				return
			}

			attestationChan <- attestation
		}(chainAttestator)
	}

	go func() {
		wg.Wait()
		close(attestationChan)
		close(errChan)
	}()

	var attestations []types.Attestation
	for attestation := range attestationChan {
		attestations = append(attestations, attestation)
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

	return attestations, nil
}

func (c *coordinator) GetAttestationForHeight(chainID string, height uint64) (types.Attestation, error) {
	var attestation types.Attestation
	if err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(heightKey(chainID, height))
		if err != nil {
			return err
		}
		return item.Value(attestation.Unmarshal)
	}); err != nil {
		return attestation, err
	}

	return attestation, nil
}

func (c *coordinator) Run(ctx context.Context) error {
	c.logger.Debug("Coordinator.Run")

	var eg errgroup.Group
	runCtx, runCtxCancel := context.WithCancel(ctx)
	for _, chainProver := range c.chainAttestators {
		c.logger.Info("Starting chain prover loop", zap.String("chain_id", chainProver.ChainID()))

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

func (c *coordinator) collectionLoop(ctx context.Context, chainProver attestator.Attestator) error {
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

func (c *coordinator) collectOnce(ctx context.Context, chainProver attestator.Attestator) {
	c.logger.Info("Collecting claims", zap.String("chain_id", chainProver.ChainID()))
	attestation, err := chainProver.CollectAttestation(ctx)
	if err != nil {
		c.logger.Error("Failed to collect claims", zap.String("chain_id", chainProver.ChainID()), zap.Error(err))
		return
	}
	c.logger.Info("Collected attestation for chain",
		zap.String("chain_id", chainProver.ChainID()),
		zap.String("client_id", attestation.AttestedData.ClientId),
		zap.String("client_to_update", attestation.AttestedData.ClientToUpdate),
		zap.String("height", fmt.Sprint(attestation.AttestedData.Height.RevisionHeight)),
		zap.String("timestamp", attestation.AttestedData.Timestamp.String()),
		zap.Int("num_packet_commitments", len(attestation.AttestedData.PacketCommitments)),
	)

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
