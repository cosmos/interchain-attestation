package types

import (
	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterAttestator{}
	_ codectypes.UnpackInterfacesMessage = (*MsgRegisterAttestator)(nil)
)


func (msg MsgRegisterAttestator) Validate(ac address.Codec) error {
	_, err := ac.StringToBytes(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}

	if msg.AttestatorId == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("attestator id cannot be empty")
	}

	if msg.AttestationPublicKey == nil {
		return sdkerrors.ErrInvalidRequest.Wrap("public key cannot be empty")
	}

	pubKey, ok := msg.AttestationPublicKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "expecting attestation public key to be cryptotypes.PubKey, got %T", pubKey)
	}

	return nil
}

func (msg MsgRegisterAttestator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.AttestationPublicKey, &pubKey)
}