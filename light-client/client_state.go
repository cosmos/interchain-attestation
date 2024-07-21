package lightclient

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"strings"
)

var _ exported.ClientState = (*ClientState)(nil)

func NewClientState(
	chainID string,
	requiredTokenPower sdkmath.Int,
	frozenHeight clienttypes.Height,
	latestHeight clienttypes.Height,
) *ClientState {
	return &ClientState{
		ChainId:            chainID,
		RequiredTokenPower: requiredTokenPower,
		FrozenHeight:       frozenHeight,
		LatestHeight:       latestHeight,
	}
}

func (cs *ClientState) ClientType() string {
	return ModuleName
}

func (cs *ClientState) Validate() error {
	if strings.TrimSpace(cs.ChainId) == "" {
		return errorsmod.Wrap(ErrInvalidChainID, "chain id cannot be empty")
	}

	if cs.RequiredTokenPower.LTE(sdkmath.ZeroInt()) {
		return errorsmod.Wrap(ErrInvalidRequiredTokenPower, "required token power must be more than zero")
	}

	// the latest height revision number must match the chain id revision number
	if cs.LatestHeight.RevisionNumber != clienttypes.ParseChainID(cs.ChainId) {
		return errorsmod.Wrapf(ErrInvalidHeaderHeight,
			"latest height revision number must match chain id revision number (%d != %d)", cs.LatestHeight.RevisionNumber, clienttypes.ParseChainID(cs.ChainId))
	}
	if cs.LatestHeight.RevisionHeight == 0 {
		return errorsmod.Wrapf(ErrInvalidHeaderHeight, "client's latest height revision height cannot be zero")
	}

	return nil
}

// Initialize checks that the initial consensus state is an 07-tendermint consensus state and
// sets the client state, consensus state and associated metadata in the provided client store.
func (cs ClientState) Initialize(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, consState exported.ConsensusState) error {
	consensusState, ok := consState.(*ConsensusState)
	if !ok {
		return errorsmod.Wrapf(clienttypes.ErrInvalidConsensus, "invalid initial consensus state. expected type: %T, got: %T",
			&ConsensusState{}, consState)
	}

	setClientState(clientStore, cdc, &cs)
	setConsensusState(clientStore, cdc, consensusState, cs.LatestHeight)
	//setConsensusMetadata(ctx, clientStore, cs.LatestHeight)

	return nil
}

