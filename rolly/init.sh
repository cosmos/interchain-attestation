#!/bin/sh

DA=$1
DA_ADDRESS=$2
DA_RPC_URL=$3

if [ -z "$DA" ]; then
  echo "DA not found (first arg should be celestia or mock). Exiting."
  exit 1
fi

if [ -z "$DA_ADDRESS" ]; then
  echo "DA_ADDRESS not found (second argument). Exiting."
  exit 1
fi

if [ "$DA" = "celestia" ]; then
  echo "DA is Celestia"
  if [ -z "$DA_RPC_URL" ]; then
    echo "DA_RPC_URL not found (third argument). Exiting."
    exit 1
  fi
elif [ "$DA" = "mock" ]; then
  echo "DA is mock"
else
  echo "DA not found (first arg should be celestia or mock). Exiting."
  exit 1
fi


# set variables for the chain
VALIDATOR_NAME=validator1
CHAIN_ID=rolly
KEY_NAME=rolly-key
KEY_2_NAME=rolly-key-2
CHAINFLAG="--chain-id ${CHAIN_ID}"
TOKEN_AMOUNT="10000000000000000000000000stake"
STAKING_AMOUNT="1000000000stake"
RELAYER_NAME=relayer
RELAYER_MNEMONIC="element achieve battle inject taxi hard purchase merit empower tower steak balance supreme purse assault lens chair dove together danger cat essence offer peace"

# build the rolly chain with Rollkit
ignite chain build

if [ "$DA" = "celestia" ]; then
  DA_BLOCK_HEIGHT=$(curl $DA_RPC_URL/block | jq -r '.result.block.header.height')
  echo -e "\n Your DA_BLOCK_HEIGHT is $DA_BLOCK_HEIGHT \n"

  AUTH_TOKEN=$(cat /root/shared-data/auth-token)
  echo -e "\n Your DA AUTH_TOKEN is $AUTH_TOKEN \n"
fi

# reset any existing genesis/chain data
#rollyd tendermint unsafe-reset-all
rm -rf ~/.rolly

# initialize the validator with the chain ID you set
rollyd init $VALIDATOR_NAME --chain-id $CHAIN_ID

# add keys for key 1 and key 2 to keyring-backend test
rollyd keys add $KEY_NAME --keyring-backend test
rollyd keys add $KEY_2_NAME --keyring-backend test
echo "$RELAYER_MNEMONIC" | rollyd keys add $RELAYER_NAME --recover --keyring-backend=test

# add these as genesis accounts
rollyd genesis add-genesis-account $KEY_NAME $TOKEN_AMOUNT --keyring-backend test
rollyd genesis add-genesis-account $KEY_2_NAME $TOKEN_AMOUNT --keyring-backend test
rollyd genesis add-genesis-account $RELAYER_NAME $TOKEN_AMOUNT --keyring-backend test

# set the staking amounts in the genesis transaction
rollyd genesis gentx $KEY_NAME $STAKING_AMOUNT --chain-id $CHAIN_ID --keyring-backend test

# collect genesis transactions
rollyd genesis collect-gentxs

ADDRESS=$(jq -r '.address' ~/.rolly/config/priv_validator_key.json)
PUB_KEY=$(jq -r '.pub_key' ~/.rolly/config/priv_validator_key.json)
jq --argjson pubKey "$PUB_KEY" '.consensus["validators"]=[{"address": "'$ADDRESS'", "pub_key": $pubKey, "power": "1000", "name": "Rollkit Sequencer"}]' ~/.rolly/config/genesis.json > temp.json && mv temp.json ~/.rolly/config/genesis.json

# allow all incoming WSS connections
sed -i '' 's/cors_allowed_origins = \[\]/cors_allowed_origins = ["*"]/g' ~/.rolly/config/config.toml

# start the chain

if [ "$DA" = "celestia" ]; then
  echo "Starting rollkit with Celestia DA..."
  rollyd start --rollkit.aggregator true --minimum-gas-prices 0stake --rollkit.da_start_height $DA_BLOCK_HEIGHT --rollkit.da_address $DA_ADDRESS --rollkit.da_auth_token $AUTH_TOKEN --rollkit.da_namespace 00000000000000000000000000000000000000000008e5f679bf7116cb --rpc.laddr tcp://0.0.0.0:27657 --p2p.laddr tcp://0.0.0.0:27656 --api.enable --api.enabled-unsafe-cors
elif [ "$DA" = "mock" ]; then
  rollyd start --rollkit.aggregator true --minimum-gas-prices 0stake --rollkit.da_address $DA_ADDRESS0 --rpc.laddr tcp://0.0.0.0:27657 --p2p.laddr tcp://0.0.0.0:27656 --api.enable --api.enabled-unsafe-cors
else
  echo "SHOULD NOT HAPPEN! Exiting.."
  exit 1
fi