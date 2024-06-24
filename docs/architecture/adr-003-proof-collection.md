# ADR 003: Proof Collection

## Changelog

* {date}: {changelog}

## Status

DRAFT

## Abstract

TODO: Write abstract

> "If you can't explain it simply, you don't understand it well enough." Provide
> a simplified and layman-accessible explanation of the ADR.
> A short (~200 word) description of the issue being addressed.

## Context

Validators who are running full nodes (pessimistic validation) or otherwise need to 
collect proof for a Prover light client, need an efficient and asynchronous process
that does not interfere with the normal validator process. 

TODO: write more?

## Alternatives

### Proof collection in-binary

It would be possible for pessimistic validation to gather proofs directly from the same
process that the node is already running on. It could be done in an async matter that
would not interfere with normal runnings of the blockchain.

The only real benefit for doing things this way would be to avoid having validators run 
an additional binary. 

The downside is that updates to the proof collection code would have to built and distributed
by the maintaining team. This makes it much harder to provide updates or otherwise customize
the collection process.

## Decision

We will implement proof collection by creating a sidecar process that can be queried by the node
binary through GRPC. The sidecar process will collect proofs asynchronously and when the main node
process is ready for any updates it will be able to provide what it has already processed.

TODO: Explain in more detail the flow and the actual proof collection of pessimistic validation

## Consequences

> This section describes the resulting context, after applying the decision. All
> consequences should be listed here, not just the "positive" ones. A particular
> decision may have positive, negative, and neutral consequences, but all of them
> affect the team and project in the future.

### Backwards Compatibility

N/A

## References

