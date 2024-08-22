# Attestation Sidecar

The attestation sidecar is a process that runs alongside a chain with Interchain Attestation and pulls in information
that the validator can use to update the attestation light client (with ABCI++).

TODO: Add more information about the sidecar

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

