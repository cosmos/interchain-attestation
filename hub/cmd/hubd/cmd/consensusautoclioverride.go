package cmd

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	consensusv1 "cosmossdk.io/api/cosmos/consensus/v1"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
)

type ConsensusAutoCliOverride struct {

}

func (c ConsensusAutoCliOverride) IsAppModule() {}

func (c ConsensusAutoCliOverride) IsOnePerModuleType() {}

func (c ConsensusAutoCliOverride) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: consensusv1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current consensus parameters",
				},
			},
			SubCommands: map[string]*autocliv1.ServiceCommandDescriptor{
				"comet": cmtservice.CometBFTAutoCLIDescriptor,
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: consensusv1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      false, // THIS IS THE WHOLE REASON FOR THIS OVERRIDE
					Use: "update-params",
					Short: "Update the consensus parameters",
				},
			},
		},
	}
}