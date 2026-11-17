include .env
export

DB_CONTAINER := sqlc-net-plugin-db

build:
	go build -o ./bin/net-plugin main.go

gen: build
	cd test && sqlc generate
	cd test/gen && dotnet restore && dotnet build

console:
	cd test/console && dotnet run

up-dep:
	docker run -d --name $(DB_CONTAINER) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		postgres:17-alpine
	@until docker exec $(DB_CONTAINER) pg_isready -U $(POSTGRES_USER) >/dev/null 2>&1; do sleep 1; done
	@echo "postgres ready at $(DATABASE_URL)"

down-dep:
	docker rm -f $(DB_CONTAINER)

migrate:
	goose up

migrate-down:
	goose down

migrate-status:
	goose status

.PHONY: build t up-dep down-dep migrate migrate-down migrate-status
