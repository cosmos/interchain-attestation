package cmd

import (
	"github.com/spf13/cobra"
	"proversidecar/server"
)

const defaultPort = 6969

func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the proof sidecar",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := &server.Server{}
			if err := s.Serve(defaultPort); err != nil {
				return err
			}

			return nil
		},
	}

	// TODO: Add flag to configure the port

	return cmd
}
