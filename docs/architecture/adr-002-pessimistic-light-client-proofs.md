# ADR 002: Pessimistic Light Client Proofs

## Changelog

* {date}: {changelog}

## Status

DRAFT

## Abstract

This ADR proposes a method for generating and verifying proofs for optimistic rollups using a combination of 
consensus state and validator signatures. By employing a conditional light client setup, 
a standard light client updates its state only when a prover light client, run by validators, 
confirms the rollup's state and height. 

The prover light client generates a proof structure that includes the consensus state and validator signatures, 
which the standard light client then uses for state updates. This solution aims to enhance security and 
interoperability of optimistic rollups within the Cosmos ecosystem until more advanced state proof mechanisms are available.

## Context

Optimistic rollups have an inherent problem with interoperability, since the state of a rollup
is not considered final/safe until a dispute window (a period of time where someone can
provide a fraud proof) - of typically 7 days - has passed.

This ADR proposes a solution for verifying the behavior of optimistic rollups by 
enabling Proof of Stake Cosmos SDK-based blockchains to verify rollup states using validators 
who run full nodes of the rollup.

The light client setup itself is document in ADR-001, but for context, we will briefly the concept:
A standard light client with an additional "conditional" prover light client that contains the pessimistic validation.
By conditional, we mean that the light client will only update its state (`UpdateClient`) if the other light client is in agreement on the height and state
The prover light client will use the method specified in this ADR to create the proofs that the standard light client will use to update its state.

## Alternatives

TODO: Write an alternative

> This section describes alternative designs to the chosen design. This section
> is important and if an adr does not have any alternatives then it should be
> considered that the ADR was not thought through.

## Decision

We will use a combination of consensus state (details below) and signatures from validators to create a proof that 
a height and state has been observed and signed off on by the validators. The validators will run a full node of the rollup
and report the state and height back to the light client (method (ABCI++ and sidecar) described in a separate ADR, which in
turn will be queried by any light client that wants to verify the rollup state against the pessimistic validation.

The consensus state could be a different set of data for different rollups, but we will use the tendermint-07 consensus state
as the initial implementation and expend with different types of consensus states as needed.

This is the generated go struct for the consensus state:
```go
type ConsensusState struct {
	// timestamp that corresponds to the block height in which the ConsensusState
	// was stored.
	Timestamp time.Time
	// commitment root (i.e app hash)
	Root               types1.MerkleRoot                                `protobuf:"bytes,2,opt,name=root,proto3" json:"root"`
	NextValidatorsHash github_com_cometbft_cometbft_libs_bytes.HexBytes `protobuf:"bytes,3,opt,name=next_validators_hash,json=nextValidatorsHash,proto3,casttype=github.com/cometbft/cometbft/libs/bytes.HexBytes" json:"next_validators_hash,omitempty"`
}
```

We will use some variation of this struct to create the consensus state for the rollup and wrap it in signatures from the validators.
It contains all we need for a conditional tendermint light client to verify the state of the rollup. 

```go
type Proof struct {
    Height     uint64 // Or exported.Height
    ConsensusState ConsensusState // Or maybe just bytes to make it easy to generalize for other types later
    Signatures []Signature
}

type Signature struct {
    ValidatorAddress sdk.Address // Or just string, we'll see
    Power            uint64
    Signature        []byte
```

Something like this. The `Proof` struct will be created by the prover light client and sent to the standard light client, 
which will use it to update its state.

### Backwards Compatibility

Not applicable, since Pessimistic Validation is not in production yet.

## References