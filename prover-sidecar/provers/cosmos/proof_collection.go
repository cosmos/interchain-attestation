package cosmos

import (
	"context"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/gjermundgaraba/pessimistic-validation/lightclient"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"github.com/gjermundgaraba/pessimistic-validation/proversidecar/types"
	"time"
)

func (c *CosmosProver) CollectProofs(ctx context.Context) error {
	c.logger.Info("Collecting proofs for chain", zap.String("chain_id", c.chainID))

	// TODO: add locks to prevent multiple CollectProofs from running at the same time

	// TODO: Replace this with actual proof collection logic, just here now to actually query the chain for something
	commitments, err := c.QueryPacketCommitments(ctx, c.clientID)
	if err != nil {
		return errors.Wrapf(err, "failed to query packet commitments for chain id %s", c.chainID)
	}

	revHeight := int64(commitments.Height.RevisionHeight)
	blockAtHeight, err := c.rpcClient.Block(ctx, &revHeight)
	if err != nil {
		return errors.Wrapf(err, "failed to query block for chain id %s", c.chainID)
	}

	var packetCommitments [][]byte
	for _, commitment := range commitments.Commitments {
		packetCommitments = append(packetCommitments, commitment.Data)
	}

	// TODO: Figure out how to use the same type as the lightclient, something stupid proto related
	claim := types.PacketCommitmentsClaim{
		Height:            clienttypes.Height{},
		Timestamp:         blockAtHeight.Block.Time,
		PacketCommitments: packetCommitments,
	}

	signableBytes := lightclient.GetSignableBytes(c.codec.Marshaler, lightclient.PacketCommitmentsClaim{
		Height:            claim.Height,
		Timestamp:         claim.Timestamp,
		PacketCommitments: claim.PacketCommitments,
	})
	signature, err := c.signer(signableBytes)
	if err != nil {
		return errors.Wrapf(err, "failed to sign packet commitments claim for chain id %s", c.chainID)
	}

	signedClaim := types.SignedPacketCommitmentsClaim{
		AttestatorId: 			[]byte(c.attestatorID),
		PacketCommitmentsClaim: claim,
		Signature:              signature,
	}

	// TODO: Put in a database or something?
	c.latestProof = &signedClaim

	c.logger.Info("Collected proofs for chain", zap.String("chain_id", c.chainID), zap.String("client_id", c.clientID), zap.Int("num_packet_commitments", len(packetCommitments)))

	// sleep for random 1-5 seconds
	time.Sleep(time.Duration(1 + time.Now().Unix()%5) * time.Second)

	return nil
}
