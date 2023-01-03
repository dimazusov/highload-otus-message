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
	cd migrations && goose mysql "root:pass@/highload" up

migrate-down: goose-install
	cd migrations && goose mysql "root:pass@/highload" down

test-data:
	docker-compose -f deployments/docker-compose.yaml exec app sh -c "/opt/social/testdatagen --config=/etc/social/config.yaml"

swagger-init:
	go get -u github.com/swaggo/swag/cmd/swag

swagger: swagger-init
	swag init -g ./internal/server/http/router.go -o api

up-env:
	MASTER_EXTERNAL_PORT=5433 COMPOSE_PROJECT_NAME=citus docker-compose up -d

up-env:
	COMPOSE_PROJECT_NAME=citus docker-compose down -v

