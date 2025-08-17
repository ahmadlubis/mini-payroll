.PHONY: build run test clean docker-up docker-down migrate seed

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=payslip-system
MAIN_PATH=./main.go

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	$(GOCMD) run $(MAIN_PATH)

# Test the application
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Benchmark tests
bench:
	$(GOTEST) -bench=. -benchmem ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy
	$(GOMOD) vendor

# Start docker services
docker-up:
	docker-compose up -d

# Stop docker services
docker-down:
	docker-compose down

# Database migration (requires running postgres)
migrate:
	$(GOCMD) run scripts/migrate.go

# Seed database (requires running postgres)
seed:
	$(GOCMD) run scripts/seed.go

# Build docker image
docker-build:
	docker build -t payslip-system:latest .

# Run all tests including integration
test-all:
	docker-compose up -d postgres_test
	sleep 5 # Wait for database to be ready
	$(GOTEST) -v ./...
	docker-compose down

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Generate API documentation
docs:
	swag init -g main.go