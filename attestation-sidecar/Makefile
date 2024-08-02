
proto:
	@echo "Generating proto files"
	@cd proto && buf generate --template buf.gen.gogo.yaml server/claim.proto
	@cd proto && buf generate --template buf.gen.gogo.yaml server/types.proto
	@mv github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/server/* ./server
	@mv github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/types/* ./types
	@rm -rf github.com
	@echo "Generated proto files"

build:
	@echo "Building attestationsidecar binary"
	@mkdir -p build && go build -o build/ ./...
	@echo "Built attestationsidecar binary in build/ directory"

install:
	@echo "Installing attestationsidecar binary"
	@go install ./...
	@echo "Installed attestationsidecar binary"

test:
	go test -mod=readonly ./...

.PHONY: *