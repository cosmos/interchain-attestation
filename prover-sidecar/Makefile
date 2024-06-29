
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        proto/proof.proto
	mv proto/*.pb.go server/

install:
	go install ./...
	echo "Installed proversidecar binary"

.PHONY: proto