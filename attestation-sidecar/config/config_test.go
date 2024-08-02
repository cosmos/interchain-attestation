package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Validate(t *testing.T) {

	tests := []struct {
		name    string
		config Config
		expErr string
	}{
		{
			name: "valid config",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:  "chain1",
						RPC:      "http://localhost:26657",
						ClientID: "client1",
					},
					{
						ChainID:  "chain2",
						RPC:      "http://localhost:26658",
						ClientID: "client2",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "",
		},
		{
			name: "empty chains",
			config: Config{
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "at least one chain must be defined in the config",
		},
		{
			name: "empty chain id",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:  "",
						RPC:      "http://localhost:26657",
						ClientID: "client1",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "chain id cannot be empty",
		},
		{
			name: "empty rpc",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:  "chain1",
						RPC:      "",
						ClientID: "client1",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "rpc address cannot be empty",
		},
		{
			name: "empty client id",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:  "chain1",
						RPC:      "http://localhost:26657",
						ClientID: "",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "client id cannot be empty",
		},
		{
			name: "duplicate chain id",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:  "chain1",
						RPC:      "http://localhost:26657",
						ClientID: "client1",
					},
					{
						ChainID:  "chain1",
						RPC:      "http://localhost:26658",
						ClientID: "client2",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "duplicate chain id",
		},
		{
			name: "empty attestator id",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:  "chain1",
						RPC:      "http://localhost:26657",
						ClientID: "client1",
					},
					{
						ChainID:  "chain2",
						RPC:      "http://localhost:26658",
						ClientID: "client2",
					},
				},
				AttestatorID:          "",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "attestator id cannot be empty",
		},
		{
			name: "empty signing private key path",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:  "chain1",
						RPC:      "http://localhost:26657",
						ClientID: "client1",
					},
					{
						ChainID:  "chain2",
						RPC:      "http://localhost:26658",
						ClientID: "client2",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "",
			},
			expErr: "private key path cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expErr)
			}

		})
	}
}
