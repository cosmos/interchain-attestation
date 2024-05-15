package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

type Config struct {
	SrcChain ChainConfig `yaml:"src_chain"`
	DstChain ChainConfig `yaml:"dst_chain"`

	Mnemonic string `yaml:"mnemonic"`
}

type ChainConfig struct {
	ChainType string `yaml:"chain_type"`
	ChainId   string `yaml:"chain_id"`
	RpcAddr   string `yaml:"rpc_addr"`
	AccountPrefix string `yaml:"account_prefix"`
	GasPrices     string `yaml:"gas_prices"`

	ClientId string `yaml:"client_id"`
	ConnectionId string `yaml:"connection_id"`
	ChannelId string `yaml:"channel_id"`
}

const (
	ChainTypeCosmos = "cosmos"
	ChainTypeRollkit = "rollkit"
)

func ReadConfigFromFile(home string) (Config, error) {
	configPath := path.Join(home, "config.yaml")
	file, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	if config.SrcChain.ChainType == "" {
		return Config{}, fmt.Errorf("src_chain chain_type is required")
	} else if config.SrcChain.ChainType != ChainTypeCosmos && config.SrcChain.ChainType != ChainTypeRollkit {
		return Config{}, fmt.Errorf("src_chain_type must be either %s or %s", ChainTypeCosmos, ChainTypeRollkit)
	}
	if config.SrcChain.ChainId == "" {
		return Config{}, fmt.Errorf("src_chain chain_id is required")
	}
	if config.SrcChain.RpcAddr == "" {
		return Config{}, fmt.Errorf("src_chain rpc_addr is required")
	}
	if config.SrcChain.AccountPrefix == "" {
		return Config{}, fmt.Errorf("src_chain account_prefix is required")
	}
	if config.SrcChain.GasPrices == "" {
		return Config{}, fmt.Errorf("src_chain gas_prices is required")
	}

	if config.DstChain.ChainType == "" {
		return Config{}, fmt.Errorf("dst_chain chain_type is required")
	} else if config.DstChain.ChainType != ChainTypeCosmos && config.DstChain.ChainType != ChainTypeRollkit {
		return Config{}, fmt.Errorf("dst_chain_type must be either %s or %s", ChainTypeCosmos, ChainTypeRollkit)
	}
	if config.DstChain.ChainId == "" {
		return Config{}, fmt.Errorf("dst_chain chain_id is required")
	}
	if config.DstChain.RpcAddr == "" {
		return Config{}, fmt.Errorf("dst_chain rpc_addr is required")
	}
	if config.DstChain.AccountPrefix == "" {
		return Config{}, fmt.Errorf("dst_chain account_prefix is required")
	}
	if config.DstChain.GasPrices == "" {
		return Config{}, fmt.Errorf("dst_chain gas_prices is required")
	}

	if config.Mnemonic == "" {
		return Config{}, fmt.Errorf("mnemonic is required")
	}

	return config, nil
}

func UpdateConfigFile(home string, config Config) error {
	configPath := path.Join(home, "config.yaml")

	// Make sure the full path exists
	if _, err := os.Stat(configPath); err != nil {
		return err
	}

	file, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, file, 0644)
}
