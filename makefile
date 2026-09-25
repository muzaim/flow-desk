# Load environment variables from .env file if it exists
ifneq (...,$(wildcard .env))
    include .env
    export
endif

# MySQL connection URL for golang-migrate
# Format: mysql://user:password@tcp(host:port)/dbname
DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

.PHONY: dev build run migrate-up migrate-down migrate-create

# Run the server with Air (Hot Reload)
dev:
	air

# Build the application binary
build:
	go build -o bin/app main.go

# Run the built application binary
run: build
	./bin/app

# Run all UP migrations
migrate-up:
	migrate -path database/migrations -database "$(DB_URL)" -verbose up

# Roll back the last migration
migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" -verbose down 1

# Roll back all migrations (reset the database)
migrate-reset:
	migrate -path database/migrations -database "$(DB_URL)" -verbose down -all

# Create a new migration file (Usage: make migrate-create name=add_priority_to_tasks)
migrate-create:
	migrate create -ext sql -dir database/migrations -seq $(name)