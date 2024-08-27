package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type AttestatorRegistration struct {
	AttestatorID         []byte
	AttestationPublicKey *codectypes.Any
}

type AttestatorRegistrationJson struct {
	AttestatorID         string          `json:"attestator-id"`
	AttestationPublicKey json.RawMessage `json:"attestation-public-key"`
}

func ParseAndValidateAttestationRegistrationJSONFromFile(cdc codec.Codec, path string) (AttestatorRegistration, error) {
	bz, err := os.ReadFile(path)
	if err != nil {
		return AttestatorRegistration{}, err
	}

	attestatorRegistration, err := UnmarshalAttestationRegistrationJSON(cdc, bz)
	if err != nil {
		return AttestatorRegistration{}, err
	}

	if err := attestatorRegistration.Validate(); err != nil {
		return AttestatorRegistration{}, err
	}

	return attestatorRegistration, nil
}

func UnmarshalAttestationRegistrationJSON(cdc codec.Codec, bz []byte) (AttestatorRegistration, error) {
	var attestatorRegistrationJSON AttestatorRegistrationJson
	if err := json.Unmarshal(bz, &attestatorRegistrationJSON); err != nil {
		return AttestatorRegistration{}, err
	}

	var publicKey cryptotypes.PubKey
	if err := cdc.UnmarshalInterfaceJSON(attestatorRegistrationJSON.AttestationPublicKey, &publicKey); err != nil {
		return AttestatorRegistration{}, err
	}
	publicKeyAny, err := codectypes.NewAnyWithValue(publicKey)
	if err != nil {
		return AttestatorRegistration{}, err
	}

	attestatorID, err := base64.StdEncoding.DecodeString(attestatorRegistrationJSON.AttestatorID)
	if err != nil {
		return AttestatorRegistration{}, err
	}

	return AttestatorRegistration{
		AttestatorID:         attestatorID,
		AttestationPublicKey: publicKeyAny,
	}, nil
}

func (a AttestatorRegistration) Validate() error {
	pubKey, ok := a.AttestationPublicKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pubKey)
	}

	if len(a.AttestatorID) == 0 {
		return fmt.Errorf("attestator id cannot be empty")
	}

	return nil
}
