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
	if len(m.Commitments) == 0 {
		return errorsmod.Wrap(ErrInvalidCommitteeProposal, "commitments cannot be empty")
	}

	for _, commitment := range m.Commitments {
		if commitment.Timestamp.IsZero() {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "timestamp cannot be zero")
		}

		if commitment.ValidatorAddress == nil {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "validator address cannot be nil")
		}

		valAddr := sdk.ValAddress(commitment.ValidatorAddress)
		if _, err := sdk.ValAddressFromBech32(valAddr.String()); err != nil {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "invalid validator address")
		}

		if commitment.ExtensionSignature == nil {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "extension signature cannot be nil")
		}

		if commitment.CanonicalVoteExtension.Extension == nil {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "extension cannot be nil")
		}

		if commitment.CanonicalVoteExtension.Height <= 0 {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "height cannot be zero or negative")
		}

		if commitment.CanonicalVoteExtension.ChainId == "" {
			return errorsmod.Wrap(ErrInvalidCommitteeProposal, "chain id cannot be empty")
		}
	}

	return nil
}