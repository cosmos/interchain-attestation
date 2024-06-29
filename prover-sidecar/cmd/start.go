package cmd

import (
	"github.com/spf13/cobra"
	"proversidecar/server"
)

const (
	defaultPort = 6969

	flagListenAddr = "listen-addr"
)

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the proof sidecar",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr, _ := cmd.Flags().GetString(flagListenAddr)

			s := &server.Server{}
			if err := s.Serve(listenAddr); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagListenAddr, "localhost:6969", "Address for grpc server to listen on")


	return cmd
}
