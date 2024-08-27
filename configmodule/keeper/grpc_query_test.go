package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"github.com/cosmos/interchain-attestation/configmodule/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	resp, err := suite.queryClient.Params(suite.ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(types.DefaultParams(), resp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryAttestors() {
	attestatorID := suite.registerAttestator()

	resp, err := suite.queryClient.Attestators(suite.ctx, &types.QueryAttestatorsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(resp.Attestators, 1)
	suite.Require().Equal(attestatorID, resp.Attestators[0].AttestatorId)
	var actualPK secp256k1.PubKey
	err = suite.cdc.Unmarshal(resp.Attestators[0].PublicKey.GetValue(), &actualPK)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestGRPCQueryAttestator() {
	attestatorID := suite.registerAttestator()

	resp, err := suite.queryClient.Attestator(suite.ctx, &types.QueryAttestatorRequest{AttestatorId: attestatorID})
	suite.Require().NoError(err)
	suite.Require().Equal(attestatorID, resp.Attestator.AttestatorId)
	var actualPK secp256k1.PubKey
	err = suite.cdc.Unmarshal(resp.Attestator.PublicKey.GetValue(), &actualPK)
	suite.Require().NoError(err)
}
