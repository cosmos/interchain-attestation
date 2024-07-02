package pessimisticinterchaintest

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v8/relayer"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	testifysuite "github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
	"strings"
)

// Not const because we need to give them as pointers later
var (
	hubVals   = 4
	hubFull   = 0
	rollyVals = 1
	rollyFull = 0

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

	rollyChainID = "rolly"
	hubChainID = "hub"
)

type E2ETestSuite struct {
	testifysuite.Suite

	ctx     context.Context
	ic      *interchaintest.Interchain
	network string
	r 	 	ibc.Relayer
	eRep   	*testreporter.RelayerExecReporter
	initialPath string

	hub     *cosmos.CosmosChain
	rolly   *cosmos.CosmosChain
}

func (s *E2ETestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain()
	s.ic = ic
	cf := s.getChainFactory()
	chains, err := cf.Chains(s.T().Name())
	hub, rolly := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)
	s.NoError(err)
	s.hub = hub
	s.rolly = rolly

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
	s.initialPath = "ibc-path"
	ic.AddLink(interchaintest.InterchainLink{
		Chain1:  hub,
		Chain2:  rolly,
		Relayer: r,
		Path:    s.initialPath,
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
			Name:      "hub",
			ChainName: "hub",
			Version:   "local",
			ChainConfig: ibc.ChainConfig{
				Type:    "cosmos",
				Name:    "hub",
				ChainID: hubChainID,
				Images: []ibc.DockerImage{
					{
						Repository: "hub",
						Version:    "local",
						UidGid:     "1025:1025",
					},
				},
				Bin:                 "hubd",
				Bech32Prefix:        "hub",
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
						ProcessName:      "proversidecar",
						Image:            ibc.DockerImage{
							Repository: "proversidecar",
							Version:    "local",
							UidGid:     "1025:1025",
						},
						HomeDir:          "",
						Ports:            []string{"6969/tcp"},
						StartCmd:         []string{"/usr/bin/proversidecar", "start", "--home", "/home/sidecar", "--listen-addr", "0.0.0.0:6969"},
						Env:              nil,
						PreStart:         false,
						ValidatorProcess: true,
					},
				},
			},
			NumValidators: &hubVals,
			NumFullNodes:  &hubFull,
		},
		{
			Name:      "rolly",
			ChainName: "rolly",
			Version:   "local",
			ChainConfig: ibc.ChainConfig{
				Type:    "cosmos",
				Name:    "rolly",
				ChainID: rollyChainID,
				Images: []ibc.DockerImage{
					{
						Repository: "rolly",
						Version:    "local",
						UidGid:     "1025:1025",
					},
				},
				Bin:                 "rollyd",
				Bech32Prefix:        "rolly",
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
					addressBz, _, err := s.rolly.Validators[0].Exec(s.ctx, []string{"jq", "-r", ".address", "/var/cosmos-chain/rolly/config/priv_validator_key.json"}, []string{})
					if err != nil {
						return nil, err
					}
					address := strings.TrimSuffix(string(addressBz), "\n")
					pubKeyBz, _, err := s.rolly.Validators[0].Exec(s.ctx, []string{"jq", "-r", ".pub_key.value", "/var/cosmos-chain/rolly/config/priv_validator_key.json"}, []string{})
					if err != nil {
						return nil, err
					}
					pubKey := strings.TrimSuffix(string(pubKeyBz), "\n")

					newGenesis := append(genesis, cosmos.GenesisKV{
						Key: "consensus.validators",
						Value: []map[string]interface{}{
							{
								"address": address,
								"pub_key": map[string]interface{}{
									"type":  "tendermint/PubKeyEd25519",
									"value": pubKey,
								},
								"power": "1000",
								"name":  "Rollkit Sequencer",
							},
						},
					})

					name := s.rolly.Sidecars[0].HostName()
					_, _, err = s.rolly.Validators[0].Exec(s.ctx, []string{"bash", "-c", fmt.Sprintf(`echo "[rollkit]
da_address = \"http://%s:%s\"" >> /var/cosmos-chain/rolly/config/config.toml`, name, "7980")}, []string{})
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
			NumValidators: &rollyVals,
			NumFullNodes:  &rollyFull,
		},
	})
}

func getEncodingConfig() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	// whatever.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}
