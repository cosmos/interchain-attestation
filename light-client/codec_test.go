package lightclient_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/gogoproto/proto"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/gjermundgaraba/pessimistic-validation/lightclient"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCodec(t *testing.T) {
	testCases := []struct {
		name            string
		toMarshalFrom   proto.Message
		toUnmarshalInto proto.Message
	} {
		{
			"ClientState",
			initialClientState,
			&lightclient.ClientState{},
		},
		{
			"ConsensusState",
			initialConsensusState,
			&lightclient.ConsensusState{},
		},
		{
			"PessimisticClaims",
			&lightclient.PessimisticClaims{
				Claims: []lightclient.PacketCommitmentsClaim{
					{
						ValidatorAddress:  []byte("validator_address"),
						PacketCommitments: [][]byte{},
						Signature:         []byte("signature"),
						Height:            clienttypes.Height{},
					},
				},
			},
			&lightclient.PessimisticClaims{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			encodingCfg := moduletestutil.MakeTestEncodingConfig(lightclient.AppModuleBasic{})

			bz, err := encodingCfg.Codec.Marshal(tc.toMarshalFrom)
			require.NoError(t, err)
			require.NotNil(t, bz)

			err = encodingCfg.Codec.Unmarshal(bz, tc.toUnmarshalInto)
			require.NoError(t, err)

			msg, err := encodingCfg.Codec.InterfaceRegistry().Resolve(sdk.MsgTypeURL(tc.toUnmarshalInto))
			require.NoError(t, err)
			require.NotNil(t, msg)
		})
	}
}
