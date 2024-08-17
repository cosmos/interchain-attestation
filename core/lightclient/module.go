package lightclient

import (
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var _ module.AppModuleBasic = (*AppModuleBasic)(nil)
var _ appmodule.AppModule   = (*AppModule)(nil)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterLegacyAminoCodec performs a no-op. The attestation client does not support amino.
func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {}

// RegisterInterfaces registers module concrete types into protobuf Any. This allows core IBC
// to unmarshal the attestation light client types.
func (a AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes performs a no-op.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {}

// AppModule is the application module for the attestation client module
type AppModule struct {
	AppModuleBasic
	lightClientModule LightClientModule
}

// NewAppModule creates a new attestation client AppModule
func NewAppModule(lightClientModule LightClientModule) AppModule {
	return AppModule{
		lightClientModule: lightClientModule,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModuleBasic) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModuleBasic) IsOnePerModuleType() {}