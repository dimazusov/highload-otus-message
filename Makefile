up:
	go run cmd/message/main.go

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint run ./...

goose-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

migrate-up: goose-install
	cd migrations && goose postgres "host=localhost user=postgres password= dbname=postgres port=5432 sslmode=disable" up

migrate-down: goose-install
	cd migrations && goose postgres "host=localhost user=postgres password= dbname=postgres port=5432 sslmode=disable" down

test-data:
	go run cmd/test_data_init/main.go

swagger-init:
	go get -u github.com/swaggo/swag/cmd/swag

swagger: swagger-init
	swag init -g ./internal/server/http/router.go -o api

env-up:
	MASTER_EXTERNAL_PORT=5433 COMPOSE_PROJECT_NAME=citus docker-compose up -d

env-stop:
	COMPOSE_PROJECT_NAME=citus docker-compose down -v