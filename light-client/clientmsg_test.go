package lightclient_test

import "github.com/gjermundgaraba/pessimistic-validation/lightclient"

func (s *PessimisticLightClientTestSuite) TestGetSignableBytes() {


	for i := 0; i < 10; i++ {
		clientMsg := generateClientMsg(s.encCfg.Codec, s.mockAttestators, i)

		expected := lightclient.GetSignableBytes(s.encCfg.Codec, clientMsg.Claims[0].PacketCommitmentsClaim)
		for j := 0; j < 10; j++ {
			for _, claim := range clientMsg.Claims {
				bz := lightclient.GetSignableBytes(s.encCfg.Codec, claim.PacketCommitmentsClaim)
				s.Require().NotNil(bz)

				// verify bytes are the same every time
				s.Require().Equal(expected, bz)

				pubKey, err := s.mockAttestatorsHandler.GetPublicKey(s.ctx, claim.AttestatorId)
				s.Require().NoError(err)
				verified := pubKey.VerifySignature(bz, claim.Signature);
				s.Require().True(verified)
			}
		}

	}
}
