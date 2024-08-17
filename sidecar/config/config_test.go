package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Validate(t *testing.T) {

	tests := []struct {
		name   string
		config Config
		expErr string
	}{
		{
			name: "valid config",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:        "chain1",
						RPC:            "http://localhost:26657",
						ClientID:       "client1",
						Attestation:    true,
						ClientToUpdate: "client1",
						AddressPrefix:  "",
						KeyringBackend: "",
						KeyName:        "",
						Gas:            "",
						GasPrices:      "",
						GasAdjustment:  0,
					},
					{
						ChainID:        "chain2",
						RPC:            "http://localhost:26658",
						ClientID:       "",
						Attestation:    false,
						ClientToUpdate: "",
						AddressPrefix:  "",
						KeyringBackend: "",
						KeyName:        "",
						Gas:            "",
						GasPrices:      "",
						GasAdjustment:  0,
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
						ChainID:             "",
						RPC:                 "http://localhost:26657",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "chain id cannot be empty",
		},
		{
			name: "empty attestation client id",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:        "chain1",
						RPC:            "http://localhost:26657",
						ClientID:       "",
						Attestation:    true,
						ClientToUpdate: "client1",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "client id cannot be empty when attestation is true",
		},
		{
			name: "empty client to update",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:        "chain1",
						RPC:            "http://localhost:26657",
						ClientID:       "client1",
						Attestation:    true,
						ClientToUpdate: "",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "client to update cannot be empty when attestation is true",
		},
		{
			name: "empty rpc",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:             "chain1",
						RPC:                 "",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "rpc address cannot be empty",
		},
		{
			name: "duplicate chain id",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:        "chain1",
						RPC:            "http://localhost:26657",
						ClientID:       "client1",
						Attestation:    true,
						ClientToUpdate: "client1",
					},
					{
						ClientID: "client2",
						ChainID:  "chain1",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "duplicate chain id",
		},
		{
			name: "duplicate client to update",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:        "chain1",
						RPC:            "http://localhost:26657",
						ClientID:       "client1",
						Attestation:    true,
						ClientToUpdate: "client1",
					},
					{
						ChainID:        "chain2",
						RPC:            "http://localhost:26658",
						ClientID:       "client2",
						Attestation:    true,
						ClientToUpdate: "client1",
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "duplicate client to update",
		},
		{
			name: "empty attestator id",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:        "chain1",
						RPC:            "http://localhost:26657",
						ClientID:       "client1",
						Attestation:    true,
						ClientToUpdate: "client1",
					},
					{
						ChainID:     "chain2",
						RPC:         "http://localhost:26658",
						Attestation: false,
					},
				},
				AttestatorID:          "",
				SigningPrivateKeyPath: "some/path",
			},
			expErr: "attestator id cannot be empty if any chains have attestation true",
		},
		{
			name: "empty signing private key path if any chain have attestation true",
			config: Config{
				CosmosChains: []CosmosChainConfig{
					{
						ChainID:        "chain1",
						RPC:            "http://localhost:26657",
						ClientID:       "client1",
						Attestation:    true,
						ClientToUpdate: "client1",
					},
					{
						ChainID:     "chain2",
						RPC:         "http://localhost:26658",
						Attestation: false,
					},
				},
				AttestatorID:          "test-attestator-id",
				SigningPrivateKeyPath: "",
			},
			expErr: "private key path cannot be empty if any chains have attestation true",
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
