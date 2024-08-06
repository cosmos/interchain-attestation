
build:
	@echo "Building attestationsidecar binary"
	@mkdir -p build && go build -o build/attestationsidecar
	@echo "Built attestationsidecar binary in build/ directory"

test:
	go test -mod=readonly ./...

.PHONY: *