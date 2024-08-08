package cosmos

import (
	"context"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/gjermundgaraba/pessimistic-validation/core/lightclient"
	"github.com/gjermundgaraba/pessimistic-validation/core/types"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
)

func (c *CosmosAttestor) CollectAttestation(ctx context.Context) (types.Attestation, error) {
	c.logger.Info("Collecting attestationData for chain", zap.String("chain_id", c.chainID))

	// TODO: add locks to prevent multiple CollectAttestation from running at the same time

	commitments, err := c.queryPacketCommitments(ctx, c.clientID)
	if err != nil {
		return types.Attestation{}, errors.Wrapf(err, "failed to query packet commitments for chain id %s", c.chainID)
	}

	revHeight := int64(commitments.Height.RevisionHeight)
	blockAtHeight, err := c.rpcClient.Block(ctx, &revHeight)
	if err != nil {
		return types.Attestation{}, errors.Wrapf(err, "failed to query block for chain id %s", c.chainID)
	}

	var packetCommitments [][]byte
	for _, commitment := range commitments.Commitments {
		packetCommitments = append(packetCommitments, commitment.Data)
	}

	attestationData := types.IBCData{
		Height:            clienttypes.Height{},
		Timestamp:         blockAtHeight.Block.Time,
		PacketCommitments: packetCommitments,
	}

	signableBytes := lightclient.GetSignableBytes(c.codec.Marshaler, attestationData)
	signature, err := c.signer(signableBytes)
	if err != nil {
		return types.Attestation{}, errors.Wrapf(err, "failed to sign attestation data for chain id %s", c.chainID)
	}

	attestation := types.Attestation{
		AttestatorId:           []byte(c.attestatorID),
		AttestedData: attestationData,
		Signature:              signature,
	}

	c.logger.Info("Collected attestation for chain", zap.String("chain_id", c.chainID), zap.String("client_id", c.clientID), zap.Int("num_packet_commitments", len(packetCommitments)))

	return attestation, nil
}
