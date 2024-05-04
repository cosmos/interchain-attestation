package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)
var _ exported.ClientMessage = (*CommitteeProposal)(nil)

func (m *ClientState) ClientType() string {
	return ClientType
}

func (m *ClientState) Validate() error {
	if m.DependentClientId == "" {
		return errorsmod.Wrap(clienttypes.ErrInvalidClient, "dependent client id cannot be empty")
	}

	return nil
}

func (m *CommitteeProposal) ClientType() string {
	return ClientType
}

func (m *CommitteeProposal) ValidateBasic() error {
	if m.Height.RevisionNumber < 0 {
		return errorsmod.Wrap(ErrInvalidCommitteeProposal, "revision number cannot be negative")
	}

	if m.Height.RevisionHeight < 0 {
		return errorsmod.Wrap(ErrInvalidCommitteeProposal, "revision height cannot be negative")
	}

	if len(m.Commitments) == 0 {
		return errorsmod.Wrap(ErrInvalidCommitteeProposal, "commitments cannot be empty")
	}

	for _, commitment := range m.Commitments {
		if commitment == nil {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "commitment cannot be nil")
		}

		if commitment.Height.RevisionNumber != m.Height.RevisionNumber {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "commitment revision number must match proposal revision number")
		}

		if commitment.Height.RevisionHeight != m.Height.RevisionHeight {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "commitment revision height must match proposal revision height")
		}

		if _, err := sdk.ValAddressFromBech32(commitment.ValidatorAddr); err != nil {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "invalid validator address")
		}

		if commitment.ClientId == "" {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "client id cannot be empty")
		}

		if commitment.Signature == nil {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "signature cannot be nil")
		}

		if commitment.Timestamp.IsZero() {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "timestamp cannot be zero")
		}
	}

	return nil
}

func (m *Commitment) Data() []byte {
	data := []byte(m.ClientId)
	data = append(data, sdk.Uint64ToBigEndian(m.Height.RevisionNumber)...)
	data = append(data, sdk.Uint64ToBigEndian(m.Height.RevisionHeight)...)
	data = append(data, sdk.Uint64ToBigEndian(uint64(m.Timestamp.Unix()))...)
	data = append(data, m.ValidatorAddr...)

	return data
}
