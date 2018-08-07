gobuild:
	@go build -o kothakexample cmd/real/main.go

gorun:
	make gobuild
	@./kothakexample