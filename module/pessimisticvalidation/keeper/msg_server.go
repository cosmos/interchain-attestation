package keeper

import (
	"context"
	"cosmossdk.io/errors"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the pessimisticvalidation MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// Check if the provided authority is the same as the keeper's authority
	if m.authority != msg.Authority {
		return nil,  errors.Wrapf(types.ErrUnauthorized, "invalid authority; expected %s, got %s", m.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := m.Params.Set(ctx, msg.Params); err != nil {
		return nil, errors.Wrapf(err, "failed to update params")
	}

	return &types.MsgUpdateParamsResponse{}, nil
}