BINARY := item-app

help:
	@echo "usage: make <command>"
	@echo
	@echo "commands:"
	@echo "	build	compiles and build the app"
	@echo "	start	starts the app"
	@echo "	test	execute the tests"
	@echo "	bench	execute the benchmark tests"

start:
	@./bin/${BINARY}

build:
	@go build -o bin/${BINARY} cmd/main.go

test:
	@go test

bench:
	@go test -bench=.