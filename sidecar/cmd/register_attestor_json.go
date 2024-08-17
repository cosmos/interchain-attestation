package cmd

import (
	"encoding/base64"
	"github.com/cometbft/cometbft/libs/json"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestators"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestators/cosmos"
	"github.com/spf13/cobra"
	"os"
)

const registerAttestatorJSONFileName = "register-attestator.json"

func GenerateRegisterAttestatorJSONCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-register-attestator-json",
		Short: "generate a json file used to register an attestator on-chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := GetConfig(cmd)

			attestatorID := config.AttestatorID
			attestatorIDBase64 := base64.StdEncoding.EncodeToString([]byte(attestatorID))

			cdc := cosmos.NewCodecConfig().Marshaler
			attestatorSigningKey, err := attestators.AttestatorSigningKeyFromConfig(cdc, config)
			if err != nil {
				return err
			}

			pubKey, err := attestatorSigningKey.PubKeyJSON(cdc)
			if err != nil {
				return err
			}

			attestatorRegistrationJSON := types.AttestatorRegistrationJson {
				AttestatorID:           attestatorIDBase64,
				AttestationPublicKey: pubKey,
			}

			jsonBz, err := json.Marshal(attestatorRegistrationJSON)
			if err != nil {
				return err
			}

			if err := os.WriteFile(registerAttestatorJSONFileName, jsonBz, 0644); err != nil {
				return err
			}

			cmd.Printf("Generated %s\n", registerAttestatorJSONFileName)

			return nil
		},
	}

	return cmd
}
