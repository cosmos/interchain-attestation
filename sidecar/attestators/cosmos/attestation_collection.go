package cosmos

import (
	"context"

	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"

	chantypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"

	"github.com/cosmos/interchain-attestation/core/types"
)

func (c *Attestator) CollectAttestation(ctx context.Context) (types.Attestation, error) {
	c.logger.Info("Collecting attestationData for chain", zap.String("chain_id", c.config.ChainID), zap.String("client_id", c.config.ClientID))

	// TODO: add locks to prevent multiple CollectAttestation from running at the same time

	commitments, err := c.queryPacketCommitments(ctx, c.config.ClientID)
	if err != nil || commitments.Height.RevisionHeight == 0 {
		c.logger.Info("Failed to query packet commitments, but to keep the client updated, we will return empty list of commitments", zap.Error(err))

		resp, err := c.cometClient.ABCIInfo(ctx)
		if err != nil {
			return types.Attestation{}, errors.Errorf("failed to query status for chain id %s: %w", c.config.ChainID, err)
		}
		commitments = &chantypes.QueryPacketCommitmentsResponse{
			Commitments: []*chantypes.PacketState{},
			Height:      c.config.GetClientHeight(uint64(resp.Response.LastBlockHeight - 1)),
		}
		// return types.Attestation{}, errors.Errorf("failed to query packet commitments for client id %s on chain id %s: %w", c.config.ClientID, c.config.ChainID, err)
	}

	height := commitments.Height
	revHeight := int64(height.GetRevisionHeight())
	blockAtHeight, err := c.cometClient.Block(ctx, &revHeight)
	if err != nil {
		return types.Attestation{}, errors.Errorf("failed to query block for client id %s (height %d) on chain id %s: %w", c.config.ClientID, revHeight, c.config.ChainID, err)
	}

	var packetCommitments [][]byte
	for _, commitment := range commitments.Commitments {
		packetCommitments = append(packetCommitments, commitment.Data)
	}

	attestationData := types.IBCData{
		ChainId:           c.config.ChainID,
		ClientId:          c.config.ClientID,
		ClientToUpdate:    c.config.ClientToUpdate,
		Height:            height,
		Timestamp:         blockAtHeight.Block.Time,
		PacketCommitments: packetCommitments,
	}

	c.logger.Debug("Generated attestation data",
		zap.String("chain_id", c.config.ChainID),
		zap.String("client_id", c.config.ClientID),
		zap.String("client_to_update", c.config.ClientToUpdate),
		zap.Int64("height", revHeight),
		zap.Time("timestamp", blockAtHeight.Block.Time),
		zap.Int("packet_commitments", len(packetCommitments)),
	)

	attestation := types.Attestation{
		AttestatorId: []byte(c.attestatorID),
		AttestedData: attestationData,
	}

	return attestation, nil
}
