# Interchain Attestation

![Interchain Attestion Logo.png](docs/images/Interchain%20Attestion%20Logo.png)

(Previously known as Pessimistic Validation)

Interchain Attestation is a project to enable IBC everywhere. In particular, IBC for chains that can't/don't have a 
light client implementation that can be used with IBC. It enables IBC for any chain that can implement IBC, and
let another chain safely validate it (for instance by running a full node). This then includes optimistic rollups,
Ethereum, Solana, and more.

The project is partially funded by the Dorahacks ATOM Economic Zone Quadratic Grant rounds.
You can find project information and contribute to the project here: https://dorahacks.io/aez

## Current status
The project is under development and is not yet ready for production use.

For a more detailed roadmap, see the [GitHub project board](https://github.com/orgs/cosmos/projects/35)

## Documentation

The documentation contains information about the project, how to build and test it, as well as how to use it.

You can find the documentation here: [Interchain Attestation Documentation](https://interchain-attestation.io)
 
## Background

This project was originally built by Gjermund Garaba (https://github.com/gjermundgaraba/, https://twitter.com/gjermundgaraba)
for Celestia's Infinite Space Bazaar to solve the problem of waiting for the dispute period to pass when bridging assets from an optimistic rollup to a receiving chain.

You can find the original working proof of concept code here:
https://github.com/cosmos/interchain-attestation/tree/9bc691c585697921b84c5467b13996389e6d119f