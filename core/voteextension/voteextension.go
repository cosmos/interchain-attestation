package voteextension

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	ve "vote-extensions.dev"

	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/json"

	attestationlightclient "github.com/cosmos/interchain-attestation/core/lightclient"
	"github.com/cosmos/interchain-attestation/core/types"
)

var _ ve.HasVoteExtension = AppModule{}

// TODO: Add checks for ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight

// ExtendVote asks sidecar for ibc data and return vote extension
// This is called for each validator and is non-deteministic (each validator may have different data)
func (a AppModule) ExtendVote(ctx sdk.Context, vote *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	ctx.Logger().Info("AttestationVoteExtension: ExtendVote")

	sidecarAddress := a.GetSidecarAddress(ctx)
	if sidecarAddress == "" {
		ctx.Logger().Info("AttestationVoteExtension: ExtendVote (no sidecar address set)")
		return &abci.ResponseExtendVote{}, nil
	}

	if a.sidecarGrpcClient == nil || a.sidecarGrpcClient.GetState() != connectivity.Ready {
		var err error
		a.sidecarGrpcClient, err = grpc.NewClient(sidecarAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			ctx.Logger().Error("AttestationVoteExtension: ExtendVote (failed to create client)", "sidecarAddress", sidecarAddress, "error", err)
			return &abci.ResponseExtendVote{}, nil // TODO: Should this return the error or not? We need to check what the correct handling is
		}
	}

	sidecarClient := types.NewSidecarClient(a.sidecarGrpcClient)
	resp, err := sidecarClient.GetIBCData(ctx, &types.GetIBCDataRequest{})
	if err != nil {
		ctx.Logger().Error("AttestationVoteExtension: ExtendVote (failed to get ibc data from sidecar)", "sidecarAddress", sidecarAddress, "error", err)
		return &abci.ResponseExtendVote{}, nil // TODO: Should this return the error or not? We need to check what the correct handling is
	}

	ctx.Logger().Info("AttestationVoteExtension: ExtendVote (got ibc data)", "num_ibc_data", len(resp.IbcData))
	for i, attestation := range resp.IbcData {
		ctx.Logger().Info("AttestationVoteExtension: ExtendVote (attestation)",
			"ibc_data #", i+1,
			"height", attestation.Height.RevisionHeight,
			"client_id", attestation.ClientId,
			"timestamp", attestation.Timestamp,
			"client_to_update", attestation.ClientToUpdate,
			"num_packet_commitments", len(attestation.PacketCommitments),
		)
	}

	voteExtension := &VoteExtension{
		IbcData: resp.IbcData,
	}

	voteExtensionBz, err := a.cdc.Marshal(voteExtension)
	if err != nil {
		ctx.Logger().Error("AttestationVoteExtension: ExtendVote (failed to marshal vote extension)", "error", err)
		return &abci.ResponseExtendVote{}, nil // TODO: Should this return the error or not? We need to check what the correct handling is
	}

	return &abci.ResponseExtendVote{
		VoteExtension: voteExtensionBz,
	}, nil
}

// VerifyVote verifies a vote extension
// It should primarily focus on basic sanity checking of the data in the vote extension
func (a AppModule) VerifyVote(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	ctx.Logger().Info("AttestationVoteExtension: VerifyVote (doing nothing atm)")

	// TODO: implement me

	return &abci.ResponseVerifyVoteExtension{
		Status: abci.ResponseVerifyVoteExtension_ACCEPT,
	}, nil
}

// PrepareProposal prepares a proposal by adding adding a "fake tx" to the proposal with the vote extensions
// Contract: We do not need to check if vote extensions are enabled as that is dealt with by the top level handler
func (a AppModule) PrepareProposal(ctx sdk.Context, proposal *abci.RequestPrepareProposal, bytes []byte) (*abci.ResponsePrepareProposal, error) {
	ctx.Logger().Info("AttestationVoteExtension: PrepareProposal", "num_votes", len(proposal.LocalLastCommit.Votes))

	extInfo := proposal.LocalLastCommit

	specialTxBz, err := a.cdc.Marshal(&extInfo)
	if err != nil {
		ctx.Logger().Error("failed to marshal client updates", "error", err)
		return nil, err
	}

	// TODO: Deal with resizing the proposal.Txs slice if too big

	return &abci.ResponsePrepareProposal{
		Txs: append([][]byte{specialTxBz}, proposal.Txs...),
	}, nil
}

