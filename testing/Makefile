###############################################################################
###                                 Build                                   ###
###############################################################################

install-simapps:
	@echo "===========     Installing simapp     ============"
	@cd simapp && make install

###############################################################################
###                                 Docker                                  ###
###############################################################################

docker-images: simapp-image sidecar-image

simapp-image:
	@echo "Building ghcr.io/cosmos/interchain-attestation-simapp:local docker image"
	docker build -t ghcr.io/cosmos/interchain-attestation-simapp:local -f simapp.Dockerfile ../

sidecar-image:
	@echo "Building ghcr.io/cosmos/interchain-attestation-sidecar:local docker image"
	docker build -t ghcr.io/cosmos/interchain-attestation-sidecar:local -f sidecar.Dockerfile ../

###############################################################################
###                                Testing                                  ###
###############################################################################

interchaintest:
	@echo "Running interchaintest tests"
	@cd interchaintest && go test -race -count=1 -v ./...

.PHONY: *
