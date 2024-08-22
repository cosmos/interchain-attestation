package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/cometbft/cometbft/libs/json"
	"github.com/gjermundgaraba/interchain-attestation/core/types"
	"github.com/gjermundgaraba/interchain-attestation/sidecar/attestators"
	"github.com/gjermundgaraba/interchain-attestation/sidecar/attestators/cosmos"
	"github.com/spf13/cobra"
)

func GenerateRegisterAttestatorJSONCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-register-attestator-json",
		Short: "generate json used to register an attestator on-chain",
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

			fmt.Println(string(jsonBz))

			return nil
		},
	}

	return cmd
}
