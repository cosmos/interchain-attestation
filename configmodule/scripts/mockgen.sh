#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=types/expected_keepers.go -package testutil -destination testutil/expected_keepers_mocks.go