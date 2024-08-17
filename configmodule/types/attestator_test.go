package types_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAttestatorValidate(t *testing.T) {
	consPubKey := ed25519.GenPrivKey().PubKey()
	require.NotNil(t, consPubKey)
	consPubKeyAny, err := codectypes.NewAnyWithValue(consPubKey)
	require.NoError(t, err)

	pubKey := secp256k1.GenPrivKey().PubKey()
	require.NotNil(t, pubKey)
	pubKeyAny, err := codectypes.NewAnyWithValue(pubKey)
	require.NoError(t, err)

	notConsKeyAny, err := codectypes.NewAnyWithValue(&types.MsgRegisterAttestator{
		ValidatorAddress: "tt",
		AttestatorId: []byte("tt"),
		AttestationPublicKey: pubKeyAny,
	})

	testCases := []struct {
		name     string
		attestator types.Attestator
		expErr string
	} {
		{
			"valid attestator",
			types.Attestator{
				AttestatorId:      []byte("attestator id"),
				PublicKey:       consPubKeyAny,
				ConsensusPubkey: consPubKeyAny,
			},
			"",
		},
		{
			"attestator id is nil",
			types.Attestator{
				AttestatorId:      nil,
				PublicKey:       consPubKeyAny,
				ConsensusPubkey: consPubKeyAny,
			},
			"attestator id cannot be empty: invalid attestator",
		},
		{
			"public key is nil",
			types.Attestator{
				AttestatorId:      []byte("attestator id"),
				PublicKey:       nil,
				ConsensusPubkey: consPubKeyAny,
			},
			"public key cannot be empty: invalid attestator",
		},
		{
			"invalid pub key",
			types.Attestator{
				AttestatorId:      []byte("attestator id"),
				PublicKey:       notConsKeyAny,
				ConsensusPubkey: consPubKeyAny,
			},
			"expecting attestation pubkey to be cryptotypes.PubKey, got <nil>: invalid type",
		},
		{
			"consensus pubkey is nil",
			types.Attestator{
				AttestatorId: []byte("attestator id"),
				PublicKey: consPubKeyAny,
				ConsensusPubkey: nil,
			},
			"consensus pubkey cannot be empty: invalid attestator",
		},
		{
			"invalid cons pubkey",
			types.Attestator{
				AttestatorId: []byte("attestator id"),
				PublicKey: consPubKeyAny,
				ConsensusPubkey: notConsKeyAny,
			},
			"expecting consensus pubkey to be cryptotypes.PubKey, got <nil>: invalid type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.attestator.Validate()
			if tc.expErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expErr)
			}
		})
	}
}
