package lightclient

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"hub/x/pessimist/keeper"
	"hub/x/pessimist/types"
)

var _ exported.LightClientModule = (*LightClientModule)(nil)

type LightClientModule struct {
	cdc           codec.BinaryCodec
	storeProvider exported.ClientStoreProvider
	keeper        keeper.Keeper
}

func NewLightClientModule(cdc codec.BinaryCodec, keeper keeper.Keeper) LightClientModule {
	return LightClientModule{
		cdc:    cdc,
		keeper: keeper,
	}
}

func (l *LightClientModule) RegisterStoreProvider(storeProvider exported.ClientStoreProvider) {
	l.storeProvider = storeProvider
}

// TODO: Test this
func (l *LightClientModule) Initialize(ctx sdk.Context, clientID string, clientStateBz, consensusStateBz []byte) error {
	var clientState types.ClientState
	if err := l.cdc.Unmarshal(clientStateBz, &clientState); err != nil {
		return err
	}

	if err := clientState.Validate(); err != nil {
		return err
	}

	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	setClientState(clientStore, l.cdc, &clientState)

	return nil
}

// TODO: Test this
func (l *LightClientModule) VerifyClientMessage(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) error {
	clientStore := l.storeProvider.ClientStore(ctx, clientID)
	clientState, found := getClientState(clientStore, l.cdc)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	committeeProposal, ok := clientMsg.(*types.CommitteeProposal)
	if !ok {
		return clienttypes.ErrInvalidClientType
	}
	if err := committeeProposal.ValidateBasic(); err != nil {
		return err
	}

	// TODO: Call dependent light client to check the height is good

	latestHeight := clientState.LatestHeight.ToIBCHeight()
	proposedHeight := committeeProposal.Height.ToIBCHeight()

	if proposedHeight.LTE(latestHeight) {
		return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "proposed height is less than or equal to the latest height")
	}

	if proposedHeight.EQ(latestHeight) {
		return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "proposed height is equal to the latest height")
	}

	// Not even sure how one would go about supporting an incremented revision number, so just won't support it for now
	if proposedHeight.GetRevisionNumber() != latestHeight.GetRevisionNumber() {
		return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "proposed height revision number must match latest height revision number")
	}

	validationObjective, found := l.keeper.GetValidationObjective(ctx, clientState.DependentClientId)
	if !found {
		return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "validation objective not found")
	}

	var signedValidators []*types.Validator
	for _, commitment := range committeeProposal.Commitments {
		if commitment.ClientId != clientState.DependentClientId {
			return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "commitment client id must match dependent client id")
		}

		// Validate basic has already verified that all the heights are the same as the top level proposal height

		var validator *types.Validator
		for _, v := range validationObjective.Validators {
			if v.ValidatorAddr == commitment.ValidatorAddr {
				validator = v
				break
			}
		}
		if validator == nil {
			return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "validator not found in validator set")
		}

		data := commitment.Data()
		pubKey, ok := validator.PubKey.GetCachedValue().(cryptotypes.PubKey)
		if !ok {
			return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "validator public key is not cryptotypes.PubKey")
		}

		if verified := pubKey.VerifySignature(data, commitment.Signature); !verified {
			return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "signature verification failed")
		}

		signedValidators = append(signedValidators, validator)
	}

	power := l.keeper.GetValidatorPower(ctx, signedValidators)
	if power.LT(math.NewIntFromUint64(validationObjective.RequiredPower)) {
		return errorsmod.Wrap(types.ErrInvalidCommitteeProposal, "insufficient power by signed members")
	}

	return nil
}

func (l *LightClientModule) CheckForMisbehaviour(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) bool {
	//TODO implement me
	return false
}

func (l *LightClientModule) UpdateStateOnMisbehaviour(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) {
	//TODO implement me
}

func (l *LightClientModule) UpdateState(ctx sdk.Context, clientID string, clientMsg exported.ClientMessage) []exported.Height {
	//TODO implement me
	panic("implement me")
}

func (l *LightClientModule) VerifyMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path, value []byte) error {
	//TODO implement me
	return nil
}

func (l *LightClientModule) VerifyNonMembership(ctx sdk.Context, clientID string, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path) error {
	//TODO implement me
	return nil
}

func (l *LightClientModule) Status(ctx sdk.Context, clientID string) exported.Status {
	//TODO implement me
	panic("implement me")
}

func (l *LightClientModule) LatestHeight(ctx sdk.Context, clientID string) exported.Height {
	//TODO implement me
	panic("implement me")
}

func (l *LightClientModule) TimestampAtHeight(ctx sdk.Context, clientID string, height exported.Height) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (l *LightClientModule) RecoverClient(ctx sdk.Context, clientID, substituteClientID string) error {
	//TODO implement me
	panic("implement me")
}

func (l *LightClientModule) VerifyUpgradeAndUpdateState(ctx sdk.Context, clientID string, newClient []byte, newConsState []byte, upgradeClientProof, upgradeConsensusStateProof []byte) error {
	//TODO implement me
	panic("implement me")
}
