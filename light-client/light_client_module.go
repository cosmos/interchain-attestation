package lightclient

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.LightClientModule = (*LightClientModule)(nil)

// LightClientModule implements the core IBC api.LightClientModule interface.
type LightClientModule struct {
	cdc                codec.BinaryCodec
	storeProvider      exported.ClientStoreProvider
	attestatorsHandler AttestatorsHandler
}

// NewLightClientModule creates and returns a new pessimistic LightClientModule.
func NewLightClientModule(cdc codec.BinaryCodec, attestatorsHandler AttestatorsHandler) LightClientModule {
	return LightClientModule{
		cdc:                cdc,
		attestatorsHandler: attestatorsHandler,
	}
}

// RegisterStoreProvider is called by core IBC when a LightClientModule is added to the router.
// It allows the LightClientModule to set a ClientStoreProvider which supplies isolated prefix client stores
// to IBC light client instances.
func (l *LightClientModule) RegisterStoreProvider(storeProvider exported.ClientStoreProvider) {
	l.storeProvider = storeProvider
}

// Initialize unmarshals the provided client and consensus states and performs basic validation.
// It then initializes the client state with the provided consensus state (and stores it in the client store).A
//
// CONTRACT: clientID is validated in 02-client router, thus clientID is assumed here to have the correct format.
func (l *LightClientModule) Initialize(ctx sdk.Context, clientID string, clientStateBz, consensusStateBz []byte) error {
	var clientState ClientState
	if err := l.cdc.Unmarshal(clientStateBz, &clientState); err != nil {
		return err
	}
	if err := clientState.Validate(); err != nil {
		return err
	}

	var consensusState ConsensusState
	if err := l.cdc.Unmarshal(consensusStateBz, &consensusState); err != nil {
		return err
	}
	if err := consensusState.ValidateBasic(); err != nil {
		return err
	}

	clientStore := l.storeProvider.ClientStore(ctx, clientID)

	return clientState.Initialize(ctx, l.cdc, clientStore, &consensusState)
}

// VerifyClientMessage obtains the client state associated with the client identifier and calls into the clientState.VerifyClientMessage method.
//
// CONTRACT: clientID is validated in 02-client router, thus clientID is assumed here to have the format 07-tendermint-{n}.
func (l *LightClientModule) VerifyClientMessage(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) error {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	return clientState.VerifyClientMessage(ctx, l.cdc, l.attestatorsHandler, clientMsg)
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) CheckForMisbehaviour(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) bool {
	return false
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) UpdateStateOnMisbehaviour(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) {
}

// UpdateState obtains the client state associated with the client identifier and calls into the clientState.UpdateState method.
func (l *LightClientModule) UpdateState(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		panic(errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID))
	}

	return clientState.UpdateState(ctx, l.cdc, clientStore, clientMsg)
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) VerifyMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path, value []byte) error {
	return nil
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) VerifyNonMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path) error {
	return nil
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) Status(ctx sdk.Context, clientID string) exported.Status {
	return exported.Active
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) LatestHeight(ctx sdk.Context, clientID string) exported.Height {
	return clienttypes.NewHeight(0, 0)
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) TimestampAtHeight(ctx sdk.Context, clientID string, height exported.Height) (uint64, error) {
	return 0, nil
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) RecoverClient(ctx sdk.Context, clientID, substituteClientID string) error {
	return fmt.Errorf("not implemented")
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) VerifyUpgradeAndUpdateState(ctx sdk.Context, clientID string, newClient []byte, newConsState []byte, upgradeClientProof, upgradeConsensusStateProof []byte) error {
	return fmt.Errorf("not implemented")
}
