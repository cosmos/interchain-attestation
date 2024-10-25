package voteextension

import (
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"cosmossdk.io/core/appmodule"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cometbft/cometbft/libs/json"

	lightclient "github.com/cosmos/interchain-attestation/core/lightclient"
)

const ModuleName = "attestationvoteextension"

type SidecarConfig struct {
	SidecarAddress string `json:"sidecar_address"`
}

var (
	_ module.AppModuleBasic = (*AppModuleBasic)(nil)
	_ appmodule.AppModule   = (*AppModule)(nil)
)

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
	sidecarAddress          string
	trustedUpdateClientFunc lightclient.TrustedClientUpdateFunc
	cdc                     codec.Codec

	// Create lazily
	sidecarGrpcClient *grpc.ClientConn
}

// NewAppModule creates a new attestation vote extension AppModule
func NewAppModule(trustedUpdateClientFunc lightclient.TrustedClientUpdateFunc, cdc codec.Codec) AppModule {
	sidecarAddress := os.Getenv(SidecarAddressEnv)

	return AppModule{
		sidecarAddress:          sidecarAddress,
		trustedUpdateClientFunc: trustedUpdateClientFunc,
		cdc:                     cdc,
	}
}

func (a AppModule) GetSidecarAddress(ctx sdk.Context) string {
	if a.sidecarAddress == "" {
		ctx.Logger().Info("GetSidecarAddress: no sidecar address set")
		sidecarConfigPath := os.Getenv(SidecarConfigPathEnv)
		if sidecarConfigPath == "" {
			ctx.Logger().Info("GetSidecarAddress: no sidecar config path set")
			return ""
		}

		sidecarConfigBz, err := os.ReadFile(sidecarConfigPath)
		if err != nil {
			ctx.Logger().Error("GetSidecarAddress: failed to read sidecar config", "path", sidecarConfigPath, "error", err)
			return ""
		}
		var sidecarConfig SidecarConfig
		if err = json.Unmarshal(sidecarConfigBz, &sidecarConfig); err != nil {
			ctx.Logger().Error("GetSidecarAddress: failed to unmarshal sidecar config", "error", err)
			return ""
		}
		a.sidecarAddress = sidecarConfig.SidecarAddress
	}

	return a.sidecarAddress
}

// IsAppModule implements the appmodule.AppModule interface.
func (AppModuleBasic) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModuleBasic) IsOnePerModuleType() {}
