package lightclient

import (
	errorsmod "cosmossdk.io/errors"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"time"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (m *ConsensusState) ClientType() string {
	return ModuleName
}

func NewConsensusState(
	timestamp time.Time,
	packetCommitments [][]byte,
) *ConsensusState {
	return &ConsensusState{
		Timestamp:         timestamp,
		PacketCommitments: packetCommitments,
	}
}

// GetTimestamp returns the timestamp (in nanoseconds) of the consensus state
func (m *ConsensusState) GetTimestamp() uint64 {
	return uint64(m.Timestamp.UnixNano())
}

func (m *ConsensusState) ValidateBasic() error {
	if m.Timestamp.Unix() <= 0 {
		return errorsmod.Wrap(clienttypes.ErrInvalidConsensus, "timestamp must be a positive Unix time")
	}

	return nil
}

