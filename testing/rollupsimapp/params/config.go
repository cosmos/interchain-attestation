package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gjermundgaraba/pessimistic-validation/rollupsimapp"
)

func InitSDKConfig() {
	// Set prefixes
	accountPubKeyPrefix := rollupsimapp.AccountAddressPrefix + "pub"
	validatorAddressPrefix := rollupsimapp.AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := rollupsimapp.AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := rollupsimapp.AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := rollupsimapp.AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(rollupsimapp.AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}