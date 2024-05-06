package keeper

import (
	"fmt"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"hub/x/pessimist/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		stakingKeeper  types.StakingKeeper
		getIBCKeeperFn func() *ibckeeper.Keeper

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	stakingKeeper types.StakingKeeper,
	getIBCKeeperFn func() *ibckeeper.Keeper,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	if storeService == nil {
		panic("store service cannot be nil")
	}

	if logger == nil {
		panic("logger cannot be nil")
	}

	if stakingKeeper == nil {
		panic("staking keeper cannot be nil")
	}

	if authority == "" {
		panic("authority address cannot be empty")
	}

	return Keeper{
		cdc:            cdc,
		storeService:   storeService,
		logger:         logger,
		stakingKeeper:  stakingKeeper,
		getIBCKeeperFn: getIBCKeeperFn,
		authority:      authority,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetClientKeeper() types.ClientKeeper {
	return k.getIBCKeeperFn().ClientKeeper
}
