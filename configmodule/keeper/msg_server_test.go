package keeper_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gjermundgaraba/interchain-attestation/configmodule/types"
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

func (suite *KeeperTestSuite) TestMsgRegisterAttestator() {
	pubKey := secp256k1.GenPrivKey().PubKey()
	suite.Require().NotNil(pubKey)
	pubKeyAny, err := codectypes.NewAnyWithValue(pubKey)
	suite.Require().NoError(err)

	attestatorID := []byte("attestator id")

	testCases := []struct {
		name   string
		input  *types.MsgRegisterAttestator
		expErr string
	}{
		{
			"valid message",
			&types.MsgRegisterAttestator{
				ValidatorAddress:     testValidatorAddress,
				AttestatorId:           attestatorID,
				AttestationPublicKey: pubKeyAny,
			},
			"",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			resp, err := suite.msgSrvr.RegisterAttestator(suite.ctx, tc.input)
			if tc.expErr != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expErr)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(resp)

				storedAttestator, err := suite.keeper.Attestators.Get(suite.ctx, tc.input.AttestatorId)
				suite.Require().NoError(err)

				var actualPK secp256k1.PubKey
				err = suite.cdc.Unmarshal(storedAttestator.PublicKey.GetValue(), &actualPK)
				suite.Require().NoError(err)
				suite.Require().Equal(pubKey.Bytes(), actualPK.Bytes())

				suite.Require().Equal(tc.input.AttestatorId, storedAttestator.AttestatorId)
				suite.Require().NotNil(storedAttestator.ConsensusPubkey.GetValue())
				suite.Require().Equal(suite.mockValidator.ConsensusPubkey.GetValue(), storedAttestator.ConsensusPubkey.GetValue())
			}
		})
	}
}
