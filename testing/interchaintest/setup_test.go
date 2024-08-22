package attestationinterchaintest

import (
	"bytes"
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	"fmt"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/config"
	"github.com/pelletier/go-toml/v2"
	"path"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	attestationve "github.com/gjermundgaraba/pessimistic-validation/core/voteextension"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/chain/ethereum"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	testifysuite "github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

// Not const because we need to give them as pointers later
var (
	simappVals = 1 // TODO: Set to more than one, once we figure out a good coordination/aggregation strategy
	simappFullNodes = 0
	rollupsimappVals = 1
	rollupsimappFullNodes = 0

	votingPeriod     = "15s"
	maxDepositPeriod = "10s"

	genesis = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: votingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: maxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: "stake",
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.amount",
			Value: "1",
		},
	}

	simappChainID = "simapp-1"
	rollupsimappChainID = "rollupsimapp-1"
)

const (
	relayerKeyName = "relayer"
	relayerMnemonic = "worry enable range three surprise skull arctic flame swear crush bunker panel stumble nature strike candy mango junior jealous add sea title unaware alpha"

	nodeSidecarConfigFileName = "sidecar-config.json"
)

type E2ETestSuite struct {
	testifysuite.Suite

	ctx     context.Context
	ic      *interchaintest.Interchain
	network string
	r 	 	ibc.Relayer
	eRep   	*testreporter.RelayerExecReporter

	ibcPath string

	simapp *cosmos.CosmosChain
	rollupsimapp *cosmos.CosmosChain
	eth   *ethereum.EthereumChain
}

func TestE2ETestSuite(t *testing.T) {
	testifysuite.Run(t, new(E2ETestSuite))
}

func (s *E2ETestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain()
	s.ic = ic
	cf := s.getChainFactory()
	chains, err := cf.Chains(s.T().Name())
	s.Require().NoError(err)
	s.simapp, s.rollupsimapp, s.eth = chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain), chains[2].(*ethereum.EthereumChain)

	for _, chain := range chains {
		ic.AddChain(chain)
	}

	client, network := interchaintest.DockerSetup(s.T())

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(s.T())
	s.eRep = eRep


	err = ic.Build(s.ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         s.T().Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	})
	s.Require().NoError(err)

	s.setupSidecars()

	// Create relayer users
	var userFunds = math.NewInt(10_000_000_000)
	_, err = interchaintest.GetAndFundTestUserWithMnemonic(s.ctx, relayerKeyName, relayerMnemonic, userFunds, s.rollupsimapp)
	s.Require().NoError(err)
	_, err = interchaintest.GetAndFundTestUserWithMnemonic(s.ctx, relayerKeyName, relayerMnemonic, userFunds, s.simapp)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		_ = ic.Close()
	})
}

