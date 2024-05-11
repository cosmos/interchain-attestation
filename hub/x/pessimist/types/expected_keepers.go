package types

import (
	"context"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}

type StakingKeeper interface {
	Validator(context.Context, sdk.ValAddress) (stakingtypes.ValidatorI, error)
}

type ClientKeeper interface {
	GetClientStatus(ctx sdk.Context, clientID string) exported.Status
	CreateClient(ctx sdk.Context, clientType string, clientState, consensusState []byte) (string, error)
	GetClientLatestHeight(ctx sdk.Context, clientID string) clienttypes.Height
	Route(clientID string) (exported.LightClientModule, bool)
}
