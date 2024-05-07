
build-docker: build-hub-docker build-rolly-docker build-mock-da-docker

build-hub-docker:
	@cd hub && ignite chain build --skip-proto --output build && docker build -t hub:local .

build-rolly-docker:
	@cd rolly && ignite chain build --skip-proto --output build && docker build -t rolly:local .

build-mock-da-docker:
	@cd mock-da && docker build -t mock-da:local .