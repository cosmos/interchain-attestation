package types

import (
	errorsmod "cosmossdk.io/errors"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ codectypes.UnpackInterfacesMessage = (*Attestator)(nil)

func (a Attestator) Validate() error {
	if a.AttestatorId == nil {
		return errorsmod.Wrap(ErrInvalidAttestator, "attestator id cannot be empty")
	}

	if a.PublicKey == nil {
		return errorsmod.Wrap(ErrInvalidAttestator, "public key cannot be empty")
	}

	pubKey, ok := a.PublicKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "expecting attestation pubkey to be cryptotypes.PubKey, got %T", pubKey)
	}

	if a.ConsensusPubkey == nil {
		return errorsmod.Wrap(ErrInvalidAttestator, "consensus pubkey cannot be empty")
	}

	consPubKey, ok := a.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "expecting consensus pubkey to be cryptotypes.PubKey, got %T", consPubKey)
	}

	if a.ConsensusPubkey == nil {
		return errorsmod.Wrap(ErrInvalidAttestator, "consensus pubkey cannot be empty")
	}

	return nil
}

func (a Attestator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	if err := unpacker.UnpackAny(a.PublicKey, &pubKey); err != nil {
		return err
	}

	var consPubKey cryptotypes.PubKey
	if err := unpacker.UnpackAny(a.ConsensusPubkey, &consPubKey); err != nil {
		return err
	}

	return nil
}
