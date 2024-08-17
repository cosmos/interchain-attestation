package relayer

import (
	"context"
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/light"
	"github.com/cosmos/cosmos-sdk/client"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclientutils "github.com/cosmos/ibc-go/v9/modules/core/02-client/client/utils"
	clienttypes "github.com/cosmos/ibc-go/v9/modules/core/02-client/types"
	commitmenttypes "github.com/cosmos/ibc-go/v9/modules/core/23-commitment/types"
	"github.com/cosmos/ibc-go/v9/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v9/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v9/testing"
	attestationlightclient "github.com/gjermundgaraba/pessimistic-validation/core/lightclient"
	"github.com/gjermundgaraba/pessimistic-validation/sidecar/config"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"time"
)

type ClientType int

const (
	TENDERMINT ClientType = iota
	ATTESTATION
)

var defaultUpgradePath = []string{"upgrade", "upgradedIBCState"}

func (r *Relayer) CreateClients(ctx context.Context, chainConfig config.CosmosChainConfig, clientType ClientType, counterpartyChainConfig config.CosmosChainConfig, counterpartyClientType ClientType) (string, string, error) {
	clientID, err := r.CreateSingleClient(ctx, chainConfig, clientType, counterpartyChainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to create client on chain %s: %w", chainConfig.ChainID, err)
	}

	counterpartyClientID, err := r.CreateSingleClient(ctx, counterpartyChainConfig, counterpartyClientType, chainConfig)
	if err != nil {
		return "", "", errors.Errorf("failed to create client on chain %s: %w", counterpartyChainConfig.ChainID, err)
	}

	return clientID, counterpartyClientID, nil
}

func (r *Relayer) CreateSingleClient(ctx context.Context, chainConfig config.CosmosChainConfig, clientType ClientType, counterpartyChainConfig config.CosmosChainConfig) (string, error) {
	clientState, consensusState, err := r.CreateClientAndConsensusState(ctx, clientType, counterpartyChainConfig)

	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return "", err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	msg, err := clienttypes.NewMsgCreateClient(clientState, consensusState, clientCtx.From)
	if err != nil {
		return "", err
	}

	txResp, err := r.sendTx(clientCtx, txf, msg)
	if err != nil {
		return "", err
	}

	return ibctesting.ParseClientIDFromEvents(txResp.Events)
}

/* For later: Eureka
func (r *Relayer) ProvideCounterparty(ctx context.Context, chainConfig config.CosmosChainConfig, clientID string, counterpartyClientID string) error {
	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return err
	}

	txf := r.createTxFactory(clientCtx, chainConfig)

	//merklePath := commitmenttypes.NewMerklePath([]byte(exported.StoreKey), host.KeyClientStorePrefix, []byte(clientID))
	msg := clienttypes.NewMsgProvideCounterparty(clientCtx.From, clientID, counterpartyClientID) // , &merklePath)

	if _, err := r.sendTx(clientCtx, txf, msg); err != nil {
		return err
	}

	return nil
}*/

func (r *Relayer) GetUnbondingPeriod(ctx context.Context, chainConfig config.CosmosChainConfig) (time.Duration, error) {
	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return 0, err
	}

	queryClient := stakingtypes.NewQueryClient(clientCtx)
	res, err := queryClient.Params(clientCtx.CmdContext, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return 0, err
	}

	return res.Params.UnbondingTime, nil
}

