package main

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"go.uber.org/zap"
	"os"
	"time"
)

const KeyName = "relayer"

type Relayer struct {
	ctx     context.Context
	cdc     codec.Codec
	config Config
	homedir string
	kr 	keyring.Keyring
}

func NewRelayer(ctx context.Context, cdc codec.Codec, config Config, homedir string) (*Relayer, error) {
	kr, err := keyring.New("sidecar-relayer", "test", homedir, os.Stdin, cdc)
	if err != nil {
		return nil, err
	}

	if _, err := kr.Key(KeyName); err != nil {
		fmt.Println("Error getting key from keyring")
		return nil, err
	}

	return &Relayer{
		ctx:     ctx,
		cdc:     cdc,
		config:  config,
		homedir: homedir,
		kr: kr,
	}, nil
}

func (r *Relayer) createClientCtx(chain ChainConfig) client.Context {
	cfg := sdk.GetConfig()
	accountPubKeyPrefix := chain.AccountPrefix + "pub"
	validatorAddressPrefix := chain.AccountPrefix + "valoper"
	validatorPubKeyPrefix := chain.AccountPrefix + "valoperpub"
	consNodeAddressPrefix := chain.AccountPrefix + "valcons"
	consNodePubKeyPrefix := chain.AccountPrefix + "valconspub"
	cfg.SetBech32PrefixForAccount(chain.AccountPrefix, accountPubKeyPrefix)
	cfg.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	cfg.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	//cfg.Seal()

	txCfg := authtx.NewTxConfig(r.cdc, authtx.DefaultSignModes)

	from, err := r.kr.Key(KeyName)
	if err != nil {
		panic(err)
	}
	fromAddr, err := from.GetAddress()
	if err != nil {
		panic(err)
	}

	rpcClient, err := client.NewClientFromNode(chain.RpcAddr)
	if err != nil {
		panic(err)
	}

	return client.Context{}.
		WithCmdContext(r.ctx).
		WithCodec(r.cdc).
		WithInterfaceRegistry(r.cdc.InterfaceRegistry()).
		WithTxConfig(txCfg).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(r.homedir).
		WithChainID(chain.ChainId).
		WithKeyring(r.kr).
		WithOffline(false).
		WithNodeURI(chain.RpcAddr).
		WithFromName(KeyName).
		WithFromAddress(fromAddr).
		WithFrom(fromAddr.String()).
		WithClient(rpcClient).
		WithBroadcastMode("sync").
		WithViper("")
}

func (r *Relayer) createTxFactory(clientCtx client.Context, chain ChainConfig) tx.Factory {

	gasSetting, err := flags.ParseGasSetting("auto")
	if err != nil {
		panic(err)
	}

	return tx.Factory{}.
		WithTxConfig(clientCtx.TxConfig).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithKeybase(clientCtx.Keyring).
		WithChainID(clientCtx.ChainID).
		WithFromName(clientCtx.FromName).
		WithGas(gasSetting.Gas).
		WithGasPrices(chain.GasPrices).
		WithSimulateAndExecute(gasSetting.Simulate).
		WithGasAdjustment(2.0).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
}

// A good amount of this is copied from cosmos sdk client code
func (r *Relayer) sendTx(clientCtx client.Context, txf tx.Factory, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	for _, msg := range msgs {
		m, ok := msg.(sdk.HasValidateBasic)
		if !ok {
			continue
		}

		if err := m.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	var err error
	txf, err = txf.Prepare(clientCtx)
	if err != nil {
		fmt.Println("Failed to prepare tx", err)
		return nil, err
	}

	adjusted, err := calculateGas(clientCtx, txf, msgs...)
	if err != nil {
		return nil, err
	}
	txf = txf.WithGas(adjusted)

	builtTx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	if err = tx.Sign(clientCtx.CmdContext, txf, clientCtx.FromName, builtTx, true); err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(builtTx.GetTx())
	if err != nil {
		return nil, err
	}

	// broadcast to a CometBFT node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	if res.Code != 0 {
		return nil, fmt.Errorf(res.RawLog)
	}

	if _, err = clientCtx.Codec.MarshalJSON(res); err != nil {
		return nil, err
	}
	var msgTypes []string
	for _, msg := range msgs {
		msgTypes = append(msgTypes, sdk.MsgTypeURL(msg))

	}
	fmt.Println("Successfully broadcast tx", msgTypes)

	return r.waitForTX(clientCtx, res.TxHash)
}

func calculateGas(clientCtx gogogrpc.ClientConn, txf tx.Factory, msgs ...sdk.Msg) (uint64, error) {
	txBytes, err := txf.BuildSimTx(msgs...)
	if err != nil {
		return 0, err
	}

	txSvcClient := txtypes.NewServiceClient(clientCtx)
	simRes, err := txSvcClient.Simulate(context.Background(), &txtypes.SimulateRequest{
		TxBytes: txBytes,
	})
	if err != nil {
		return 0, err
	}

	return uint64(txf.GasAdjustment() * float64(simRes.GasInfo.GasUsed)), nil
}

func (r *Relayer) waitForTX(clientCtx client.Context, txHash string) (*sdk.TxResponse, error) {
	fmt.Println("Starting to wait for tx", zap.String("tx_hash", txHash))
	try := 1
	maxTries := 25
	for {
		txResp, err := authtx.QueryTx(clientCtx, txHash)
		if err != nil {
			if try == maxTries {
				err2 := fmt.Errorf("transaction with hash %s exceeded max retry limit of %d with error %s", txHash, try, err)
				fmt.Println("Transaction not found", err, err2)
				return nil, err2
			}

			fmt.Println("Waiting for transaction", zap.String("tx_hash", txHash), zap.Int("try", try), zap.Error(err))
			time.Sleep(500 * time.Millisecond)
			try++
			continue
		}

		if txResp.Code != 0 {
			return nil, fmt.Errorf("transaction failed: %s", txResp.RawLog)
		}

		fmt.Println("Transaction succeeded", zap.String("tx_hash", txHash))
		return txResp, nil
	}
}

func (r *Relayer) getLatestHeight(chain ChainConfig) (int64, error) {
	clientCtx := r.createClientCtx(chain)
	stat, err := clientCtx.Client.Status(clientCtx.CmdContext)
	if err != nil {
		return -1, err
	} else if stat.SyncInfo.CatchingUp {
		return -1, fmt.Errorf("node at %s running chain %s not caught up", clientCtx.NodeURI, clientCtx.ChainID)
	}

	return stat.SyncInfo.LatestBlockHeight, nil
}