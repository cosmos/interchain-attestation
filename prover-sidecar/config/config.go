package config

import (
	"github.com/pelletier/go-toml/v2"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"os"
	"path"
)

const (
	configFileName = "config.toml"
)

type Config struct {
	CosmosChains []CosmosChainConfig `toml:"cosmos_chain"`
}

type CosmosChainConfig struct {
	ChainID string `toml:"chain_id"`
	RPC    string `toml:"rpc"`
	ClientID string `toml:"client_id"`
}

func (c Config) Validate() error {
	if len(c.CosmosChains) == 0 {
		return errors.New("at least one chain must be defined in the config")
	}

	seenChainIDs := make(map[string]bool)
	for _, chain := range c.CosmosChains {
		if chain.ChainID == "" {
			return errors.New("chain id cannot be empty")
		}

		if chain.RPC == "" {
			return errors.New("rpc address cannot be empty")
		}

		if chain.ClientID == "" {
			return errors.New("client id cannot be empty")
		}

		if _, ok := seenChainIDs[chain.ChainID]; ok {
			return errors.New("duplicate chain id")
		}
		seenChainIDs[chain.ChainID] = true
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

	return config, true, nil
}

func InitConfig(logger *zap.Logger, homedir string, force bool) error {
	configFilePath := getConfigFilePath(homedir)

	if !force {
		_, err := os.Stat(configFilePath)
		if !os.IsNotExist(err) {
			return errors.Errorf("config file already exists at %s", configFilePath)
		}
	}

	logger.Debug("InitConfig", zap.String("configFilePath", configFilePath))

	if err := os.MkdirAll(homedir, os.ModePerm); err != nil {
		return err
	}

	config := Config{
		CosmosChains: []CosmosChainConfig{
			{
				ChainID: "example-1",
				RPC:    "http://localhost:26657",
				ClientID: "example-1-client",
			},
		},
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	if err := toml.NewEncoder(f).Encode(config); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func getConfigFilePath(homedir string) string {
	return path.Join(homedir, configFileName)
}

