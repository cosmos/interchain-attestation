package cmd

import (
	"fmt"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestators"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestators/cosmos"
	"github.com/spf13/cobra"
	"os"
	"path"
)

const (
	signingKeysCommandName = "signing-keys"
)

func SigningKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   signingKeysCommandName,
		Short: fmt.Sprintf("%s subcommands", signingKeysCommandName),
	}

	cmd.AddCommand(CreateSigningKeyCmd())

	return cmd
}

func CreateSigningKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: fmt.Sprintf("create a new signing key and writes it to %s and %s", attestators.DefaultSigningPubKeyFileName, attestators.DefaultSigningPrivKeyFileName),
		RunE: func(cmd *cobra.Command, args []string) error {
			homedir := GetHomedir(cmd)
			config := GetConfig(cmd)

			signingKey, err := attestators.GenerateAttestatorSigningKey()
			if err != nil {
				return err
			}

			encodingConfig := cosmos.NewCodecConfig()
			cdc := encodingConfig.Marshaler

			// json encoding
			pubKeyJSON, err := signingKey.PubKeyJSON(cdc)
			if err != nil {
				return err
			}
			privKeyJSON, err := signingKey.PrivKeyJSON(cdc)
			if err != nil {
				return err
			}

			// write to files
			pubKeyFullPath := path.Join(homedir, attestators.DefaultSigningPubKeyFileName)
			privKeyFullPath := path.Join(homedir, attestators.DefaultSigningPrivKeyFileName)

			if err := os.WriteFile(pubKeyFullPath, pubKeyJSON, 0644); err != nil {
				return err
			}
			if err := os.WriteFile(privKeyFullPath, privKeyJSON, 0644); err != nil {
				return err
			}
			config.SigningPrivateKeyPath = privKeyFullPath
			if err := config.Save(); err != nil {
				return err
			}

			fmt.Printf("Generated %s and %s, as well as updated the config file with the path to the private key\n", pubKeyFullPath, privKeyFullPath)

			return nil
		},
	}

	return cmd
}
