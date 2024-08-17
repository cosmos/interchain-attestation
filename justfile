
tidy:
    @echo "Running go mod tidy in all packages"
    cd core && go mod tidy
    cd configmodule && go mod tidy
    cd sidecar && go mod tidy
    cd testing/simapp && go mod tidy
    cd testing/rollupsimapp && go mod tidy
    cd testing/interchaintest && go mod tidy
