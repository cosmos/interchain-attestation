package relayer

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/gjermundgaraba/interchain-attestation/sidecar/config"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"os"
	"time"
)

func (r *Relayer) createClientCtx(ctx context.Context, chainConfig config.CosmosChainConfig) (client.Context, error) {
	cfg := sdk.GetConfig()
	accountPubKeyPrefix := chainConfig.AddressPrefix + "pub"
	validatorAddressPrefix := chainConfig.AddressPrefix + "valoper"
	validatorPubKeyPrefix := chainConfig.AddressPrefix + "valoperpub"
	consNodeAddressPrefix := chainConfig.AddressPrefix + "valcons"
	consNodePubKeyPrefix := chainConfig.AddressPrefix + "valconspub"
	cfg.SetBech32PrefixForAccount(chainConfig.AddressPrefix, accountPubKeyPrefix)
	cfg.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	cfg.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	//cfg.Seal()

	kr, err := keyring.New("attestation-sidecar", chainConfig.KeyringBackend, r.homedir, os.Stdin, r.cdc)
	if err != nil {
		return client.Context{}, err
	}

	txCfg := authtx.NewTxConfig(r.cdc, authtx.DefaultSignModes)

	from, err := kr.Key(chainConfig.KeyName)
	if err != nil {
		panic(err)
	}
	fromAddr, err := from.GetAddress()
	if err != nil {
		panic(err)
	}

	rpcClient, err := client.NewClientFromNode(chainConfig.RPC)
	if err != nil {
		panic(err)
	}

	r.logger.Debug("Creating client context",
		zap.String("chain_id", chainConfig.ChainID),
		zap.String("rpc", chainConfig.RPC),
		zap.String("address_prefix", chainConfig.AddressPrefix),
		zap.String("key_name", chainConfig.KeyName),
		zap.String("from_address", fromAddr.String()),
	)

	return client.Context{}.
		WithCmdContext(ctx).
		WithCodec(r.cdc).
		WithInterfaceRegistry(r.cdc.InterfaceRegistry()).
		WithTxConfig(txCfg).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithHomeDir(r.homedir).
		WithChainID(chainConfig.ChainID).
		WithKeyring(kr).
		WithOffline(false).
		WithNodeURI(chainConfig.RPC).
		WithFromName(chainConfig.KeyName).
		WithFromAddress(fromAddr).
		WithFrom(fromAddr.String()).
		WithClient(rpcClient).
		WithBroadcastMode("sync").
		WithViper(""), nil
}

// TODO: return error
func (r *Relayer) createTxFactory(clientCtx client.Context, chainConfig config.CosmosChainConfig) tx.Factory {
	gasSetting, err := flags.ParseGasSetting(chainConfig.Gas)
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
		WithGasPrices(chainConfig.GasPrices).
		WithSimulateAndExecute(gasSetting.Simulate).
		WithGasAdjustment(chainConfig.GasAdjustment).
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
			return nil, errors.Errorf("failed to validate msg: %w", err)
		}
	}

	var err error
	txf, err = txf.Prepare(clientCtx)
	if err != nil {
		return nil, errors.Errorf("failed to prepare transaction: %w", err)
	}

	adjusted, err := calculateGas(clientCtx, txf, msgs...)
	if err != nil {
		return nil, errors.Errorf("failed to calculate gas: %w", err)
	}
	txf = txf.WithGas(adjusted)
	r.logger.Debug("Estimated gas", zap.Uint64("Gas", txf.Gas()))

	builtTx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, errors.Errorf("failed to build unsigned tx: %w", err)
	}

	if err = tx.Sign(clientCtx.CmdContext, txf, clientCtx.FromName, builtTx, true); err != nil {
		return nil, errors.Errorf("failed to sign tx: %w", err)
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(builtTx.GetTx())
	if err != nil {
		return nil, errors.Errorf("failed to encode tx: %w", err)
	}

	// broadcast to a CometBFT node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, errors.Errorf("failed to broadcast tx: %w", err)
	}

	if res.Code != 0 {
		return nil, errors.Errorf("transaction failed with code != 0: %s", res.RawLog)
	}

	if _, err = clientCtx.Codec.MarshalJSON(res); err != nil {
		return nil, errors.Errorf("failed to marshal tx response: %w", err)
	}
	var msgTypes []string
	for _, msg := range msgs {
		msgTypes = append(msgTypes, sdk.MsgTypeURL(msg))

	}
	r.logger.Info("Successfully broadcast tx", zap.Strings("msgs", msgTypes))

	return r.waitForTX(clientCtx, res.TxHash)
}

func calculateGas(clientCtx gogogrpc.ClientConn, txf tx.Factory, msgs ...sdk.Msg) (uint64, error) {
	txBytes, err := txf.BuildSimTx(msgs...)
	if err != nil {
		return 0, errors.Errorf("failed to build sim tx: %w", err)
	}

	txSvcClient := txtypes.NewServiceClient(clientCtx)

	var gas uint64
	if err := WaitUntilCondition(10*time.Second, 2*time.Second, func() (bool, error) {
		simRes, err := txSvcClient.Simulate(context.Background(), &txtypes.SimulateRequest{
			TxBytes: txBytes,
		})
		if err != nil {
			return false, nil
		}

		gas = uint64(txf.GasAdjustment() * float64(simRes.GasInfo.GasUsed))
		return true, nil
	}); err != nil {
		return 0, errors.Errorf("failed to wait for simulation: %w", err)
	}

	return gas, nil
}

func (r *Relayer) waitForTX(clientCtx client.Context, txHash string) (*sdk.TxResponse, error) {
	r.logger.Debug("Starting to wait for tx", zap.String("tx_hash", txHash))
	try := 1
	maxTries := 25
	for {
		txResp, err := authtx.QueryTx(clientCtx, txHash)
		if err != nil {
			if try == maxTries {
				return nil, errors.Errorf("failed to wait for tx after max tries: %d: %w", maxTries, err)
			}

			r.logger.Debug("Waiting for transaction", zap.String("tx_hash", txHash), zap.Int("try", try), zap.Error(err))
			time.Sleep(500 * time.Millisecond)
			try++
			continue
		}

		if txResp.Code != 0 {
			return nil, errors.Errorf("transaction failed: %s", txResp.RawLog)
		}

		r.logger.Info("Transaction succeeded on chain", zap.String("tx_hash", txHash), zap.String("chain_id", clientCtx.ChainID))
		time.Sleep(5 * time.Second) // for good measure
		return txResp, nil
	}
}
