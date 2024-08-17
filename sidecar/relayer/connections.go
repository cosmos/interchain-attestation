package relayer

import (
	"context"
	"fmt"
	"github.com/cometbft/cometbft/rpc/client"
	connectiontypes "github.com/cosmos/ibc-go/v9/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v9/modules/core/24-host"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
	ibcexported "github.com/cosmos/ibc-go/v9/modules/core/exported"
	ibctesting "github.com/cosmos/ibc-go/v9/testing"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/config"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
)

func (r *Relayer) CreateConnections(ctx context.Context, chainConfig config.CosmosChainConfig, counterpartyChainConfig config.CosmosChainConfig) (string, string, error) {
	// Init connection on chain
	connectionID, err := r.InitConnection(ctx, chainConfig, counterpartyChainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to init connection on %s: %w", chainConfig.ChainID, err)
	}

	updatedClientStateHeight, err := r.UpdateClient(ctx, counterpartyChainConfig, chainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to update client after init on %s with client id %s: %w", counterpartyChainConfig.ChainID, counterpartyChainConfig.ClientID, err)
	}

	// Open try on counterparty
	counterpartyConnectionID, err := r.ConnectionOpenTry(ctx, counterpartyChainConfig, chainConfig, connectionID, updatedClientStateHeight)
	if err != nil {
		return "", "", errors.Errorf("failed to open try connection on %s: %w", counterpartyChainConfig.ChainID, err)
	}

	updatedClientStateHeight, err = r.UpdateClient(ctx, chainConfig, counterpartyChainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to update client after opentry on %s with client id %s: %w", chainConfig.ChainID, chainConfig.ClientID, err)
	}

	// Open ack on chain
	if err := r.ConnectionOpenAck(ctx, chainConfig, connectionID, counterpartyChainConfig, counterpartyConnectionID, updatedClientStateHeight); err != nil {
		return "", "", errors.Errorf("failed to open ack connection on %s: %w", chainConfig.ChainID, err)
	}

	updatedClientStateHeight, err = r.UpdateClient(ctx, counterpartyChainConfig, chainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to update client after openack on %s with client id %s: %w", counterpartyChainConfig.ChainID, counterpartyChainConfig.ClientID, err)
	}

	if err := r.ConnectionOpenConfirm(ctx, counterpartyChainConfig, counterpartyConnectionID, chainConfig, connectionID, updatedClientStateHeight); err != nil {
		return "", "", errors.Errorf("failed to open confirm connection on %s: %w", counterpartyChainConfig.ChainID, err)
	}

	return connectionID, counterpartyConnectionID, nil
}

func (r *Relayer) InitConnection(ctx context.Context, chainConfig config.CosmosChainConfig, counterpartyChainConfig config.CosmosChainConfig) (string, error) {
	r.logger.Debug("Init connection", zap.String("chain_id", chainConfig.ChainID), zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID))

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return "", err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	var version *connectiontypes.Version // Can be nil? Not sure.
	merklePrefix := commitmenttypes.NewMerklePrefix([]byte(ibcexported.StoreKey))
	initMsg := connectiontypes.NewMsgConnectionOpenInit(
		chainConfig.ClientID,
		counterpartyChainConfig.ClientID,
		merklePrefix,
		version,
		0,
		clientCtx.From,
	)

	txResp, err := r.sendTx(clientCtx, txf, initMsg)
	if err != nil {
		return "", err
	}

	connectionID, err := ibctesting.ParseConnectionIDFromEvents(txResp.Events)
	if err != nil {
		return "", err
	}

	return connectionID, nil
}

