# Interchain Attestation

(Previously known as Pessimistic Validation)

Interchain Attestation is a project to enable IBC everywhere. In particular, IBC for chains that can't/don't have a 
light client implementation that can be used with IBC. It enables IBC for any chain that can implement IBC, and
let another chain safely validate it (for instance by running a full node). This then includes optimistic rollups,
Ethereum, Solana, and more.

The project is partially funded by the Dorahacks ATOM Economic Zone Quadratic Grant rounds.
You can find project information and contribute to the project here: https://dorahacks.io/aez


## Current status
The project is under development and is not yet ready for production use.

## Interchain Attestation

Interchain Attestation enables IBC connectivity (with no intermediary chains) with any chain where you can't/don't have a light client.

Interchain Attestation solves the problem where you can't, for whatever reason, trust the counterparty with a "normal" light client. 
Instead, it allows a chain's validators to attest to the state of the counterparty - moving the security to someone you already trust.

![No IBC for you.png](docs/images/No%20IBC%20for%20you.png)

This enables any chain to connect with IBC, as long as it can implement the IBC protocol (e.g. smart contracts), 
and the validators using Interchain Attestation are attesting to the state of the counterparty IBC implementation.

![Attestion enables IBC.png](docs/images/Attestion%20enables%20IBC.png)

Interchain Attestion is based on using validators with existing economic security to attest to the state of the counterparty chain.
We move the security assumption over to the receiving validator set (e.g. Cosmos Hub/Osmosis/whatever), away from the one we can't trust (like a single sequencer).

![Attest.png](docs/images/Attest.png)

In addition, the Attestation light client is based on attesting to IBC packets, rather than full state.
This makes it much easier to implement new chains and support consensus algorithms that doesn't (yet?) have a light client implementation.

The Attestation light client verifies the signatures of the attestors (validators) and stores the packet commitments to be able to verify the packet later.

![Packet commitments.png](docs/images/Packet%20commitments.png)

Not all scenarios warrant all validators running a full node and attesting to every chain it connects to, 
so Interchain Attestation also has a config module that allows for configurable security requirements.

The architecture is using a combination of a sidecar process, ABCI++ Vote Extensions, and a light client to enable the attestation process.

![High level architecture.png](docs/images/High%20level%20architecture.png)

A talk about the project can be found here: https://www.youtube.com/watch?v=loNyUjSgR8M

## Background

This project was originally built for Celestia's Infinite Space Bazaar to solve the problem of 
waiting for the dispute period to pass when bridging assets from an optimistic rollup to a receiving chain.

You can find the original working proof of concept code here:
https://github.com/gjermundgaraba/pessimistic-validation/tree/9bc691c585697921b84c5467b13996389e6d119f