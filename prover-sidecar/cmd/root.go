package cmd

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"proversidecar/server"
	"proversidecar/utils"
)

const (
	ContextKeyConfig  = "config"
	ContextKeyHomedir = "homedir"
	ContextKeyLogger  = "logger"

	flagHome = "home"
	flagVerbose = "verbose"

	defaultAppFolderName = ".prover-sidecar"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prover-sidecar",
		Short: "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Context() == nil {
				cmd.SetContext(context.Background())
			}

			verbose, err := cmd.Flags().GetBool(flagVerbose)
			if err != nil {
				return err
			}
			logger := utils.CreateLogger(verbose)
			cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyLogger, logger))

			homedir, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyHomedir, homedir))

			config, found, err := server.ReadConfig(homedir)
			if err != nil {
				return err
			}

			if !found {
				if err := server.InitConfig(homedir); err != nil {
					return err
				}
				return errors.New("config.toml was not found, example config.toml created")
			}

			cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyConfig, config))

			return nil
		},
	}

	cmd.AddCommand(
		StartCmd(),
	)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultHomedir := filepath.Join(userHomeDir, defaultAppFolderName)

	cmd.PersistentFlags().String(flagHome, defaultHomedir, "home directory for prover-sidecar config and data")
	cmd.PersistentFlags().Bool(flagVerbose, false, "enable verbose output")

	return cmd
}
