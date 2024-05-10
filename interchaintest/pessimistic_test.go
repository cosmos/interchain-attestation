package pessimisticinterchaintest

import (
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	testifysuite "github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
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

	authorityUser, err := interchaintest.GetAndFundTestUserWithMnemonic(s.ctx, "authority", "copy horror distance stick flock tortoise talk robust grape alter quality call climb dumb arrive leopard digital panel scale decide regret digital humble dust", userFunds, s.hub)
	s.NoError(err)

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

	rollyHostName := s.rolly.GetNode().HostName()
	config := PessimisticValidationConfig{
		ChainsToValidate: map[string]struct {
			RPC string `yaml:"rpc"`
		}{
			"07-tendermint-0": {
				RPC: fmt.Sprintf("http://%s:26657", rollyHostName),
			},
		},
	}
	configBz, err := yaml.Marshal(config)
	s.NoError(err)
	s.NoError(s.hub.Validators[0].WriteFile(s.ctx, configBz, "config/pessimist.yaml"))
	s.NoError(s.hub.Validators[1].WriteFile(s.ctx, configBz, "config/pessimist.yaml"))
	s.NoError(s.hub.Validators[2].WriteFile(s.ctx, configBz, "config/pessimist.yaml"))

	currentHeight, err := s.hub.Height(s.ctx)
	s.NoError(err)
	enableHeight := currentHeight + 5

	// The only reason this is possible is because I've overwritten the consensus store with another authority address and enabled the update-params cli command
	s.NoError(s.SendTx(s.hub.GetNode(), authorityUser.KeyName(),
		"consensus",
		"update-params",
		"--block", "{\"max_bytes\": \"22020096\", \"max_gas\": \"-1\"}",
		"--evidence", "{\"max_age_num_blocks\": \"100000\", \"max_age_duration\": \"172800s\", \"max_bytes\": \"1048576\"}",
		"--validator", "{ \"pub_key_types\": [\"ed25519\"] }",
		"--abci", fmt.Sprintf("{\"vote_extensions_enable_height\": \"%d\"}", enableHeight)))

	s.NoError(testutil.WaitForBlocks(s.ctx, 5, s.hub))

	s.NoError(testutil.WaitForBlocks(s.ctx, 5, s.hub))
}



