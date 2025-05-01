# Define the default target
.PHONY: default
default: help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  seed               Seed the database"
	@echo "  run-air            Run the app with hot reloading using Air"

# application repository and binary file name
NAME=Bougette

# application repository path
REPOSITORY=github.com/projects/${NAME}

install:
	go mod download

run-dev:
	echo "Starting Application In Development Mode"
	go run ./cmd/api

run-air:
	echo "Starting Application with Hot Reloading using Air"
	~/go/bin/air

migrate-up:
	echo "Running migrations up"
	go run ./internal/database/migrate.go -direction=up

migrate-down:
	echo "Running migrations up"
	go run ./internal/database/migrate.go -direction=down

migrate-fresh:
	echo "Running fresh migrations"
	go run ./internal/database/migrate.go -direction=down
	go run ./internal/database/migrate.go -direction=up

seed:
	echo "Seeding Database"
	go run "./cmd/seed.go"
	echo "Seeding Completed"

test:
	gotest ./tests/... -v