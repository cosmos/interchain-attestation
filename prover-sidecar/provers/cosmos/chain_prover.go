package cosmos

import (
	clientwrapper "github.com/strangelove-ventures/cometbft-client/client"
	"go.uber.org/zap"
	"time"

	"proversidecar/provers"
)

var _ provers.ChainProver = &CosmosProver{}

type CosmosProver struct {
	logger *zap.Logger

	rpcClient *clientwrapper.Client
	codec     Codec

	chainID string
	rpcAddr string

	clientID string
}

func NewCosmosProver(logger *zap.Logger, chainID, rpcAddr, clientID string) (*CosmosProver, error) {
	cometClient, err := clientwrapper.NewClient(rpcAddr, time.Second * 30) // TODO: Make timeout configurable per chain
	if err != nil {
		return nil, err
	}

	codec := newCodec()

	return &CosmosProver{
		logger:    logger,
		rpcClient: cometClient,
		codec:     codec,
		chainID:   chainID,
		rpcAddr:   rpcAddr,
		clientID: clientID,
	}, nil
}

func (c *CosmosProver) ChainID() string {
	return c.chainID
}
