###############################################################################
###                                 Build                                   ###
###############################################################################

install-simapps:
	@echo "===========     Installing simapp     ============"
	@cd simapp && make install
	@echo "===========     Installing rollupsimapp     ============"
	@cd rollupsimapp && make install

# Generate the `Counter.json` file containing the ABI of the Counter contract
# Requires `forge` to be installed on the system
# Requires `jq` to be installed on the system
# Requires `abigen` to be installed on the system to generate the go bindings for e2e tests
generate-abi:
	@cd contracts && forge install && forge build
	jq '.abi' contracts/out/Counter.sol/Counter.json > contracts/abi/Counter.json
	@echo "ABI file created at 'contracts/abi/Counter.json'"
	@echo "Generating go bindings for the end-to-end tests..."
	abigen --abi contracts/abi/Counter.json --pkg counter --type Contract --out interchaintest/types/counter/contract.go
	@echo "Done."

###############################################################################
###                                 Docker                                  ###
###############################################################################

docker-images: simapp-image rollupsimapp-image sidecar-image mock-da-image

simapp-image:
	@echo "Building simapp:local docker image"
	docker build -t simapp:local -f simapp.Dockerfile ../

rollupsimapp-image:
	@echo "Building rollupsimapp:local docker image"
	docker build -t rollupsimapp:local -f rollupsimapp.Dockerfile ../

sidecar-image:
	@echo "Building attestationsidecar:local docker image"
	docker build -t attestationsidecar:local -f sidecar.Dockerfile ../

mock-da-image:
	@echo "Building mock-da:local docker image"
	docker build -t mock-da:local -f mock-da.Dockerfile .

###############################################################################
###                                Testing                                  ###
###############################################################################

interchaintest:
	@echo "Running interchaintest tests"
	@cd interchaintest && go test -race -count=1 -v -run E2ETestSuite/TestIBC ./...
	@sleep 5
	@cd interchaintest && go test -race -count=1 -v -run E2ETestSuite/TestSidecar ./...

###############################################################################
###                                  Serve                                  ###
###############################################################################

serve: kill-all
	@echo "===========     Serve     ============"
	./scripts/serve-simapp.sh
	./scripts/serve-rollupsimapp.sh

kill-all:
	@echo "Killing simappd"
	-@pkill simappd 2>/dev/null
	@echo "Killing rollupsimappd"
	-@pkill rollupsimappd 2>/dev/null

.PHONY: *