package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

const ClientType = "69-pessimist"

var _ exported.ClientState = (*ClientState)(nil)

func (m *ClientState) ClientType() string {
	return ClientType
}

func (m *ClientState) Validate() error {
	if m.DependentClientId == "" {
		return errorsmod.Wrap(clienttypes.ErrInvalidClient, "dependent client id cannot be empty")
	}

	return nil
}

func (m *ClientState) GetLatestHeight() exported.Height {
	return clienttypes.NewHeight(0, uint64(m.LatestHeight))
}

func (m *ClientState) Status(ctx sdk.Context, clientStore types.KVStore, cdc codec.BinaryCodec) exported.Status {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) ExportMetadata(clientStore types.KVStore) []exported.GenesisMetadata {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) ZeroCustomFields() exported.ClientState {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) GetTimestampAtHeight(ctx sdk.Context, clientStore types.KVStore, cdc codec.BinaryCodec, height exported.Height) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) Initialize(ctx sdk.Context, cdc codec.BinaryCodec, clientStore types.KVStore, consensusState exported.ConsensusState) error {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) VerifyMembership(ctx sdk.Context, clientStore types.KVStore, cdc codec.BinaryCodec, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) VerifyNonMembership(ctx sdk.Context, clientStore types.KVStore, cdc codec.BinaryCodec, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path) error {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) VerifyClientMessage(ctx sdk.Context, cdc codec.BinaryCodec, clientStore types.KVStore, clientMsg exported.ClientMessage) error {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) CheckForMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore types.KVStore, clientMsg exported.ClientMessage) bool {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) UpdateStateOnMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore types.KVStore, clientMsg exported.ClientMessage) {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) UpdateState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore types.KVStore, clientMsg exported.ClientMessage) []exported.Height {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) CheckSubstituteAndUpdateState(ctx sdk.Context, cdc codec.BinaryCodec, subjectClientStore, substituteClientStore types.KVStore, substituteClient exported.ClientState) error {
	//TODO implement me
	panic("implement me")
}

func (m *ClientState) VerifyUpgradeAndUpdateState(ctx sdk.Context, cdc codec.BinaryCodec, store types.KVStore, newClient exported.ClientState, newConsState exported.ConsensusState, upgradeClientProof, upgradeConsensusStateProof []byte) error {
	//TODO implement me
	panic("implement me")
}
