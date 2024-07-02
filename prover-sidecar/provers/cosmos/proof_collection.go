package cosmos

import (
	"context"
	"fmt"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"time"
)

func (c *CosmosProver) CollectProofs(ctx context.Context) error {
	c.logger.Info("Collecting proofs for chain", zap.String("chain_id", c.chainID))

	// TODO: add locks to prevent multiple CollectProofs from running at the same time

	// TODO: Replace this with actual proof collection logic, just here now to actually query the chain for something
	status, err := c.rpcClient.Status(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to query status for chain id %s", c.chainID)
	}

	c.tmpProof = fmt.Sprintf("app hash %s for block %d on chain %s", status.SyncInfo.LatestAppHash.String(), status.SyncInfo.LatestBlockHeight, c.chainID)

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
