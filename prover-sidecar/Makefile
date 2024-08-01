
proto:
	@echo "Generating proto files"
	@cd proto && buf generate --template buf.gen.gogo.yaml server/proof.proto
	@cd proto && buf generate --template buf.gen.gogo.yaml server/types.proto
	@mv github.com/gjermundgaraba/pessimistic-validation/proversidecar/server/* ./server
	@mv github.com/gjermundgaraba/pessimistic-validation/proversidecar/types/* ./types
	@rm -rf github.com
	@echo "Generated proto files"

build:
	@echo "Building proversidecar binary"
	@mkdir -p build && go build -o build/ ./...
	@echo "Built proversidecar binary in build/ directory"

install:
	@echo "Installing proversidecar binary"
	@go install ./...
	@echo "Installed proversidecar binary"

test:
	go test -mod=readonly ./...

.PHONY: *