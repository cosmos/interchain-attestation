package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"hub/x/pessimist/types"
	"time"
)

func (k Keeper) CreateNewValidationObjective(ctx sdk.Context, clientIDToValidate string, requiredPower uint64) error {
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
	dependentClientModule, found := k.GetClientKeeper().Route(clientIDToValidate)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrInvalidClientType, "dependent client not found")
	}
	clientStatus := dependentClientModule.Status(ctx, clientIDToValidate)
	if clientStatus != exported.Active {
		return errorsmod.Wrapf(types.ErrClientNotActive, "client %s is not active: %s", clientIDToValidate, clientStatus)
	}

	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ValidatorObjectiveKeyPrefix))
	store.Set(types.ValidatorObjectiveKey(clientIDToValidate), k.cdc.MustMarshal(objective))

	return nil
}

func (k Keeper) GetAllValidationObjectives(ctx sdk.Context) []types.ValidationObjective {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.ValidatorObjectiveKeyPrefix))
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var objectives []types.ValidationObjective
	for ; iterator.Valid(); iterator.Next() {
		var objective types.ValidationObjective
		k.cdc.MustUnmarshal(iterator.Value(), &objective)
		objectives = append(objectives, objective)
	}

	return objectives
}

func (k Keeper) GetValidationObjective(ctx sdk.Context, clientID string) (types.ValidationObjective, bool) {
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

func (k Keeper) GetValidatorForObjective(ctx sdk.Context, validatorAddress string, validationObjective types.ValidationObjective) (types.ValidatorPower, bool, error) {
	for _, v := range validationObjective.Validators {
		if v.ValidatorAddr == validatorAddress || v.ConsAddr == validatorAddress {
			addr, err := sdk.ValAddressFromBech32(v.ValidatorAddr)
			if err != nil {
				return types.ValidatorPower{}, false, err
			}
			validator, err := k.stakingKeeper.Validator(ctx, addr)
			if err != nil {
				return types.ValidatorPower{}, false, err
			}

			consPubKey, err := validator.ConsPubKey()
			if err != nil {
				return types.ValidatorPower{}, false, err
			}
			var pkAny *codectypes.Any
			if pkAny, err = codectypes.NewAnyWithValue(consPubKey); err != nil {
				return types.ValidatorPower{}, false, err
			}
			return types.ValidatorPower{
				Validator: types.Validator{
					ValidatorAddr: validatorAddress,
					PubKey:       pkAny,
				},
				Power:     validator.GetBondedTokens().Uint64(),
			}, true, nil
		}
	}

	return types.ValidatorPower{}, false, nil
}

func (k Keeper) GetValidatorPower(ctx sdk.Context, validators []*types.Validator) math.Int {
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

func (k Keeper) AddValidatorToObjective(ctx sdk.Context, clientID string, validator *types.Validator) error {
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

		clientState := types.ClientState{
			DependentClientId: objective.ClientIdToValidate,
			LatestHeight:      0,
		}
		clientStateBz := k.cdc.MustMarshal(&clientState)

		consensusState := tmclient.ConsensusState{
			Timestamp:          time.Time{},
			Root:               commitmenttypes.MerkleRoot{},
			NextValidatorsHash: nil,
		}
		consensusStateBz := k.cdc.MustMarshal(&consensusState)
		newClientID, err := k.GetClientKeeper().CreateClient(ctx, types.ClientType, clientStateBz, consensusStateBz)
		if err != nil {
			return err
		}

		objective.ClientIdToNotify = newClientID
		store.Set(types.ValidatorObjectiveKey(clientID), k.cdc.MustMarshal(&objective))
	}

	return nil
}
