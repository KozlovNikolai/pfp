include .env
LOCAL_BIN:=$(CURDIR)/bin

# Линтер:
install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

# Миграции:
install-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

# make create-new-migration name=fill_users
create-new-migration:
ifndef name
	$(error name is not set)
endif
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} create $(name) sql
local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

# 
.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out -tags="!mock" ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: mockgen
mockgen:
	mockgen --source=internal/app/transport/httpserver/interfaces.go \
	--destination=internal/app/transport/httpserver/mocks/mock_interfaces.go

.PHONY: mockery
mockery:
	mockery --name=IUserService \
	--dir=./internal/chat/transport/httpserver \
	--output=./internal/chat/transport/httpserver/mocks

test20:
	go test -v -count=20 ./...

