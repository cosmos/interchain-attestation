# NOTE: This needs to be run from the root context of the repo

FROM golang:1.22-alpine3.20 AS builder

RUN set -eux; apk add --no-cache git libusb-dev linux-headers gcc musl-dev make;

# Set necessary environmet for Go module download
ENV GOPATH=""
ENV GOMODCACHE="/go/pkg/mod"

# Copy relevant files before go mod download. Replace directives to local paths break if local
# files are not copied before go mod download.
COPY testing/rollupsimapp testing/rollupsimapp
COPY module module
COPY light-client light-client

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd module && go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd light-client && go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd testing/rollupsimapp && go mod download

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd testing/rollupsimapp && make build

FROM alpine:3.18

RUN apk add --no-cache jq

COPY --from=builder /go/testing/rollupsimapp/build/rollupsimappd /bin/rollupsimappd

ENTRYPOINT ["rollupsimappd"]
