package configmodule

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/interchain-attestation/configmodule/keeper"
	"github.com/cosmos/interchain-attestation/configmodule/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := k.Params.Set(ctx, *data.Params); err != nil {
		panic(err)
	}
	for _, attestator := range data.Attestators {
		if err := k.Attestators.Set(ctx, attestator.AttestatorId, attestator); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	iterator, err := k.Attestators.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}
	attestators, err := iterator.Values()
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params:      &params,
		Attestators: attestators,
	}
}
