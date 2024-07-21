package lightclient

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	modulev1 "github.com/gjermundgaraba/pessimistic-validation/lightclient/api/pessimisticvalidation/lightclient/module/v1"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var _ module.AppModuleBasic = (*AppModuleBasic)(nil)
var _ appmodule.AppModule   = (*AppModule)(nil)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterLegacyAminoCodec performs a no-op. The pessimistic client does not support amino.
func (a AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {}

// RegisterInterfaces registers module concrete types into protobuf Any. This allows core IBC
// to unmarshal the pessimistic light client types.
func (a AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes performs a no-op.
func (a AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {}

// AppModule is the application module for the pessimistic client module
type AppModule struct {
	AppModuleBasic
	lightClientModule LightClientModule
}

// NewAppModule creates a new pessimistic client AppModule
func NewAppModule(lightClientModule LightClientModule) AppModule {
	return AppModule{
		lightClientModule: lightClientModule,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModuleBasic) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModuleBasic) IsOnePerModuleType() {}

// ----------------------------------------------------------------------------
// App Wiring Setup
// ----------------------------------------------------------------------------

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *modulev1.Module
	Cdc          codec.Codec
	StoreService store.KVStoreService
}

type ModuleOutputs struct {
	depinject.Out

	Module            appmodule.AppModule
	LightClientModule LightClientModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	lightClientModule := NewLightClientModule(in.Cdc)
	m := NewAppModule(lightClientModule)

	return ModuleOutputs{
		Module: m,
		LightClientModule: lightClientModule,
	}
}