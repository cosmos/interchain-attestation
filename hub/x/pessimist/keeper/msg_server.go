package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"hub/x/pessimist/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// TODO: Test this
func (k msgServer) CreateValidationObjective(goCtx context.Context, msg *types.MsgCreateValidationObjective) (*types.MsgCreateValidationObjectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CreateNewValidationObjective(ctx, msg.ClientId, msg.RequiredPower); err != nil {
		return nil, err
	}

	return &types.MsgCreateValidationObjectiveResponse{}, nil
}

// TODO: Test this
func (k msgServer) SignUpForObjective(goCtx context.Context, msg *types.MsgSignUpForObjective) (*types.MsgSignUpForObjectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to decode validator address %s", msg.ValidatorAddress)
	}

	val, err := k.stakingKeeper.Validator(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	if val.GetStatus() != stakingtypes.Bonded {
		return nil, types.ErrValidatorNotBonded
	}

	pubKey, err := val.ConsPubKey()
	if err != nil {
		return nil, err
	}

	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	if err := k.AddValidatorToObjective(ctx, msg.ClientId, &types.Validator{
		ValidatorAddr: msg.ValidatorAddress,
		PubKey:        pkAny,
	}); err != nil {
		return nil, err
	}

	return &types.MsgSignUpForObjectiveResponse{}, nil
}
