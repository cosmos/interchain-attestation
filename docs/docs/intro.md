---
sidebar_position: 1
---

# Introduction

Interchain Attestation is a project to enable IBC everywhere. In particular, IBC for chains that can't/don't have a
light client implementation that can be used with IBC. It enables IBC for any chain that can implement IBC, and
let another chain safely validate it (for instance by running a full node). This then includes optimistic rollups,
Ethereum, Solana, and more.

The system enables any chain to connect with IBC, as long as it can implement the IBC protocol (e.g. smart contracts),
and the validators using Interchain Attestation are attesting to the state of the counterparty IBC implementation.

![Attestation enables IBC](../static/img/attestation-enables-ibc.png)