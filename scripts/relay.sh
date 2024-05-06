#!/bin/bash

rm -rf ~/.relayer
rly config init
cp relay-config.yaml ~/.relayer/config/config.yaml

rly keys restore hub default "element achieve battle inject taxi hard purchase merit empower tower steak balance supreme purse assault lens chair dove together danger cat essence offer peace"
rly keys restore rolly default "element achieve battle inject taxi hard purchase merit empower tower steak balance supreme purse assault lens chair dove together danger cat essence offer peace"

rly paths new hub rolly theoptimists
rly tx link theoptimists
rly start