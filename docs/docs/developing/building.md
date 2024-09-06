---
sidebar_position: 2
---

# Building

## Building all the modules

To build all the modules, you can use the following command:

```bash
$ just build
```

## Linting

To lint all the modules, you can use the following command:

```bash
$ just lint
```

## Proto

To generate go code for all the proto files, you can use the following command:

```bash
$ just proto-gen
```

## Docker images (used for testing)

There is a set of docker images used for the e2e tests. To build these images, you can use the following command:

```bash
$ just build-docker-images
```

## Simapps

If you want to build and install the simapp binaries locally, you can use the following command:

```bash
just install-simapps
```

## Running locally

There is a command that spins up a local test environment with the following components:
* Simapp (Cosmos SDK chain with Interchain Attestion integrated)
* Rollupsimapp (Rollkit rollup)
* Mock DA service (for the rollup)
* Sidecar
* Configuration
  * Light clients
  * Connections
  * Channels
  * Validator registered and wired up to Interchain Attestation

To run this environment, you can use the following command:
```bash
just serve
```
