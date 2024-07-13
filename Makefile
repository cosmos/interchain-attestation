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

docker-images: simapp-image rollupsimapp-image prover-sidecar-image mock-da-image

simapp-image:
	@echo "Building simapp:local docker image"
	docker build -t simapp:local -f simapp.Dockerfile .

rollupsimapp-image:
	@echo "Building rollupsimapp:local docker image"
	docker build -t rollupsimapp:local -f rollupsimapp.Dockerfile .

proversidecar-image:
	@echo "Building proversidecar:local docker image"
	docker build -t proversidecar:local -f proversidecar.Dockerfile .

mock-da-image:
	@echo "Building mock-da:local docker image"
	docker build -t mock-da:local -f mock-da.Dockerfile .

# TODO: REMOVE
hub-docker:
	@cd hub && ignite chain build --skip-proto --output build && docker build -t hub:local .

# TODO: REMOVE
rolly-docker:
	@cd rolly && ignite chain build --skip-proto --output build && docker build -t rolly:local .

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
