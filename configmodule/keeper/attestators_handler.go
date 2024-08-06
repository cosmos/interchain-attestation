package keeper

import (
	"context"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/gjermundgaraba/pessimistic-validation/core/lightclient"
)

var _ lightclient.AttestatorsController = Keeper{}

func (k Keeper) GetPublicKey(ctx context.Context, attestatorId []byte) (cryptotypes.PubKey, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) SufficientAttestations(ctx context.Context, attestatorIds [][]byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}
