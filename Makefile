build_server:
	@go build -o bin/server ./cmd/server/main.go

run_server:
	@go build -o bin/server ./cmd/server/main.go
	@./bin/server

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
