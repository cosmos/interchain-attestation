package cli

import (
	"cosmossdk.io/core/address"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/types"
	coretypes "github.com/gjermundgaraba/pessimistic-validation/core/types"
	"github.com/spf13/cobra"
	"strings"
)

func TxCmd(valAddrCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Attestation config module subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(RegisterAttestatorCmd(valAddrCodec))

	return cmd
}

func RegisterAttestatorCmd(valAddrCodec address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-attestator [path/to/attestator-registration.json]",
		Short: "Register an attestator",
		Long:  "Register an attestator with the module by submitting a JSON file with the attestator details. See example",
		Example: strings.TrimSpace(
			fmt.Sprintf(`
$ %s tx %s register-attestator path/to/attestator-registration.json

Where the JSON file contains the following:

{
	"attestator-id": "base64encodedattestatorid",
	"attestation-public-key": {
		"@type": "/cosmos.crypto.secp256k1.PubKey",
		"key": "base64encodedpubkey"
	}
}

The attestation-public-key should be the same as the one generated using the sidecar cli.
It is the public key that matches with the sidecar private key used to sign attestations.
`, version.AppName, types.ModuleName)),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			valAddr := clientCtx.GetFromAddress()
			valStr, err := valAddrCodec.BytesToString(valAddr)
			if err != nil {
				return err
			}

			attestationRegistration, err := coretypes.ParseAndValidateAttestationRegistrationJSONFromFile(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			msg := &types.MsgRegisterAttestator{
				ValidatorAddress:     valStr,
				AttestatorId:           attestationRegistration.AttestatorID,
				AttestationPublicKey: attestationRegistration.AttestationPublicKey,
			}

			if err := msg.Validate(valAddrCodec); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}