package pessimisticinterchaintest

import (
	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

func (s *E2ETestSuite) TestIBC() {
	s.NotNil(s.ic)

	var userFunds = math.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(s.T(), s.ctx, s.T().Name(), userFunds, s.rollupsimapp, s.simapp)
	rollyUser, hubUser := users[0], users[1]

	// This works because we assume only 1 client, 1 connection, and 1 channel
	initialChannel, err := ibc.GetTransferChannel(s.ctx, s.r, s.eRep, s.rollupsimapp.Config().ChainID, s.simapp.Config().ChainID)
	s.NoError(err)

	IBCTransferWorksTest(s.T(), s.ctx, s.r, s.eRep, s.ibcPath, s.rollupsimapp, s.simapp, rollyUser, hubUser, initialChannel.ChannelID, initialChannel.Counterparty.ChannelID)
}