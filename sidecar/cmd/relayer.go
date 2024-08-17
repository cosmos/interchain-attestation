package cmd

import (
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/attestators/cosmos"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/relayer"
	"github.com/spf13/cobra"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
)

func RelayerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer",
		Short: "relayer sub commands",
	}

	cmd.AddCommand(
		CreateCommand(),
		TransferCmd(),
	)

	return cmd
}

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create subcommands",
	}

	cmd.AddCommand(
		CreateClientCmd(),
		CreateConnectionsCmd(),
	)

	return cmd
}

func CreateClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "clients [chain-id] [client-type] [counterparty-chain-id] [counterparty-client-type]",
		Short: "create client for a chain",
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID := args[0]
			clientTypeStr := args[1]
			counterPartyChainID := args[2]
			counterPartyClientTypeStr := args[3]

			clientType, err := relayer.ConvertClientType(clientTypeStr)
			if err != nil {
				return err
			}

			counterPartyClientType, err := relayer.ConvertClientType(counterPartyClientTypeStr)
			if err != nil {
				return err
			}

			logger := GetLogger(cmd)
			homedir := GetHomedir(cmd)
			cfg := GetConfig(cmd)
			codecConfig := cosmos.NewCodecConfig()

			r := relayer.NewRelayer(logger, codecConfig.Marshaler, homedir)

			chainConfig, found := cfg.GetChain(chainID)
			if !found {
				return errors.Errorf("chain with id %s not found in config", chainID)
			}

			counterpartyChainConfig, found := cfg.GetChain(counterPartyChainID)
			if !found {
				return errors.Errorf("counterparty chain with id %s not found in config", counterPartyChainID)
			}

			newClientID, newCounterpartyClientID, err := r.CreateClients(cmd.Context(), chainConfig, clientType, counterpartyChainConfig, counterPartyClientType)
			if err != nil {
				return err
			}

			logger.Info("client created", zap.String("new_client_id", newClientID), zap.String("new_counterparty_client_id", newCounterpartyClientID))

			// TODO: Update config with the new values!

			return nil
		},
	}

	return cmd
}

func CreateConnectionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connections [chain-id] [counterparty-chain-id]",
		Short: "create connections between two chains",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID := args[0]
			counterPartyChainID := args[1]

			logger := GetLogger(cmd)
			homedir := GetHomedir(cmd)
			cfg := GetConfig(cmd)
			codecConfig := cosmos.NewCodecConfig()

			r := relayer.NewRelayer(logger, codecConfig.Marshaler, homedir)

			chainConfig, found := cfg.GetChain(chainID)
			if !found {
				return errors.Errorf("chain with id %s not found in config", chainID)
			}

			counterpartyChainConfig, found := cfg.GetChain(counterPartyChainID)
			if !found {
				return errors.Errorf("counterparty chain with id %s not found in config", counterPartyChainID)
			}

			connectionID, counterpartyConnectionID, err := r.CreateConnections(cmd.Context(), chainConfig, counterpartyChainConfig)
			if err != nil {
				return err
			}

			logger.Info("connections created", zap.String("connection_id", connectionID), zap.String("counterparty_connection_id", counterpartyConnectionID))

			return nil
		},
	}

	return cmd
}

func TransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [from-chain-id] [to-chain-id] [to] [amount]",
		Short: "do an ibc transfer",
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			fromChainID := args[0]
			toChainID := args[1]
			to := args[2]
			amount := args[3]

			logger := GetLogger(cmd)
			homedir := GetHomedir(cmd)
			cfg := GetConfig(cmd)
			codecConfig := cosmos.NewCodecConfig()

			r := relayer.NewRelayer(logger, codecConfig.Marshaler, homedir)

			toChainConfig, found := cfg.GetChain(toChainID)
			if !found {
				return errors.Errorf("to chain with id %s not found in config", toChainID)
			}

			fromChainConfig, found := cfg.GetChain(fromChainID)
			if !found {
				return errors.Errorf("from chain with id %s not found in config", fromChainID)
			}

			packet, err := r.Transfer(cmd.Context(), fromChainConfig, toChainConfig, to, amount)
			if err != nil {
				return err
			}

			logger.Info("transfer packet created", zap.Uint64("sequence", packet.Sequence))

			return nil
		},
	}


	return cmd
}