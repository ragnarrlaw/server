build:
	@go build -o bin/json_api

run:
	@./bin/json_api

test:
	@go test -v ./...
