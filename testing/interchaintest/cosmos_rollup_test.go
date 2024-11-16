package interchaintest

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	transfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v9/testing"

	"github.com/strangelove-ventures/interchaintest/v8"

	"github.com/cosmos/interchain-attestation/core/types"
)

func (s *E2ETestSuite) TestCosmosRollupAttestation() {
	s.Require().NotNil(s.ic)

	// Create light clients
	// TODO: Get the client ids from output and update config files
	stdOut, stdErr, err := s.simapp.Validators[0].Sidecars[0].Exec(s.ctx, []string{
		"attestation-sidecar",
		"relayer",
		"create",
		"clients",
		s.rollupsimapp.Config().ChainID,
		"tendermint", // light client type to be created on rollupsimapp
		s.simapp.Config().ChainID,
		"attestation", // light client type to be created on simapp
		"--home",
		"/home/sidecar",
		"--verbose",
	}, []string{})
	s.Require().NoError(err, string(stdOut), string(stdErr))

	// Create connections
	// TODO: Get connection ids from output (used in the next command)
	stdOut, stdErr, err = s.simapp.Validators[0].Sidecars[0].Exec(s.ctx, []string{
		"attestation-sidecar",
		"relayer",
		"create",
		"connections",
		s.rollupsimapp.Config().ChainID,
		s.simapp.Config().ChainID,
		"--home",
		"/home/sidecar",
		"--verbose",
	}, []string{})
	s.Require().NoError(err, string(stdOut), string(stdErr))

	// Create channels
	stdOut, stdErr, err = s.simapp.Validators[0].Sidecars[0].Exec(s.ctx, []string{
		"attestation-sidecar",
		"relayer",
		"create",
		"channels",
		s.rollupsimapp.Config().ChainID,
		"connection-0",
		"transfer",
		"ics20-1",
		s.simapp.Config().ChainID,
		"connection-0",
		"transfer",
		"--home",
		"/home/sidecar",
		"--verbose",
	}, []string{})
	s.Require().NoError(err, string(stdOut), string(stdErr))

	// Create an IBC transfer
	users := interchaintest.GetAndFundTestUsers(s.T(), s.ctx, "ibcuser", math.NewInt(1_000_000), s.rollupsimapp, s.simapp)
	relayUsers := interchaintest.GetAndFundTestUsers(s.T(), s.ctx, "relayuser", math.NewInt(1_000_000), s.rollupsimapp, s.simapp)
	rollupUser, simappUser := users[0], users[1]
	_, simappRelayUser := relayUsers[0], relayUsers[1]
	transferAmount := math.NewInt(1_000)
	timeoutTimestamp := uint64(time.Now().Add(10 * time.Minute).UnixNano())

	resp, err := s.BroadcastMessages(s.ctx, s.rollupsimapp, rollupUser, 200_000, &transfertypes.MsgTransfer{
		SourcePort:       "transfer",
		SourceChannel:    "channel-0",
		Sender:           rollupUser.FormattedAddress(),
		Receiver:         simappUser.FormattedAddress(),
		TimeoutHeight:    clienttypes.Height{},
		TimeoutTimestamp: timeoutTimestamp,
		Memo:             "",
		Tokens:           sdk.NewCoins(sdk.NewCoin(s.rollupsimapp.Config().Denom, transferAmount)),
		Forwarding:       &transfertypes.Forwarding{},
	})
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(resp.Events)
	s.Require().NoError(err)

	// Wait for validators to catch up
	time.Sleep(5 * time.Second)

	// Check that all sidecars have the attestation
	for i, val := range s.simapp.Validators {
		s.Require().Len(val.Sidecars, 1)
		sidecar := val.Sidecars[0]

		hostPorts, err := sidecar.GetHostPorts(s.ctx, "6969/tcp")
		s.Require().NoError(err)
		s.Require().Len(hostPorts, 1)

		client, err := grpc.NewClient(hostPorts[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
		s.Require().NoError(err)
		defer client.Close()

		sidecarClient := types.NewSidecarClient(client)
		resp, err := sidecarClient.GetAttestations(s.ctx, &types.GetAttestationsRequest{})
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.Require().Len(resp.Attestations, 1)
		s.Require().Equal([]byte(fmt.Sprintf("attestator-%d", i)), resp.Attestations[0].AttestatorId)
		s.Require().Len(resp.Attestations[0].AttestedData.PacketCommitments, 1)
	}

	// Receive packet
	_, err = s.BroadcastMessages(s.ctx, s.simapp, simappRelayUser, 200_000, &channeltypes.MsgRecvPacket{
		Packet:          packet,
		ProofCommitment: []byte("not-used"),
		ProofHeight:     clienttypes.Height{},
		Signer:          simappRelayUser.FormattedAddress(),
	})
	s.Require().NoError(err)

	// Check balance on simapp
	denomOnSimapp := transfertypes.NewDenom(s.rollupsimapp.Config().Denom, transfertypes.NewHop("transfer", "channel-0"))
	balanceResp, err := GRPCQuery[banktypes.QueryBalanceResponse](s.ctx, s.simapp, &banktypes.QueryBalanceRequest{
		Address: simappUser.FormattedAddress(),
		Denom:   denomOnSimapp.IBCDenom(),
	})
	s.Require().NoError(err)
	s.Require().NotNil(balanceResp.Balance)
	s.Require().Equal(transferAmount.Int64(), balanceResp.Balance.Amount.Int64())
}
