package lightclient

import (
	"bytes"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v9/modules/core/exported"

	"github.com/cosmos/interchain-attestation/core/types"
)

// VerifyClientMessage checks if the clientMessage is the correct type and verifies the message
func (cs *ClientState) VerifyClientMessage(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	attestatorsHandler AttestatorsController,
	clientMsg exported.ClientMessage,
) error {
	attestationClaim, ok := clientMsg.(*AttestationClaim)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "invalid client message type %T", clientMsg)
	}

	return cs.verifyAttestationClaim(ctx, cdc, attestatorsHandler, attestationClaim)
}

// verifyAttestationClaim verifies that the provided attestation claims are valid, all the same and valid signatures from enough validators
func (cs *ClientState) verifyAttestationClaim(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	attestatorsHandler AttestatorsController,
	attestationClaim *AttestationClaim,
) error {
	if len(attestationClaim.Attestations) == 0 {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "empty attestations")
	}

	seenAttestators := make(map[string]bool)
	var attestatorsSignedOff [][]byte
	for _, attestation := range attestationClaim.Attestations {
		attestator := string(attestation.AttestatorId)

		// check that all attestators are unqiue
		_, ok := seenAttestators[attestator]
		if ok {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "duplicate attestation from %s", attestator)
		}
		seenAttestators[attestator] = true
		attestatorsSignedOff = append(attestatorsSignedOff, attestation.AttestatorId)
	}

	// check that enough attestators have signed off
	sufficient, err := attestatorsHandler.SufficientAttestations(ctx, attestatorsSignedOff)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "failed to check sufficient attestations: %s", err)
	}
	if !sufficient {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "not enough attestations")
	}

	// check that all attestations have packet commitments that are unique
	seenPacketCommitments := make(map[string]bool)
	for _, packetCommitements := range attestationClaim.Attestations[0].AttestedData.PacketCommitments {
		_, ok := seenPacketCommitments[string(packetCommitements)]
		if ok {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "duplicate packet commitment %s", string(packetCommitements))
		}
		seenPacketCommitments[string(packetCommitements)] = true
	}

	// Used to check against all the other attestations to make sure they match
	firstAttestationBytes := types.GetDeterministicAttestationBytes(cdc, attestationClaim.Attestations[0].AttestedData)

	// check that the attestations are all the same
	for i, attestation := range attestationClaim.Attestations {
		// we are going to equals check against the first one, so we skip it here
		if i == 0 {
			continue
		}

		attestationBytes := types.GetDeterministicAttestationBytes(cdc, attestation.AttestedData)
		if !bytes.Equal(firstAttestationBytes, attestationBytes) {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "attestations must all be the same")
		}
	}

	return nil
}

func (cs *ClientState) UpdateState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) []exported.Height {
	attestationClaim, ok := clientMsg.(*AttestationClaim)
	if !ok {
		panic(errorsmod.Wrapf(ErrInvalidClientMsg, "invalid client message type %T", clientMsg))
	}

	if len(attestationClaim.Attestations) == 0 {
		// perform no-op
		return []exported.Height{}
	}

	height := attestationClaim.Attestations[0].AttestedData.Height
	timestamp := attestationClaim.Attestations[0].AttestedData.Timestamp
	packetCommitements := attestationClaim.Attestations[0].AttestedData.PacketCommitments

	// TODO: Pruning

	// check for duplicate update
	if _, found := getConsensusState(clientStore, cdc, height); found {
		// perform no-op
		return []exported.Height{height}
	}

	if height.GT(cs.LatestHeight) {
		cs.LatestHeight = height
	}

	consensusState := NewConsensusState(timestamp)

	setClientState(clientStore, cdc, cs)
	setConsensusState(clientStore, cdc, consensusState, height)
	setPacketCommitmentState(clientStore, packetCommitements)

	return []exported.Height{height}
}