func (s *E2ETestSuite) setupSidecars() {
	chainConfigs := []config.CosmosChainConfig{
		{
			ChainID:        rollupsimappChainID,
			RPC:            s.rollupsimapp.GetRPCAddress(),
			ClientID:       "07-tendermint-0", // TODO: All the client IDs should come from the creation of the light clients
			Attestation:    true,
			ClientToUpdate: "10-attestation-0",
			AddressPrefix:  "rollup",
			KeyringBackend: "test",
			KeyName:        relayerKeyName,
			Gas:            "auto",
			GasPrices:      "0.025stake",
			GasAdjustment:  1.5,
		},
		{
			ChainID:        simappChainID,
			RPC:            s.simapp.GetRPCAddress(),
			ClientID:       "10-attestation-0",
			Attestation:    false,
			ClientToUpdate: "", // Not needed
			AddressPrefix:  "simapp",
			KeyringBackend: "test",
			KeyName:        relayerKeyName,
			Gas:            "auto",
			GasPrices:      "0.025stake",
			GasAdjustment:  1.5,
		},
	}

	for i, val := range s.simapp.Validators {
		s.Require().Len(val.Sidecars, 1)
		sidecar := val.Sidecars[0]

		attestorID := fmt.Sprintf("attestator-%d", i)
		sidecarConfig := config.Config{
			CosmosChains:          chainConfigs,
			AttestatorID:          attestorID,
			SigningPrivateKeyPath: "dummy", // Will be updated when we create the key
		}

		byteWriter := new(bytes.Buffer)
		err := toml.NewEncoder(byteWriter).Encode(sidecarConfig)
		s.Require().NoError(err)
		err = sidecar.WriteFile(s.ctx, byteWriter.Bytes(), "config.toml")
		s.Require().NoError(err)

		stdOut, stdErr, err := sidecar.Exec(s.ctx, []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("echo %s | attestation-sidecar keys add %s --recover --keyring-backend test --home /home/sidecar --address-prefix rollup",
				relayerMnemonic,
				relayerKeyName),
		}, []string{})
		s.Require().NoError(err, string(stdOut), string(stdErr))

		// Create the attestator signing key
		stdOut, stdErr, err = sidecar.Exec(s.ctx, []string{"attestation-sidecar", "signing-keys", "create", "--home", "/home/sidecar"}, []string{})
		s.Require().NoError(err, string(stdOut), string(stdErr))

		// Generate registration file
		stdOut, stdErr, err = sidecar.Exec(s.ctx, []string{
			"attestation-sidecar",
			"generate-register-attestator-json",
			"--home",
			"/home/sidecar",
		}, []string{})
		s.Require().NoError(err, string(stdOut), string(stdErr))
		registrationFileName := fmt.Sprintf("register-attestator-%d.json", i)
		err = val.WriteFile(s.ctx, stdOut, registrationFileName)
		s.Require().NoError(err)

		_, err = val.ExecTx(s.ctx, "validator", "attestationconfig", "register-attestator", path.Join(val.HomeDir(), registrationFileName))
		s.Require().NoError(err)

		err = sidecar.CreateContainer(s.ctx)
		s.Require().NoError(err)

		err = sidecar.StartContainer(s.ctx)
		s.Require().NoError(err)

		nodeSidecarConfig := attestationve.SidecarConfig{
			SidecarAddress: fmt.Sprintf("%s:6969", sidecar.HostName()),
		}
		nodeSidecarConfigBz, err := json.Marshal(nodeSidecarConfig)
		s.Require().NoError(err)
		err = val.WriteFile(s.ctx, nodeSidecarConfigBz, nodeSidecarConfigFileName)
		s.Require().NoError(err)
	}
}

func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e test suite")
	if s.ic != nil {
		_ = s.ic.Close()
	}
}

