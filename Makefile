mainprogram=projectbackend

install-deps:
	@./scripts/install_dependencies.sh

gobuild:
	@go build -v -o $(mainprogram) cmd/project/*.go 

gorun:
	make gobuild
	@./$(mainprogram) \
		-config_file=./project.config.toml \
		-env_file=./project.env.toml

testconfig:
	make gobuild
	@./$(mainprogram) \
		--config_file=./project.config.toml \
		--env_file=./project.env.toml \
		--debug=testconfig=1