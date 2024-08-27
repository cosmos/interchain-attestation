package cmd

import (
	"path"

	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/cosmos/interchain-attestation/sidecar/attestators"
	"github.com/cosmos/interchain-attestation/sidecar/server"
)

const (
	flagListenAddr = "listen-addr"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the attestation sidecar",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr, _ := cmd.Flags().GetString(flagListenAddr)

			logger := GetLogger(cmd)

			coordinator, err := setUpCoordinator(cmd, logger)
			if err != nil {
				return err
			}

			s := server.NewServer(logger, coordinator)

			var eg errgroup.Group

			eg.Go(func() error {
				if err := s.Serve(listenAddr); err != nil {
					logger.Error("server.Serve crashed", zap.Error(err))
					return err
				}

				return nil
			})

			eg.Go(func() error {
				if err := coordinator.Run(cmd.Context()); err != nil {
					logger.Error("coordinator.Run crashed", zap.Error(err))
					return err
				}

				return nil
			})

			return eg.Wait()
		},
	}

	cmd.Flags().String(flagListenAddr, "localhost:6969", "Address for grpc server to listen on")

	return cmd
}

func setUpCoordinator(cmd *cobra.Command, logger *zap.Logger) (attestators.Coordinator, error) {
	sidecarConfig := GetConfig(cmd)
	homedir := GetHomedir(cmd)

	dbPath := path.Join(homedir, "db")
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}

	coordinator, err := attestators.NewCoordinator(logger, db, sidecarConfig)
	if err != nil {
		return nil, err
	}

	return coordinator, nil
}
