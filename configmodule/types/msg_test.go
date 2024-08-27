package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/interchain-attestation/configmodule/types"
)

func TestMsgRegisterAttestatorValidate(t *testing.T) {
	pubKey := secp256k1.GenPrivKey().PubKey()
	require.NotNil(t, pubKey)
	pubKeyAny, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)

	notAPubKey, err := codectypes.NewAnyWithValue(&types.MsgRegisterAttestator{
		ValidatorAddress:     "tt",
		AttestatorId:         []byte("tt"),
		AttestationPublicKey: pubKeyAny,
	})
	require.NoError(t, err)

	testCases := []struct {
		name     string
		msg      types.MsgRegisterAttestator
		expError string
	}{
		{
			"valid message",
			types.MsgRegisterAttestator{
				ValidatorAddress:     "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8",
				AttestatorId:         []byte("attestator id"),
				AttestationPublicKey: pubKeyAny,
			},
			"",
		},
		{
			"validator is empty",
			types.MsgRegisterAttestator{
				ValidatorAddress:     "",
				AttestatorId:         []byte("attestator id"),
				AttestationPublicKey: pubKeyAny,
			},
			"invalid validator address: empty address string is not allowed: invalid address",
		},
		{
			"invalid validator address",
			types.MsgRegisterAttestator{
				ValidatorAddress:     "invalid",
				AttestatorId:         []byte("attestator id"),
				AttestationPublicKey: pubKeyAny,
			},
			"invalid validator address: decoding bech32 failed: invalid bech32 string length 7: invalid address",
		},
		{
			"nil attestator id",
			types.MsgRegisterAttestator{
				ValidatorAddress:     "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8",
				AttestatorId:         nil,
				AttestationPublicKey: pubKeyAny,
			},
			"attestator id cannot be empty: invalid request",
		},
		{
			"nil attestation public key",
			types.MsgRegisterAttestator{
				ValidatorAddress:     "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8",
				AttestatorId:         []byte("attestator id"),
				AttestationPublicKey: nil,
			},
			"public key cannot be empty: invalid request",
		},
		{
			"invalid attestation public key",
			types.MsgRegisterAttestator{
				ValidatorAddress:     "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8",
				AttestatorId:         []byte("attestator id"),
				AttestationPublicKey: notAPubKey,
			},
			"expecting attestation public key to be cryptotypes.PubKey, got <nil>: invalid type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.Validate(addresscodec.NewBech32Codec(sdk.Bech32PrefixValAddr))
			if tc.expError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expError)
			}
		})
	}
}
