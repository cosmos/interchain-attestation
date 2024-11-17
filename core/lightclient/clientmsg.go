package lightclient

import (
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
)

var _ exported.ClientMessage = (*AttestationTally)(nil)

func (m *AttestationTally) ClientType() string {
	return ModuleName
}

func (m *AttestationTally) ValidateBasic() error {
	// TODO: implement me
	panic("implement me")
}
