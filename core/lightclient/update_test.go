package lightclient_test

import (
	"fmt"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
	tmclienttypes "github.com/cosmos/ibc-go/v9/modules/light-clients/07-tendermint"
	"github.com/gjermundgaraba/pessimistic-validation/core/lightclient"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
)

func (s *PessimisticLightClientTestSuite) TestVerifyClientMessage() {
	var attestators []mockAttestator
	var clientMsg exported.ClientMessage
	var attestatorsHandler mockAttestatorsHandler

	tests := []struct {
		name     string
		numberOfAttestator        int
		numberOfPacketCommitments int
		malleate                  func(*types.Attestation)
		expError string
	}{
		{
			"valid pessimistic attestations",
			10,
			5,
			func(_ *types.Attestation) {},
			"",
		},
		{
			"valid pessimistic attestations: single attestator",
			1,
			5,
			func(_ *types.Attestation) {},
			"",
		},
		{
			"valid pessimistic attestations: single packet commitment",
			10,
			1,
			func(_ *types.Attestation) {},
			"",
		},
		{
			"valid pessimistic attestations: zero commitments",
			10,
			0,
			func(_ *types.Attestation) {},
			"",
		},
		{
			"valid pessimistic attestations: many attestators",
			100,
			5,
			func(_ *types.Attestation) {},
			"",
		},
		{
			"valid pessimistic attestations: many packet commitments",
			10,
			100,
			func(_ *types.Attestation) {},
			"",
		},
		{
			"valid pessimistic attestations: many attestators and packet commitments",
			100,
			100,
			func(_ *types.Attestation) {},
			"",
		},
		{
			"invalid client message: type",
			10,
			5,
			func(_ *types.Attestation) {
				clientMsg = &tmclienttypes.Header{}
			},
			"invalid client message type",
		},
		{
			"invalid client message: zero attestations",
			10,
			5,
			func(_ *types.Attestation) {
				clientMsg = &lightclient.AttestationClaim{}
			},
			"empty attestations",
		},
		{
			"invalid client message: different heights",
			10,
			5,
			func(attestation *types.Attestation) {
				attestation.AttestedData.Height = clienttypes.Height{
					RevisionNumber: 1,
					RevisionHeight: 100000,
				}
				attestatorsHandler.reSignAttestation(s.encCfg.Codec, attestation)
			},
			"attestations must all be the same",
		},
		{
			"invalid client message: different timestamps",
			10,
			5,
			func(attestation *types.Attestation) {
				attestation.AttestedData.Timestamp = attestation.AttestedData.Timestamp.Add(10)
				attestatorsHandler.reSignAttestation(s.encCfg.Codec, attestation)
			},
			"attestations must all be the same",
		},
		{
			"invalid client message: different packet commitments",
			10,
			5,
			func(attestation *types.Attestation) {
				attestation.AttestedData.PacketCommitments[0] = []byte{0x01}
				attestatorsHandler.reSignAttestation(s.encCfg.Codec, attestation)
			}                                                                                                                   ,
			"attestations must all be the same",
		},
		{
			"invalid client message: different amount of packet commitments",
			10,
			5,
			func(attestation *types.Attestation) {
				attestation.AttestedData.PacketCommitments = append(attestation.AttestedData.PacketCommitments, []byte{0x01})
				attestatorsHandler.reSignAttestation(s.encCfg.Codec, attestation)
			},
			"attestations must all be the same",
		},
		{
			"invalid client message: different chain id",
			10,
			5,
			func(attestation *types.Attestation) {
				attestation.AttestedData.ChainId = "different chain id"
				attestatorsHandler.reSignAttestation(s.encCfg.Codec, attestation)
			},
			"attestations must all be the same",
		},
		{
			"invalid client message: different client id",
			10,
			5,
			func(attestation *types.Attestation) {
				attestation.AttestedData.ClientId = "different client id"
				attestatorsHandler.reSignAttestation(s.encCfg.Codec, attestation)
			},
			"attestations must all be the same",
		},
		{
			"invalid client message: duplicate packet commitment",
			10,
			5,
			func(_ *types.Attestation) {
				for _, attestation := range clientMsg.(*lightclient.AttestationClaim).Attestations {
					attestation.AttestedData.PacketCommitments[1] = attestation.AttestedData.PacketCommitments[0]
					attestatorsHandler.reSignAttestation(s.encCfg.Codec, &attestation)
				}
			},
			"duplicate packet commitment",
		},
		{
			"invalid client message: duplicate attestator",
			10,
			5,
			func(attestation *types.Attestation) {
				clientMsg.(*lightclient.AttestationClaim).Attestations = append(clientMsg.(*lightclient.AttestationClaim).Attestations, *attestation)
			},
			"duplicate attestation from",
		},
		{
			"invalid client message: invalid signature over different bytes",
			10,
			5,
			func(attestation *types.Attestation) {
				var err error
				attestation.Signature, err = attestators[0].privateKey.Sign([]byte("different bytes"))
				s.Require().NoError(err)
			},
			"invalid signature from attestator",
		},
		{
			"invalid client message: invalid signature",
			10,
			5,
			func(attestation *types.Attestation) {
				attestation.Signature = []byte{0x01}
			},
			"invalid signature from attestator",
		},
		{
			"insufficient number of attestators in claim",
			10,
			5,
			func(attestation *types.Attestation) {
				attestatorsHandler.sufficientAttestations = func() (bool, error) {
					return false, nil
				}
			},
			"not enough attestations",
		},
		{
			"sufficient attestators handler error",
			10,
			5,
			func(attestation *types.Attestation) {
				attestatorsHandler.sufficientAttestations = func() (bool, error) {
					return false, fmt.Errorf("handler error")
				}
			},
			"failed to check sufficient attestations",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			attestators = generateAttestators(tt.numberOfAttestator)
			attestatorsHandler = NewMockAttestatorsHandler(attestators)

			for i := 0; i < tt.numberOfAttestator; i++ {
				clientMsg = generateClientMsg(s.encCfg.Codec, attestators, tt.numberOfPacketCommitments)
				tt.malleate(&clientMsg.(*lightclient.AttestationClaim).Attestations[i])

				err := initialClientState.VerifyClientMessage(s.ctx, s.encCfg.Codec, attestatorsHandler, clientMsg)
				if tt.expError != "" {
					s.Require().Error(err)
					s.Require().Contains(err.Error(), tt.expError)
				} else {
					s.Require().NoError(err)
				}
			}
		})
	}
}
