# NOTE: This needs to be run from the root context of the repo

FROM golang:1.23-alpine3.20 AS builder

RUN set -eux; apk add --no-cache git libusb-dev linux-headers gcc musl-dev make;

# Set necessary environmet for Go module download
ENV GOPATH=""
ENV GOMODCACHE="/go/pkg/mod"

# Copy relevant files before go mod download. Replace directives to local paths break if local
# files are not copied before go mod download.
COPY core core
COPY configmodule configmodule
COPY testing/rollupsimapp testing/rollupsimapp

RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd core && go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd configmodule && go mod download
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod cd testing/rollupsimapp && make build

FROM alpine:3.20

RUN apk add --no-cache jq

COPY --from=builder /go/testing/rollupsimapp/build/rollupsimappd /bin/rollupsimappd

ENTRYPOINT ["rollupsimappd"]
