package keeper_test

import (
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
)

func (suite *KeeperTestSuite) TestMsgUpdateParams() {
	authority := suite.keeper.GetAuthority()
	params := types.DefaultParams()

	testCases := []struct {
		name      string
		input     *types.MsgUpdateParams
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "valid: default params",
			input:     &types.MsgUpdateParams{
				Authority: authority,
				Params: params,
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name: "valid: custom params",
			input: &types.MsgUpdateParams{
				Authority: authority,
				Params: types.Params{
					MinimumPower: 30_000_000,
				},
			},
			expErr:    false,
			expErrMsg: "",
		},
		{
			name:      "invalid: invalid authority",
			input:     &types.MsgUpdateParams{
				Authority: "invalid",
				Params: params,
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.msgSrvr.UpdateParams(suite.ctx, tc.input)
			if tc.expErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErrMsg)
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
