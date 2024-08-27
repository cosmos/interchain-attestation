
lint:
    @echo "Running golangci-lint in all packages"
    cd core && golangci-lint run -c ../.golangci.yml
    cd configmodule && golangci-lint run -c ../.golangci.yml
    cd sidecar && golangci-lint run -c ../.golangci.yml
    cd testing/simapp && golangci-lint run -c ../../.golangci.yml
    cd testing/rollupsimapp && golangci-lint run -c ../../.golangci.yml
    cd testing/interchaintest && golangci-lint run -c ../../.golangci.yml

tidy:
    @echo "Running go mod tidy in all packages"
    cd core && go mod tidy
    cd configmodule && go mod tidy
    cd sidecar && go mod tidy
    cd testing/simapp && go mod tidy
    cd testing/rollupsimapp && go mod tidy
    cd testing/interchaintest && go mod tidy

proto-gen:
    @echo "Generating proto files in all packages"
    cd core && make proto-gen
    cd configmodule && make proto-gen

test-unit:
    @echo "Running unit tests in all packages"
    cd core && make test
    cd configmodule && make test
    cd sidecar && make test

test-e2e:
    cd testing && make interchaintest