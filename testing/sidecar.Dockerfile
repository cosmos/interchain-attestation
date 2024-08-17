# NOTE: This needs to be run from the root context of the repo

FROM golang:1.23-alpine AS builder

WORKDIR /code

RUN apk add --no-cache build-base

# Set necessary environmet for Go module download
ENV GOPATH=""
ENV GOMODCACHE="/go/pkg/mod"

COPY core core
COPY sidecar sidecar

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd core && go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd sidecar && go mod download

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd sidecar && make build

FROM alpine:3.20
COPY --from=builder /code/sidecar/build/attestationsidecar /usr/bin/attestationsidecar

EXPOSE 6969

ENTRYPOINT ["/usr/bin/attestationsidecar"]