// ProcessProposal primarily focuses on validating the proposal and checking if it has enough votes and that all the signatures are correct (and that the signatures come from the correct validators)
// Contract: We do not need to check if vote extensions are enabled as that is dealt with by the top level handler
func (a AppModule) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal, i int) (*abci.ResponseProcessProposal, error) {
	ctx.Logger().Info("AttestationVoteExtension: ProcessProposal (doing nothing atm)")

	// Validate the proposal (enough votes, signatures, etc)

	// TODO: implement me

	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
}

// PreBlocker takes apart the "fake tx" (extended commit info) and constructs client updates from it
// Contract: We do not need to check if vote extensions are enabled as that is dealt with by the top level handler
func (a AppModule) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock, txIndex int) error {
	// If something panics here, don't panic the app
	defer func() {
		if r := recover(); r != nil {
			ctx.Logger().Error("AttestationVoteExtension: PreBlocker (panic recovered)", "panic", r)
		}
	}()

	// Extract "fake tx" and send an update client msg to the light client
	ctx.Logger().Info("AttestationVoteExtension: PreBlocker")

	if len(req.Txs) == 0 {
		ctx.Logger().Info("AttestationVoteExtension: PreBlocker doing nothing (no txs)")
		return nil
	}

	specialTxBz := req.Txs[0] // TODO: use txIndex to get our correct one
	var extInfo abci.ExtendedCommitInfo
	if err := a.cdc.Unmarshal(specialTxBz, &extInfo); err != nil {
		ctx.Logger().Error("failed to unmarshal client updates", "error", err)
		return nil
	}

	clientClaims := make(map[string]*attestationlightclient.AttestationClaim)
	for _, vote := range extInfo.Votes {
		if vote.VoteExtension == nil {
			continue
		}

		var ext map[string][]byte
		if err := json.Unmarshal(vote.VoteExtension, &ext); err != nil {
			ctx.Logger().Error("failed to handler unmarshal vote extension", "error", err)
			return nil
		}

		var voteExtension VoteExtension
		if err := a.cdc.Unmarshal(ext[ModuleName], &voteExtension); err != nil {
			ctx.Logger().Error("failed to unmarshal vote extension", "error", err)
			return nil
		}

		if len(voteExtension.IbcData) == 0 {
			ctx.Logger().Error("AttestationVoteExtension: PreBlocker (no attestations in vote extension)")
			continue
		}

		for _, ibcData := range voteExtension.IbcData {
			claim, ok := clientClaims[ibcData.ClientToUpdate]
			if !ok {
				claim = &attestationlightclient.AttestationClaim{}
				clientClaims[ibcData.ClientToUpdate] = claim
			}

			claim.Attestations = append(claim.Attestations, types.Attestation{
				ValidatorAddress: vote.Validator.Address,
				AttestedData:     ibcData,
			})
		}
	}
	if len(clientClaims) == 0 {
		ctx.Logger().Info("AttestationVoteExtension: PreBlocker (no client claims)")
		return nil
	}

	clientUpdates := ClientUpdates{
		ClientUpdates: make([]ClientUpdate, len(clientClaims)),
	}
	i := 0
	for clientID, claim := range clientClaims {
		ctx.Logger().Info("AttestationVoteExtension: PrepareProposal (adding client update)",
			"client_id", clientID,
			"num_attestations", len(claim.Attestations),
		)
		clientUpdates.ClientUpdates[i] = ClientUpdate{
			ClientToUpdate:   clientID,
			AttestationClaim: *claim,
		}
		i++
	}
	for _, clientUpdate := range clientUpdates.ClientUpdates {
		if err := a.trustedUpdateClientFunc(ctx, clientUpdate.ClientToUpdate, &clientUpdate.AttestationClaim); err != nil {
			ctx.Logger().Error("failed to update client", "error", err, "client_id", clientUpdate.ClientToUpdate)
		}

		ctx.Logger().Info("AttestationVoteExtension: PreBlocker (updated client)", "client_id", clientUpdate.ClientToUpdate)
	}

	return nil
}
