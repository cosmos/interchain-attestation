package keeper

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = queryServer{}

type queryServer struct{ k Keeper }

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