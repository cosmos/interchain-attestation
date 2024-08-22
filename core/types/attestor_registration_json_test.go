package types_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/gjermundgaraba/interchain-attestation/core/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnmarshalAttestationJSON(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	validJSON := []byte(`{
    "attestator-id":"aGVsbG8=",
    "attestation-public-key": {
		"@type": "/cosmos.crypto.secp256k1.PubKey",
        "key":"AkMIdx2z1dKWFXSIIKMa6UEWw0qrtDnYmPp5MMi1PUFQ"
	}
}`)
	attestationRegistration, err := types.UnmarshalAttestationRegistrationJSON(cdc, validJSON)
	require.NoError(t, err)
	require.NotNil(t, attestationRegistration.AttestatorID)
	require.NotNil(t, attestationRegistration.AttestationPublicKey)

	err = attestationRegistration.Validate()
	require.NoError(t, err)

	// TODO: Test some error scenarios
}

// TODO: Test validate attestator separately with both valid and invalid inputs