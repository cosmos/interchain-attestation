package cmd

import "github.com/spf13/cobra"

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prover-sidecar",
		Short: "",
	}

	cmd.AddCommand(
		StartCmd(),
	)

	return cmd
}
