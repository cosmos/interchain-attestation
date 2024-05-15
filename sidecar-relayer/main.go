package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	tmclient "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"time"
)

func main() {
	var home string
	flag.StringVar(&home, "home", "", "home dir for config files, keys and so on")

	cdc := SetupCodec()
	ctx := context.Background()

	for {
		fmt.Println("Trying to read config file...")

		config, err := ReadConfigFromFile(home)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			fmt.Println("Trying again in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		if config.SrcChain.ClientId == "" {
			panic("src_chain client_id is required because client creation is not supported yet")
		}
		if config.DstChain.ClientId == "" {
			panic("dst_chain client_id is required because client creation is not supported yet")
		}


		r, err := NewRelayer(ctx, cdc, config, home)
		if err != nil {
			panic(err)
		}

		if config.SrcChain.ConnectionId == "" {
			srcConnectionId, err := r.InitConnection()
			if err != nil {
				panic(err)
			}

			config.SrcChain.ConnectionId = srcConnectionId
			if err := UpdateConfigFile(home, config); err != nil {
				panic(err)
			}
		}

		if config.DstChain.ConnectionId == "" {
			dstConnectionId, err := r.OpenTryConnection()
			if err != nil {
				panic(err)
			}

			config.DstChain.ConnectionId = dstConnectionId
			if err := UpdateConfigFile(home, config); err != nil {
				panic(err)
			}
		}

		srcHeight, err := r.getLatestHeight(config.SrcChain)
		if err != nil {
			panic(err)
		}
		connection, err := r.QueryConnection(ctx, config.SrcChain, srcHeight, config.SrcChain.ConnectionId)
		if err != nil {
			panic(err)
		}
		fmt.Println("Connection state:", connection.Connection.State)

		if connection.Connection.State == ibcconnectiontypes.INIT {
			
		}

		fmt.Println("Loop done, sleeping 5 seconds before going again...")
		time.Sleep(5 * time.Second)
	}
}

func SetupCodec() codec.Codec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	ibcclienttypes.RegisterInterfaces(interfaceRegistry)
	ibcconnectiontypes.RegisterInterfaces(interfaceRegistry)
	ibcchanneltypes.RegisterInterfaces(interfaceRegistry)
	tmclient.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	return cdc
}