func (s *E2ETestSuite) getChainFactory() *interchaintest.BuiltinChainFactory {
	return interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(s.T()), []*interchaintest.ChainSpec{
		{
			Name:      "simapp",
			ChainName: "simapp",
			Version:   "local",
			ChainConfig: ibc.ChainConfig{
				Type:    "cosmos",
				Name:    "simapp",
				ChainID: simappChainID,
				Images: []ibc.DockerImage{
					{
						Repository: "simapp",
						Version:    "local",
						UidGid:     "1025:1025",
					},
				},
				Bin:                 "simappd",
				Bech32Prefix:        "simapp",
				Denom:               "stake",
				CoinType:            "118",
				GasPrices:           "0stake",
				GasAdjustment:       2.0,
				TrustingPeriod:      "112h",
				NoHostMount:         false,
				ConfigFileOverrides: nil,
				EncodingConfig:      getEncodingConfig(),
				ModifyGenesis:       cosmos.ModifyGenesis(genesis),
				Env: []string{
					fmt.Sprintf("%s=/var/cosmos-chain/simapp/%s", attestationve.SidecarConfigPathEnv, nodeSidecarConfigFileName),
				},
				SidecarConfigs: []ibc.SidecarConfig{
					{
						ProcessName: "attestationsidecar",
						Image: ibc.DockerImage{
							Repository: "attestationsidecar",
							Version:    "local",
							UidGid:     "1025:1025",
						},
						HomeDir:          "",
						Ports:            []string{"6969/tcp"},
						StartCmd:         []string{"/usr/bin/attestation-sidecar", "--verbose", "start", "--home", "/home/sidecar", "--listen-addr", "0.0.0.0:6969"},
						Env:              nil,
						PreStart:         false,
						ValidatorProcess: true,
					},
				},
			},
			NumValidators: &simappVals,
			NumFullNodes:  &simappFullNodes,
		},
		{
			Name:      "rollupsimapp",
			ChainName: "rollupsimapp",
			Version:   "local",
			ChainConfig: ibc.ChainConfig{
				Type:    "cosmos",
				Name:    "rollupsimapp",
				ChainID: rollupsimappChainID,
				Images: []ibc.DockerImage{
					{
						Repository: "rollupsimapp",
						Version:    "local",
						UidGid:     "1025:1025",
					},
				},
				Bin:                 "rollupsimappd",
				Bech32Prefix:        "rollup",
				Denom:               "stake",
				CoinType:            "118",
				GasPrices:           "0stake",
				GasAdjustment:       2.0,
				TrustingPeriod:      "112h",
				NoHostMount:         false,
				ConfigFileOverrides: nil,
				EncodingConfig:      getEncodingConfig(),
				ModifyGenesisAmounts: func(_ int) (sdk.Coin, sdk.Coin) {
					return sdk.NewInt64Coin("stake", 10_000_000_000_000), sdk.NewInt64Coin("stake", 1_000_000_000)
				},
				ModifyGenesis: func(config ibc.ChainConfig, bytes []byte) ([]byte, error) {
					addressBz, _, err := s.rollupsimapp.Validators[0].Exec(s.ctx, []string{"jq", "-r", ".address", "/var/cosmos-chain/rollupsimapp/config/priv_validator_key.json"}, []string{})
					if err != nil {
						return nil, err
					}
					address := strings.TrimSuffix(string(addressBz), "\n")
					pubKeyBz, _, err := s.rollupsimapp.Validators[0].Exec(s.ctx, []string{"jq", "-r", ".pub_key.value", "/var/cosmos-chain/rollupsimapp/config/priv_validator_key.json"}, []string{})
					if err != nil {
						return nil, err
					}
					pubKey := strings.TrimSuffix(string(pubKeyBz), "\n")
					pubKeyValueBz, _, err := s.rollupsimapp.Validators[0].Exec(s.ctx, []string{"jq", "-r", ".pub_key .value", "/var/cosmos-chain/rollupsimapp/config/priv_validator_key.json"}, []string{})
					if err != nil {
						return nil, err
					}
					pubKeyValue := strings.TrimSuffix(string(pubKeyValueBz), "\n")

					newGenesis := []cosmos.GenesisKV{
						{
							Key: "consensus.validators",
							Value: []map[string]interface{}{
								{
									"address": address,
									"pub_key": map[string]interface{}{
										"type":  "tendermint/PubKeyEd25519",
										"value": pubKey,
									},
									"power": "1",
									"name":  "Rollkit Sequencer",
								},
							},
						},
						{
							Key: "app_state.sequencer.sequencers",
							Value: []map[string]interface{}{
								{
									"name": "test-1",
									"consensus_pubkey": map[string]interface{}{
										"@type": "/cosmos.crypto.ed25519.PubKey",
										"key":   pubKeyValue,
									},
								},
							},
						},
					}

					name := s.rollupsimapp.Sidecars[0].HostName()
					_, _, err = s.rollupsimapp.Validators[0].Exec(s.ctx, []string{"sh", "-c", fmt.Sprintf(`echo "[rollkit]
da_address = \"http://%s:%s\""`+" >> /var/cosmos-chain/rollupsimapp/config/config.toml", name, "7980")}, []string{})
					if err != nil {
						return nil, err
					}

					return cosmos.ModifyGenesis(newGenesis)(config, bytes)
				},
				AdditionalStartArgs: []string{"--rollkit.aggregator", "true", "--api.enable", "--api.enabled-unsafe-cors", "--rpc.laddr", "tcp://0.0.0.0:26657"},
				SidecarConfigs: []ibc.SidecarConfig{
					{
						ProcessName: "mock-da",
						Image: ibc.DockerImage{
							Repository: "mock-da",
							Version:    "local",
							UidGid:     "1025:1025",
						},
						HomeDir:          "",
						Ports:            []string{"7980/tcp"},
						StartCmd:         []string{"/usr/bin/mock-da", "-listen-all"},
						Env:              nil,
						PreStart:         true,
						ValidatorProcess: false,
					},
				},
			},
			NumValidators: &rollupsimappVals,
			NumFullNodes:  &rollupsimappFullNodes,
		},
		// -- ETH --
		{
			ChainConfig: ethereum.DefaultEthereumAnvilChainConfig("ethereum"),
		},
	})
}

func getEncodingConfig() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	// whatever.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}
