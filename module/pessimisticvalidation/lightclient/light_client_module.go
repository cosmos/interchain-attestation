package lightclient

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.LightClientModule = (*LightClientModule)(nil)

type LightClientModule struct {
	cdc           codec.BinaryCodec
	storeProvider exported.ClientStoreProvider
}

func NewLightClientModule(cdc codec.BinaryCodec) LightClientModule {
	return LightClientModule{
		cdc:    cdc,
	}
}

func (l *LightClientModule) RegisterStoreProvider(storeProvider exported.ClientStoreProvider) {
	l.storeProvider = storeProvider
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) Initialize(ctx sdk.Context, clientID string, clientStateBz, consensusStateBz []byte) error {
	return nil
}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) VerifyClientMessage(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) error {
	return nil
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
func (l *LightClientModule) UpdateStateOnMisbehaviour(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) {}

// TODO: implement this
// TODO: test this
// TODO: godoc
func (l *LightClientModule) UpdateState(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	return []exported.Height{
		clienttypes.NewHeight(0, 0),
	}
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
