package lightclient_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/interchain-attestation/core/lightclient"
)

func TestConsensusState_ValidateBasic(t *testing.T) {
	testCases := []struct {
		name           string
		consensusState *lightclient.ConsensusState
		expError       string
	}{
		{
			"valid: initial consensus state",
			initialConsensusState,
			"",
		},
		{
			"invalid: zero timestamp",
			&lightclient.ConsensusState{
				Timestamp: time.Time{},
			},
			"timestamp must be a positive Unix time",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.consensusState.ValidateBasic()
			if tc.expError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
