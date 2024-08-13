build:
	@go build -o bin/json_api main.go

run:
	@./bin/json_api

test:
	@go test -v ./...
