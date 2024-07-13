# NOTE: This needs to be run from the root context of the repo

FROM golang:1.22-alpine as builder

WORKDIR /code

RUN apk add --no-cache build-base

# Set necessary environmet for Go module download
ENV GOPATH=""
ENV GOMODCACHE="/go/pkg/mod"

ADD prover-sidecar/go.mod prover-sidecar/go.sum ./

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download

COPY ./prover-sidecar .

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod make build

FROM alpine:3.16
COPY --from=builder /code/build/proversidecar /usr/bin/proversidecar

EXPOSE 6969

ENTRYPOINT ["/usr/bin/proversidecar"]
