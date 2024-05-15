# Pessimistic Validation

Main goal: to allow optimistic rollups to IBC without having to wait for the dispute period.

The problem of bridging in reasonable time from optimistic rollups are required to be solved in order for them to be useful
in an interchain context where token bridging over IBC is one of the main use cases.

This repo is a proof of concept for the second solution: pessimistic validation.

## Pessimistic Validation

Pessimistic validation is where you validate all the blocks on the whole rollup chain, rather than optimistic trust or waiting for dispute periods.
This is essentially the way most full L1 blockchains works today, but with a twist: the rollup can still run as an optimistic rollup and the receiver chain can
pessimistically validate the rollup chain with a partial set from their own validators.

For an optimistic rollup, to bridge assets you normally need to wait for the dispute period to pass before you can trust the rollup
and the assets bridged over. This is because the rollup chain can be rolled back if a fraud proof is submitted within the dispute period.

Anyone can however validate the rollup themselves with a full node, but it doesn't help the receiving chain unless it can either validate it itself or a have a trusted party validate it.

The goal of pessimistic validation is to allow a receving chain to have a partial set of its validators validate the rollup chain and sign off on the heights they trust. 
This way the receiving chain can trust the rollup chain without having to wait for the dispute period to pass, and assets can be bridged over quickly.

![Pessimistic validation overview](pessimistic-validation-overview.png)

## Implementation

The implementation of pessimist validation is done by combining the following:
* A custom SDK module that has "validation objectives" where validators can sign up to validate a rollup chain.
* A custom light client that can be updated with new rollup heights that have been signed off by the validators.
* Implementation of ABCI++ VoteExtensions to allow the light client to be updated with new rollup heights easily (and quickly).

The system sketch below shows how the system works:
1. A new validation objective is created on the receiving chain with a given "required consensus power" requirement
    - The objective also needs a "dependent light client" that is used for the pessimistic light client to prove memberships and stuff like that
2. Validators sign up for the validation objective
3. Once enough validator power has signed up, the validation objective is started and a new pessimistic light client is created
4. The validators set up a config for reading the headers from a full node they control
5. The validators automatically read headers from their full node and sign off on the height by putting them in an ABCI++ VoteExtensions message
6. VoteExtensions are validated (this part is not implemented much in this PoC, just some basic stuff)
7. VoteExtensions are sent to the Pessimistic Light Client in the form of an UpdateClient with a ClientMessage
8. Connections, channels and packets can now be set up between the rollup chain and the receiving chain

> Side note: The VoteExtension implementation could easily be extended to also perform all the relaying on both chains directly.

![Pessimistic validation system sketch](pessimistic-validation-system-sketch.png)

The whole flow can also be seen in the E2E tests in the `interchaintest` folder.

The following folders are included in this repo:
* `hub/` - Standard Cosmos SDK chain with
* `hub/app/app.go` - ABCI++ VoteExtensions implementation
* `hub/x/pessimist` - Custom SDK module for validation objectives
* `hub/x/pessimist/lightclient` - Custom light client for pessimistic validation
* `rolly/` - A basic Rollkit chain that is used to test the solution
* `mock-da/` - Just a dockerfile for a mock data availability layer used in the tests
* `interchaintest/` - E2E tests for the solution

Some other repos are used where custom forks are made:
* https://github.com/gjermundgaraba/interchaintest/tree/rollkit-celestia-example
* https://github.com/gjermundgaraba/relayer/tree/pessimistic-rollkit
* https://github.com/gjermundgaraba/mock-da

And some custom docker builds are used, which can be found here: https://github.com/gjermundgaraba?tab=packages

In addition to the above we are also using a main-line commit of ibc-go that has the new LightClientModule interface (see hub go.mod).

### Shortcomings of the solution
* The dependent light client could still be used to create channels and send packets, but would not have the same security guarantees as the pessimistic validation light client.
  * This could be solved by adding a simple IBC middleware that just implements the channel handshake callbacks to stop them from being created. Or just drop packets on that client.
* Requires a new channel between source light client and the new pessimistic validation light client.
  * Might not be an issue, but if there are existing channels this could be a problem.
  * Could be solved by IBC adding support for dependent light clients and therefore flipping the dependency relationship.
* Requires a new SDK module and light client to be added to the chain.
  * Could be solved by rewriting the light client in CosmWasm for 08-WASM and implement the module as a CosmWasm contract as well (or maybe it could be embedded in 08-WASM client contract?).

## Economic models and considerations
I've not implemented any economic models in this proof of concept, but I have some ideas:
* Whoever creates a validation objective can set a bounty for validators to claim
* Leverage ibc fees somehow?
* ...

Incentives to behave:
* Validators can be slashed if they sign off on a rollup height that is later proven to be invalid

## Alternative implementation using IBC middleware
I also considered another alternative implementation that does not involve a new light client,
but instead uses an IBC middleware that can block packets for heights not signed off by the validators.

![Alternative IBC Middleware implementation sketch](alternative-middleware-implementation.png)

The drawback of this solution is that it requires all applications routes to be wrapped with the middleware and might not work for all light clients.
It would however be a simpler and perhaps more straight-forward solution. **It could also be combined with an eIBC/escrow type of solution.**

## Work to be done

- [x] Create validation objective
- [x] Signing up for a validation objective
- [x] Implement a simple "sign-off" light client that validators can use to sign off on rollup heights.
- [x] Starting a validation objective (including creating light client)
- [x] E2E Test setup
- [x] Implement VoteExtension to update light client
- [x] Test the solution end-to-end
- [x] Document the solution, shortcomings and hacks

Some next steps:
- [ ] Add back signature verification in the light client for the VE commitments
- [ ] Add tests for less happy paths

## Some issues I'm aware of
* Rotated validator keys are not handled
* Creating validation objective is permissionless, should maybe be done by gov/authority
* Not sure how well shifting validator power would work. Might need some kind of margin before starting the light client.
* No slashing or payments implemented
* No way to back out of a validation objective
* Test coverage is shit
* Lots of hacks

## Another potential solution to the problem

Another solution that could work (besides changing to ZK rollups) is that which has been proposed by Dymension 
in their [eIBC paper](https://eibc.dymension.xyz) (and implementation) where incoming token packets are first put
on hold (not processed) until either the fraud detection window passes (7 days or whatever) or someone pays the 
escrow (minus a fee) giving the user the funds right away (again, minus the fee) and redirects the waiting IBC packet to
the payer of the escrow. 

This is also essentially a form of pessimistic validation, but requires only the escrow payer to validate the rollup.

1. Allow anyone to give the funds (minus a fee) to the user on the other side of the bridge and redirecting the waiting IBC message to the person who paid the escrow.
2. Let someone you trust (your validators) to pessimistically validate the rollup and sign off on heights we trust.