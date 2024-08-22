package lightclient_test

import (
	sdkmath "cosmossdk.io/math"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"github.com/gjermundgaraba/interchain-attestation/core/lightclient"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClientState_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		clientState *lightclient.ClientState
		expError string
	} {
		{
			"valid: initial client state",
			initialClientState,
			"",
		},
		{
			"invalid: empty chain id",
			&lightclient.ClientState{
				ChainId:            "",
				RequiredTokenPower: sdkmath.NewInt(100),
				FrozenHeight:       clienttypes.Height{},
				LatestHeight:       clienttypes.NewHeight(1, 42),
			},
			"chain id cannot be empty",
		},
		{
			"invalid: empty chain id with spaces",
			&lightclient.ClientState{
				ChainId:            "  ",
				RequiredTokenPower: sdkmath.NewInt(100),
				FrozenHeight:       clienttypes.Height{},
				LatestHeight:       clienttypes.NewHeight(1, 42),
			},
			"chain id cannot be empty",
		},
		{
			"invalid: required token power is zero",
			&lightclient.ClientState{
				ChainId:            "testchain-1",
				RequiredTokenPower: sdkmath.NewInt(0),
				FrozenHeight:       clienttypes.Height{},
				LatestHeight:       clienttypes.NewHeight(1, 42),
			},
			"required token power must be more than zero",
		},
		{
			"invalid: required token power is negative",
			&lightclient.ClientState{
				ChainId:            "testchain-1",
				RequiredTokenPower: sdkmath.NewInt(-1),
				FrozenHeight:       clienttypes.Height{},
				LatestHeight:       clienttypes.NewHeight(1, 42),
			},
			"required token power must be more than zero",
		},
		{
			"invalid: latest height revision number does not match chain id revision number",
			&lightclient.ClientState{
				ChainId:            "testchain-2",
				RequiredTokenPower: sdkmath.NewInt(100),
				FrozenHeight:       clienttypes.Height{},
				LatestHeight:       clienttypes.NewHeight(1, 42),
			},
			"latest height revision number must match chain id revision number",
		},
		{
			"invalid: latest height revision height is zero",
			&lightclient.ClientState{
				ChainId:            "testchain-1",
				RequiredTokenPower: sdkmath.NewInt(100),
				FrozenHeight:       clienttypes.Height{},
				LatestHeight:       clienttypes.NewHeight(1, 0),
			},
			"client's latest height revision height cannot be zero",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.clientState.Validate()
			if tc.expError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}



