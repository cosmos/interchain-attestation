package lightclient_test

import (
	sdkmath "cosmossdk.io/math"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/gjermundgaraba/pessimistic-validation/lightclient"
	"time"
)

func (s *PessimisticLightClientTestSuite) TestLightClientModule_Initialize() {
	testCases := []struct {
		name string
		clientState *lightclient.ClientState
		consensusState *lightclient.ConsensusState
		expError string
	} {
		{
			"valid client and consensus state",
			initialClientState,
			initialConsensusState,
			"",
		},
		{
			"invalid client state",
			&lightclient.ClientState{
				ChainId: "testchain-1",
				RequiredTokenPower: sdkmath.NewInt(0),
				FrozenHeight: clienttypes.Height{},
				LatestHeight: clienttypes.Height{},
			},
			initialConsensusState,
			"invalid required token power",
		},
		{
			"invalid consensus state",
			initialClientState,
			&lightclient.ConsensusState{
				Timestamp: time.Time{},
				PacketCommitments: [][]byte{},
			},
			"timestamp must be a positive Unix time",
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			s.SetupTest() // to reset the store and such

			clientID := createClientID(0)
			clientStateBz := s.encCfg.Codec.MustMarshal(tc.clientState)
			consensusStateBz := s.encCfg.Codec.MustMarshal(tc.consensusState)

			err := s.lightClientModule.Initialize(s.ctx, clientID, clientStateBz, consensusStateBz)

			clientStore := s.storeProvider.ClientStore(s.ctx, clientID)
			storedClientState := getClientState(clientStore, s.encCfg.Codec)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(storedClientState)
			} else {
				s.Require().NoError(err)

				// verify client state is stored
				s.Require().Equal(tc.clientState, storedClientState)
			}
		})
	}
}
