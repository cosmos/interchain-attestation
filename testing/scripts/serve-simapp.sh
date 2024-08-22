#!/bin/bash
set -eE -o functrace

failure() {
  local lineno=$1
  local msg=$2
  echo "Failed at $lineno: $msg"
}
trap 'failure ${LINENO} "$BASH_COMMAND"' ERR

source scripts/common-serve-env.sh

BINARY=simappd
CHAIN_ID=simapp-1
ROOT_DIR=/tmp
CHAIN_DIR=$ROOT_DIR/$CHAIN_ID
LOG_FILE_PATH=$CHAIN_DIR/$CHAIN_ID.log

P2P_PORT=26656
RPC_PORT=26657
REST_PORT=1317
ROSETTA_PORT=8080
GRPC_PORT=8090

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

echo "Initializing $CHAIN_ID..."
$BINARY init test --home $CHAIN_DIR --chain-id=$CHAIN_ID --default-denom stake > /dev/null 2>&1

echo "Adding genesis accounts..."
echo "$ALICE_MNEMONIC" | $BINARY keys add alice --home $CHAIN_DIR --recover --keyring-backend=test > /dev/null 2>&1
echo "$BOB_MNEMONIC" | $BINARY keys add bob --home $CHAIN_DIR --recover --keyring-backend=test > /dev/null 2>&1
echo "$VALIDATOR_MNEMONIC" | $BINARY keys add validator --home $CHAIN_DIR --recover --keyring-backend=test > /dev/null 2>&1

$BINARY genesis add-genesis-account $($BINARY --home $CHAIN_DIR keys show alice --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY genesis add-genesis-account $($BINARY --home $CHAIN_DIR keys show bob --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR
$BINARY genesis add-genesis-account $($BINARY --home $CHAIN_DIR keys show validator --keyring-backend test -a) 100000000000stake  --home $CHAIN_DIR

$BINARY genesis gentx validator 7000000000stake --home $CHAIN_DIR --chain-id $CHAIN_ID --keyring-backend test > /dev/null 2>&1
$BINARY genesis collect-gentxs --home $CHAIN_DIR > /dev/null 2>&1

sed -i -e 's/stake/stake/g' $CHAIN_DIR/config/genesis.json
sed -i -e 's/"voting_period": "172800s"/"voting_period": "30s"/g' $CHAIN_DIR/config/genesis.json

echo "Changing defaults and ports in app.toml and config.toml files..."
sed -i -e 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2P_PORT"'"#g' $CHAIN_DIR/config/config.toml
sed -i -e 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPC_PORT"'"#g' $CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/config/config.toml
sed -i -e 's/enable = false/enable = true/g' $CHAIN_DIR/config/app.toml
sed -i -e 's/swagger = false/swagger = true/g' $CHAIN_DIR/config/app.toml
sed -i -e 's#"tcp://0.0.0.0:1317"#"tcp://0.0.0.0:'"$REST_PORT"'"#g' $CHAIN_DIR/config/app.toml
sed -i -e 's#":8080"#":'"$ROSETTA_PORT"'"#g' $CHAIN_DIR/config/app.toml
sed -i -e 's#"localhost:9090"#"0.0.0.0:'"$GRPC_PORT"'"#g' $CHAIN_DIR/config/app.toml
sed -i -e 's/enable-unsafe-cors = false/enable-unsafe-cors = true/g' $CHAIN_DIR/config/app.toml
sed -i -e 's/enabled-unsafe-cors = false/enable-unsafe-cors = true/g' $CHAIN_DIR/config/app.toml
sed -i.bak -e "s/^minimum-gas-prices *=.*/minimum-gas-prices = \"0.025stake\"/" $CHAIN_DIR/config/app.toml


echo "Starting $CHAIN_ID in $CHAIN_DIR..."
echo "Creating log file at $LOG_FILE_PATH"
ATTESTATION_SIDECAR_ADDRESS=localhost:6969 $BINARY start --home $CHAIN_DIR --pruning=nothing --rpc.unsafe --grpc.address="0.0.0.0:$GRPC_PORT" --state-sync.snapshot-interval 10 --state-sync.snapshot-keep-recent 2 > $LOG_FILE_PATH 2>&1 &

sleep 3

echo ""
echo "----------- Config -------------"
echo "RPC: tcp://0.0.0.0:$RPC_PORT"
echo "REST: tcp://0.0.0.0:$REST_PORT"
echo "GRPC: tcp://0.0.0.0:$GRPC_PORT"
echo "chain-id: $CHAIN_ID"
echo ""

if ! $BINARY --home $CHAIN_DIR --node tcp://:$RPC_PORT status; then
  echo "Chain failed to start"
  exit 1
fi

echo "-------- Chain started! --------"
