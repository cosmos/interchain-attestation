FROM golang:1.22-alpine3.20 as builder

RUN set -eux; apk add --no-cache git libusb-dev linux-headers gcc musl-dev make;

ENV GOPATH=""

# Copy relevant files before go mod download. Replace directives to local paths break if local
# files are not copied before go mod download.
ADD simapp simapp
ADD module module

COPY contrib/devtools/Makefile contrib/devtools/Makefile

RUN cd module && go mod download
RUN cd simapp && go mod download

RUN cd simapp && make build

FROM alpine:3.18

COPY --from=builder /go/simapp/build/simd /bin/simd

ENTRYPOINT ["simd"]
