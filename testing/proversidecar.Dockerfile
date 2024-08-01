# NOTE: This needs to be run from the root context of the repo

FROM golang:1.22-alpine AS builder

WORKDIR /code

RUN apk add --no-cache build-base

# Set necessary environmet for Go module download
ENV GOPATH=""
ENV GOMODCACHE="/go/pkg/mod"

COPY light-client light-client
COPY prover-sidecar prover-sidecar

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd light-client && go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd prover-sidecar && go mod download

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd prover-sidecar && make build

FROM alpine:3.16
COPY --from=builder /code/prover-sidecar/build/proversidecar /usr/bin/proversidecar

EXPOSE 6969

ENTRYPOINT ["/usr/bin/proversidecar"]
