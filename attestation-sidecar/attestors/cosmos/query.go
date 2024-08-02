package cosmos

import (
	"context"
	querytypes "github.com/cosmos/cosmos-sdk/types/query"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"gitlab.com/tozd/go/errors"
	"time"
)

const PaginationDelay = 10 * time.Millisecond

func (c *CosmosAttestor) queryLatestHeight(ctx context.Context) (int64, error) {
	stat, err := c.rpcClient.Status(ctx)
	if err != nil {
		return -1, err
	} else if stat.SyncInfo.CatchingUp {
		return -1, errors.Errorf("node at %s running chain %s not caught up", c.rpcAddr, c.chainID)
	}
	return stat.SyncInfo.LatestBlockHeight, nil
}

func (c *CosmosAttestor) queryConnectionsForClient(ctx context.Context, clientID string) ([]string, error) {
	qc := connectiontypes.NewQueryClient(c)
	connections, err := qc.ClientConnections(ctx, &connectiontypes.QueryClientConnectionsRequest{
		ClientId: clientID,
	})
	if err != nil {
		return nil, err
	}

	var connectionPaths []string
	for _, connectionID := range connections.ConnectionPaths {
		connectionPaths = append(connectionPaths, connectionID)
	}

	return connectionPaths, nil
}

func (c *CosmosAttestor) queryPacketCommitments(ctx context.Context, clientID string) (*chantypes.QueryPacketCommitmentsResponse, error) {
	// TODO: Check if the client is in the correct state

	connections, err := c.queryConnectionsForClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	qc := chantypes.NewQueryClient(c)

	var channels []*chantypes.IdentifiedChannel
	p := defaultPageRequest()
	for _, connectionID := range connections {
		res, err := qc.ConnectionChannels(ctx, &chantypes.QueryConnectionChannelsRequest{
			Connection: connectionID,
			Pagination: p,
		})
		if err != nil {
			return nil, err
		}

		for _, channel := range res.Channels {
			channels = append(channels, channel)
		}

		next := res.GetPagination().GetNextKey()
		if len(next) == 0 {
			break
		}
		time.Sleep(PaginationDelay)
		p.Key = next
	}

	commitments := &chantypes.QueryPacketCommitmentsResponse{}
	p = defaultPageRequest()
	for _, channel := range channels {
		for {
			if channel.State != chantypes.OPEN {
				break
			}

			res, err := qc.PacketCommitments(ctx, &chantypes.QueryPacketCommitmentsRequest{
				PortId:     channel.PortId,
				ChannelId:  channel.ChannelId,
				Pagination: p,
			})
			if err != nil {
				return nil, err
			}

			commitments.Commitments = append(commitments.Commitments, res.Commitments...)
			commitments.Height = res.Height
			next := res.GetPagination().GetNextKey()
			if len(next) == 0 {
				break
			}
			time.Sleep(PaginationDelay)
			p.Key = next
		}
	}

	return commitments, nil
}

func defaultPageRequest() *querytypes.PageRequest {
	return &querytypes.PageRequest{
		Key:        []byte(""),
		Offset:     0,
		Limit:      1000,
		CountTotal: false,
	}
}