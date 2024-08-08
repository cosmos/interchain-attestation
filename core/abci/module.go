package abci

import (
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const ModuleName = "attestationvoteextension"

var _ module.AppModuleBasic = (*AppModuleBasic)(nil)
var _ appmodule.AppModule   = (*AppModule)(nil)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterLegacyAminoCodec performs a no-op. The attestation vote extension does not support amino.
func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {}

// RegisterInterfaces registers module concrete types into protobuf Any.
func (a AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
}

// RegisterGRPCGatewayRoutes performs a no-op.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {}

// AppModule is the application module for the attestation vote extension module
type AppModule struct {
	AppModuleBasic
}

// NewAppModule creates a new attestation vote extension AppModule
func NewAppModule() AppModule {
	return AppModule{}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModuleBasic) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModuleBasic) IsOnePerModuleType() {}
