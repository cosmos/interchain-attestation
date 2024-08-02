
tidy:
	cd light-client && go mod tidy
	cd module && go mod tidy
	cd attestation-sidecar && go mod tidy
	cd testing/simapp && go mod tidy
	cd testing/rollupsimapp && go mod tidy
	cd testing/interchaintest && go mod tidy