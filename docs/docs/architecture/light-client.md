---
sidebar_position: 5
---

# Attestation Light Client

The Attestation Light Client is an ibc-go light client implementation that uses attestation data to verify validator consensus
and update the client state and consensus state with packet commitments and header information.

TODO: Briefly explain the light client interface from ibc-go and how it works (including consensus state)

TODO: Explain packet commitments

TODO: Add illustration

## Consensus State

TODO: Document how packet commitments and stuff are stored in the consensus state.

## Verify Membership

The light client uses packet commitments that it receiver in the attestation data to verify that the packet was included in the chain.

## Verify non-membership

