#!/bin/bash

set -eE -o functrace

failure() {
  local lineno=$1
  local msg=$2
  echo "Failed at $lineno: $msg"
}
trap 'failure ${LINENO} "$BASH_COMMAND"' ERR

BINARY=rollupsimappd
CHAIN_ID=rollupsimapp-1
ROOT_DIR=/tmp
CHAIN_DIR=$ROOT_DIR/$CHAIN_ID
LOG_FILE_PATH=$CHAIN_DIR/$CHAIN_ID.log
VALIDATOR_NAME=validator1

P2P_PORT=36656
RPC_PORT=36657
REST_PORT=2317
ROSETTA_PORT=9080

ALICE_MNEMONIC="clock post desk civil pottery foster expand merit dash seminar song memory figure uniform spice circle try happy obvious trash crime hybrid hood cushion"
BOB_MNEMONIC="angry twist harsh drastic left brass behave host shove marriage fall update business leg direct reward object ugly security warm tuna model broccoli choice"
VALIDATOR_MNEMONIC="banner spread envelope side kite person disagree path silver will brother under couch edit food venture squirrel civil budget number acquire point work mass"

# Stop if it is already running
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    pkill $BINARY
    sleep 5 # To avoid removing the folder to be any issue
fi

if [ -d $CHAIN_DIR ]; then
  echo "Removing previous data..."
  rm -rf $CHAIN_DIR &> /dev/null
fi

# Add directories for chain(s), exit if an error occurs
if ! mkdir -p $CHAIN_DIR 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

# Rollkit specific
# Mocha
DA_BLOCK_HEIGHT=$(curl public-celestia-mocha4-consensus.numia.xyz:26657/block |jq -r '.result.block.header.height')
AUTH_TOKEN=$(celestia light auth write --p2p.network mocha)
# Arabica
#DA_BLOCK_HEIGHT=$(curl https://rpc.celestia-arabica-11.com/block |jq -r '.result.block.header.height')
#AUTH_TOKEN=$(celestia light auth write --p2p.network arabica)
echo -e "\n Your DA_BLOCK_HEIGHT is $DA_BLOCK_HEIGHT \n"
echo -e "\n Your DA AUTH_TOKEN is $AUTH_TOKEN \n"

echo "Initializing $CHAIN_ID..."
$BINARY init $VALIDATOR_NAME --home $CHAIN_DIR --chain-id=$CHAIN_ID --default-denom stake > /dev/null 2>&1

echo "Adding genesis accounts..."
echo "$ALICE_MNEMONIC" | $BINARY keys add alice --home $CHAIN_DIR --recover --keyring-backend=test > /dev/null 2>&1
echo "$BOB_MNEMONIC" | $BINARY keys add bob --home $CHAIN_DIR --recover --keyring-backend=test > /dev/null 2>&1
echo "$VALIDATOR_MNEMONIC" | $BINARY keys add $VALIDATOR_NAME --home $CHAIN_DIR --recover --keyring-backend=test > /dev/null 2>&1

$BINARY genesis add-genesis-account $($BINARY --home $CHAIN_DIR keys show alice --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY genesis add-genesis-account $($BINARY --home $CHAIN_DIR keys show bob --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY genesis add-genesis-account $($BINARY --home $CHAIN_DIR keys show $VALIDATOR_NAME --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR

$BINARY genesis gentx $VALIDATOR_NAME 7000000000stake --home $CHAIN_DIR --chain-id $CHAIN_ID --keyring-backend test > /dev/null 2>&1
$BINARY genesis collect-gentxs --home $CHAIN_DIR > /dev/null 2>&1

sed -i -e 's/"voting_period": "172800s"/"voting_period": "30s"/g' $CHAIN_DIR/config/genesis.json
sed -i -e 's/"expedited_voting_period": "86400s"/"expedited_voting_period": "20s"/g' $CHAIN_DIR/config/genesis.json

echo "Changing defaults and ports in app.toml and config.toml files..."
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2P_PORT"'"#g' $CHAIN_DIR/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPC_PORT"'"#g' $CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/config/config.toml
sed -i -e 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$REST_PORT"'"#g' $CHAIN_DIR/config/app.toml
sed -i -e 's#":8080"#":'"$ROSETTA_PORT"'"#g' $CHAIN_DIR/config/app.toml
sed -i -e 's/enable-unsafe-cors = false/enable-unsafe-cors = true/g' $CHAIN_DIR/config/app.toml
sed -i -e 's/enabled-unsafe-cors = false/enable-unsafe-cors = true/g' $CHAIN_DIR/config/app.toml
sed -i.bak -e "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"0.025stake\"/" $CHAIN_DIR/config/app.toml

# Rollkit specific
# copy centralized sequencer address into genesis.json
# Note: validator and sequencer are used interchangeably here
ADDRESS=$(jq -r '.address' $CHAIN_DIR/config/priv_validator_key.json)
PUB_KEY=$(jq -r '.pub_key' $CHAIN_DIR/config/priv_validator_key.json)
jq --argjson pubKey "$PUB_KEY" '.consensus["validators"]=[{"address": "'$ADDRESS'", "pub_key": $pubKey, "power": "1", "name": "'$VALIDATOR_NAME'"}]' $CHAIN_DIR/config/genesis.json > temp.json && mv temp.json $CHAIN_DIR/config/genesis.json
PUB_KEY_VALUE=$(jq -r '.pub_key .value' $CHAIN_DIR/config/priv_validator_key.json)
jq --arg pubKey $PUB_KEY_VALUE '.app_state .sequencer["sequencers"]=[{"name": "'$VALIDATOR_NAME'", "consensus_pubkey": {"@type": "/cosmos.crypto.ed25519.PubKey","key":$pubKey}}]' $CHAIN_DIR/config/genesis.json >temp.json && mv temp.json $CHAIN_DIR/config/genesis.json

echo "Starting $CHAIN_ID in $CHAIN_DIR..."
echo "Creating log file at $LOG_FILE_PATH"

$BINARY genesis validate --home $CHAIN_DIR
$BINARY start --home $CHAIN_DIR --rollkit.aggregator --rollkit.da_auth_token=$AUTH_TOKEN --rollkit.da_namespace 00000000000000000000000000000000000000000008e5f679bf7116cb --rollkit.da_start_height $DA_BLOCK_HEIGHT --minimum-gas-prices="0stake" --api.enable --api.enabled-unsafe-cors > $LOG_FILE_PATH 2>&1 &
# $BINARY start --log_format json --home $CHAIN_DIR --pruning=nothing --rpc.unsafe --grpc.address="0.0.0.0:$GRPC_PORT" --state-sync.snapshot-interval 10 --state-sync.snapshot-keep-recent 2 > $LOG_FILE_PATH 2>&1 &

sleep 3

echo ""
echo "----------- Config -------------"
echo "RPC: tcp://0.0.0.0:$RPC_PORT"
echo "REST: tcp://0.0.0.0:$REST_PORT"
echo "chain-id: $CHAIN_ID"
echo ""

if ! $BINARY --home $CHAIN_DIR --node tcp://:$RPC_PORT status; then
  echo "Chain failed to start"
  exit 1
fi

echo "-------- Chain started! --------"
