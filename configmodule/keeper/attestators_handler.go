package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/interchain-attestation/core/lightclient"
)

type AttestatorHandler struct{ k Keeper }

var _ lightclient.AttestatorsController = AttestatorHandler{}

func NewAttestatorHandler(k Keeper) lightclient.AttestatorsController {
	return AttestatorHandler{k: k}
}

func (a AttestatorHandler) SufficientAttestations(ctx context.Context, validatorAddresses [][]byte) (bool, error) {
	for _, validatorAddress := range validatorAddresses {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		val, found := a.k.stakingKeeper.GetValidatorByConsAddr(sdkCtx, validatorAddress)
		if !found {
			sdkCtx.Logger().Warn("validator not found", "validatorAddress", validatorAddress)
			continue // we just ignore unknown validators, but should really not happen
		}

		consensusPower := val.ConsensusPower(a.k.stakingKeeper.PowerReduction(sdkCtx))
	}

	return true, nil
}
