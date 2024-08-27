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

// ExtendVote asks sidecar for attestations and return vote extension
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
	resp, err := sidecarClient.GetAttestations(ctx, &types.GetAttestationsRequest{})
	if err != nil {
		ctx.Logger().Error("AttestationVoteExtension: ExtendVote (failed to get attestations from sidecar)", "sidecarAddress", sidecarAddress, "error", err)
		return &abci.ResponseExtendVote{}, nil // TODO: Should this return the error or not? We need to check what the correct handling is
	}

	ctx.Logger().Info("AttestationVoteExtension: ExtendVote (got attestations)", "num_attestations", len(resp.Attestations))
	for i, attestation := range resp.Attestations {
		ctx.Logger().Info("AttestationVoteExtension: ExtendVote (attestation)",
			"attestation #", i+1,
			"attestator_id", attestation.AttestatorId,
			"height", attestation.AttestedData.Height.RevisionHeight,
			"client_id", attestation.AttestedData.ClientId,
			"timestamp", attestation.AttestedData.Timestamp,
			"client_to_update", attestation.AttestedData.ClientToUpdate,
			"num_packet_commitments", len(attestation.AttestedData.PacketCommitments),
		)
	}

	voteExtension := &VoteExtension{
		Attestations: resp.Attestations,
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

// TODO: Document
func (a AppModule) VerifyVote(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	// Verify vote extension

	ctx.Logger().Info("AttestationVoteExtension: VerifyVote (doing nothing atm)")

	// TODO implement me

	return &abci.ResponseVerifyVoteExtension{
		Status: abci.ResponseVerifyVoteExtension_ACCEPT,
	}, nil
}

// TODO: Document
// Contract: We do not need to check if vote extensions are enabled as that is dealt with by the top level handler
func (a AppModule) PrepareProposal(ctx sdk.Context, proposal *abci.RequestPrepareProposal, bytes []byte) (*abci.ResponsePrepareProposal, error) {
	// Extract vote extensions and add "fake tx" to the proposal

	ctx.Logger().Info("AttestationVoteExtension: PrepareProposal", "num_votes", len(proposal.LocalLastCommit.Votes))

	clientClaims := make(map[string]*attestationlightclient.AttestationClaim)
	for _, vote := range proposal.LocalLastCommit.Votes {
		if vote.VoteExtension == nil {
			continue
		}

		var ext map[string][]byte
		if err := json.Unmarshal(vote.VoteExtension, &ext); err != nil {
			ctx.Logger().Error("failed to handler unmarshal vote extension", "error", err)
			return nil, err // TODO: Should we return errors here or not?
		}

		var voteExtension VoteExtension
		if err := a.cdc.Unmarshal(ext[ModuleName], &voteExtension); err != nil {
			ctx.Logger().Error("failed to unmarshal vote extension", "error", err)
			return nil, err // TODO: Should we return errors here or not?
		}

		if len(voteExtension.Attestations) == 0 {
			ctx.Logger().Error("AttestationVoteExtension: PrepareProposal (no attestations in vote extension)")
			continue
		}

		for _, attestation := range voteExtension.Attestations {
			claim, ok := clientClaims[attestation.AttestedData.ClientToUpdate]
			if !ok {
				claim = &attestationlightclient.AttestationClaim{}
				clientClaims[attestation.AttestedData.ClientToUpdate] = claim
			}

			claim.Attestations = append(claim.Attestations, attestation)
		}
	}

	if len(clientClaims) == 0 {
		ctx.Logger().Info("AttestationVoteExtension: PrepareProposal (no client claims)")
		return &abci.ResponsePrepareProposal{
			Txs: proposal.Txs,
		}, nil
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

	specialTxBz, err := a.cdc.Marshal(&clientUpdates)
	if err != nil {
		ctx.Logger().Error("failed to marshal client updates", "error", err)
		return nil, err
	}

	ctx.Logger().Info("AttestationVoteExtension: PrepareProposal (adding special tx) with client updates", "num_client_updates", len(clientUpdates.ClientUpdates))

	return &abci.ResponsePrepareProposal{
		Txs: append([][]byte{specialTxBz}, proposal.Txs...),
	}, nil
}

// TODO: Document
// Contract: We do not need to check if vote extensions are enabled as that is dealt with by the top level handler
func (a AppModule) ProcessProposal(ctx sdk.Context, req *abci.RequestProcessProposal, i int) (*abci.ResponseProcessProposal, error) {
	ctx.Logger().Info("AttestationVoteExtension: ProcessProposal (doing nothing atm)")

	// ?

	// TODO implement me

	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
}

// TODO: Document
// Contract: We do not need to check if vote extensions are enabled as that is dealt with by the top level handler
func (a AppModule) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock, i int) error {
	// Extract "fake tx" and send an update client msg to the light client

	ctx.Logger().Info("AttestationVoteExtension: PreBlocker")

	if len(req.Txs) == 0 {
		ctx.Logger().Info("AttestationVoteExtension: PreBlocker doing nothing (no txs)")
		return nil
	}

	specialTxBz := req.Txs[0] // TODO: use i to get our correct one
	var clientUpdates ClientUpdates
	if err := a.cdc.Unmarshal(specialTxBz, &clientUpdates); err != nil {
		ctx.Logger().Error("failed to unmarshal client updates", "error", err)
		return nil
	}

	for _, clientUpdate := range clientUpdates.ClientUpdates {
		if err := a.clientKeeper.UpdateClient(ctx, clientUpdate.ClientToUpdate, &clientUpdate.AttestationClaim); err != nil {
			ctx.Logger().Error("failed to update client", "error", err, "client_id", clientUpdate.ClientToUpdate)
		}

		ctx.Logger().Info("AttestationVoteExtension: PreBlocker (updated client)", "client_id", clientUpdate.ClientToUpdate)
	}

	return nil
}
