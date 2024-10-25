package keeper_test

import (
	"github.com/cosmos/interchain-attestation/configmodule/types"
)

func (suite *KeeperTestSuite) TestMsgUpdateParams() {
	authority := suite.keeper.GetAuthority()

	testCases := []struct {
		name   string
		input  *types.MsgUpdateParams
		expErr string
	}{
		{
			"valid: default params",
			&types.MsgUpdateParams{
				Authority: authority,
				Params:    types.DefaultParams(),
			},
			"",
		},
		{
			"valid: custom params",
			&types.MsgUpdateParams{
				Authority: authority,
				Params:    types.Params{},
			},
			"",
		},
		{
			"invalid: invalid authority",
			&types.MsgUpdateParams{
				Authority: "invalid",
				Params:    types.DefaultParams(),
			},
			"invalid authority",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.msgSrvr.UpdateParams(suite.ctx, tc.input)
			if tc.expErr != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErr)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)

				// Check if the params were updated
				updatedParams, err := suite.keeper.Params.Get(suite.ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(tc.input.Params, updatedParams)
			}
		})
	}
}
