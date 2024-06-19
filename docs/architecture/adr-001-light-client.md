# ADR 001: Light Client

## Changelog

* {date}: Initial draft

## Status

DRAFT

## Abstract

TODO: Write the abstract

> "If you can't explain it simply, you don't understand it well enough." Provide
> a simplified and layman-accessible explanation of the ADR.
> A short (~200 word) description of the issue being addressed.

## Context

The modular thesis is expected to bring many new rollups to the Cosmos ecosystem.
Until state proofs (ZK) are a viable option, we will see many (most?) of the rollups
using so-called optimistic rollups. State proofs are undoubtedly the future for rollups,
but even with all the advancements in ZK tech, provable state transitions will take time
to be available everywhere.

Optimistic rollups have an inherent problem with interoperability, since the state of a rollup
is not considered final/safe until a dispute window (a period of time where someone can 
provide a fraud proof) - of typically 7 days - has passed.

The Pessimistic Rollup project was conceived to solve this problem by allowing any Proof of Stake
Cosmos SDK-based blockchain to leverage its economic security and verify that a given rollup
is behaving correctly. The validators do this by running a full node of the rollup and report
the state and height back to its own chain. With this information a chain can safely accept and verify
IBC packets from the rollup.

The initial prototype, developed during the Celestia Infinite Space Bazaar hackathon, used
two light clients (standard tendermint + a custom light client updated with ABCI++ Vote Extensions)
and a custom Cosmos SDK module to achieve Pessimistic Validation for a rollup. The design was
functional, but had some downsides that makes adjustments necessary. In particular, the design
required a standard tendermint light client to be running, which could be used directly by anyone.
While the denoms coming across a connection with that light client would not be canonical, it was
an unnecessary confusing to have an operational light client that should not be used. It could
have been solved by requiring an IBC middleware to stop packets or channel creation, but this
would require more complexity for integrators and the protocol itself.

## Alternatives

TODO: Move the chosen alternative down to decision.

### IBC Middleware

TODO: Add diagram

### Conditional Tendermint Light Client

TODO: Add diagram

## Decision

> This section describes our response to these forces. It is stated in full
> sentences, with active voice. "We will ..."
> {decision body}

## Consequences

> This section describes the resulting context, after applying the decision. All
> consequences should be listed here, not just the "positive" ones. A particular
> decision may have positive, negative, and neutral consequences, but all of them
> affect the team and project in the future.

### Backwards Compatibility

Not applicable, since Pessimistic Validation is not in production yet.

### Positive

> {positive consequences}

### Negative

> {negative consequences}

### Neutral

> {neutral consequences}

## Further Discussions

> While an ADR is in the DRAFT or PROPOSED stage, this section should contain a
> summary of issues to be solved in future iterations (usually referencing comments
> from a pull-request discussion).
> 
> Later, this section can optionally list ideas or improvements the author or
> reviewers found during the analysis of this ADR.

## References

* {reference link}
