###############################################################################
###                                 Build                                   ###
###############################################################################

install-simapps:
	@echo "===========     Installing simapp     ============"
	@cd simapp && make install
	@echo "===========     Installing rollupsimapp     ============"
	@cd rollupsimapp && make install
	# TODO: Install sidecar

###############################################################################
###                                 Docker                                  ###
###############################################################################

docker-images: simapp-image rollupsimapp-image proversidecar-image mock-da-image

simapp-image:
	@echo "Building simapp:local docker image"
	docker build -t simapp:local -f simapp.Dockerfile ../

rollupsimapp-image:
	@echo "Building rollupsimapp:local docker image"
	docker build -t rollupsimapp:local -f rollupsimapp.Dockerfile ../

proversidecar-image:
	@echo "Building proversidecar:local docker image"
	docker build -t proversidecar:local -f proversidecar.Dockerfile ../

mock-da-image:
	@echo "Building mock-da:local docker image"
	docker build -t mock-da:local -f mock-da.Dockerfile .

###############################################################################
###                                Testing                                  ###
###############################################################################

interchaintest:
	@echo "Running interchaintest tests"
	@cd interchaintest && go test -race -count=1 -v -run E2ETestSuite/TestSidecar ./...

###############################################################################
###                                  Serve                                  ###
###############################################################################

serve: kill-all
	@echo "===========     Serve     ============"
	./scripts/serve-simapp.sh
	./scripts/serve-rollupsimapp.sh

kill-all:
	@echo "Killing simappd"
	-@pkill simappd 2>/dev/null
	@echo "Killing rollupsimappd"
	-@pkill rollupsimappd 2>/dev/null

.PHONY: *