package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"hub/x/pessimist/types"
)

func (k Keeper) ValidationObjective(goCtx context.Context, req *types.QueryValidationObjectiveRequest) (*types.QueryValidationObjectiveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	objective, found := k.GetValidationObjective(ctx, req.ClientId)
	if !found {
		return nil, status.Error(codes.NotFound, "objective not found")
	}

	return &types.QueryValidationObjectiveResponse{
		ValidationObjective: &objective,
	}, nil
}
