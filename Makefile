# Load environment variables from .env file 
include .env
export $(shell sed 's/=.*//' .env)

# Build the application
all: build test

migrate-up:
	@echo "Migration up with DB_HOST=$(DB_HOST) and DB_PORT=$(EXTERNAL_DB_PORT)"
	@goose -dir ./migrations postgres "user=postgres password=postgres port=$(EXTERNAL_DB_PORT) host=localhost dbname=golang_db sslmode=disable" up

migrate-down:
	@echo "Migration down with DB_HOST=$(DB_HOST) and DB_PORT=$(EXTERNAL_DB_PORT)"
	@goose -dir ./migrations postgres "user=postgres password=postgres port=$(EXTERNAL_DB_PORT) host=localhost dbname=golang_db sslmode=disable" down-to 0

build:
	@echo "Building..."
	
	
	@go build -o main main.go

# Run the application
run:
	@go run main.go

# Start DB container in detached mode
docker-up:
	@if docker compose up --build -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build -d; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests/... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./integration_tests/... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest
