package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"rolly/x/pessimist/types"
)

var _ clienttypes.ConsensusHost = (*CustomConsensusHost)(nil)

type CustomConsensusHost struct {
	TMConsensusHost clienttypes.ConsensusHost
}

func (c CustomConsensusHost) GetSelfConsensusState(ctx sdk.Context, height exported.Height) (exported.ConsensusState, error) {
	return c.TMConsensusHost.GetSelfConsensusState(ctx, height)
}

func (c CustomConsensusHost) ValidateSelfClient(ctx sdk.Context, clientState exported.ClientState) error {
	switch clientState.(type) {
	case *tmclient.ClientState:
		return c.TMConsensusHost.ValidateSelfClient(ctx, clientState)
	case *types.ClientState:
		return clientState.Validate()
	default:
		return clienttypes.ErrInvalidClientType
	}
}
