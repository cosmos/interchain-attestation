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
	GetLatestIBCData() ([]types.IBCData, error)
	GetIBCDataForHeight(chainID string, height uint64) (types.IBCData, error)
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

func (c *coordinator) GetLatestIBCData() ([]types.IBCData, error) {
	var wg sync.WaitGroup
	ibcDataChan := make(chan types.IBCData, len(c.chainAttestators))
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

			var ibcData types.IBCData
			if err := proto.Unmarshal(bz, &ibcData); err != nil {
				errChan <- err
				return
			}

			ibcDataChan <- ibcData
		}(chainAttestator)
	}

	go func() {
		wg.Wait()
		close(ibcDataChan)
		close(errChan)
	}()

	var ibcData []types.IBCData
	for attestation := range ibcDataChan {
		ibcData = append(ibcData, attestation)
	}

	if err := <-errChan; err != nil {
		return nil, err
	}

	return ibcData, nil
}

func (c *coordinator) GetIBCDataForHeight(chainID string, height uint64) (types.IBCData, error) {
	var ibcData types.IBCData
	if err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(heightKey(chainID, height))
		if err != nil {
			return err
		}
		return item.Value(ibcData.Unmarshal)
	}); err != nil {
		return ibcData, err
	}

	return ibcData, nil
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
	c.logger.Info("Collecting ibc data", zap.String("chain_id", chainProver.ChainID()))
	ibcData, err := chainProver.CollectIBCData(ctx)
	if err != nil {
		c.logger.Error("Failed to collect ibc data", zap.String("chain_id", chainProver.ChainID()), zap.Error(err))
		return
	}
	c.logger.Info("Collected ibc data for chain",
		zap.String("chain_id", chainProver.ChainID()),
		zap.String("client_id", ibcData.ClientId),
		zap.String("client_to_update", ibcData.ClientToUpdate),
		zap.String("height", fmt.Sprint(ibcData.Height.RevisionHeight)),
		zap.String("timestamp", ibcData.Timestamp.String()),
		zap.Int("num_packet_commitments", len(ibcData.PacketCommitments)),
	)

	if err := c.db.Update(func(txn *badger.Txn) error {
		ibcDataBz, err := ibcData.Marshal()
		if err != nil {
			return err
		}
		height := ibcData.Height.RevisionHeight
		if err := txn.Set(heightKey(chainProver.ChainID(), height), ibcDataBz); err != nil {
			return err
		}
		if err := txn.Set(latestKey(chainProver.ChainID()), ibcDataBz); err != nil {
			return err
		}

		return nil
	}); err != nil {
		c.logger.Error("Failed to store ibc data", zap.String("chain_id", chainProver.ChainID()), zap.Error(err))
		return
	}
}

func heightKey(chainID string, height uint64) []byte {
	return []byte(fmt.Sprintf("%s/ibcdata/%d", chainID, height))
}

func latestKey(chainID string) []byte {
	return []byte(fmt.Sprintf("%s/ibcdata/latest", chainID))
}
