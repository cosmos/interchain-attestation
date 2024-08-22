package relayer

import (
	"context"
	"fmt"
	"github.com/cometbft/cometbft/rpc/client"
	channeltypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v9/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v9/modules/core/exported"
	ibctesting "github.com/cosmos/ibc-go/v9/testing"
	"github.com/gjermundgaraba/interchain-attestation/sidecar/config"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
)

func (r *Relayer) CreateChannels(
	ctx context.Context,
	chainConfig config.CosmosChainConfig,
	connectionID string,
	portID string,
	version string,
	counterpartyChainConfig config.CosmosChainConfig,
	counterpartyConnectionID string,
	counterpartyPortID string,
) (string, string, error) {
	// Init channel on chain
	channelID, err := r.InitChannel(ctx, chainConfig, connectionID, portID, version, counterpartyPortID)
	if err != nil {
		return "", "", errors.Errorf("failed to init channel on %s: %w", chainConfig.ChainID, err)
	}

	updatedClientStateHeight, err := r.UpdateClient(ctx, counterpartyChainConfig, chainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to update client after init on %s with client id %s: %w", counterpartyChainConfig.ChainID, counterpartyChainConfig.ClientID, err)
	}

	// Open try on counterparty
	counterpartyChannelID, err := r.ChannelOpenTry(ctx, counterpartyChainConfig, counterpartyPortID, version, counterpartyConnectionID, chainConfig, portID, version, channelID, updatedClientStateHeight)
	if err != nil {
		return "", "", errors.Errorf("failed to open try channel on %s (proof height %d): %w", counterpartyChainConfig.ChainID, updatedClientStateHeight, err)
	}

	updatedClientStateHeight, err = r.UpdateClient(ctx, chainConfig, counterpartyChainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to update client after opentry on %s with client id %s: %w", chainConfig.ChainID, chainConfig.ClientID, err)
	}

	// Open ack on chain
	if err := r.ChannelOpenAck(ctx, chainConfig, portID, channelID, version, counterpartyChainConfig, portID, counterpartyChannelID, updatedClientStateHeight); err != nil {
		return "", "", errors.Errorf("failed to open ack channel on %s (proof height %d): %w", chainConfig.ChainID, updatedClientStateHeight, err)
	}

	updatedClientStateHeight, err = r.UpdateClient(ctx, counterpartyChainConfig, chainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to update client after openack on %s with client id %s: %w", counterpartyChainConfig.ChainID, counterpartyChainConfig.ClientID, err)
	}

	// Open confirm on counterparty
	if err := r.ChannelOpenConfirm(ctx, counterpartyChainConfig, counterpartyPortID, counterpartyChannelID, chainConfig, portID, channelID, updatedClientStateHeight); err != nil {
		return "", "", errors.Errorf("failed to open confirm channel on %s (proof height %d): %w", counterpartyChainConfig.ChainID, updatedClientStateHeight, err)
	}

	return channelID, counterpartyChannelID, nil
}

func (r *Relayer) InitChannel(ctx context.Context, chainConfig config.CosmosChainConfig, connectionID string, portID string, version string, counterpartyPortID string) (string, error) {
	r.logger.Debug("Init channel", zap.String("chain_id", chainConfig.ChainID), zap.String("connection_id", connectionID), zap.String("port_id", portID), zap.String("version", version), zap.String("counterparty_port_id", counterpartyPortID))

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return "", err
	}
	txf := r.createTxFactory(clientCtx, chainConfig)

	initMsg := channeltypes.NewMsgChannelOpenInit(
		portID,
		version,
		channeltypes.UNORDERED,
		[]string{connectionID},
		counterpartyPortID,
		clientCtx.From,
	)
	txResp, err := r.sendTx(clientCtx, txf, initMsg)
	if err != nil {
		return "", err
	}

	channelID, err := ibctesting.ParseChannelIDFromEvents(txResp.Events)
	if err != nil {
		return "", err
	}

	return channelID, nil
}

