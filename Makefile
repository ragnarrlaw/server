build_server:
	@go build -o bin/api ./cmd/api/main.go

run_api:
	@go build -o bin/api ./cmd/api/main.go
	@./bin/api

build_migrate:
	@go build -o bin/migrate ./cmd/migrate/main.go

run_migrate:
	@go build -o bin/migrate ./cmd/migrate/main.go
	@./bin/migrate

build_seed:
	@go build -o bin/migrate ./cmd/seed/main.go

run_seed:
	@go build -o bin/seed ./cmd/seed/main.go
	@./bin/seed

test:
	@go test -v ./...
