package lightclient

import (
	"time"

	errorsmod "cosmossdk.io/errors"

	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (m *ConsensusState) ClientType() string {
	return ModuleName
}

func NewConsensusState(
	timestamp time.Time,
) *ConsensusState {
	return &ConsensusState{
		Timestamp: timestamp,
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
