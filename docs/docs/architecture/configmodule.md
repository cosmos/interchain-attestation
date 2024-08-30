---
sidebar_position: 1
---

# Attestation Config Module

The Attestation Config Module is a Cosmos SDK module that allows for configurable security requirements for the Interchain Attestation system.

It has the following responsibilities:
* Registering and keeping track of validator attestation signing keys
* Registering and keeping track of chains/clients that are registered for attestation with configuration details
* Providing the light client with the data and integration points to the chain it needs:
  * Validator signing keys
  * Validating required security backing of an attestation (i.e. is there enough voting power behind the attestation)
  * Any chain native capabilities like slashing or incentives

TODO: Add illustration with actors and interactions