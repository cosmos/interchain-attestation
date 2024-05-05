package keeper

import (
	"context"

	"hub/x/pessimist/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: Test this
func (k msgServer) CreateValidationObjective(goCtx context.Context, msg *types.MsgCreateValidationObjective) (*types.MsgCreateValidationObjectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Check that the client exists

	k.CreateNewValidationObjective(ctx, msg.ClientId, msg.RequiredPower)

	return &types.MsgCreateValidationObjectiveResponse{}, nil
}
