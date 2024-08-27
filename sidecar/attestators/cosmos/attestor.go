package cosmos

import (
	"time"

	clientwrapper "github.com/strangelove-ventures/cometbft-client/client"
	"go.uber.org/zap"

	"github.com/cosmos/interchain-attestation/sidecar/attestators/attestator"
	"github.com/cosmos/interchain-attestation/sidecar/config"
)

var _ attestator.Attestator = &Attestator{}

type Attestator struct {
	logger *zap.Logger

	clientConn  *ClientConn
	cometClient *clientwrapper.Client
	codec       CodecConfig

	attestatorID string
	config       config.CosmosChainConfig
	signer       func(msg []byte) ([]byte, error)
}

func NewCosmosAttestator(logger *zap.Logger, attestatorID string, config config.CosmosChainConfig, signer func(msg []byte) ([]byte, error)) (*Attestator, error) {
	cometClient, err := clientwrapper.NewClient(config.RPC, time.Second*30) // TODO: Make timeout configurable per chain
	if err != nil {
		return nil, err
	}

	codec := NewCodecConfig()

	clientConn := &ClientConn{
		cometClient: cometClient,
		codec:       codec,
	}

	return &Attestator{
		logger: logger,

		clientConn:  clientConn,
		cometClient: cometClient,
		codec:       codec,

		attestatorID: attestatorID,
		config:       config,
		signer:       signer,
	}, nil
}

func (c *Attestator) ChainID() string {
	return c.config.ChainID
}
