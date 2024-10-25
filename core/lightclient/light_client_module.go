package lightclient

import (
	"bytes"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	v2 "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types/v2"
	host "github.com/cosmos/ibc-go/v9/modules/core/24-host"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
)

var _ exported.LightClientModule = (*LightClientModule)(nil)

type TrustedClientUpdateFunc func(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height

// LightClientModule implements the core IBC api.LightClientModule interface.
type LightClientModule struct {
	cdc                codec.BinaryCodec
	storeProvider      clienttypes.StoreProvider
	attestatorsHandler AttestatorsController
}

// NewLightClientModule creates and returns a new attestation LightClientModule.
func NewLightClientModule(cdc codec.BinaryCodec, storeProvider clienttypes.StoreProvider, attestatorsHandler AttestatorsController) (LightClientModule, TrustedClientUpdateFunc) {
	lcm := LightClientModule{
		cdc:                cdc,
		storeProvider:      storeProvider,
		attestatorsHandler: attestatorsHandler,
	}
	return lcm, lcm.trustedUpdateState
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

	return clientState.Initialize(l.cdc, clientStore, &consensusState)
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

// UpdateState will always return an error, because it is only supposed to be called by validators in the trustedUpdateState call
func (l *LightClientModule) UpdateState(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	panic(ErrInvalidUpdateMethod)
}

// TODO: Document
func (l *LightClientModule) trustedUpdateState(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		panic(errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID))
	}

	return clientState.UpdateState(ctx, l.cdc, clientStore, clientMsg)
}

// VerifyMembership uses the packet commitment bytes (value) to verify the membership proof.
// The client module has all the packet commitments stored and will just look for their existence.
func (l *LightClientModule) VerifyMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path, value []byte) error {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	merklePath, ok := path.(v2.MerklePath)
	if ok {
		// TODO deal with other stores
		if bytes.Equal(merklePath.KeyPath[0], []byte("ibc")) {
			split := strings.Split(string(merklePath.KeyPath[1]), "/")
			if split[0] == host.KeyConnectionPrefix ||
				split[0] == host.KeyChannelEndPrefix {
				// TODO: Verify membership using merkleroot
				// For now, to get things moving, we just return true here :O
				ctx.Logger().Info("we are not verifying merkle root for connection/channel keys yet, this might be dangerous, we just accept it", "key path", merklePath.KeyPath)
				return nil
			} else {
				ctx.Logger().Info("key path is not supported for merkle root verification, so it better be a packet commitment!", "key path", merklePath.KeyPath)
			}
		}
	}

	return clientState.VerifyMembership(clientStore, value)
}

// VerifyNonMembership is currently not possible as we need a packet commitment to check if it exists in the store
func (l *LightClientModule) VerifyNonMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path) error {
	return fmt.Errorf("not implemented")
}

// Status returns the status of the light client with the given clientID.
// TODO: Implement frozen client
func (l *LightClientModule) Status(ctx sdk.Context, clientID string) exported.Status {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	_, found := getClientState(clientStore, l.cdc)
	if !found {
		return exported.Unknown
	}

	return exported.Active
}

// LatestHeight returns the latest height of the light client with the given clientID.
func (l *LightClientModule) LatestHeight(ctx sdk.Context, clientID string) exported.Height {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return clienttypes.ZeroHeight()
	}

	return clientState.LatestHeight
}

// TimestampAtHeight returns the timestamp associated with the given height.
func (l *LightClientModule) TimestampAtHeight(ctx sdk.Context, clientID string, height exported.Height) (uint64, error) {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return 0, errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	return clientState.getTimestampAtHeight(clientStore, l.cdc, height)
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