func (r *Relayer) UpdateClient(ctx context.Context, chainConfig config.CosmosChainConfig, counterpartyChainConfig config.CosmosChainConfig) (uint64, error) {
	r.logger.Debug("updating client", zap.String("chain_id", chainConfig.ChainID), zap.String("client_id", chainConfig.ClientID), zap.String("counterparty_chain_id", counterpartyChainConfig.ChainID))

	clientID := chainConfig.ClientID

	clientType, err := GetClientType(clientID)
	if err != nil {
		return 0, err
	}

	var updatedClientStateHeight uint64
	switch clientType {
	case TENDERMINT:
		// TODO: are there any more efficient ways to do this??
		counterpartyClientCtx, err := r.createClientCtx(ctx, counterpartyChainConfig)
		if err != nil {
			return 0, err
		}
		header, _, err := ibcclientutils.QueryTendermintHeader(counterpartyClientCtx)
		if err != nil {
			return 0, errors.Errorf("failed to query tendermint header: %w", err)
		}
		updatedClientStateHeight = header.GetHeight().GetRevisionHeight()

		currentClientState, err := r.QueryClientState(ctx, chainConfig)
		if err != nil {
			return 0, err
		}
		currentTendermintClientState, ok := currentClientState.(*tmclient.ClientState)
		if !ok {
			return 0, errors.Errorf("queried client state for chain id %s and client id %s is not a tendermint client state", chainConfig.ChainID, chainConfig.ClientID)
		}
		header.TrustedHeight = currentTendermintClientState.LatestHeight

		counterpartyClientCtx.Height = int64(currentTendermintClientState.LatestHeight.RevisionHeight)
		trustedHeader, _, err := ibcclientutils.QueryTendermintHeader(counterpartyClientCtx)
		if err != nil {
			return 0, errors.Errorf("failed to query trusted tendermint header: %w", err)
		}

		header.TrustedValidators = trustedHeader.ValidatorSet

		clientMsg := &header

		clientCtx, err := r.createClientCtx(ctx, chainConfig)
		if err != nil {
			return 0, err
		}

		txf := r.createTxFactory(clientCtx, chainConfig)
		msg, err := clienttypes.NewMsgUpdateClient(clientID, clientMsg, clientCtx.From)
		if err != nil {
			return 0, err
		}

		if _, err := r.sendTx(clientCtx, txf, msg); err != nil {
			return 0, errors.Errorf("failed to send tx: %w", err)
		}
		time.Sleep(5 * time.Second)

		clientStateAfter, err := r.QueryClientState(ctx, chainConfig)
		if err != nil {
			return 0, err
		}
		tmClientStateAfter, ok := clientStateAfter.(*tmclient.ClientState)
		if !ok {
			return 0, errors.Errorf("queried client state for chain id %s and client id %s is not a tendermint client state", chainConfig.ChainID, chainConfig.ClientID)
		}

		if tmClientStateAfter.LatestHeight.RevisionHeight == currentTendermintClientState.LatestHeight.RevisionHeight {
			return 0, errors.Errorf("tendermint client %s on chain id %s has not been updated (h: %d/%d)", chainConfig.ClientID, chainConfig.ChainID, currentTendermintClientState.LatestHeight.RevisionHeight, tmClientStateAfter.LatestHeight.RevisionHeight)
		}
	case ATTESTATION:
		clientStateBefore, err := r.QueryClientState(ctx, chainConfig)
		if err != nil {
			return 0, err
		}
		attClientStateBefore, ok := clientStateBefore.(*attestationlightclient.ClientState)
		if !ok {
			return 0, errors.Errorf("queried client state for chain id %s and client id %s is not an attestation client state", chainConfig.ChainID, chainConfig.ClientID)
		}

		r.logger.Info("waiting for the attestation client to be updated")
		time.Sleep(10 * time.Second)

		clientState, err := r.QueryClientState(ctx, chainConfig)
		if err != nil {
			return 0, err
		}
		attClientState, ok := clientState.(*attestationlightclient.ClientState)
		if !ok {
			return 0, errors.Errorf("queried client state for chain id %s and client id %s is not an attestation client state", chainConfig.ChainID, chainConfig.ClientID)
		}

		if attClientStateBefore.LatestHeight.RevisionHeight == attClientState.LatestHeight.RevisionHeight {
			return 0, errors.Errorf("attestation client %s on chain id %s has not been updated (h: %d/%d)", chainConfig.ClientID, chainConfig.ChainID, attClientStateBefore.LatestHeight.RevisionHeight, attClientState.LatestHeight.RevisionHeight)
		}

		return attClientState.LatestHeight.RevisionHeight, nil
	}

	return updatedClientStateHeight, nil
}

func (r *Relayer) QueryClientState(ctx context.Context, chainConfig config.CosmosChainConfig) (exported.ClientState, error) {
	clientCtx, err := r.createClientCtx(ctx, chainConfig)
	if err != nil {
		return nil, err
	}

	queryClient := clienttypes.NewQueryClient(clientCtx)
	res, err := queryClient.ClientState(clientCtx.CmdContext, &clienttypes.QueryClientStateRequest{ClientId: chainConfig.ClientID})
	if err != nil {
		return nil, errors.Errorf("failed to query client state for chain id %s and client id %s: %w", chainConfig.ChainID, chainConfig.ClientID, err)
	}

	clientStateUnpacked, err := clienttypes.UnpackClientState(res.ClientState)
	if err != nil {
		return nil, errors.Errorf("failed to unpack client state for chain id %s and client id %s: %w", chainConfig.ChainID, chainConfig.ClientID, err)
	}

	return clientStateUnpacked, nil
}

