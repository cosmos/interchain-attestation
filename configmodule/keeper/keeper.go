package keeper

import (
	"cosmossdk.io/collections"
	addresscodec "cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/types"
)

type Keeper struct {
	storeService          store.KVStoreService
	cdc                   codec.BinaryCodec
	validatorAddressCodec addresscodec.Codec

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string

	stakingKeeper types.StakingKeeper

	Schema    collections.Schema
	Params    collections.Item[types.Params]
	Attestators collections.Map[[]byte, types.Attestator]
}

func NewKeeper(
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	validatorAddressCodec addresscodec.Codec,
	authority string,
	stakingKeeper types.StakingKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:          storeService,
		cdc:                   cdc,
		validatorAddressCodec: validatorAddressCodec,
		authority:             authority,
		stakingKeeper:         stakingKeeper,
		Params:                collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Attestators:             collections.NewMap(sb, types.AttestatorsKey, "attestators", collections.BytesKey, codec.CollValue[types.Attestator](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the configmodule module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
