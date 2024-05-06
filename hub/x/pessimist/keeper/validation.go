package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"hub/x/pessimist/types"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateNewValidationObjective(ctx context.Context, clientIDToValidate string, requiredPower uint64) error {
	objective := &types.ValidationObjective{
		ClientIdToValidate: clientIDToValidate,
		RequiredPower:      requiredPower,
		Validators:         nil,
		Activated:          false,
		ClientIdToNotify:   "",
	}
	// TODO: Validate it

	if k.GetClientKeeper() == nil {
		panic("client keeper is nil!!!")
	}
	clientStatus := k.GetClientKeeper().GetClientStatus(ctx.(sdk.Context), clientIDToValidate)
	if clientStatus != exported.Active {
		return errorsmod.Wrapf(types.ErrClientNotActive, "client %s is not active: %s", clientIDToValidate, clientStatus)
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ValidatorObjectiveKeyPrefix))
	store.Set(types.ValidatorObjectiveKey(clientIDToValidate), k.cdc.MustMarshal(objective))

	return nil
}

func (k Keeper) GetValidationObjective(ctx context.Context, clientID string) (types.ValidationObjective, bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ValidatorObjectiveKeyPrefix))
	validatorObjective := store.Get(types.ValidatorObjectiveKey(clientID))
	if validatorObjective == nil {
		return types.ValidationObjective{}, false
	}

	var objective types.ValidationObjective
	k.cdc.MustUnmarshal(validatorObjective, &objective)
	return objective, true
}

func (k Keeper) GetValidatorPower(ctx context.Context, validators []*types.Validator) math.Int {
	totalPower := math.ZeroInt()
	for _, v := range validators {
		addr, err := sdk.ValAddressFromBech32(v.ValidatorAddr)
		if err != nil {
			panic(err)
		}
		validator, err := k.stakingKeeper.Validator(ctx, addr)
		if err != nil {
			panic(err)
		}

		totalPower = totalPower.Add(validator.GetBondedTokens())
	}

	return totalPower
}

func (k Keeper) AddValidatorToObjective(ctx context.Context, clientID string, validator *types.Validator) error {
	objective, ok := k.GetValidationObjective(ctx, clientID)
	if !ok {
		return errorsmod.Wrapf(types.ErrObjectiveNotFound, "objective not found for client %s", clientID)
	}

	objective.Validators = append(objective.Validators, validator)
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ValidatorObjectiveKeyPrefix))
	store.Set(types.ValidatorObjectiveKey(clientID), k.cdc.MustMarshal(&objective))

	if objective.Activated {
		return nil
	}

	// Check if the required power is reached
	totalPower := k.GetValidatorPower(ctx, objective.Validators)
	if totalPower.GTE(math.NewInt(int64(objective.RequiredPower))) {
		objective.Activated = true
		store.Set(types.ValidatorObjectiveKey(clientID), k.cdc.MustMarshal(&objective))

		clientState := types.ClientState{
			DependentClientId: objective.ClientIdToValidate,
			LatestHeight:      types.Height{},
		}
		clientStateBz := k.cdc.MustMarshal(&clientState)
		newClientID, err := k.GetClientKeeper().CreateClient(ctx.(sdk.Context), types.ClientType, clientStateBz, nil)
		if err != nil {
			return err
		}

		objective.ClientIdToNotify = newClientID
	}

	return nil
}
