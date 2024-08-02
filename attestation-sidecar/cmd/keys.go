package cmd

import (
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path"
)

const keysCommandName = "keys"

func KeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   keysCommandName,
		Short: "keys subcommands",
	}

	cmd.AddCommand(CreateCmd())

	return cmd
}

func CreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a new signing key and writes it to a priv.key file",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Use a similar system to the one in cosmos where the key is in a json format with a base64 encoded key

			privKey := secp256k1.GenPrivKey()
			bz := privKey.Bytes()
			base64EncodedKey := make([]byte, base64.StdEncoding.EncodedLen(len(bz)))
			base64.StdEncoding.Encode(base64EncodedKey, bz)

			// write to file
			filename := "priv.key"
			homedir := GetHomedir(cmd)
			fullPath := path.Join(homedir, filename)
			if err := os.WriteFile(fullPath, base64EncodedKey, 0600); err != nil {
				return err
			}

			logger := GetLogger(cmd)
			logger.Info("private key written to file", zap.String("path", fullPath))

			return nil
		},
	}

	return cmd
}