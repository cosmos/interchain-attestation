package keeper

import (
	"context"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"hub/x/pessimist/types"
)

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
	var totalPower math.Int
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
