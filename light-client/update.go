package lightclient

import (
	"bytes"
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

// VerifyClientMessage checks if the clientMessage is the correct type and verifies the message
func (cs *ClientState) VerifyClientMessage(
	ctx sdk.Context, cdc codec.BinaryCodec, attestatorsHandler AttestatorsHandler,
	clientMsg exported.ClientMessage,
) error {
	pessimisticClaims, ok := clientMsg.(*PessimisticClaims)
	if !ok {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "invalid client message type %T", clientMsg)
	}

	return cs.verifyPessimisticClaims(ctx, cdc, attestatorsHandler, pessimisticClaims)
}

// verifyPessimisticClaims verifies that the provided pessimistic claims are valid, all the same and valid signatures from enough validators
func (cs *ClientState) verifyPessimisticClaims(
	ctx sdk.Context, cdc codec.BinaryCodec, attestatorsHandler AttestatorsHandler,
	pessimisticClaims *PessimisticClaims,
) error {
	if len(pessimisticClaims.Claims) == 0 {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "empty claims")
	}

	seenAttestators := make(map[string]bool)
	var attestatorsSignedOff [][]byte
	for _, claim := range pessimisticClaims.Claims {
		attestator := string(claim.AttestatorId)

		// check that all attestators are unqiue
		_, ok := seenAttestators[attestator]
		if ok {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "duplicate attestation from %s", attestator)
		}
		seenAttestators[attestator] = true
		attestatorsSignedOff = append(attestatorsSignedOff, claim.AttestatorId)
	}

	// check that enough attestators have signed off
	sufficient, err := attestatorsHandler.SufficientAttestations(ctx, attestatorsSignedOff)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "failed to check sufficient attestations: %s", err)
	}
	if !sufficient {
		return errorsmod.Wrapf(ErrInvalidClientMsg, "not enough attestations")
	}

	// check that all claims are packet commitments are unique
	seenPacketCommitments := make(map[string]bool)
	for _, packetCommitements := range pessimisticClaims.Claims[0].PacketCommitmentsClaim.PacketCommitments {
		_, ok := seenPacketCommitments[string(packetCommitements)]
		if ok {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "duplicate packet commitment %s", string(packetCommitements))
		}
		seenPacketCommitments[string(packetCommitements)] = true
	}

	// check that the claims are all the same
	for i, claim := range pessimisticClaims.Claims {
		attestator := string(claim.AttestatorId)

		// verify signature
		signBytes := GetSignableBytes(cdc, claim.PacketCommitmentsClaim)
		pubKey, err := attestatorsHandler.GetPublicKey(ctx, claim.AttestatorId)
		if err != nil {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "failed to get public key for attestator %s: %s", attestator, err)
		}
		if verified := pubKey.VerifySignature(signBytes, claim.Signature); !verified {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "invalid signature from attestator %s", attestator)
		}

		// for the rest we just verify against the first claim, so we skip the first one
		if i == 0 {
			continue
		}

		// check that all claims have the same height
		if !claim.PacketCommitmentsClaim.Height.EQ(pessimisticClaims.Claims[0].PacketCommitmentsClaim.Height) {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "claims must all have the same height")
		}

		// check that all claims have the same timestamp
		if !claim.PacketCommitmentsClaim.Timestamp.Equal(pessimisticClaims.Claims[0].PacketCommitmentsClaim.Timestamp) {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "claims must all have the same timestamp")
		}

		// check that all claims have the same packet commitments
		if !byteSlicesAreEqual(claim.PacketCommitmentsClaim.PacketCommitments, pessimisticClaims.Claims[0].PacketCommitmentsClaim.PacketCommitments) {
			return errorsmod.Wrapf(ErrInvalidClientMsg, "claims must all have the same packet commitments")
		}
	}

	return nil
}

func (cs *ClientState) UpdateState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) []exported.Height {
	pessimisticClaims, ok := clientMsg.(*PessimisticClaims)
	if !ok {
		panic(errorsmod.Wrapf(ErrInvalidClientMsg, "invalid client message type %T", clientMsg))
	}

	if len(pessimisticClaims.Claims) == 0 {
		// perform no-op
		return []exported.Height{}
	}

	height := pessimisticClaims.Claims[0].PacketCommitmentsClaim.Height
	timestamp := pessimisticClaims.Claims[0].PacketCommitmentsClaim.Timestamp
	packetCommitements := pessimisticClaims.Claims[0].PacketCommitmentsClaim.PacketCommitments

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

func byteSlicesAreEqual(slice1, slice2 [][]byte) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if !bytes.Equal(slice1[i], slice2[i]) {
			return false
		}
	}
	return true
}
