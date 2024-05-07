FROM ignitehq/cli:v28.3.0 AS builder

USER root

WORKDIR /code

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN ignite chain build --skip-proto --output build
RUN cp build/hubd /usr/bin/hubd

# TODO: Put it into a smaller image (alpine or whatever)

# rest server, tendermint p2p, tendermint rpc
EXPOSE 1317 26656 26657

ENTRYPOINT ["/bin/sh"]
