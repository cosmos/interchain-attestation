
build-docker: build-hub-docker build-rolly-docker build-mock-da-docker build-simapp-docker

build-simapp-docker:
	@docker build -t simapp:local .

build-hub-docker:
	@cd hub && ignite chain build --skip-proto --output build && docker build -t hub:local .

build-rolly-docker:
	@cd rolly && ignite chain build --skip-proto --output build && docker build -t rolly:local .

build-mock-da-docker:
	@cd mock-da && docker build -t mock-da:local .

# Generate the `Counter.json` file containing the ABI of the Counter contract
# Requires `jq` to be installed on the system
# Requires `abigen` to be installed on the system to generate the go bindings for e2e tests
generate-abi:
	@cd contracts && forge install && forge build
	jq '.abi' contracts/out/Counter.sol/Counter.json > contracts/abi/Counter.json
	@echo "ABI file created at 'contracts/abi/Counter.json'"
	@echo "Generating go bindings for the end-to-end tests..."
	abigen --abi contracts/abi/Counter.json --pkg counter --type Contract --out interchaintest/types/counter/contract.go
	@echo "Done."
