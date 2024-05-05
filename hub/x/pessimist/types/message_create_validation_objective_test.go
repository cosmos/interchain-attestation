package types

import (
	"testing"

	"hub/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateValidationObjective_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateValidationObjective
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateValidationObjective{
				Creator:       "invalid_address",
				ClientId:      "tt",
				RequiredPower: 1,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateValidationObjective{
				Creator:       sample.AccAddress(),
				ClientId:      "tt",
				RequiredPower: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
