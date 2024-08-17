package keeper

import (
	"context"
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the configmodule MsgServer interface
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

func (m msgServer) RegisterAttestator(ctx context.Context, msg *types.MsgRegisterAttestator) (*types.MsgRegisterAttestatorResponse, error) {
	valAddr, err := m.validatorAddressCodec.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if err := msg.Validate(m.validatorAddressCodec); err != nil {
		return nil, err
	}

	validator, err := m.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	attestator := types.Attestator{
		AttestatorId:      msg.AttestatorId,
		PublicKey:       msg.AttestationPublicKey,
		ConsensusPubkey: validator.ConsensusPubkey,
	}

	if err := m.Keeper.SetNewAttestator(ctx, attestator); err != nil {
		return nil, err
	}

	return &types.MsgRegisterAttestatorResponse{}, nil
}