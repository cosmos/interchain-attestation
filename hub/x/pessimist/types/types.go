package types

import (
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

func (h *Height) ToIBCHeight() exported.Height {
	return clienttypes.NewHeight(h.RevisionNumber, h.RevisionHeight)
}
