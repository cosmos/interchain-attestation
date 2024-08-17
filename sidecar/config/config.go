package config

import (
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"github.com/pelletier/go-toml/v2"
	"gitlab.com/tozd/go/errors"
	"os"
	"path"
)

const (
	configFileName = "config.toml"
)

type Config struct {
	AttestatorID          string              `toml:"attestator_id"`
	SigningPrivateKeyPath string              `toml:"private_key_path"`
	CosmosChains          []CosmosChainConfig `toml:"cosmos_chain"`

	configFilePath string
}

// TODO: Document the config properly in the readme with examples
type CosmosChainConfig struct {
	ChainID string `toml:"chain_id"`
	RPC      string `toml:"rpc"`
	ClientID string `toml:"client_id"`

	// Attestation related stuff
	Attestation         bool   `toml:"attestation"`
	ClientToUpdate      string `toml:"client_to_update"`

	// Relaying and tx related stuff
	AddressPrefix  string  `toml:"address_prefix"`
	KeyringBackend string  `toml:"keyring_backend"`
	KeyName        string  `toml:"key_name"`
	Gas            string  `toml:"gas"`
	GasPrices      string  `toml:"gas_prices"`
	GasAdjustment  float64 `toml:"gas_adjustment"`
}

func (c Config) Validate() error {
	if len(c.CosmosChains) == 0 {
		return errors.New("at least one chain must be defined in the config")
	}

	anyAttestationChains := false
	seenChainIDs := make(map[string]bool)
	seenClientsToUpdate := make(map[string]bool)
	for _, chain := range c.CosmosChains {
		if chain.ChainID == "" {
			return errors.New("chain id cannot be empty")
		}

		if _, ok := seenChainIDs[chain.ChainID]; ok {
			return errors.New("duplicate chain id")
		}
		seenChainIDs[chain.ChainID] = true

		if chain.RPC == "" {
			return errors.New("rpc address cannot be empty")
		}

		if chain.Attestation {
			anyAttestationChains = true

			if chain.ClientID == "" {
				return errors.New("client id cannot be empty when attestation is true")
			}

			if chain.ClientToUpdate == "" {
				return errors.New("client to update cannot be empty when attestation is true")
			}

			if _, ok := seenClientsToUpdate[chain.ClientToUpdate]; ok {
				return errors.New("duplicate client to update")
			}
			seenClientsToUpdate[chain.ClientToUpdate] = true
		}
	}

	if anyAttestationChains {
		if c.AttestatorID == "" {
			return errors.New("attestator id cannot be empty if any chains have attestation true")
		}

		if c.SigningPrivateKeyPath == "" {
			return errors.New("private key path cannot be empty if any chains have attestation true")
		}
	}

	return nil
}

func ReadConfig(homedir string) (Config, bool, error) {
	configFilePath := getConfigFilePath(homedir)

	// Check if config file exists
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return Config{}, false, nil
	}

	f, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, false, err
	}
	var config Config
	if err := toml.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, false, err
	}
	if err := f.Close(); err != nil {
		return Config{}, false, err
	}

	config.configFilePath = configFilePath

	return config, true, nil
}

func InitConfig(homedir string, force bool) (string, error) {
	configFilePath := getConfigFilePath(homedir)

	if !force {
		_, err := os.Stat(configFilePath)
		if !os.IsNotExist(err) {
			return "", errors.Errorf("config file already exists at %s", configFilePath)
		}
	}

	config := Config{
		AttestatorID:          "your-attestator-id",
		SigningPrivateKeyPath: "/path/to/your/private/signing/key.json",
		CosmosChains: []CosmosChainConfig{
			{
				ChainID:        "chain-to-attest-1",
				RPC:            "http://localhost:26657",
				Attestation:    true,
				ClientID:       "example-1-client",
				ClientToUpdate: "client-id-to-update",
				AddressPrefix:  "",
				KeyringBackend: "",
				KeyName:        "",
				Gas:            "",
				GasPrices:      "",
				GasAdjustment:  0,
			},
			{
				ChainID:        "non-attestation-chain-1",
				RPC:            "http://localhost:36657",
				Attestation:    false,
				ClientID:       "",
				ClientToUpdate: "",
				AddressPrefix:  "",
				KeyringBackend: "",
				KeyName:        "",
				Gas:            "",
				GasPrices:      "",
				GasAdjustment:  0,
			},
		},
	}

	config.configFilePath = configFilePath
	if err := config.Save(); err != nil {
		return "", err
	}

	return configFilePath, nil
}

func (c Config) Save() error {
	f, err := os.Create(c.configFilePath)
	if err != nil {
		return err
	}
	if err := toml.NewEncoder(f).Encode(c); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (c Config) GetChain(chainID string) (CosmosChainConfig, bool) {
	var chain CosmosChainConfig
	for _, chain := range c.CosmosChains {
		if chain.ChainID == chainID {
			return chain, true
		}
	}

	return chain, false
}

func getConfigFilePath(homedir string) string {
	return path.Join(homedir, configFileName)
}

func (c CosmosChainConfig) GetClientHeight(height uint64) clienttypes.Height {
	revisionNumber := clienttypes.ParseChainID(c.ChainID)

	return clienttypes.NewHeight(revisionNumber, height)
}