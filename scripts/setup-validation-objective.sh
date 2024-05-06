#!/bin/bash

hubd tx pessimist create-validation-objective 07-tendermint-0 1 --from bob --chain-id hub
sleep 1
hubd tx pessimist sign-up-for-objective 07-tendermint-0 --from alice --chain-id hub