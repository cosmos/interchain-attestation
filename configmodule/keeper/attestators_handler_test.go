package keeper_test

import "github.com/cosmos/interchain-attestation/configmodule/keeper"

func (s *KeeperTestSuite) TestSufficientAttestations() {
	attestatorsHandler := keeper.NewAttestatorHandler(s.keeper)

	// TODO: Test when implemented properly, right now it just always returns true
	s.Require().True(attestatorsHandler.SufficientAttestations(s.ctx, nil))
}
