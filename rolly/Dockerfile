FROM ghcr.io/gjermundgaraba/ignitecli:v28.3.0-go1.22 AS builder

USER root

RUN apt-get update -yq \
    && apt-get install --no-install-recommends -yq \
    jq

WORKDIR /code

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN ignite chain build --skip-proto --output build
RUN cp build/rollyd /usr/bin/rollyd

# TODO: Put it into a smaller image (alpine or whatever)

# rest server, tendermint p2p, tendermint rpc
EXPOSE 1317 26656 26657

ENTRYPOINT ["/bin/sh"]