func (r *Relayer) CreateClientAndConsensusState(ctx context.Context, clientType ClientType, counterpartyChainConfig config.CosmosChainConfig) (exported.ClientState, exported.ConsensusState, error) {
	counterpartyClientCtx, err := r.createClientCtx(ctx, counterpartyChainConfig)
	if err != nil {
		return nil, nil, err
	}

	var clientState exported.ClientState
	var consensusState exported.ConsensusState

	switch clientType {
	case TENDERMINT:
		clientState, consensusState, err = generateTendermintStates(ctx, counterpartyClientCtx, counterpartyChainConfig)
		if err != nil {
			return nil, nil, err
		}
	case ATTESTATION:
		clientState, consensusState, err = generateAttestationStates(ctx, counterpartyClientCtx, counterpartyChainConfig)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, errors.New("unknown client type")
	}

	return clientState, consensusState, nil
}

func generateTendermintStates(ctx context.Context, counterpartyClientCtx client.Context, counterpartyChainConfig config.CosmosChainConfig) (exported.ClientState, exported.ConsensusState, error) {
	stakingQueryClient := stakingtypes.NewQueryClient(counterpartyClientCtx)
	resp, err := stakingQueryClient.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, nil, errors.Errorf("failed to query staking params: %w", err)
	}

	consensusState, height, err := generateTendermintConsensusState(counterpartyClientCtx)
	if err != nil {
		return nil, nil, err
	}

	clientState := &tmclient.ClientState{
		ChainId:         counterpartyChainConfig.ChainID,
		TrustLevel:      tmclient.NewFractionFromTm(light.DefaultTrustLevel),
		TrustingPeriod:  (resp.Params.UnbondingTime / 3) * 2,
		UnbondingPeriod: resp.Params.UnbondingTime,
		MaxClockDrift:   10 * time.Minute,
		FrozenHeight:    clienttypes.ZeroHeight(),
		LatestHeight:    counterpartyChainConfig.GetClientHeight(height),
		ProofSpecs:      commitmenttypes.GetSDKSpecs(),
		UpgradePath:     defaultUpgradePath,
	}

	return clientState, consensusState, nil
}

func generateTendermintConsensusState(clientCtx client.Context) (exported.ConsensusState, uint64, error) {
	header, height, err := ibcclientutils.QueryTendermintHeader(clientCtx)
	if err != nil {
		return nil, 0, errors.Errorf("failed to get latest light block for tm consensus state: %w", err)
	}

	consensusState := &tmclient.ConsensusState{
		Timestamp:          header.SignedHeader.Header.GetTime(),
		Root:               commitmenttypes.NewMerkleRoot(header.SignedHeader.Header.AppHash),
		NextValidatorsHash: header.SignedHeader.Header.NextValidatorsHash,
	}

	return consensusState, uint64(height), nil
}

func generateAttestationStates(ctx context.Context, counterpartyClientCtx client.Context, counterpartyChainConfig config.CosmosChainConfig) (exported.ClientState, exported.ConsensusState, error) {
	cometClient := counterpartyClientCtx.Client
	status, err := cometClient.Status(ctx)
	if err != nil {
		return nil, nil, errors.Errorf("failed to query status for attestation states: %w", err)
	}

	height := status.SyncInfo.LatestBlockHeight
	timestamp := status.SyncInfo.LatestBlockTime

	clientState := &attestationlightclient.ClientState{
		ChainId:            counterpartyChainConfig.ChainID,
		RequiredTokenPower: math.NewInt(1), // TODO: take this in or something
		FrozenHeight:       clienttypes.ZeroHeight(),
		LatestHeight:       counterpartyChainConfig.GetClientHeight(uint64(height)),
	}

	consensusState := &attestationlightclient.ConsensusState{
		Timestamp: timestamp,
	}

	return clientState, consensusState, nil
}
