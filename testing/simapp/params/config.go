package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gjermundgaraba/pessimistic-validation/simapp"
)

func InitSDKConfig() {
	// Set prefixes
	accountPubKeyPrefix := simapp.AccountAddressPrefix + "pub"
	validatorAddressPrefix := simapp.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := simapp.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := simapp.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := simapp.AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(simapp.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}