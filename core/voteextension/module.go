package voteextension

import (
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	clientkeeper "github.com/cosmos/ibc-go/v9/modules/core/02-client/keeper"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"os"
)

const ModuleName = "attestationvoteextension"

var _ module.AppModuleBasic = (*AppModuleBasic)(nil)
var _ appmodule.AppModule = (*AppModule)(nil)

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

	// TODO: Should we just move this stuff into a keeper, or is it fine here?
	sidecarAddress string
	clientKeeper   *clientkeeper.Keeper
	cdc            codec.Codec

	// Create lazily
	sidecarGrpcClient *grpc.ClientConn
}

// NewAppModule creates a new attestation vote extension AppModule
func NewAppModule(clientKeeper *clientkeeper.Keeper, cdc codec.Codec) AppModule {
	sidecarAddress := os.Getenv(SidecarAddressEnv)

	return AppModule{
		sidecarAddress: sidecarAddress,
		clientKeeper:   clientKeeper,
		cdc:            cdc,
	}
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModuleBasic) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModuleBasic) IsOnePerModuleType() {}
