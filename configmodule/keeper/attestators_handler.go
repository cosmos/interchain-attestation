package keeper

import (
	"context"

	"github.com/cosmos/interchain-attestation/core/lightclient"
)

type AttestatorHandler struct{ k Keeper }

var _ lightclient.AttestatorsController = AttestatorHandler{}

func NewAttestatorHandler(k Keeper) lightclient.AttestatorsController {
	return AttestatorHandler{k: k}
}

// TODO: Implement properly
func (a AttestatorHandler) SufficientAttestations(ctx context.Context, attestatorIds [][]byte) (bool, error) {
	// TODO implement me
	// Just return true for now until we implement the actual logic
	return true, nil
}
