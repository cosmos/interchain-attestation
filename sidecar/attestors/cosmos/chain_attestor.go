package cosmos

import (
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestors/chainattestor"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
	clientwrapper "github.com/strangelove-ventures/cometbft-client/client"
	"go.uber.org/zap"
	"time"
)

var _ chainattestor.ChainAttestor = &CosmosAttestor{}

type CosmosAttestor struct {
	logger *zap.Logger

	rpcClient *clientwrapper.Client
	codec     Codec

	attestatorID string
	chainID string
	rpcAddr string
	clientID string
	signer func(msg []byte) ([]byte, error)

	latestAttestation *types.Attestation
}

func NewCosmosAttestor(logger *zap.Logger, attestatorID string, chainID, rpcAddr, clientID string, signer func(msg []byte) ([]byte, error)) (*CosmosAttestor, error) {
	cometClient, err := clientwrapper.NewClient(rpcAddr, time.Second * 30) // TODO: Make timeout configurable per chain
	if err != nil {
		return nil, err
	}

	codec := newCodec()

	return &CosmosAttestor{
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

func (c *CosmosAttestor) ChainID() string {
	return c.chainID
}

func (c *CosmosAttestor) GetLatestAttestation() *types.Attestation {
	// TODO: Fetch from database?
	return c.latestAttestation
}
