package keeper

import (
	"context"
	"fmt"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/gjermundgaraba/interchain-attestation/core/lightclient"
)

type AttestatorHandler struct{ k Keeper }

var _ lightclient.AttestatorsController = AttestatorHandler{}

func NewAttestatorHandler(k Keeper) lightclient.AttestatorsController {
	return AttestatorHandler{k: k}
}

func (a AttestatorHandler) GetPublicKey(ctx context.Context, attestatorId []byte) (cryptotypes.PubKey, error) {
	attestator, err := a.k.Attestators.Get(ctx, attestatorId)
	if err != nil {
		return nil, err
	}

	pk, ok := attestator.PublicKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type for for attestator %s, got %T", attestatorId, pk)
	}

	return pk, nil
}

func (a AttestatorHandler) SufficientAttestations(ctx context.Context, attestatorIds [][]byte) (bool, error) {
	//TODO implement me
	// Just return true for now until we implement the actual logic
	return true, nil
}
