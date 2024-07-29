package lightclient_test

import (
	sdkmath "cosmossdk.io/math"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/gjermundgaraba/pessimistic-validation/lightclient"
	"time"
)

func (s *PessimisticLightClientTestSuite) TestLightClientModule_Initialize() {
	testCases := []struct {
		name           string
		clientState    *lightclient.ClientState
		consensusState *lightclient.ConsensusState
		expError       string
	}{
		{
			"valid client and consensus state",
			initialClientState,
			initialConsensusState,
			"",
		},
		{
			"invalid client state",
			&lightclient.ClientState{
				ChainId:            "testchain-1",
				RequiredTokenPower: sdkmath.NewInt(0),
				FrozenHeight:       clienttypes.Height{},
				LatestHeight:       clienttypes.Height{},
			},
			initialConsensusState,
			"invalid required token power",
		},
		{
			"invalid consensus state",
			initialClientState,
			&lightclient.ConsensusState{
				Timestamp:         time.Time{},
				PacketCommitments: nil,
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
			storedConsensusState := getConsensusState(clientStore, s.encCfg.Codec, defaultHeight)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(storedClientState)
				s.Require().Nil(storedConsensusState)
			} else {
				s.Require().NoError(err)

				// verify state is stored
				s.Require().Equal(tc.clientState, storedClientState)
				s.Require().Equal(tc.consensusState.Timestamp.Unix(), storedConsensusState.Timestamp.Unix())
				s.Require().Equal(tc.consensusState.PacketCommitments, storedConsensusState.PacketCommitments)
			}
		})
	}
}

func (s *PessimisticLightClientTestSuite) TestLightClientModule_VerifyClientMessage() {
	clientID := createClientID(0)
	clientStateBz := s.encCfg.Codec.MustMarshal(initialClientState)
	consensusStateBz := s.encCfg.Codec.MustMarshal(initialConsensusState)

	err := s.lightClientModule.Initialize(s.ctx, clientID, clientStateBz, consensusStateBz)
	s.Require().NoError(err)

	clientMsg := generateClientMsg(s.encCfg.Codec, s.mockAttestators, 5)

	// test happy path
	err = s.lightClientModule.VerifyClientMessage(s.ctx, clientID, clientMsg)
	s.Require().NoError(err)

	// test client not found
	err = s.lightClientModule.VerifyClientMessage(s.ctx, "non-existent-client", clientMsg)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "client not found")
}

func (s *PessimisticLightClientTestSuite) TestLightClientModule_UpdateState() {
	clientID := createClientID(0)
	clientStateBz := s.encCfg.Codec.MustMarshal(initialClientState)
	consensusStateBz := s.encCfg.Codec.MustMarshal(initialConsensusState)

	err := s.lightClientModule.Initialize(s.ctx, clientID, clientStateBz, consensusStateBz)
	s.Require().NoError(err)

	expectedHeight := clienttypes.NewHeight(1, defaultHeight.RevisionHeight+1)
	expectedTimestamp := time.Now()
	clientMsg := generateClientMsg(s.encCfg.Codec, s.mockAttestators, 5, func(claim *lightclient.PacketCommitmentsClaim) {
		claim.Height = expectedHeight
		claim.Timestamp = expectedTimestamp
	})

	heights := s.lightClientModule.UpdateState(s.ctx, clientID, clientMsg)
	s.Require().Equal([]exported.Height{expectedHeight}, heights)

	clientStore := s.storeProvider.ClientStore(s.ctx, clientID)
	storedClientState := getClientState(clientStore, s.encCfg.Codec)
	s.Require().NotNil(storedClientState)
	s.Require().Equal(expectedHeight, storedClientState.LatestHeight)

	storedConsensusState := getConsensusState(clientStore, s.encCfg.Codec, expectedHeight)
	s.Require().NotNil(storedConsensusState)
	s.Require().Equal(expectedTimestamp.Unix(), storedConsensusState.Timestamp.Unix())
	s.Require().Equal(clientMsg.Claims[0].PacketCommitmentsClaim.PacketCommitments, storedConsensusState.PacketCommitments)
}
