package lightclient

import (
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"hub/x/pessimist/types"
)

func getClientState(store storetypes.KVStore, cdc codec.BinaryCodec) (types.ClientState, bool) {
	bz := store.Get(host.ClientStateKey())
	if len(bz) == 0 {
		return types.ClientState{}, false
	}

	clientStateI := clienttypes.MustUnmarshalClientState(cdc, bz)
	var clientState *types.ClientState
	clientState, ok := clientStateI.(*types.ClientState)
	if !ok {
		panic(fmt.Errorf("cannot convert %T to %T", clientStateI, clientState))
	}

	return *clientState, true
}

func getConsensusState(store storetypes.KVStore, cdc codec.BinaryCodec, height exported.Height) (*tmclient.ConsensusState, bool) {
	bz := store.Get(host.ConsensusStateKey(height))
	if len(bz) == 0 {
		return nil, false
	}

	consensusStateI := clienttypes.MustUnmarshalConsensusState(cdc, bz)
	var consensusState *tmclient.ConsensusState
	consensusState, ok := consensusStateI.(*tmclient.ConsensusState)
	if !ok {
		panic(fmt.Errorf("cannot convert %T to %T", consensusStateI, consensusState))
	}

	return consensusState, true

}

// sets the client state to the store
func setClientState(store storetypes.KVStore, cdc codec.BinaryCodec, clientState *types.ClientState) {
	bz := clienttypes.MustMarshalClientState(cdc, clientState)
	store.Set(host.ClientStateKey(), bz)
}

// setConsensusState stores the consensus state at the given height.
func setConsensusState(clientStore storetypes.KVStore, cdc codec.BinaryCodec, consensusState *tmclient.ConsensusState, height exported.Height) {
	key := host.ConsensusStateKey(height)
	val := clienttypes.MustMarshalConsensusState(cdc, consensusState)
	clientStore.Set(key, val)
}

func (l *LightClientModule) processClientMessage(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) (int64, tmclient.ConsensusState, error) {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return 0, tmclient.ConsensusState{}, errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	committeeProposal, ok := clientMsg.(*types.CommitteeProposal)
	if !ok {
		return 0, tmclient.ConsensusState{}, clienttypes.ErrInvalidClientType
	}
	if err := committeeProposal.ValidateBasic(); err != nil {
		return 0, tmclient.ConsensusState{}, err
	}

	validationObjective, found := l.keeper.GetValidationObjective(ctx, clientState.DependentClientId)
	if !found {
		return 0, tmclient.ConsensusState{}, errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "validation objective not found")
	}

	var votes []ValidationVoteWithCommitment
	for _, commitment := range committeeProposal.Commitments {
		var voteExt types.VoteExtension
		if err := l.cdc.Unmarshal(commitment.CanonicalVoteExtension.Extension, &voteExt); err != nil {
			return 0, tmclient.ConsensusState{}, err
		}

		for _, validationVote := range voteExt.ValidationVotes {
			valAddr := sdk.ConsAddress(commitment.ValidatorAddress)
			validator, valFound, err := l.keeper.GetValidatorForObjective(ctx, valAddr.String(), validationObjective)
			if err != nil {
				return 0, tmclient.ConsensusState{}, err
			}

			// TODO: Check signature here!

			if valFound && validationVote.ClientIdToUpdate == clientID && validationVote.ClientIdToValidate == clientState.DependentClientId {
				votes = append(votes, ValidationVoteWithCommitment{
					ValidationVote: validationVote,
					Commitment: commitment,
					Validator: validator,
				})
			} else {
				ctx.Logger().Info("Vote not relevant for client", "client_id", clientID, "dependent_client_id", clientState.DependentClientId, "validator", valAddr.String(), "validation_vote", validationVote, "validator_found", valFound)
			}
		}
	}

	latestHeight := clientState.LatestHeight

	heightPowers := make(map[int64]uint64)
	heightConsensus := make(map[int64]tmclient.ConsensusState)
	for _, vote := range votes {

		heightPowers[vote.ValidationVote.Height] += vote.Validator.Power
		heightConsensus[vote.ValidationVote.Height] = tmclient.ConsensusState{
			Timestamp:          vote.ValidationVote.Timestamp,
			Root:               commitmenttypes.MerkleRoot{
				Hash: vote.ValidationVote.MerkleRoot.Hash,
			},
			NextValidatorsHash: vote.ValidationVote.NextValidatorHash,
		}

		// Increase votes for any heights less than the current height
		for height, _ := range heightPowers {
			if vote.ValidationVote.Height > height {
				heightPowers[height] += vote.Validator.Power
			}
		}
	}

	proposedHeight := int64(0)
	for height, power := range heightPowers {
		if power >= validationObjective.RequiredPower && height > latestHeight && height > proposedHeight {
			proposedHeight = height
		}
	}

	// Possible thing to do here: check that the new proposed height is not insanely high compared to the current height of the dependent client

	if proposedHeight == 0 {
		ctx.Logger().Error("No valid heights found", "heights", heightPowers)
		return 0, tmclient.ConsensusState{}, errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "no valid heights found")
	}

	proposedConsensusState, ok := heightConsensus[proposedHeight]
	if !ok{
		return 0, tmclient.ConsensusState{}, errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "proposed consensus state not found")
	}

	return proposedHeight, proposedConsensusState, nil
}