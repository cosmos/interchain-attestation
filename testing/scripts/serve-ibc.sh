#!/bin/bash

source scripts/common-serve-env.sh

echo "Adding hermes keys, an error is expected if the key already exists, the script will not fail"
echo "$ALICE_MNEMONIC" | hermes --config ./scripts/config.toml keys add --chain simapp-1 --mnemonic-file /dev/stdin
echo "$ALICE_MNEMONIC" | hermes --config ./scripts/config.toml keys add --chain rollupsimapp-1 --mnemonic-file /dev/stdin

set -eE -o functrace

failure() {
  local lineno=$1
  local msg=$2
  echo "Failed at $lineno: $msg"
}
trap 'failure ${LINENO} "$BASH_COMMAND"' ERR

#hermes --config ./scripts/config.toml create client --host-chain simapp-1 --reference-chain rollupsimapp-1
#hermes --config ./scripts/config.toml create client --host-chain rollupsimapp-1 --reference-chain simapp-1
hermes --config ./scripts/config.toml create connection --a-chain simapp-1 --b-chain rollupsimapp-1
hermes --config ./scripts/config.toml create channel --a-chain simapp-1 --a-connection connection-0 --a-port transfer --b-port transfer