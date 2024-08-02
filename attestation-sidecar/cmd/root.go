package cmd

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/config"
)

const (
	ContextKeyConfig  = "config"
	ContextKeyHomedir = "homedir"
	ContextKeyLogger  = "logger"

	flagHome = "home"
	flagVerbose = "verbose"

	defaultAppFolderName = ".attestation-sidecar"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attestation-sidecar",
		Short: "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Context() == nil {
				cmd.SetContext(context.Background())
			}

			verbose, err := cmd.Flags().GetBool(flagVerbose)
			if err != nil {
				return err
			}

			logger := CreateLogger(verbose)
			cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyLogger, logger))

			homedir, err := cmd.Flags().GetString(flagHome)
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyHomedir, homedir))

			if cmd.Parent() != nil && cmd.Parent().Name() != configCommandName && cmd.Parent().Name() != keysCommandName {
				sidecarConfig, found, err := config.ReadConfig(homedir)
				if err != nil {
					return err
				}

				if !found {
					if err := config.InitConfig(logger, homedir, false); err != nil {
						return err
					}
					return errors.New("config.toml was not found, example config.toml created")
				}

				if err := sidecarConfig.Validate(); err != nil {
					return err
				}

				cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyConfig, sidecarConfig))
			}

			return nil
		},
	}

	cmd.AddCommand(
		StartCmd(),
		ConfigCmd(),
		KeysCmd(),
	)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultHomedir := filepath.Join(userHomeDir, defaultAppFolderName)

	cmd.PersistentFlags().String(flagHome, defaultHomedir, "home directory for attestation-sidecar config and data")
	cmd.PersistentFlags().Bool(flagVerbose, false, "enable verbose output")

	return cmd
}

func CreateLogger(verbose bool) *zap.Logger {
	logLevel := zapcore.InfoLevel
	if verbose {
		logLevel = zapcore.DebugLevel
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	loggerConfig.Encoding = "console"
	loggerConfig.Level = zap.NewAtomicLevelAt(logLevel)

	// Create the logger from the core
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func GetLogger(cmd *cobra.Command) *zap.Logger {
	return cmd.Context().Value(ContextKeyLogger).(*zap.Logger)
}

func GetHomedir(cmd *cobra.Command) string {
	return cmd.Context().Value(ContextKeyHomedir).(string)
}

func GetConfig(cmd *cobra.Command) config.Config {
	return cmd.Context().Value(ContextKeyConfig).(config.Config)
}