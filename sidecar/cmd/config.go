package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/interchain-attestation/sidecar/config"
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

			configFilePath, err := config.InitConfig(homedir, force)
			if err != nil {
				return err
			}

			fmt.Printf("Example config file created at %s\n", configFilePath)
			fmt.Println("Please update it with your configuration")

			return nil
		},
	}

	cmd.Flags().Bool(flagForceInitConfig, false, "force overwrite of existing config file")

	return cmd
}
