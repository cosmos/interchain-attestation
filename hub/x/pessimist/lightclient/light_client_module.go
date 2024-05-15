package lightclient

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"hub/x/pessimist/keeper"
	"hub/x/pessimist/types"
)

var _ exported.LightClientModule = (*LightClientModule)(nil)

type ValidationVoteWithCommitment struct {
	ValidationVote types.ValidationVote
	Commitment types.Commitment
	Validator types.ValidatorPower
}

type LightClientModule struct {
	cdc           codec.BinaryCodec
	storeProvider exported.ClientStoreProvider
	keeper        keeper.Keeper
}

func NewLightClientModule(cdc codec.BinaryCodec, keeper keeper.Keeper) LightClientModule {
	return LightClientModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

func (l *LightClientModule) RegisterStoreProvider(storeProvider exported.ClientStoreProvider) {
	l.storeProvider = storeProvider
}

// TODO: Test this
func (l *LightClientModule) Initialize(ctx sdk.Context, clientID string, clientStateBz, consensusStateBz []byte) error {
	var clientState types.ClientState
	if err := l.cdc.Unmarshal(clientStateBz, &clientState); err != nil {
		return err
	}
	if err := clientState.Validate(); err != nil {
		return err
	}

	var consensusState tmclient.ConsensusState
	if err := l.cdc.Unmarshal(consensusStateBz, &consensusState); err != nil {
		return err
	}

	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	setClientState(clientStore, l.cdc, &clientState)
	setConsensusState(clientStore, l.cdc, &consensusState, clienttypes.NewHeight(0, uint64(clientState.LatestHeight)))

	return nil
}

// TODO: Test this
func (l *LightClientModule) VerifyClientMessage(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) error {
	if _, _, err := l.processClientMessage(ctx, clientID, clientMsg); err != nil {
		return err
	}

	return nil
}

func (l *LightClientModule) CheckForMisbehaviour(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) bool {
	//TODO implement me
	return false
}

func (l *LightClientModule) UpdateStateOnMisbehaviour(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) {
	//TODO implement me
}

func (l *LightClientModule) UpdateState(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		panic("client state not found") // Should not happen
	}

	height, consensusState, err := l.processClientMessage(ctx, clientID, clientMsg)
	if err != nil {
		panic(err) // Should not happen
	}
	ibcHeight := &clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: uint64(height),
	}

	ctx.Logger().Info("Updating client state", "clientID", clientID, "height", height)
	clientState.LatestHeight = height
	setClientState(clientStore, l.cdc, &clientState)
	setConsensusState(clientStore, l.cdc, &consensusState, ibcHeight)

	return []exported.Height{ibcHeight}
}

func (l *LightClientModule) VerifyMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path, value []byte) error {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	if clientState.LatestHeight < int64(height.GetRevisionHeight()) {
		return errorsmod.Wrap(clienttypes.ErrInvalidHeight, "client state height is less than the height of the proof")
	}

	dependentClientModule, found := l.keeper.GetClientKeeper().Route(clientState.DependentClientId)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrInvalidClientType, "dependent client not found")
	}

	return dependentClientModule.VerifyMembership(ctx, clientState.DependentClientId, height, delayTimePeriod, delayBlockPeriod, proof, path, value)
}

func (l *LightClientModule) VerifyNonMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path) error {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	if clientState.LatestHeight < int64(height.GetRevisionHeight()) {
		return errorsmod.Wrap(clienttypes.ErrInvalidHeight, "client state height is less than the height of the proof")
	}

	dependentClientModule, found := l.keeper.GetClientKeeper().Route(clientState.DependentClientId)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrInvalidClientType, "dependent client not found")
	}

	return dependentClientModule.VerifyNonMembership(ctx, clientState.DependentClientId, height, delayTimePeriod, delayBlockPeriod, proof, path)
}

func (l *LightClientModule) Status(ctx sdk.Context, clientID string) exported.Status {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		ctx.Logger().Error("client state not found", "clientID", clientID)
		return exported.Unknown
	}
	dependentClientModule, found := l.keeper.GetClientKeeper().Route(clientState.DependentClientId)
	if !found {
		ctx.Logger().Error("dependent client not found", "clientID", clientID)
		return exported.Unknown
	}

	return dependentClientModule.Status(ctx, clientState.DependentClientId)
}

func (l *LightClientModule) LatestHeight(ctx sdk.Context, clientID string) exported.Height {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		panic("client state not found") // Should not happen
	}

	return clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: uint64(clientState.LatestHeight),
	}
}

func (l *LightClientModule) TimestampAtHeight(ctx sdk.Context, clientID string, height exported.Height) (uint64, error) {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	consensusState, found := getConsensusState(clientStore, l.cdc, height)
	if !found {
		return 0, errorsmod.Wrap(clienttypes.ErrConsensusStateNotFound, "consensus state not found")
	}

	return consensusState.GetTimestamp(), nil
}

func (l *LightClientModule) RecoverClient(ctx sdk.Context, clientID, substituteClientID string) error {
	return fmt.Errorf("not implemented")
}

func (l *LightClientModule) VerifyUpgradeAndUpdateState(ctx sdk.Context, clientID string, newClient []byte, newConsState []byte, upgradeClientProof, upgradeConsensusStateProof []byte) error {
	return fmt.Errorf("not implemented")
}
