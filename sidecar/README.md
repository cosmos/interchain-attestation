# Prover Sidecar

The prover sidecar is a service that runs alongside a chain with pessimistic-validation and pulls in information
that the validator can use to update the pessimistic light client (with ABCI++).

TODO: Add more information on how the prover sidecar works and what it does.

## Install

```bash
$ make install
```

## Start the sidecar

```bash
$ attestation-sidecar start
```

If you start the sidecar without a config file existing, it will generate one for you in `~/.attestation-sidecar/config.toml` 
with an example chain configuration.

TODO: Show example

### Flags

`--listen-addr` - the address for the grpc server to listen on. Defaults to localhost:6969

```bash
$ attestation-sidecar start --listen-addr 0.0.0.0:1337
```

## Configuration

TODO: Document each field in the config file once it's more stable.

