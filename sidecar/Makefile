
install: build
	@echo "Installing attestation-sidecar binary to $(HOME)/go/bin/"
	@cp build/attestation-sidecar $(HOME)/go/bin/

build:
	@echo "Building attestation-sidecar binary"
	@mkdir -p build && go build -o build/attestation-sidecar

test:
	go test -mod=readonly ./...

.PHONY: *