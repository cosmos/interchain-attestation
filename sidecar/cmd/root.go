package cmd

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestators/cosmos"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/config"
	"github.com/spf13/cobra"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"slices"
)

const (
	ContextKeyConfig  = "config"
	ContextKeyHomedir = "homedir"
	ContextKeyLogger  = "logger"

	flagVerbose = "verbose"
	flagAddressPrefix = "address-prefix"

	defaultAppFolderName = ".attestation-sidecar"
)

var parentCommandsToSkipConfigSetup = []string{
	configCommandName,
}

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attestation-sidecar",
		Short: "", // TODO: Write more in both Short and Long
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Context() == nil {
				cmd.SetContext(context.Background())
			}

			verbose, _ := cmd.Flags().GetBool(flagVerbose)

			logger := CreateLogger(verbose)
			cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyLogger, logger))

			homedir, err := cmd.Flags().GetString(flags.FlagHome)
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyHomedir, homedir))

			// Create the home directory if it doesn't exist
			if _, err := os.Stat(homedir); os.IsNotExist(err) {
				if err := os.MkdirAll(homedir, os.ModePerm); err != nil {
					return err
				}
			}

			if cmd.Parent() != nil && !slices.Contains(parentCommandsToSkipConfigSetup, cmd.Parent().Name()) {
				sidecarConfig, found, err := config.ReadConfig(homedir)
				if err != nil {
					return err
				}

				if !found {
					configFilePath, err := config.InitConfig(homedir, false)
					if err != nil {
						return err
					}
					return errors.Errorf("config file was not found, example created at %s", configFilePath)
				}

				if err := sidecarConfig.Validate(); err != nil {
					return err
				}

				cmd.SetContext(context.WithValue(cmd.Context(), ContextKeyConfig, sidecarConfig))
			}

			// To avoid caching of address conversation (because prefixes change based on chain)
			sdk.SetAddrCacheEnabled(false)

			// The rest here is mostly for cosmos sdk stuff (e.g. the keyring commands (keys))
			if cmd.Parent() != nil && cmd.Parent().Name() == "keys" {
				addressPrefix, _ := cmd.Flags().GetString(flagAddressPrefix)
				if addressPrefix == "" {
					return errors.Errorf("address prefix cannot be empty")
				}

				cfg := sdk.GetConfig()
				accountPubKeyPrefix := addressPrefix + "pub"
				validatorAddressPrefix := addressPrefix + "valoper"
				validatorPubKeyPrefix := addressPrefix + "valoperpub"
				consNodeAddressPrefix := addressPrefix + "valcons"
				consNodePubKeyPrefix := addressPrefix + "valconspub"
				cfg.SetBech32PrefixForAccount(addressPrefix, accountPubKeyPrefix)
				cfg.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
				cfg.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
				//cfg.Seal()
				codecConfig := cosmos.NewCodecConfig()
				keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
				kr, err := keyring.New("attestation-sidecar", keyringBackend, homedir, os.Stdin, codecConfig.Marshaler)
				if err != nil {
					return err
				}

				clientCtx := client.Context{}.
					WithCmdContext(cmd.Context()).
					WithCodec(codecConfig.Marshaler).
					WithInput(os.Stdin).
					WithAccountRetriever(authtypes.AccountRetriever{}).
					WithHomeDir(homedir).
					WithKeyring(kr).
					WithViper("")
				cmd.SetContext(context.WithValue(cmd.Context(), client.ClientContextKey, &clientCtx))
			}

			return nil
		},
	}

	keysCmd := keys.Commands()
	keysCmd.PersistentFlags().String(flagAddressPrefix, "", "Bech32 address prefix")

	cmd.AddCommand(
		keysCmd,
		StartCmd(),
		ConfigCmd(),
		SigningKeysCmd(),
		RelayerCmd(),
		GenerateRegisterAttestatorJSONCmd(),
	)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultHomedir := filepath.Join(userHomeDir, defaultAppFolderName)

	cmd.PersistentFlags().String(flags.FlagHome, defaultHomedir, "home directory for attestation-sidecar config and data")
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
