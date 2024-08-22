package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/gjermundgaraba/interchain-attestation/configmodule/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type queryServer struct{ k Keeper }

var _ types.QueryServer = queryServer{}


func NewQueryServer(k Keeper) types.QueryServer {
	return queryServer{k: k}
}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	params, err := q.k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Attestators(ctx context.Context, req *types.QueryAttestatorsRequest) (*types.QueryAttestatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	attestators, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Attestators,
		req.Pagination,
		func(_ []byte, value types.Attestator) (types.Attestator, error) {
			return value, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryAttestatorsResponse{
		Attestators: attestators,
		Pagination: pageRes,
	}, nil
}

func (q queryServer) Attestator(ctx context.Context, req *types.QueryAttestatorRequest) (*types.QueryAttestatorResponse, error) {
	if req == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("empty request")
	}

	attestator, err := q.k.Attestators.Get(ctx, req.AttestatorId)
	if err != nil {
		return nil, sdkerrors.ErrNotFound
	}

	return &types.QueryAttestatorResponse{Attestator: attestator}, nil
}