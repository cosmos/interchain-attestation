package attestators

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/gjermundgaraba/interchain-attestation/sidecar/config"
	"os"
)

const (
	DefaultSigningPubKeyFileName  = "signing_pub_key.json"
	DefaultSigningPrivKeyFileName = "signing_priv_key.json"
)

type AttestatorSigningKey struct {
	PubKey cryptotypes.PubKey
	PrivKey cryptotypes.PrivKey
}

func GenerateAttestatorSigningKey() (AttestatorSigningKey, error) {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()

	return AttestatorSigningKey{
		PrivKey: privKey,
		PubKey: pubKey,
	}, nil
}

func AttestatorSigningKeyFromConfig(cdc codec.Codec, cfg config.Config) (AttestatorSigningKey, error) {
	signingPrivKeyJSON, err := os.ReadFile(cfg.SigningPrivateKeyPath)
	if err != nil {
		return AttestatorSigningKey{}, err
	}

	return AttestatorSigningKeyFromJSON(cdc, signingPrivKeyJSON)
}

func AttestatorSigningKeyFromJSON(cdc codec.Codec, privKeyJSON []byte) (AttestatorSigningKey, error) {
	var privKeyAny codectypes.Any
	err := cdc.UnmarshalJSON(privKeyJSON, &privKeyAny)
	if err != nil {
		return AttestatorSigningKey{}, err
	}

	var privKey cryptotypes.PrivKey
	err = cdc.UnpackAny(&privKeyAny, &privKey)
	if err != nil {
		return AttestatorSigningKey{}, err
	}

	pubKey := privKey.PubKey()

	return AttestatorSigningKey{
		PrivKey: privKey,
		PubKey: pubKey,
	}, nil
}

func (a AttestatorSigningKey) PubKeyJSON(cdc codec.Codec) ([]byte, error) {
	pubKeyAny, err := codectypes.NewAnyWithValue(a.PubKey)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSON(pubKeyAny)
}

func (a AttestatorSigningKey) PrivKeyJSON(cdc codec.Codec) ([]byte, error) {
	privKeyAny, err := codectypes.NewAnyWithValue(a.PrivKey)
	if err != nil {
		return nil, err
	}

	return cdc.MarshalJSON(privKeyAny)
}

func (a AttestatorSigningKey) Sign(msg []byte) ([]byte, error) {
	return a.PrivKey.Sign(msg)
}
