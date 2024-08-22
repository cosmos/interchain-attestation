package lightclient_test

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
	"github.com/gjermundgaraba/interchain-attestation/core/lightclient"
	"github.com/gjermundgaraba/interchain-attestation/core/types"
	"time"
)

func (s *AttestationLightClientTestSuite) TestLightClientModule_Initialize() {
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
				Timestamp: time.Time{},
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
			}
		})
	}
}

func (s *AttestationLightClientTestSuite) TestLightClientModule_VerifyClientMessage() {
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

func (s *AttestationLightClientTestSuite) TestLightClientModule_UpdateState() {
	clientID := createClientID(0)
	clientStateBz := s.encCfg.Codec.MustMarshal(initialClientState)
	consensusStateBz := s.encCfg.Codec.MustMarshal(initialConsensusState)

	err := s.lightClientModule.Initialize(s.ctx, clientID, clientStateBz, consensusStateBz)
	s.Require().NoError(err)

	expectedHeight := clienttypes.NewHeight(1, defaultHeight.RevisionHeight+1)
	expectedTimestamp := time.Now()

	for i := 0; i < 25; i++ {
		clientMsg := generateClientMsg(s.encCfg.Codec, s.mockAttestators, i, func(attestedData *types.IBCData) {
			attestedData.Height = expectedHeight
			attestedData.Timestamp = expectedTimestamp
		})

		heights := s.lightClientModule.UpdateState(s.ctx, clientID, clientMsg)
		s.Require().Equal([]exported.Height{expectedHeight}, heights)

		s.assertClientState(clientID, expectedHeight, expectedTimestamp)
		s.assertPacketCommitmentStored(clientID, clientMsg)

		expectedHeight = clienttypes.NewHeight(1, clientMsg.Attestations[0].AttestedData.Height.RevisionHeight+1)
		expectedTimestamp = expectedTimestamp.Add(2 * time.Second)
	}

	for i := 25; i != 0; i-- {
		clientMsg := generateClientMsg(s.encCfg.Codec, s.mockAttestators, i, func(attestedData *types.IBCData) {
			attestedData.Height = expectedHeight
			attestedData.Timestamp = expectedTimestamp
		})

		heights := s.lightClientModule.UpdateState(s.ctx, clientID, clientMsg)
		s.Require().Equal([]exported.Height{expectedHeight}, heights)

		s.assertClientState(clientID, expectedHeight, expectedTimestamp)
		s.assertPacketCommitmentStored(clientID, clientMsg)

		expectedHeight = clienttypes.NewHeight(1, clientMsg.Attestations[0].AttestedData.Height.RevisionHeight+1)
		expectedTimestamp = expectedTimestamp.Add(2 * time.Second)
	}
}

func (s *AttestationLightClientTestSuite) TestLightClientModule_VerifyMembership() {
	clientID := createClientID(0)
	clientStateBz := s.encCfg.Codec.MustMarshal(initialClientState)
	consensusStateBz := s.encCfg.Codec.MustMarshal(initialConsensusState)

	err := s.lightClientModule.Initialize(s.ctx, clientID, clientStateBz, consensusStateBz)
	s.Require().NoError(err)

	clientMsg := generateClientMsg(s.encCfg.Codec, s.mockAttestators, 5, func(attestedData *types.IBCData) {
		attestedData.Height = clienttypes.NewHeight(1, defaultHeight.RevisionHeight+1)
	})
	s.lightClientModule.UpdateState(s.ctx, clientID, clientMsg)

	for _, packetCommitment := range clientMsg.Attestations[0].AttestedData.PacketCommitments {
		err = s.lightClientModule.VerifyMembership(s.ctx, clientID, nil, 0, 0, nil, nil, packetCommitment)
		s.Require().NoError(err)
	}

	err = s.lightClientModule.VerifyMembership(s.ctx, clientID, nil, 0, 0, nil, nil, []byte("non-existent-packet-commitment"))
	s.Require().Error(err)

	oldPacketCommitments := clientMsg.Attestations[0].AttestedData.PacketCommitments

	// Update state with no packet commitments
	clientMsg = generateClientMsg(s.encCfg.Codec, s.mockAttestators, 0, func(attestedData *types.IBCData) {
		attestedData.Height = clienttypes.NewHeight(1, clientMsg.Attestations[0].AttestedData.Height.RevisionHeight+1)
	})
	s.lightClientModule.UpdateState(s.ctx, clientID, clientMsg)

	for _, packetCommitment := range oldPacketCommitments {
		err = s.lightClientModule.VerifyMembership(s.ctx, clientID, nil, 0, 0, nil, nil, packetCommitment)
		s.Require().Error(err)
	}
}

func (s *AttestationLightClientTestSuite) assertClientState(clientID string, expectedHeight clienttypes.Height, expectedTimestamp time.Time) {
	clientStore := s.storeProvider.ClientStore(s.ctx, clientID)
	storedClientState := getClientState(clientStore, s.encCfg.Codec)
	s.Require().NotNil(storedClientState)
	s.Require().Equal(expectedHeight, storedClientState.LatestHeight)

	storedConsensusState := getConsensusState(clientStore, s.encCfg.Codec, expectedHeight)
	s.Require().NotNil(storedConsensusState)
	s.Require().Equal(expectedTimestamp.UnixNano(), storedConsensusState.Timestamp.UnixNano())

	// Assert latest height and timestamp at height
	latestHeight := s.lightClientModule.LatestHeight(s.ctx, clientID)
	s.Require().Equal(expectedHeight, latestHeight)

	timestampAtHeight, err := s.lightClientModule.TimestampAtHeight(s.ctx, clientID, expectedHeight)
	s.Require().NoError(err)
	s.Require().Equal(uint64(expectedTimestamp.UnixNano()), timestampAtHeight)
}

func (s *AttestationLightClientTestSuite) assertPacketCommitmentStored(clientID string, clientMsg *lightclient.AttestationClaim) {
	clientStore := s.storeProvider.ClientStore(s.ctx, clientID)
	packetCommitmentStore := prefix.NewStore(clientStore, []byte(lightclient.PacketCommitmentStoreKey))

	// verify packet commitments are stored
	for _, packetCommitment := range clientMsg.Attestations[0].AttestedData.PacketCommitments {
		hasPacketCommitment := packetCommitmentStore.Has(packetCommitment)
		s.Require().True(hasPacketCommitment)
	}

	numberOfPacketsStored := 0
	iterator := packetCommitmentStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		numberOfPacketsStored++
	}
	s.Require().Equal(len(clientMsg.Attestations[0].AttestedData.PacketCommitments), numberOfPacketsStored)
}
