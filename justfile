
lint:
    @echo "Running golangci-lint in all packages"
    cd core && golangci-lint run
    cd configmodule && golangci-lint run
    cd sidecar && golangci-lint run
    cd testing/simapp && golangci-lint run
    cd testing/rollupsimapp && golangci-lint run
    cd testing/interchaintest && golangci-lint run

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

test-e2e image-version="local":
    if [[ "{{image-version}}" = "local" ]]; then just build-docker-images; fi
    cd testing && DOCKER_IMAGE_VERSION={{image-version}} make interchaintest

build-docker-images:
    cd testing && make docker-images
