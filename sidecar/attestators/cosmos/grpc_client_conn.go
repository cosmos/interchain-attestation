package cosmos

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	gogogrpc "github.com/cosmos/gogoproto/grpc"
	sltypes "github.com/strangelove-ventures/cometbft-client/abci/types"
	clientwrapper "github.com/strangelove-ventures/cometbft-client/client"
	slclient "github.com/strangelove-ventures/cometbft-client/rpc/client"
	slcoretypes "github.com/strangelove-ventures/cometbft-client/rpc/core/types"
	"gitlab.com/tozd/go/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/rootmulti"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	legacyerrors "github.com/cosmos/cosmos-sdk/types/errors"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
)

// The code below is borrowed from the Go Realyer and modified a bit

type ClientConn struct {
	cometClient *clientwrapper.Client
	codec       CodecConfig
}

var _ gogogrpc.ClientConn = &ClientConn{}

var protoCodec = encoding.GetCodec(proto.Name)

func (c *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	if reflect.ValueOf(args).IsNil() {
		return errors.New("request cannot be nil")
	}

	inMd, _ := metadata.FromOutgoingContext(ctx)
	abciRes, outMd, err := c.RunGRPCQuery(ctx, method, args, inMd)
	if err != nil {
		return err
	}

	if err = protoCodec.Unmarshal(abciRes.Value, reply); err != nil {
		return err
	}

	for _, callOpt := range opts {
		header, ok := callOpt.(grpc.HeaderCallOption)
		if !ok {
			continue
		}

		*header.HeaderAddr = outMd
	}

	if c.codec.InterfaceRegistry != nil {
		return codectypes.UnpackInterfaces(reply, c.codec.Marshaler)
	}

	return nil
}

// RunGRPCQuery runs a gRPC query from the clientCtx, given all necessary
// arguments for the gRPC method, and returns the ABCI response. It is used
// to factorize code between client (Invoke) and server (RegisterGRPCServer)
// gRPC handlers.
func (c *ClientConn) RunGRPCQuery(ctx context.Context, method string, req interface{}, md metadata.MD) (abci.ResponseQuery, metadata.MD, error) {
	reqBz, err := protoCodec.Marshal(req)
	if err != nil {
		return abci.ResponseQuery{}, nil, err
	}

	// parse height header
	if heights := md.Get(grpctypes.GRPCBlockHeightHeader); len(heights) > 0 {
		height, err := strconv.ParseInt(heights[0], 10, 64)
		if err != nil {
			return abci.ResponseQuery{}, nil, err
		}
		if height < 0 {
			return abci.ResponseQuery{}, nil, errors.Errorf("client.Context.Invoke: height (%d) from %q must be >= 0", height, grpctypes.GRPCBlockHeightHeader)
		}

	}

	height, err := GetHeightFromMetadata(md)
	if err != nil {
		return abci.ResponseQuery{}, nil, err
	}

	prove, err := GetProveFromMetadata(md)
	if err != nil {
		return abci.ResponseQuery{}, nil, err
	}

	abciReq := abci.RequestQuery{
		Path:   method,
		Data:   reqBz,
		Height: height,
		Prove:  prove,
	}

	abciRes, err := c.QueryABCI(ctx, abciReq)
	if err != nil {
		return abci.ResponseQuery{}, nil, err
	}

	// Create header metadata. For now the headers contain:
	// - block height
	// We then parse all the call options, if the call option is a
	// HeaderCallOption, then we manually set the value of that header to the
	// metadata.
	md = metadata.Pairs(grpctypes.GRPCBlockHeightHeader, strconv.FormatInt(abciRes.Height, 10))

	return abciRes, md, nil
}

// QueryABCI performs an ABCI query and returns the appropriate response and error sdk error code.
func (c *ClientConn) QueryABCI(ctx context.Context, req abci.RequestQuery) (abci.ResponseQuery, error) {
	opts := slclient.ABCIQueryOptions{
		Height: req.Height,
		Prove:  req.Prove,
	}

	slRes, err := c.cometClient.ABCIQueryWithOptions(ctx, req.Path, req.Data, opts)
	if err != nil {
		return abci.ResponseQuery{}, err
	}

	res := convertResultABCIQuery(slRes)

	if !res.Response.IsOK() {
		return abci.ResponseQuery{}, sdkErrorToGRPCError(res.Response)
	}

	// data from trusted node or subspace query doesn't need verification
	if !opts.Prove || !isQueryStoreWithProof(req.Path) {
		return res.Response, nil
	}

	return res.Response, nil
}

func GetHeightFromMetadata(md metadata.MD) (int64, error) {
	height := md.Get(grpctypes.GRPCBlockHeightHeader)
	if len(height) == 1 {
		return strconv.ParseInt(height[0], 10, 64)
	}
	return 0, nil
}

func GetProveFromMetadata(md metadata.MD) (bool, error) {
	prove := md.Get("x-cosmos-query-prove")
	if len(prove) == 1 {
		return strconv.ParseBool(prove[0])
	}
	return false, nil
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

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType must be "store" and subpath must be "key" to require a proof.
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

func convertResultABCIQuery(res *slcoretypes.ResultABCIQuery) *coretypes.ResultABCIQuery {
	var ops *crypto.ProofOps
	if res.Response.ProofOps != nil {
		ops = convertProofOps(res.Response.ProofOps)
	}

	return &coretypes.ResultABCIQuery{
		Response: abci.ResponseQuery{
			Code:      res.Response.Code,
			Log:       res.Response.Log,
			Info:      res.Response.Info,
			Index:     res.Response.Index,
			Key:       res.Response.Key,
			Value:     res.Response.Value,
			ProofOps:  ops,
			Height:    res.Response.Height,
			Codespace: res.Response.Codespace,
		},
	}
}

func convertProofOps(proofOps *sltypes.ProofOps) *crypto.ProofOps {
	ops := make([]crypto.ProofOp, len(proofOps.Ops))
	for i, op := range proofOps.Ops {
		ops[i] = crypto.ProofOp{
			Type: op.Type,
			Key:  op.Key,
			Data: op.Data,
		}
	}

	return &crypto.ProofOps{Ops: ops}
}

func (c *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// TODO implement me
	panic("implement me")
}
