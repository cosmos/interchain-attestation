package keeper_test

import "github.com/gjermundgaraba/pessimistic-validation/configmodule/keeper"

func (suite *KeeperTestSuite) TestGetPublicKey() {
	attestatorsHandler := keeper.NewAttestatorHandler(suite.keeper)

	attestatorID := suite.registerAttestator()

	key, err := attestatorsHandler.GetPublicKey(suite.ctx, attestatorID)
	suite.Require().NoError(err)
	suite.Require().NotNil(key)
}