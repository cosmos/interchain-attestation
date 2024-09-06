---
sidebar_position: 3
---

# Attestation Sidecar

The Attestation Sidecar is a process that runs alongside a chain node, and is responsible for generating and serving attestations to the chain node. 
The sidecar is used by the chain node to fetch attestations during vote extensions.

The sidecar is configured by the validator to connect to counterparty chains and will run in the background to continuously fetch and generate attestations.

The sidecar is implemented as a CLI binary that can perform multiple tasks, such as setting up signing keys, creating validator registration files, and starting the attestation process.

In addition, the sidecar also enables relaying of IBC packets (exact capabilities TBD).

The Attestation Sidecar is responsible for the following:
* Fetching attestation data from the counterparty chain
  * Packet commitments
  * Header information
* Signing attestations
* Serving attestations to the chain node via a GRPC server

## Configuration

TODO: Document the configuration (or should it be under a separate "usage" section of some kind?)

## Relaying

Currently, the sidecar has only one-off commands for creating clients, connections and channels, but the plan is to enable the sidecar to relay IBC packets
as part of its running process.

## CLI

TODO: Document the commands