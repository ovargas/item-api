BINARY := item-app

build:
	@go build -o bin/${BINARY} cmd/main.go