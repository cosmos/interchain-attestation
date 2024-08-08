package abci

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ve "vote-extensions.dev"
)

var _ ve.HasVoteExtension = AppModule{}

// TODO: Document
func (a AppModule) ExtendVote(ctx sdk.Context, vote *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
	// Ask sidecar for attestation and return vote extension
	// TODO: Check how slinky does sidecar communication

	ctx.Logger().Info("AttestationVoteExtension: ExtendVote (doing nothing atm)")

	// TODO implement me

	return &abci.ResponseExtendVote{}, nil
}

// TODO: Document
func (a AppModule) VerifyVote(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
	// Verify vote extension

	ctx.Logger().Info("AttestationVoteExtension: VerifyVote (doing nothing atm)")

	// TODO implement me

	return &abci.ResponseVerifyVoteExtension{}, nil
}

// TODO: Document
func (a AppModule) PrepareProposal(ctx sdk.Context, proposal *abci.RequestPrepareProposal, bytes []byte) (*abci.ResponsePrepareProposal, error) {
	// Extract vote extensions and add "fake tx" to the proposal

	ctx.Logger().Info("AttestationVoteExtension: PrepareProposal (doing nothing atm)")

	// TODO implement me

	return &abci.ResponsePrepareProposal{
		Txs: proposal.Txs,
	}, nil
}

// TODO: Document
func (a AppModule) ProcessProposal(ctx sdk.Context, proposal *abci.RequestProcessProposal, i int) (*abci.ResponseProcessProposal, error) {
	ctx.Logger().Info("AttestationVoteExtension: ProcessProposal (doing nothing atm)")

	// TODO implement me

	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
}

// TODO: Document
func (a AppModule) PreBlocker(ctx sdk.Context, block *abci.RequestFinalizeBlock, i int) error {
	// Extract "fake tx" and send an update client msg to the light client

	ctx.Logger().Info("AttestationVoteExtension: PreBlocker (doing nothing atm)")

	// TODO implement me

	return nil
}
