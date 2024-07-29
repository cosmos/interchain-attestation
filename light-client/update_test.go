package lightclient_test

import (
	"fmt"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tmclienttypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/gjermundgaraba/pessimistic-validation/lightclient"
)

func (s *PessimisticLightClientTestSuite) TestVerifyClientMessage() {
	var attestators []mockAttestator
	var clientMsg exported.ClientMessage
	var attestatorsHandler mockAttestatorsHandler

	tests := []struct {
		name     string
		numberOfAttestator        int
		numberOfPacketCommitments int
		malleate                  func(*lightclient.SignedPacketCommitmentsClaim)
		expError string
	}{
		{
			"valid pessimistic claims",
			10,
			5,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {},
			"",
		},
		{
			"valid pessimistic claims: single attestator",
			1,
			5,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {},
			"",
		},
		{
			"valid pessimistic claims: single packet commitment",
			10,
			1,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {},
			"",
		},
		{
			"valid pessimistic claims: zero commitments",
			10,
			0,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {},
			"",
		},
		{
			"valid pessimistic claims: many attestators",
			100,
			5,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {},
			"",
		},
		{
			"valid pessimistic claims: many packet commitments",
			10,
			100,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {},
			"",
		},
		{
			"valid pessimistic claims: many attestators and packet commitments",
			100,
			100,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {},
			"",
		},
		{
			"invalid client message: type",
			10,
			5,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {
				clientMsg = &tmclienttypes.Header{}
			},
			"invalid client message type",
		},
		{
			"invalid client message: zero claims",
			10,
			5,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {
				clientMsg = &lightclient.PessimisticClaims{}
			},
			"empty claims",
		},
		{
			"invalid client message: different heights",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
				claim.PacketCommitmentsClaim.Height = clienttypes.Height{
					RevisionNumber: 1,
					RevisionHeight: 100000,
				}
				attestatorsHandler.reSignClaim(s.encCfg.Codec, claim)
			},
			"claims must all have the same height",
		},
		{
			"invalid client message: different timestamps",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
				claim.PacketCommitmentsClaim.Timestamp = claim.PacketCommitmentsClaim.Timestamp.Add(10)
				attestatorsHandler.reSignClaim(s.encCfg.Codec, claim)
			},
			"claims must all have the same timestamp",
		},
		{
			"invalid client message: different packet commitments",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
				claim.PacketCommitmentsClaim.PacketCommitments[0] = []byte{0x01}
				attestatorsHandler.reSignClaim(s.encCfg.Codec, claim)
			}                                                                                                                   ,
			"claims must all have the same packet commitments",
		},
		{
			"invalid client message: different amount of packet commitments",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
				claim.PacketCommitmentsClaim.PacketCommitments = append(claim.PacketCommitmentsClaim.PacketCommitments, []byte{0x01})
				attestatorsHandler.reSignClaim(s.encCfg.Codec, claim)
			},
			"claims must all have the same packet commitments",
		},
		{
			"invalid client message: duplicate packet commitment",
			10,
			5,
			func(_ *lightclient.SignedPacketCommitmentsClaim) {
				for _, claim := range clientMsg.(*lightclient.PessimisticClaims).Claims {
					claim.PacketCommitmentsClaim.PacketCommitments[1] = claim.PacketCommitmentsClaim.PacketCommitments[0]
					attestatorsHandler.reSignClaim(s.encCfg.Codec, &claim)
				}
			},
			"duplicate packet commitment",
		},
		{
			"invalid client message: duplicate attestator",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
				clientMsg.(*lightclient.PessimisticClaims).Claims = append(clientMsg.(*lightclient.PessimisticClaims).Claims, *claim)
			},
			"duplicate attestation from",
		},
		{
			"invalid client message: invalid signature over different bytes",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
				var err error
				claim.Signature, err = attestators[0].privateKey.Sign([]byte("different bytes"))
				s.Require().NoError(err)
			},
			"invalid signature from attestator",
		},
		{
			"invalid client message: invalid signature",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
				claim.Signature = []byte{0x01}
			},
			"invalid signature from attestator",
		},
		{
			"not enough attestations",
			10,
			5,
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
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
			func(claim *lightclient.SignedPacketCommitmentsClaim) {
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
				tt.malleate(&clientMsg.(*lightclient.PessimisticClaims).Claims[i])

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
