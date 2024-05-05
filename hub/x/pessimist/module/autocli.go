package pessimist

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "hub/api/hub/pessimist"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "ValidationObjective",
					Use:            "validation-objective [client-id]",
					Short:          "Query validation-objective",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "clientId"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateValidationObjective",
					Use:            "create-validation-objective [client-id] [required-power]",
					Short:          "Send a create-validation-objective tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "clientId"}, {ProtoField: "requiredPower"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
