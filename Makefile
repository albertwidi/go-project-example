install-deps:
	@./scripts/install_dependencies.sh

gobuild:
	@go build -o project cmd/*.go 

gorun:
	make gobuild
	@./kothakexample