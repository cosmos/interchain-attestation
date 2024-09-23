package voteextension

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v9/modules/core/exported"
)

type ClientKeeper interface {
	UpdateClient(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) error
}
