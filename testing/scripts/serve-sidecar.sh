#!/bin/bash
#set -e
set -eE -o functrace

failure() {
  local lineno=$1
  local msg=$2
  echo "Failed at $lineno: $msg"
}
trap 'failure ${LINENO} "$BASH_COMMAND"' ERR

source scripts/common-serve-env.sh

BINARY=attestation-sidecar
ROOT_DIR=/tmp
SIDECAR_DIR=$ROOT_DIR/sidecar
LOG_FILE_PATH=$SIDECAR_DIR/sidecar.log

# Stop if it is already running
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    pkill $BINARY
    sleep 5 # To avoid removing the folder to be any issue
fi

if [ -d $SIDECAR_DIR ]; then
  echo "Removing previous data..."
  rm -rf $SIDECAR_DIR &> /dev/null
fi

# Add directories for sidecar, exit if an error occurs
if ! mkdir -p $SIDECAR_DIR 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

echo "Creating sidecar config"
cat <<EOF > $SIDECAR_DIR/config.toml
attestator_id = 'your-attestator-id'

[[cosmos_chain]]
chain_id = 'rollupsimapp-1'
client_id = '07-tendermint-0' # TODO: REMOVE AFTER CREATE CLIENT ADDS THIS AUTOMATICALLY
attestation = true
client_to_update = '10-attestation-0'
rpc = 'http://localhost:36657'
address_prefix = 'rollup'
keyring_backend = 'test'
key_name = 'alice'
gas = 'auto'
gas_prices = '0.025stake'
gas_adjustment = 1.5

[[cosmos_chain]]
chain_id = 'simapp-1'
client_id = '10-attestation-0' # TODO: REMOVE AFTER CREATE CLIENT ADDS THIS AUTOMATICALLY
rpc = 'http://localhost:26657'
attestation = false
address_prefix = 'simapp'
keyring_backend = 'test'
key_name = 'alice'
gas = 'auto'
gas_prices = '0.025stake'
gas_adjustment = 1.5

EOF

echo "Setting up keys on sidecar keychain"
echo "$ALICE_MNEMONIC" | $BINARY keys add alice --recover --home $SIDECAR_DIR --keyring-backend test --address-prefix simapp

echo "Creating light clients"
$BINARY relayer create clients rollupsimapp-1 tendermint simapp-1 attestation --home $SIDECAR_DIR

echo "Starting sidecar"
$BINARY start --home $SIDECAR_DIR > $LOG_FILE_PATH 2>&1 &

sleep 3

echo "Creating connections"
$BINARY relayer create connections rollupsimapp-1 simapp-1 --verbose --home $SIDECAR_DIR

echo "Creating channels"
$BINARY relayer create channels rollupsimapp-1 connection-0 transfer ics20-1 simapp-1 connection-0 transfer --home $SIDECAR_DIR
