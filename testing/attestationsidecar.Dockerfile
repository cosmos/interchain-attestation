# NOTE: This needs to be run from the root context of the repo

FROM golang:1.22-alpine AS builder

WORKDIR /code

RUN apk add --no-cache build-base

# Set necessary environmet for Go module download
ENV GOPATH=""
ENV GOMODCACHE="/go/pkg/mod"

COPY light-client light-client
COPY attestation-sidecar attestation-sidecar

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd light-client && go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd attestation-sidecar && go mod download

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd attestation-sidecar && make build

FROM alpine:3.16
COPY --from=builder /code/attestation-sidecar/build/attestationsidecar /usr/bin/attestationsidecar

EXPOSE 6969

ENTRYPOINT ["/usr/bin/attestationsidecar"]
