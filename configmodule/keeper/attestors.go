package keeper

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/configmodule/types"
)

func (k Keeper) SetNewAttestator(ctx context.Context, attestator types.Attestator) error {
	has, err := k.Attestators.Has(ctx, attestator.AttestatorId)
	if err != nil {
		return err
	}

	if has {
		return types.ErrAttestatorAlreadyExists
	}

	return k.Attestators.Set(ctx, attestator.AttestatorId, attestator)
}
