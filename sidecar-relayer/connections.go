package main

import (
	"context"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/store/rootmulti"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	client2 "github.com/cometbft/cometbft/rpc/client"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	legacyerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	conntypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (r *Relayer) InitConnection() (string, error) {
	clientCtx := r.createClientCtx(r.config.SrcChain)
	txf := r.createTxFactory(clientCtx, r.config.SrcChain)

	merklePrefix := commitmenttypes.NewMerklePrefix([]byte(ibcexported.StoreKey))
	initMsg := connectiontypes.NewMsgConnectionOpenInit(
		r.config.SrcChain.ClientId,
		r.config.DstChain.ClientId,
		merklePrefix,
		nil,
		0,
		clientCtx.From,
	)

	txResp, err := r.sendTx(clientCtx, txf, initMsg)
	if err != nil {
		return "", err
	}

	connectionID, err := parseConnectionIDFromEvents(txResp.Events)
	if err != nil {
		return "", err
	}

	return connectionID, nil
}

func (r *Relayer) OpenTryConnection() (string, error) {
	clientCtx := r.createClientCtx(r.config.DstChain)
	txf := r.createTxFactory(clientCtx, r.config.DstChain)

	latestHeight, err := r.getLatestHeight(r.config.SrcChain)
	if err != nil {
		return "", err
	}

	clientState, clientStateProof, consensusProof, connectionProof, connectionProofHeight, err := r.GenerateConnHandshakeProof(r.ctx, r.config.SrcChain, latestHeight, r.config.SrcChain.ClientId, r.config.SrcChain.ConnectionId)

	merklePrefix := commitmenttypes.NewMerklePrefix([]byte(ibcexported.StoreKey))
	initMsg := connectiontypes.NewMsgConnectionOpenTry(
		r.config.DstChain.ChainId,
		r.config.SrcChain.ConnectionId,
		r.config.SrcChain.ClientId,
		clientState,
		merklePrefix,
		[]*connectiontypes.Version{},
		0,
		connectionProof,
		clientStateProof,
		consensusProof,
		connectionProofHeight.(clienttypes.Height),
		clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: uint64(latestHeight),
		},
		clientCtx.From,
	)

	txResp, err := r.sendTx(clientCtx, txf, initMsg)
	if err != nil {
		return "", err
	}

	connectionID, err := parseConnectionIDFromEvents(txResp.Events)
	if err != nil {
		return "", err
	}

	return connectionID, nil
}

func (r *Relayer) GenerateConnHandshakeProof(ctx context.Context, chain ChainConfig, height int64, clientId, connId string) (clientState ibcexported.ClientState, clientStateProof []byte, consensusProof []byte, connectionProof []byte, connectionProofHeight ibcexported.Height, err error) {
	var (
		clientStateRes     *clienttypes.QueryClientStateResponse
		consensusStateRes  *clienttypes.QueryConsensusStateResponse
		connectionStateRes *conntypes.QueryConnectionResponse
		eg                 = new(errgroup.Group)
	)

	// query for the client state for the proof and get the height to query the consensus state at.
	clientStateRes, err = r.QueryClientStateResponse(ctx, chain, height, clientId)
	if err != nil {
		return nil, nil, nil, nil, clienttypes.Height{}, err
	}

	clientState, err = clienttypes.UnpackClientState(clientStateRes.ClientState)
	if err != nil {
		return nil, nil, nil, nil, clienttypes.Height{}, err
	}

	eg.Go(func() error {
		var err error
		consensusStateRes, err = r.QueryClientConsensusState(ctx, chain, height, clientId, clientStateRes.ProofHeight)
		return err
	})
	eg.Go(func() error {
		var err error
		connectionStateRes, err = r.QueryConnection(ctx, chain, height, connId)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, nil, nil, clienttypes.Height{}, err
	}

	return clientState, clientStateRes.Proof, consensusStateRes.Proof, connectionStateRes.Proof, connectionStateRes.ProofHeight, nil
}

func (r *Relayer) QueryClientStateResponse(ctx context.Context, chain ChainConfig, height int64, srcClientId string) (*clienttypes.QueryClientStateResponse, error) {
	key := host.FullClientStateKey(srcClientId)

	value, proofBz, proofHeight, err := r.QueryTendermintProof(ctx, chain, height, key)
	if err != nil {
		return nil, err
	}

	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(clienttypes.ErrClientNotFound, srcClientId)
	}

	cdc := codec.NewProtoCodec(r.cdc.InterfaceRegistry())

	clientState, err := clienttypes.UnmarshalClientState(cdc, value)
	if err != nil {
		return nil, err
	}

	anyClientState, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	return &clienttypes.QueryClientStateResponse{
		ClientState: anyClientState,
		Proof:       proofBz,
		ProofHeight: proofHeight,
	}, nil
}

