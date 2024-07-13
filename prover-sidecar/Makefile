
proto:
	@echo "Generating proto files"
	@protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/proof.proto
	@mv proto/*.pb.go server/
	@echo "Generated proto files"

build:
	@echo "Building proversidecar binary"
	@mkdir -p build && go build -o build/ ./...
	@echo "Built proversidecar binary in build/ directory"

install:
	@echo "Installing proversidecar binary"
	@go install ./...
	@echo "Installed proversidecar binary"

.PHONY: *