package server

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"path"
)

const configFileName = "config.toml"

type Config struct {
	CometBFTChains []CometBFTChain `toml:"comet_bft_chains"`
}

type CometBFTChain struct {
	ChainID string `toml:"chain_id"`
	RPC    string `toml:"rpc"`
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

func InitConfig(homedir string) error {
	if err := os.MkdirAll(homedir, os.ModePerm); err != nil {
		return err
	}

	configFilePath := getConfigFilePath(homedir)

	config := Config{
		CometBFTChains: []CometBFTChain{
			{
				ChainID: "example-1",
				RPC:    "http://localhost:26657",
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