func (r *Relayer) ConnectionOpenTry(
	ctx context.Context,
	chainConfig config.CosmosChainConfig,
	counterpartyChainConfig config.CosmosChainConfig,
	counterpartyConnectionID string,
	counterpartyProofHeight uint64,
) (string, error) {
	r.logger.Debug("Open try connection", zap.String("chain_id", chainConfig.ChainID), zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID), zap.String("counterparty_connection_id", counterpartyConnectionID), zap.Uint64("counterparty_proof_height", counterpartyProofHeight))
	counterpartyPrefix := commitmenttypes.NewMerklePrefix([]byte(exported.StoreKey))

	initProof, err := r.GenerateConnectionHandshakeProof(ctx, counterpartyChainConfig, counterpartyConnectionID, counterpartyProofHeight)
	if err != nil {
		return "", errors.Errorf("failed to generate connection handshake proof: %w", err)
	}

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return "", err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	msg := connectiontypes.NewMsgConnectionOpenTry(
		chainConfig.ClientID,
		counterpartyConnectionID,
		counterpartyChainConfig.ClientID,
		counterpartyPrefix,
		connectiontypes.GetCompatibleVersions(),
		0,
		initProof,
		counterpartyChainConfig.GetClientHeight(counterpartyProofHeight),
		clientCtx.From,
	)

	txResp, err := r.sendTx(clientCtx, txf, msg)
	if err != nil {
		return "", errors.Errorf("failed to send tx: %w", err)
	}

	connectionID, err := ibctesting.ParseConnectionIDFromEvents(txResp.Events)
	if err != nil {
		return "", err
	}

	return connectionID, nil
}

func (r *Relayer) ConnectionOpenAck(
	ctx context.Context,
	chainConfig config.CosmosChainConfig,
	connectionID string,
	counterpartyChainConfig config.CosmosChainConfig,
	counterpartyConnectionID string,
	counterpartyProofHeight uint64,
) error {
	r.logger.Debug("Open ack connection", zap.String("chain_id", chainConfig.ChainID), zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID), zap.String("connection_id", connectionID), zap.String("counterparty_connection_id", counterpartyConnectionID), zap.Uint64("counterparty_proof_height", counterpartyProofHeight))

	tryProof, err := r.GenerateConnectionHandshakeProof(ctx, counterpartyChainConfig, counterpartyConnectionID, counterpartyProofHeight)
	if err != nil {
		return errors.Errorf("failed to generate connection handshake proof: %w", err)
	}

	version := connectiontypes.GetCompatibleVersions()[0]

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	ackMsg := connectiontypes.NewMsgConnectionOpenAck(
		connectionID,
		counterpartyConnectionID,
		tryProof,
		counterpartyChainConfig.GetClientHeight(counterpartyProofHeight),
		version,
		clientCtx.From,
	)

	if _, err := r.sendTx(clientCtx, txf, ackMsg); err != nil {
		return errors.Errorf("failed to send tx %v: %w", ackMsg, err)
	}

	return nil
}

func (r *Relayer) ConnectionOpenConfirm(ctx context.Context, chainConfig config.CosmosChainConfig, connectionID string, counterpartyChainConfig config.CosmosChainConfig, counterpartyConnectionID string, counterpartyProofHeight uint64) error {
	r.logger.Debug("Open confirm connection", zap.String("chain_id", chainConfig.ChainID), zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID), zap.String("connection_id", connectionID), zap.String("counterparty_connection_id", counterpartyConnectionID), zap.Uint64("counterparty_proof_height", counterpartyProofHeight))

	ackProof, err := r.GenerateConnectionHandshakeProof(ctx, counterpartyChainConfig, counterpartyConnectionID, counterpartyProofHeight)
	if err != nil {
		return errors.Errorf("failed to generate connection handshake proof: %w", err)
	}

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	msg := connectiontypes.NewMsgConnectionOpenConfirm(
		connectionID,
		ackProof,
		counterpartyChainConfig.GetClientHeight(counterpartyProofHeight),
		clientCtx.From,
	)

	if _, err := r.sendTx(clientCtx, txf, msg); err != nil {
		return errors.Errorf("failed to send tx: %w", err)
	}

	return nil
}

func (r *Relayer) GenerateConnectionHandshakeProof(ctx context.Context, chainConfig config.CosmosChainConfig, connectionID string, proofHeight uint64) ([]byte, error) {
	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("store/%s/key", ibcexported.StoreKey)
	key := host.ConnectionKey(connectionID)
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

func (r *Relayer) QueryConnection(ctx context.Context, chainConfig config.CosmosChainConfig, connectionID string) (*connectiontypes.ConnectionEnd, error) {
	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return nil, err
	}

	queryClient := connectiontypes.NewQueryClient(clientCtx)
	req := &connectiontypes.QueryConnectionRequest{
		ConnectionId: connectionID,
	}

	res, err := queryClient.Connection(clientCtx.CmdContext, req)
	if err != nil {
		return nil, err
	}

	return res.Connection, nil
}
