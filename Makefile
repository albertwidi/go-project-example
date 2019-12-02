mainprogram=projectbackend

install-deps:
	@./scripts/install_dependencies.sh

gobuild:
	@go build -v -o $(mainprogram) cmd/project/*.go 

gorun:
	make gobuild
	@./$(mainprogram) \
		-config_file="./project.config.toml" \
		-env_file="./project.env.toml" \
		-tz="Asia/Jakarta"

testconfig:
	make gobuild
	@./$(mainprogram) \
		-config_file=./project.config.toml \
		-env_file=./project.env.toml \
		-debug=-testconfig=1-devserver=1 \
		-tz=Asia/Jakarta

dbup:
	@cd database && ./setup.sh create database.yml

dbdown:
	@cd database && ./setup.sh drop database.yml