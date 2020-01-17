mainprogram=projectbackend
build_commit=$(shell git rev-parse HEAD)
build_version=$(shell git describe --tags 2> /dev/null || echo "dev-$(shell git rev-parse HEAD)")

.PHONY: install-deps
install-deps:
	@./scripts/install_dependencies.sh

.PHONY: version
version:
	echo $(build_version)
	echo $(build_commit)

.PHONY: build
build:
	@go build -v \
		-ldflags "-X main.buildVersion=$(build_version) \
		-X main.buildCommit=$(build_commit)" \
		-race \
		-o $(mainprogram) cmd/project/*.go

.PHONY: run
run:
	make build 
	@./$(mainprogram) \
		-config_file="./project.config.toml" \
		-env_file="./project.env.toml" \
		-tz="Asia/Jakarta"

.PHONY: test
test:
	@go generate -v ./...
	@go test -race -v ./...

.PHONY: testconfig
testconfig:
	make build
	@./$(mainprogram) \
		-config_file=./project.config.toml \
		-env_file=./project.env.toml \
		-debug=-testconfig=1 \
			-devserver=1 \
		-tz=Asia/Jakarta

.PHONY: dbup
dbup:
	@cd database && ./setup.sh create database.yml

.PHONY: dbdown
dbdown:
	@cd database && ./setup.sh drop database.yml