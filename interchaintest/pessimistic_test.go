package pessimisticinterchaintest

import (
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	testifysuite "github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

func TestE2ETestSuite(t *testing.T) {
	testifysuite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) TestTheKitchenSink() {
	s.NotNil(s.ic)

	var userFunds = math.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(s.T(), s.ctx, s.T().Name(), userFunds, s.rolly, s.hub)
	rollyUser, hubUser := users[0], users[1]

	s.NoError(s.r.StartRelayer(s.ctx, s.eRep, s.initialPath))

	s.T().Cleanup(
		func() {
			err := s.r.StopRelayer(s.ctx, s.eRep)
			if err != nil {
				s.T().Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	s.NoError(testutil.WaitForBlocks(s.ctx, 5, s.rolly, s.hub))

	IBCTransferWorksTest(s.T(), s.ctx, s.rolly, s.hub, rollyUser, hubUser, s.r, s.eRep)

	clients, err := s.r.GetClients(s.ctx, s.eRep, "hub")
	s.NoError(err)
	s.Len(clients, 2)
	var tendermintClient string
	for _, client := range clients {
		if client.ClientID == "07-tendermint-0" {
			tendermintClient = client.ClientID
			break
		}
	}

	poolRespBz, _, err := s.hub.GetNode().ExecQuery(s.ctx, "staking", "pool")
	s.NoError(err)
	var poolResponse stakingtypes.QueryPoolResponse
	s.NoError(json.Unmarshal(poolRespBz, &poolResponse))
	requiredSecurity := poolResponse.Pool.BondedTokens.QuoRaw(2).AddRaw(1) //just over 50% of the bonded tokens

	s.NoError(s.SendTx(s.hub.Validators[0], "validator", "pessimist", "create-validation-objective", tendermintClient, requiredSecurity.String()))

	s.NoError(s.SendTx(s.hub.Validators[0], "validator", "pessimist", "sign-up-for-objective", tendermintClient))
	s.NoError(s.SendTx(s.hub.Validators[1], "validator", "pessimist", "sign-up-for-objective", tendermintClient))

	objectiveBz, _, err := s.hub.GetNode().ExecQuery(s.ctx, "pessimist", "validation-objective", tendermintClient)
	s.NoError(err)
	var objectiveResp ValidationObjective
	s.NoError(json.Unmarshal(objectiveBz, &objectiveResp))
	s.False(objectiveResp.ValidationObjective.Activated)
	s.Len(objectiveResp.ValidationObjective.Validators, 2)
	s.Equal(strconv.FormatInt(requiredSecurity.Int64(), 10), objectiveResp.ValidationObjective.RequiredPower)

	s.NoError(s.SendTx(s.hub.Validators[2], "validator", "pessimist", "sign-up-for-objective", tendermintClient))

	objectiveBzAfter, _, err := s.hub.GetNode().ExecQuery(s.ctx, "pessimist", "validation-objective", tendermintClient)
	s.NoError(err)
	var objectiveRespAfter ValidationObjective
	s.NoError(json.Unmarshal(objectiveBzAfter, &objectiveRespAfter))
	s.True(objectiveRespAfter.ValidationObjective.Activated)
	s.Len(objectiveRespAfter.ValidationObjective.Validators, 3)
	s.Equal(strconv.FormatInt(requiredSecurity.Int64(), 10), objectiveRespAfter.ValidationObjective.RequiredPower)
}

func (s *E2ETestSuite) SendTx(node *cosmos.ChainNode, keyName string, command ...string) error {
	command = append(command, "--keyring-dir", "/var/cosmos-chain/hub") // Something off with the hubd binary
	txHash, err := node.ExecTx(s.ctx, keyName, command...)
	if err != nil {
		return err
	}

	txRespBz, _, err := node.ExecQuery(s.ctx, "tx", txHash)
	if err != nil {
		return err
	}
	fmt.Println("tx resp", string(txRespBz))
	var txResp TxResponse
	if err := json.Unmarshal(txRespBz, &txResp); err != nil {
		return err
	}
	if txResp.Code != 0 {
		return fmt.Errorf("tx failed with code %d: %s", txResp.Code, txResp.RawLog)
	}
	return nil
}

func IBCTransferWorksTest(
	t *testing.T,
	ctx context.Context,
	srcChain *cosmos.CosmosChain,
	dstChain *cosmos.CosmosChain,
	srcUser ibc.Wallet,
	dstUser ibc.Wallet,
	r ibc.Relayer,
	eRep *testreporter.RelayerExecReporter) {
	// Wait a few blocks for relayer to start and for user accounts to be created
	err := testutil.WaitForBlocks(ctx, 5, srcChain, dstChain)
	require.NoError(t, err)

	srcUserAddr := srcUser.FormattedAddress()
	dstUserAddr := dstUser.FormattedAddress()

	// Get original account balances
	srcOrigBal, err := srcChain.GetBalance(ctx, srcUserAddr, srcChain.Config().Denom)
	require.NoError(t, err)

	// Compose an IBC transfer and send from srcChain -> dstChain
	var transferAmount = math.NewInt(1_000)
	transfer := ibc.WalletAmount{
		Address: dstUserAddr,
		Denom:   srcChain.Config().Denom,
		Amount:  transferAmount,
	}

	channel, err := ibc.GetTransferChannel(ctx, r, eRep, srcChain.Config().ChainID, dstChain.Config().ChainID)
	require.NoError(t, err)

	srcHeight, err := srcChain.Height(ctx)
	require.NoError(t, err)

	transferTx, err := srcChain.SendIBCTransfer(ctx, channel.ChannelID, srcUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, srcChain, srcHeight, srcHeight+50, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 10, srcChain)
	require.NoError(t, err)

	// Get the IBC denom for srcChain on dstChain
	srcTokenDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, srcChain.Config().Denom)
	srcIBCDenom := transfertypes.ParseDenomTrace(srcTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on srcChain and are in the user acc on dstChain
	srcUpdateBal, err := srcChain.GetBalance(ctx, srcUserAddr, srcChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, srcOrigBal.Sub(transferAmount), srcUpdateBal)

	dstUpdateBal, err := dstChain.GetBalance(ctx, dstUserAddr, srcIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, dstUpdateBal)

	// Compose an IBC transfer and send from dstChain -> srcChain
	transfer = ibc.WalletAmount{
		Address: srcUserAddr,
		Denom:   srcIBCDenom,
		Amount:  transferAmount,
	}

	dstHeight, err := dstChain.Height(ctx)
	require.NoError(t, err)

	transferTx, err = dstChain.SendIBCTransfer(ctx, channel.Counterparty.ChannelID, dstUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, dstChain, dstHeight, dstHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on srcChain and not on dstChain
	srcUpdateBal, err = srcChain.GetBalance(ctx, srcUserAddr, srcChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, srcOrigBal, srcUpdateBal)

	dstUpdateBal, err = dstChain.GetBalance(ctx, dstUserAddr, srcIBCDenom)
	require.NoError(t, err)
	require.Equal(t, int64(0), dstUpdateBal.Int64())
}
