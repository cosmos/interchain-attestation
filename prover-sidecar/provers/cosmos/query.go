package cosmos

import (
	"context"
	"gitlab.com/tozd/go/errors"
)

func (c *CosmosProver) queryLatestHeight(ctx context.Context) (int64, error) {
	stat, err := c.rpcClient.Status(ctx)
	if err != nil {
		return -1, err
	} else if stat.SyncInfo.CatchingUp {
		return -1, errors.Errorf("node at %s running chain %s not caught up", c.rpcAddr, c.chainID)
	}
	return stat.SyncInfo.LatestBlockHeight, nil
}