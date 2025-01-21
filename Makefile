build:
	@go build -o bin/simpleurl-server

run: build
	@./bin/simpleurl-server

test:
	@go test -v ./...