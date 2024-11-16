package interchaintest

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/cometbft/cometbft/abci/types"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"

	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
)

// BroadcastMessages broadcasts the provided messages to the given chain and signs them on behalf of the provided user.
// Once the broadcast response is returned, we wait for two blocks to be created on chain.
func (s *E2ETestSuite) BroadcastMessages(ctx context.Context, chain *cosmos.CosmosChain, user ibc.Wallet, gas uint64, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	sdk.GetConfig().SetBech32PrefixForAccount(chain.Config().Bech32Prefix, chain.Config().Bech32Prefix+sdk.PrefixPublic)
	sdk.GetConfig().SetBech32PrefixForValidator(
		chain.Config().Bech32Prefix+sdk.PrefixValidator+sdk.PrefixOperator,
		chain.Config().Bech32Prefix+sdk.PrefixValidator+sdk.PrefixOperator+sdk.PrefixPublic,
	)

	broadcaster := cosmos.NewBroadcaster(s.T(), chain)

	broadcaster.ConfigureClientContextOptions(func(clientContext client.Context) client.Context {
		return clientContext.
			WithCodec(chain.Config().EncodingConfig.Codec).
			WithChainID(chain.Config().ChainID).
			WithTxConfig(chain.Config().EncodingConfig.TxConfig)
	})

	broadcaster.ConfigureFactoryOptions(func(factory tx.Factory) tx.Factory {
		return factory.WithGas(gas)
	})

	resp, err := cosmos.BroadcastTx(ctx, broadcaster, user, msgs...)
	if err != nil {
		return nil, err
	}

	// wait for 2 blocks for the transaction to be included
	s.Require().NoError(testutil.WaitForBlocks(ctx, 2, chain))

	if resp.Code != 0 {
		return nil, fmt.Errorf("transaction failed with code %d: %s", resp.Code, resp.RawLog)
	}

	return &resp, nil
}

// TODO: Replace with ibc-go/v9/testing when possible
// ParsePacketFromEvents parses events emitted from a MsgRecvPacket and returns
// the first packet found.
// Returns an error if no packet is found.
func ParsePacketFromEvents(events []abci.Event) (channeltypes.Packet, error) {
	packets, err := ParsePacketsFromEvents(events)
	if err != nil {
		return channeltypes.Packet{}, err
	}
	return packets[0], nil
}

// ParsePacketsFromEvents parses events emitted from a MsgRecvPacket and returns
// all the packets found.
// Returns an error if no packet is found.
func ParsePacketsFromEvents(events []abci.Event) ([]channeltypes.Packet, error) {
	ferr := func(err error) ([]channeltypes.Packet, error) {
		return nil, fmt.Errorf("ibctesting.ParsePacketsFromEvents: %w", err)
	}
	var packets []channeltypes.Packet
	for _, ev := range events {
		if ev.Type == channeltypes.EventTypeSendPacket {
			var packet channeltypes.Packet
			for _, attr := range ev.Attributes {
				switch attr.Key {
				case channeltypes.AttributeKeyDataHex:
					data, err := hex.DecodeString(attr.Value)
					if err != nil {
						return ferr(err)
					}
					packet.Data = data
				case channeltypes.AttributeKeySequence:
					seq, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return ferr(err)
					}

					packet.Sequence = seq

				case channeltypes.AttributeKeySrcPort:
					packet.SourcePort = attr.Value

				case channeltypes.AttributeKeySrcChannel:
					packet.SourceChannel = attr.Value

				case channeltypes.AttributeKeyDstPort:
					packet.DestinationPort = attr.Value

				case channeltypes.AttributeKeyDstChannel:
					packet.DestinationChannel = attr.Value

				case channeltypes.AttributeKeyTimeoutHeight:
					height, err := clienttypes.ParseHeight(attr.Value)
					if err != nil {
						return ferr(err)
					}

					packet.TimeoutHeight = height

				case channeltypes.AttributeKeyTimeoutTimestamp:
					timestamp, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return ferr(err)
					}

					packet.TimeoutTimestamp = timestamp

				default:
					continue
				}
			}

			packets = append(packets, packet)
		}
	}
	if len(packets) == 0 {
		return ferr(fmt.Errorf("acknowledgement event attribute not found"))
	}
	return packets, nil
}

// ParseAckFromEvents parses events emitted from a MsgRecvPacket and returns the
// acknowledgement.
func ParseAckFromEvents(events []abci.Event) ([]byte, error) {
	for _, ev := range events {
		if ev.Type == channeltypes.EventTypeWriteAck {
			for _, attr := range ev.Attributes {
				if attr.Key == channeltypes.AttributeKeyAckHex {
					value, err := hex.DecodeString(attr.Value)
					if err != nil {
						return nil, err
					}

					return value, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("acknowledgement event attribute not found")
}
