package cosmos

import (
	"context"
	"time"

	"gitlab.com/tozd/go/errors"

	querytypes "github.com/cosmos/cosmos-sdk/types/query"

	connectiontypes "github.com/cosmos/ibc-go/v9/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
)

const PaginationDelay = 10 * time.Millisecond

func (c *Attestator) queryConnectionsForClient(ctx context.Context, clientID string) ([]string, error) {
	qc := connectiontypes.NewQueryClient(c.clientConn)
	connections, err := qc.ClientConnections(ctx, &connectiontypes.QueryClientConnectionsRequest{
		ClientId: clientID,
	})
	if err != nil {
		return nil, err
	}

	return connections.ConnectionPaths, nil
}

func (c *Attestator) queryPacketCommitments(ctx context.Context, clientID string) (*chantypes.QueryPacketCommitmentsResponse, error) {
	// TODO: Check if the client is in the correct state
	// TODO: Cache some of this crap
	// TODO: Add support for ibc lite (i.e. skip a bunch of this)

	connections, err := c.queryConnectionsForClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if len(connections) == 0 {
		return nil, errors.Errorf("no connections found for client id %s", clientID)
	}

	qc := chantypes.NewQueryClient(c.clientConn)

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

		channels = append(channels, res.Channels...)

		next := res.GetPagination().GetNextKey()
		if len(next) == 0 {
			break
		}
		time.Sleep(PaginationDelay)
		p.Key = next
	}

	if len(channels) == 0 {
		return nil, errors.Errorf("no channels found for client id %s", clientID)
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
