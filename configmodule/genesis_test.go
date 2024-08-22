package configmodule_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gjermundgaraba/interchain-attestation/configmodule"
	"github.com/gjermundgaraba/interchain-attestation/configmodule/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenesis(t *testing.T) {
	consPubKey := ed25519.GenPrivKey().PubKey()
	require.NotNil(t, consPubKey)
	consPubKeyAny, err := codectypes.NewAnyWithValue(consPubKey)
	require.NoError(t, err)

	pubKey := secp256k1.GenPrivKey().PubKey()
	require.NotNil(t, pubKey)
	pubKeyAny, err := codectypes.NewAnyWithValue(pubKey)

	testCases := []struct {
		name string
		genesis types.GenesisState
		// expectedError bool // TODO: Add later when this can fail
	} {
		{
			name: "default",
			genesis: *types.DefaultGenesisState(),
		},
		{
			name: "custom",
			genesis: types.GenesisState{
				Params: &types.Params{},
				Attestators: []types.Attestator{
					{
						AttestatorId:      []byte("test-attestator-id"),
						PublicKey:       pubKeyAny,
						ConsensusPubkey: consPubKeyAny,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suite := createTestSuite(t)
			app := suite.App
			ctx := app.BaseApp.NewContext(false)

			configmodule.InitGenesis(ctx, suite.Keeper, tc.genesis)

			exportedGenesis := configmodule.ExportGenesis(ctx, suite.Keeper)
			require.Equal(t, tc.genesis.Params, exportedGenesis.Params)
			require.Len(t, tc.genesis.Attestators, len(exportedGenesis.Attestators))
			attestatorMap := make(map[string]types.Attestator)
			for _, attestator := range tc.genesis.Attestators {
				attestatorMap[string(attestator.AttestatorId)] = attestator
			}
			for _, actualAttestator := range exportedGenesis.Attestators {
				foundAttestator, ok := attestatorMap[string(actualAttestator.AttestatorId)]
				require.True(t, ok)

				require.Equal(t, foundAttestator.AttestatorId, actualAttestator.AttestatorId)
				require.Equal(t, foundAttestator.PublicKey.GetValue(), actualAttestator.PublicKey.GetValue())
				require.Equal(t, foundAttestator.ConsensusPubkey.GetValue(), actualAttestator.ConsensusPubkey.GetValue())
			}
		})
	}
}
