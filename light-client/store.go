package lightclient

import (
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// getClientState retrieves the client state from the store using the provided KVStore and codec.
// It returns the unmarshaled ClientState and a boolean indicating if the state was found.
func getClientState(store storetypes.KVStore, cdc codec.BinaryCodec) (*ClientState, bool) {
	bz := store.Get(host.ClientStateKey())
	if len(bz) == 0 {
		return nil, false
	}

	clientStateI := clienttypes.MustUnmarshalClientState(cdc, bz)
	var clientState *ClientState
	clientState, ok := clientStateI.(*ClientState)
	if !ok {
		panic(fmt.Errorf("cannot convert %T into %T", clientStateI, clientState))
	}
	return clientState, true
}

// setClientState stores the client state
func setClientState(clientStore storetypes.KVStore, cdc codec.BinaryCodec, clientState *ClientState) {
	key := host.ClientStateKey()
	val := clienttypes.MustMarshalClientState(cdc, clientState)
	clientStore.Set(key, val)
}

// setConsensusState stores the consensus state at the given height.
func setConsensusState(clientStore storetypes.KVStore, cdc codec.BinaryCodec, consensusState *ConsensusState, height exported.Height) {
	key := host.ConsensusStateKey(height)
	val := clienttypes.MustMarshalConsensusState(cdc, consensusState)
	clientStore.Set(key, val)
}

