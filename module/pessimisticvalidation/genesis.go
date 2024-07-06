package pessimisticvalidation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/keeper"
	"github.com/gjermundgaraba/pessimistic-validation/pessimisticvalidation/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := k.Params.Set(ctx, *data.Params); err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}
	
	return &types.GenesisState{
		Params: &params,
	}
}
