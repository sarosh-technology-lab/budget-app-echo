# Define the default target
.PHONY: default
default: help

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  seeder FILENAME     Process the specified seeder file name"

# application repository and binary file name
NAME=Bougette

# application repository path
REPOSITORY=github.com/projects/${NAME}


install:
	go mod download

run-dev:
	echo "Starting Application In Development Mode"
	go run ./cmd/api


migrate-up:
	echo "Running migrations up"
	go run ./internal/database/migrate.go -direction=up

migrate-down:
	echo "Running migrations up"
	go run ./internal/database/migrate.go -direction=down

migrate_fresh:
	echo "Running fresh migrations"
	go run ./internal/database/migrate.go -direction=down
	go run ./internal/database/migrate.go -direction=up

#E.G make seeder FILENAME=category
.PHONY: seeder
seeder:
ifdef FILENAME
	echo "Seedinng : $(FILENAME)_seeder.go"
	go run "./internal/database/seeders/$(FILENAME)_seeder.go"
else
	echo "Error: FILENAME is not specified. Please provide the filename using 'make seeder FILENAME=<filename>'"
	exit 1
endif

test:
	 gotest ./tests/... -v