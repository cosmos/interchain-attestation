#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=voteextension/expected_keepers.go -package testutil -destination voteextension/testutil/expected_keepers_mocks.go
