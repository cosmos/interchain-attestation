package cmd

import (
	"github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/attestors"
	"github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

			sidecarConfig := GetConfig(cmd)
			logger := GetLogger(cmd)

			coordinator, err := attestors.NewCoordinator(logger, sidecarConfig)
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
