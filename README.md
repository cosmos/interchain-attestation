# Pessimistic Validation

Main goal: to allow optimistic rollups to IBC without having to wait for the dispute period.

Possible solutions to this problem:
1. Allow anyone to give the funds (minus a fee) to the user on the other side of the bridge and redirecting the waiting IBC message to the person who paid the escrow.
2. Let someone you trust (your validators) to pessimistically validate the rollup and sign off on heights we trust.

This repo is a proof of concept for the second solution.

## Implementation

I have two ideas on how to do this:
1. Have a separate light client that validators use to sign off on the heights of a another light client.
2. Extend IBC with a light client middleware that can block packets for heights not signed off by the validators.
   - Potentially: some kind of middleware that just blocks the actual packets, but allows light client updates.

## Economic models and considerations
I've not implemented any economic models in this proof of concept, but I have some ideas:
* Whoever creates a validation objective can set a bounty for validators to claim
* Leverage ibc fees somehow?
* ...

Incentives to behave:
* Validators can be slashed if they sign off on a rollup height that is later proven to be invalid

## Work to be done

- [x] Create validation objective
- [x] Signing up for a validation objective
- [ ] Implement a simple "sign-off" light client that validators can use to sign off on rollup heights.
- [ ] Starting a validation objective (including creating light client)
- [ ] Implement VoteExtension to update light client and relay packets

## Some issues I'm aware of
* Rotated validator keys are not handled