package attestators_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/interchain-attestation/sidecar/attestators"
	"github.com/cosmos/interchain-attestation/sidecar/attestators/cosmos"
)

func TestAttestatorSigningKeyJSON(t *testing.T) {
	encodingConfig := cosmos.NewCodecConfig()
	cdc := encodingConfig.Marshaler

	attestatorSigningKey, err := attestators.GenerateAttestatorSigningKey()
	require.NoError(t, err)

	pubKeyJSON, err := attestatorSigningKey.PubKeyJSON(cdc)
	require.NoError(t, err)
	require.Equal(t,
		fmt.Sprintf(
			`{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"%s"}`,
			base64.StdEncoding.EncodeToString(attestatorSigningKey.PubKey.Bytes()),
		),
		string(pubKeyJSON),
	)

	privKeyJSON, err := attestatorSigningKey.PrivKeyJSON(cdc)
	require.NoError(t, err)
	require.Equal(t,
		fmt.Sprintf(
			`{"@type":"/cosmos.crypto.secp256k1.PrivKey","key":"%s"}`,
			base64.StdEncoding.EncodeToString(attestatorSigningKey.PrivKey.Bytes()),
		),
		string(privKeyJSON),
	)

	attestatorSigningKey2, err := attestators.AttestatorSigningKeyFromJSON(cdc, privKeyJSON)
	require.NoError(t, err)
	require.Equal(t, attestatorSigningKey, attestatorSigningKey2)
}