func (r *Relayer)  QueryTendermintProof(ctx context.Context, chain ChainConfig, height int64, key []byte) ([]byte, []byte, clienttypes.Height, error) {
	// ABCI queries at heights 1, 2 or less than or equal to 0 are not supported.
	// Base app does not support queries for height less than or equal to 1.
	// Therefore, a query at height 2 would be equivalent to a query at height 3.
	// A height of 0 will query with the latest state.
	if height != 0 && height <= 2 {
		return nil, nil, clienttypes.Height{}, fmt.Errorf("proof queries at height <= 2 are not supported")
	}

	// Use the IAVL height if a valid tendermint height is passed in.
	// A height of 0 will query with the latest state.
	if height != 0 {
		height--
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", ibcexported.StoreKey),
		Height: height,
		Data:   key,
		Prove:  true,
	}

	res, err := r.QueryABCI(ctx, chain, req)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	cdc := codec.NewProtoCodec(r.cdc.InterfaceRegistry())

	proofBz, err := cdc.Marshal(&merkleProof)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	revision := clienttypes.ParseChainID(chain.ChainId)
	return res.Value, proofBz, clienttypes.NewHeight(revision, uint64(res.Height)+1), nil
}

func (r *Relayer) QueryABCI(ctx context.Context, chain ChainConfig, req abci.RequestQuery) (abci.ResponseQuery, error) {
	opts := client2.ABCIQueryOptions{
		Height: req.Height,
		Prove:  req.Prove,
	}

	result, err := r.ABCIQueryWithOptions(ctx, chain, req.Path, req.Data, opts)
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	if !result.Response.IsOK() {
		return abci.ResponseQuery{}, sdkErrorToGRPCError(result.Response)
	}

	// data from trusted node or subspace query doesn't need verification
	if !opts.Prove || !isQueryStoreWithProof(req.Path) {
		return result.Response, nil
	}

	return result.Response, nil
}

func sdkErrorToGRPCError(resp abci.ResponseQuery) error {
	switch resp.Code {
	case legacyerrors.ErrInvalidRequest.ABCICode():
		return status.Error(codes.InvalidArgument, resp.Log)
	case legacyerrors.ErrUnauthorized.ABCICode():
		return status.Error(codes.Unauthenticated, resp.Log)
	case legacyerrors.ErrKeyNotFound.ABCICode():
		return status.Error(codes.NotFound, resp.Log)
	default:
		return status.Error(codes.Unknown, resp.Log)
	}
}

func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}

	paths := strings.SplitN(path[1:], "/", 3)

	switch {
	case len(paths) != 3:
		return false
	case paths[0] != "store":
		return false
	case rootmulti.RequireProof("/" + paths[2]):
		return true
	}

	return false
}

func (r *Relayer) ABCIQueryWithOptions(
	ctx context.Context,
	chain ChainConfig,
	path string,
	data bytes.HexBytes,
	opts rpcclient.ABCIQueryOptions,
) (*coretypes.ResultABCIQuery, error) {
	rpcClient, err := client.NewClientFromNode(chain.RpcAddr)
	if err != nil {
		return nil, err
	}

	return rpcClient.ABCIQueryWithOptions(ctx, path, data, opts)
}

