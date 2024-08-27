package attestationinterchaintest

import (
	"cosmossdk.io/math"
	"fmt"
	"github.com/cosmos/interchain-attestation/core/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	users := interchaintest.GetAndFundTestUsers(s.T(), s.ctx, "ibcuser", math.NewInt(1_000_000), s.rollupsimapp, s.simapp)
	rollupUser, simAppUser := users[0], users[1]
	transferAmount := ibc.WalletAmount{
		Address: simAppUser.FormattedAddress(),
		Denom:   s.rollupsimapp.Config().Denom,
		Amount:  math.NewInt(1_000),
	}
	_, err = s.rollupsimapp.SendIBCTransfer(s.ctx, "channel-0", rollupUser.KeyName(), transferAmount, ibc.TransferOptions{})
	s.Require().NoError(err)

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
}
