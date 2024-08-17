package relayer

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v9/testing"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/config"
	"time"
)

func (r *Relayer) Transfer(ctx context.Context, fromChainConfig config.CosmosChainConfig, channelID string, to string, amount string) (channeltypes.Packet, error) {
	clientCtx, err := r.createClientCtx(ctx, fromChainConfig)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	txf := r.createTxFactory(clientCtx, fromChainConfig)

	coins, err := sdk.ParseCoinsNormalized(amount)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	// create a new transfer message
	transferMsg := &transfertypes.MsgTransfer{
		SourcePort:       transfertypes.PortID,
		SourceChannel:    channelID,
		Sender:           clientCtx.From,
		Receiver:         to,
		TimeoutTimestamp: uint64(time.Now().Add(time.Minute * 10).UnixNano()),
		Tokens:           coins,
	}

	txResp, err := r.sendTx(clientCtx, txf, transferMsg)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	packet, err := ibctesting.ParsePacketFromEvents(txResp.Events)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	return packet, nil
}