func (r *Relayer) ChannelOpenTry(
	ctx context.Context,
	chainConfig config.CosmosChainConfig,
	portID string,
	version string,
	connectionID string,
	counterpartyChainConfig config.CosmosChainConfig,
	counterpartyPortID string,
	counterpartyVersion string,
	counterpartyChannelID string,
	counterpartyProofHeight uint64,
) (string, error) {
	r.logger.Debug("Open try channel",
		zap.String("chain_id", chainConfig.ChainID),
		zap.String("port_id", portID),
		zap.String("version", version),
		zap.String("connection_id", connectionID),
		zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID),
		zap.String("counterparty_port_id", counterpartyPortID),
		zap.String("counterparty_version", counterpartyVersion),
		zap.String("counterparty_channel_id", counterpartyChannelID),
		zap.Uint64("counterparty_proof_height", counterpartyProofHeight),
	)

	initProof, err := r.GenerateChannelHandshakeProof(ctx, counterpartyChainConfig, counterpartyPortID, counterpartyChannelID, counterpartyProofHeight)
	if err != nil {
		return "", errors.Errorf("failed to generate channel handshake init proof: %w", err)
	}

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return "", err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	msg := channeltypes.NewMsgChannelOpenTry(
		portID,
		version,
		channeltypes.UNORDERED,
		[]string{connectionID},
		counterpartyPortID,
		counterpartyChannelID,
		counterpartyVersion,
		initProof,
		counterpartyChainConfig.GetClientHeight(counterpartyProofHeight),
		clientCtx.From,
	)

	txResp, err := r.sendTx(clientCtx, txf, msg)
	if err != nil {
		return "", errors.Errorf("failed to send tx: %w", err)
	}

	channelID, err := ibctesting.ParseChannelIDFromEvents(txResp.Events)
	if err != nil {
		return "", err
	}

	return channelID, nil
}

func (r *Relayer) ChannelOpenAck(
	ctx context.Context,
	chainConfig config.CosmosChainConfig,
	portID string,
	channelID string,
	version string,
	counterpartyChainConfig config.CosmosChainConfig,
	counterpartyPortID string,
	counterpartyChannelID string,
	counterpartyProofHeight uint64,
) error {
	r.logger.Debug("Open ack channel",
		zap.String("chain_id", chainConfig.ChainID),
		zap.String("port_id", portID),
		zap.String("channel_id", channelID),
		zap.String("version", version),
		zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID),
		zap.String("counterparty_port_id", counterpartyPortID),
		zap.String("counterparty_channel_id", counterpartyChannelID),
		zap.Uint64("counterparty_proof_height", counterpartyProofHeight),
	)

	tryProof, err := r.GenerateChannelHandshakeProof(ctx, counterpartyChainConfig, counterpartyPortID, counterpartyChannelID, counterpartyProofHeight)
	if err != nil {
		return errors.Errorf("failed to generate channel handshake try proof: %w", err)
	}

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return err
	}
	txf := r.createTxFactory(clientCtx, chainConfig)

	ackMsg := channeltypes.NewMsgChannelOpenAck(
		portID,
		channelID,
		counterpartyChannelID,
		version,
		tryProof,
		counterpartyChainConfig.GetClientHeight(counterpartyProofHeight),
		clientCtx.From,
	)

	_, err = r.sendTx(clientCtx, txf, ackMsg)
	if err != nil {
		return errors.Errorf("failed to send tx: %w", err)
	}

	return nil
}

func (r *Relayer) ChannelOpenConfirm(
	ctx context.Context,
	chainConfig config.CosmosChainConfig,
	portID string,
	channelID string,
	counterpartyChainConfig config.CosmosChainConfig,
	counterpartyPortID string,
	counterpartyChannelID string,
	counterpartyProofHeight uint64,
) error {
	r.logger.Debug("Open confirm channel",
		zap.String("chain_id", chainConfig.ChainID),
		zap.String("port_id", portID),
		zap.String("channel_id", channelID),
		zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID),
		zap.String("counterparty_port_id", counterpartyPortID),
		zap.String("counterparty_channel_id", counterpartyChannelID),
		zap.Uint64("counterparty_proof_height", counterpartyProofHeight),
	)

	ackProof, err := r.GenerateChannelHandshakeProof(ctx, counterpartyChainConfig, counterpartyPortID, counterpartyChannelID, counterpartyProofHeight)
	if err != nil {
		return errors.Errorf("failed to generate channel handshake ack proof: %w", err)
	}

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	confirmMsg := channeltypes.NewMsgChannelOpenConfirm(
		portID,
		channelID,
		ackProof,
		counterpartyChainConfig.GetClientHeight(counterpartyProofHeight),
		clientCtx.From,
	)

	_, err = r.sendTx(clientCtx, txf, confirmMsg)
	if err != nil {
		return errors.Errorf("failed to send tx: %w", err)
	}

	return nil
}

func (r *Relayer) GenerateChannelHandshakeProof(ctx context.Context, chainConfig config.CosmosChainConfig, portID string, channelID string, proofHeight uint64) ([]byte, error) {
	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("store/%s/key", ibcexported.StoreKey)
	key := host.ChannelKey(portID, channelID)
	resp, err := clientCtx.Client.ABCIQueryWithOptions(ctx, path, key, client.ABCIQueryOptions{
		Height: int64(proofHeight-1),
		Prove:  true,
	})

	merkleProof, err := commitmenttypes.ConvertProofs(resp.Response.ProofOps)
	if err != nil {
		return nil, err
	}

	proofBz, err := r.cdc.Marshal(&merkleProof)
	if err != nil {
		return nil, err
	}

	return proofBz, nil
}
