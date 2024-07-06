package keeper_test

import (
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	resp, err := suite.queryClient.Params(suite.ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(types.DefaultParams(), resp.Params)
}
