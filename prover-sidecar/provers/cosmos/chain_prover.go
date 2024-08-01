package cosmos

import (
	clientwrapper "github.com/strangelove-ventures/cometbft-client/client"
	"go.uber.org/zap"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/provers/chainprover"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/types"
	"time"
)

var _ chainprover.ChainProver = &CosmosProver{}

type CosmosProver struct {
	logger *zap.Logger

	rpcClient *clientwrapper.Client
	codec     Codec

	attestatorID string
	chainID string
	rpcAddr string
	clientID string
	signer func(msg []byte) ([]byte, error)

	latestProof *types.SignedPacketCommitmentsClaim
}

func NewCosmosProver(logger *zap.Logger, attestatorID string, chainID, rpcAddr, clientID string, signer func(msg []byte) ([]byte, error)) (*CosmosProver, error) {
	cometClient, err := clientwrapper.NewClient(rpcAddr, time.Second * 30) // TODO: Make timeout configurable per chain
	if err != nil {
		return nil, err
	}

	codec := newCodec()

	return &CosmosProver{
		logger:    logger,
		rpcClient: cometClient,
		codec:     codec,
		attestatorID: attestatorID,
		chainID:   chainID,
		rpcAddr:   rpcAddr,
		clientID: clientID,
		signer:    signer,
	}, nil
}

func (c *CosmosProver) ChainID() string {
	return c.chainID
}

func (c *CosmosProver) GetProof() *types.SignedPacketCommitmentsClaim {
	// TODO: Fetch from database?
	return c.latestProof
}
