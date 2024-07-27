package keeper

import (
	"context"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/gjermundgaraba/pessimistic-validation/lightclient"
)

var _ lightclient.AttestatorsHandler = Keeper{}

func (k Keeper) GetPublicKey(ctx context.Context, attestatorId []byte) (cryptotypes.PubKey, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) SufficientAttestations(ctx context.Context, attestatorIds [][]byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}