func (r *Relayer) QueryClientConsensusState(ctx context.Context, chain ChainConfig, chainHeight int64, clientid string, clientHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	key := host.FullConsensusStateKey(clientid, clientHeight)

	value, proofBz, proofHeight, err := r.QueryTendermintProof(ctx, chain, chainHeight, key)
	if err != nil {
		return nil, err
	}

	// check if consensus state exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(clienttypes.ErrConsensusStateNotFound, clientid)
	}

	cdc := codec.NewProtoCodec(r.cdc.InterfaceRegistry())

	cs, err := clienttypes.UnmarshalConsensusState(cdc, value)
	if err != nil {
		return nil, err
	}

	anyConsensusState, err := clienttypes.PackConsensusState(cs)
	if err != nil {
		return nil, err
	}

	return &clienttypes.QueryConsensusStateResponse{
		ConsensusState: anyConsensusState,
		Proof:          proofBz,
		ProofHeight:    proofHeight,
	}, nil
}

func (r *Relayer) QueryConnection(ctx context.Context, chain ChainConfig, height int64, connectionid string) (*conntypes.QueryConnectionResponse, error) {
	res, err := r.queryConnectionABCI(ctx, chain, height, connectionid)
	if err != nil && strings.Contains(err.Error(), "not found") {
		return &conntypes.QueryConnectionResponse{
			Connection: &conntypes.ConnectionEnd{
				ClientId: "client",
				Versions: []*conntypes.Version{},
				State:    conntypes.UNINITIALIZED,
				Counterparty: conntypes.Counterparty{
					ClientId:     "client",
					ConnectionId: "connection",
					Prefix:       commitmenttypes.MerklePrefix{KeyPrefix: []byte{}},
				},
				DelayPeriod: 0,
			},
			Proof:       []byte{},
			ProofHeight: clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
		}, nil
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Relayer) queryConnectionABCI(ctx context.Context, chain ChainConfig, height int64, connectionID string) (*conntypes.QueryConnectionResponse, error) {
	key := host.ConnectionKey(connectionID)

	value, proofBz, proofHeight, err := r.QueryTendermintProof(ctx, chain, height, key)
	if err != nil {
		return nil, err
	}

	// check if connection exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(conntypes.ErrConnectionNotFound, connectionID)
	}

	cdc := codec.NewProtoCodec(r.cdc.InterfaceRegistry())

	var connection conntypes.ConnectionEnd
	if err := cdc.Unmarshal(value, &connection); err != nil {
		return nil, err
	}

	return &conntypes.QueryConnectionResponse{
		Connection:  &connection,
		Proof:       proofBz,
		ProofHeight: proofHeight,
	}, nil
}

func (r *Relayer) ConnOpenAck() (string, error) {
	clientCtx := r.createClientCtx(r.config.SrcChain)
	txf := r.createTxFactory(clientCtx, r.config.SrcChain)

	latestHeight, err := r.getLatestHeight(r.config.DstChain)
	if err != nil {
		return "", err
	}

	clientState, clientStateProof, consensusProof, connectionProof, connectionProofHeight, err := r.GenerateConnHandshakeProof(r.ctx, r.config.DstChain, latestHeight, r.config.DstChain.ClientId, r.config.DstChain.ConnectionId)

	initMsg := connectiontypes.NewMsgConnectionOpenAck(
		r.config.DstChain.ConnectionId,
		r.config.SrcChain.ConnectionId,
		clientState,
		connectionProof,
		clientStateProof,
		consensusProof,
		connectionProofHeight.(clienttypes.Height),
		clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: uint64(latestHeight),
		},
		nil,
		clientCtx.From,
	)

	txResp, err := r.sendTx(clientCtx, txf, initMsg)
	if err != nil {
		return "", err
	}

	connectionID, err := parseConnectionIDFromEvents(txResp.Events)
	if err != nil {
		return "", err
	}

	return connectionID, nil
}