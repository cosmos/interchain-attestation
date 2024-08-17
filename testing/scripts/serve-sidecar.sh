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
private_key_path = 'will-be-updated'

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

echo "Creating signing key"
$BINARY signing-keys create --home $SIDECAR_DIR

echo "Setting up keys on sidecar keychain"
echo "$ALICE_MNEMONIC" | $BINARY keys add alice --recover --home $SIDECAR_DIR --keyring-backend test --address-prefix simapp

echo "Creating registration json"
$BINARY generate-register-attestator-json --home $SIDECAR_DIR

echo "Registering attestator"
TX_HASH=$(simappd tx attestationconfig register-attestator register-attestator.json --from validator --chain-id simapp-1 --keyring-backend test --home /tmp/simapp-1 --gas auto --gas-adjustment 1.5 --gas-prices 0.025stake --output json -y | jq -r ".txhash")
sleep 5
RES_CODE=$(simappd q tx $TX_HASH --output json | jq -r ".code")
if [ "$RES_CODE" != "0" ]; then
  echo "Error: Attestator registration failed: $RES_CODE"
  exit 1
fi

echo "Creating light clients"
$BINARY relayer create clients rollupsimapp-1 tendermint simapp-1 attestation --home $SIDECAR_DIR

echo "Starting sidecar"
$BINARY start --home $SIDECAR_DIR # TODO: Background and log to file

#echo "Creating connections"
echo $BINARY relayer create connections rollupsimapp-1 simapp-1 --home $SIDECAR_DIR