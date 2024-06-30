package cosmos

import (
	"context"
	"go.uber.org/zap"
	"time"
)

func (c *CosmosProver) CollectProofs(ctx context.Context) error {
	c.logger.Info("Collecting proofs for chain", zap.String("chain_id", c.chainID))

	// sleep for random 1-5 seconds
	time.Sleep(time.Duration(1 + time.Now().Unix()%5) * time.Second)

	return nil
}

/*func (c *CosmosProver) CollectProofs(ctx context.Context) error {
	queryClient := channeltypes.NewQueryClient(c)
	commitmentsResp, err := queryClient.PacketCommitments(ctx, &channeltypes.QueryPacketCommitmentsRequest{
		PortId:     c.portID,
		ChannelId:  c.channelID,
		Pagination: nil,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to query for packet commitments for chain id %s", c.chainID)
	}

	for _, commitment := range commitmentsResp.Commitments {
		fmt.Println("Found commitment with sequence", commitment.Sequence)
		fmt.Println(commitment.String())
	}

	return nil
}*/
