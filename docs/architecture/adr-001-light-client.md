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

A simpler solution to the problem (in terms of moving pieces at least) would be to not have any 
new light clients, and instead use IBC Middleware to filter packets.

IBC Middleware lives between the core IBC modules (client, channel, connection) and the application
modules (ICS20, ICS721, etc). It allows you to manipulate or even block packets as they are coming in.

If Pessimistic Validation was implemented as an IBC middleware, you could use existing light clients
(such as 07-tendermint) without any modifications. The middleware would keep track of the validation 
updates of the counterparty rollup, and only let through packets that are less than or equal to the
latest height proven to be safe.

An IBC middleware solution could implement a fee-taking system, where a small fee is redirected from ICS20 packets.

While this solution is materially simpler than any that involves custom light clients, it also has some
drawbacks that are important to note:

TODO: Look more into IBC fees
Potential issue that is not confirmed: it might be harder to force native IBC relaying fees. 

Relayers and clients would not be able to check that a light client has been proven by the pessimistic validators,
because all of that would happen in the middleware. If you simply queried the light client you would think
that a packet can be proven.


### Conditional Tendermint Light Client

TODO: Add diagram

The solution with a conditional tendermint light client is based on a new feature that landed in ibc-go
recently. It allows a light client to query another light client for 

Using a conditional tendermint light client we can avoid most of the issues from the initial prototype:
- 

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

* [https://github.com/cosmos/ibc-go/issues/5112](Conditional clients ibc-go issue)
