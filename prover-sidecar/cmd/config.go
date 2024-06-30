package cmd

import (
	"github.com/spf13/cobra"
	"proversidecar/config"
)

const (
	configCommandName = "config" // so we can reference it in the root command pre-run hook

	flagForceInitConfig = "force"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   configCommandName,
		Short: "config subcommands",
	}

	cmd.AddCommand(InitConfigCmd())

	return cmd
}

func InitConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool(flagForceInitConfig)

			homedir := GetHomedir(cmd)
			logger := GetLogger(cmd)

			return config.InitConfig(logger, homedir, force)
		},
	}

	cmd.Flags().Bool(flagForceInitConfig, false, "force overwrite of existing config file")

	return cmd
}
