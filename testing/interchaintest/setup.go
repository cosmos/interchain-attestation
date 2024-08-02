package pessimisticinterchaintest

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/chain/ethereum"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	testifysuite "github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

// Not const because we need to give them as pointers later
var (
	simappVals = 4
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

func (s *E2ETestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain()
	s.ic = ic
	cf := s.getChainFactory()
	chains, err := cf.Chains(s.T().Name())
	s.NoError(err)
	s.simapp, s.rollupsimapp, s.eth = chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain), chains[2].(*ethereum.EthereumChain)


	for _, chain := range chains {
		ic.AddChain(chain)
	}

	client, network := interchaintest.DockerSetup(s.T())

	rf := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(s.T()),
		interchaintestrelayer.CustomDockerImage("ghcr.io/gjermundgaraba/relayer", "pessimistic-rollkit", "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	)
	r := rf.Build(s.T(), client, network)
	s.r = r

	ic.AddRelayer(r, "relayer")
	s.ibcPath = "ibc-path"
	ic.AddLink(interchaintest.InterchainLink{
		Chain1:  s.simapp,
		Chain2:  s.rollupsimapp,
		Relayer: r,
		Path:    s.ibcPath,
	})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(s.T())
	s.eRep = eRep

	err = ic.Build(s.ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         s.T().Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	})
	s.NoError(err)

	s.T().Cleanup(func() {
		_ = ic.Close()
	})
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
				// TODO: Add a callback for each validator to adjust env and maybe other config things?
				Env: []string{
					"PESSIMIST_CONFIG_PATH=/var/cosmos-chain/hub/config/pessimist.yaml",
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
						StartCmd:         []string{"/usr/bin/attestationsidecar", "start", "--home", "/home/sidecar", "--listen-addr", "0.0.0.0:6969"},
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
