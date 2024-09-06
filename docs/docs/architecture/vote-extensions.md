---
sidebar_position: 4
---

# Vote Extensions

ABCI++ is a CometBFT interface that allows a chain to add more functionality to the low-level parts of their application.
In particular, for our purposes, it allows validators to communicate with each other to aggregate potentially disparate information
and come to consensus on an aggregate value. To read more about ABCI++, please refer to the official documentation [here](https://docs.cometbft.com/v0.38/spec/abci/).

In the context of Interchain Attestation, we use the ABCI++ interface to fetch attestations from the sidecar, 
aggregate them, and send them to the light client for verification and client updates.

TODO: Add illustration with all the different callbacks