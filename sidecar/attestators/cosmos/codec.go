package cosmos

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v9/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v9/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v9/modules/core/04-channel/types"
	tmclient "github.com/cosmos/ibc-go/v9/modules/light-clients/07-tendermint"
	"github.com/gjermundgaraba/pessimistic-validation/core/lightclient"
)

type CodecConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.Codec
	//TxConfig          client.TxConfig // Add if we need to do txs at some point
}

func NewCodecConfig() CodecConfig {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	ibcclienttypes.RegisterInterfaces(interfaceRegistry)
	ibcconnectiontypes.RegisterInterfaces(interfaceRegistry)
	ibcchanneltypes.RegisterInterfaces(interfaceRegistry)
	tmclient.RegisterInterfaces(interfaceRegistry)
	lightclient.RegisterInterfaces(interfaceRegistry)
	transfertypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	return CodecConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         cdc,
		//TxConfig:          nil,
	}
}