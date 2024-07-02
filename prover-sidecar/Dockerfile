FROM golang:1.22-alpine as builder

WORKDIR /code

RUN apk add --no-cache build-base

ADD go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

FROM alpine:3.16
COPY --from=builder /code/build/proversidecar /usr/bin/proversidecar

EXPOSE 6969

ENTRYPOINT ["/usr/bin/proversidecar"]
