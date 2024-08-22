package types_test

import (
	"github.com/gjermundgaraba/interchain-attestation/configmodule/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenesisValidate(t *testing.T) {
	tests := []struct {
		name string
		genesis *types.GenesisState
		expError string
	} {
		{
			"valid: default genesis state",
			types.DefaultGenesisState(),
			"",
		},
		{
			"invalid: empty params",
			&types.GenesisState{
				Params: nil,
			},
			"params cannot be nil",
		},
		// TODO: Invalid params when that exists...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.genesis.Validate()

			if tt.expError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.expError)
			}
		})
	}
}
