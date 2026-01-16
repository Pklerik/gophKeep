.PHONY: build server client test coverage bench clean help install

# Build variables
BINARY_SERVER = gophkeeper-server
BINARY_CLIENT = gophkeeper-client
BUILD_DIR = bin
VERSION = v1.0.0
BUILD_TIME = $(shell date -u '+%Y-%m-%d_%H:%M:%S')

help:
	@echo "GophKeeper - Password Manager Application"
	@echo ""
	@echo "Available targets:"
	@echo "  make build          - Build both server and client"
	@echo "  make server         - Build server binary"
	@echo "  make client         - Build client binary"
	@echo "  make test           - Run all tests"
	@echo "  make coverage       - Run tests with coverage report"
	@echo "  make bench          - Run benchmarks"
	@echo "  make clean          - Remove built binaries"
	@echo "  make install        - Install binaries"
	@echo ""

build: server client
	@echo "✓ Build complete"

server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(BUILD_DIR)/$(BINARY_SERVER) ./cmd/gophkeep

client:
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(BUILD_DIR)/$(BINARY_CLIENT) ./cmd/gophkeep

test:
	@echo "Running tests..."
	@go test -v -cover -coverprofile=coverage.out ./internal/...
	@echo "✓ Tests passed"

coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out -covermode=atomic ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated (coverage.html)"
	@go tool cover -func=coverage.out | tail -1

bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem -benchtime=100ms -run=^$$ ./internal/...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f *.db
	@echo "✓ Clean complete"

install: build
	@echo "Installing binaries..."
	@cp $(BUILD_DIR)/$(BINARY_SERVER) $(GOPATH)/bin/
	@cp $(BUILD_DIR)/$(BINARY_CLIENT) $(GOPATH)/bin/
	@echo "✓ Installation complete"

run-server: server
	@echo "Starting server..."
	@./$(BUILD_DIR)/$(BINARY_SERVER) server -a 0.0.0.0:8080

run-client: client
	@echo "Starting client..."
	@./$(BUILD_DIR)/$(BINARY_CLIENT) client

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

vet:
	@echo "Running go vet..."
	@go vet -vettool=$$(which statictest) ./...
	@echo "✓ Vet passed"

staticcheck:
	@echo "Running staticcheck..."
	@staticcheck ./...
	@echo "✓ Staticcheck passed"

cilint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...
	@echo "✓ Golangci-lint passed"

lint: fmt vet staticcheck cilint
	@echo "✓ Lint checks passed"
	
	go tool pprof -http :9000 profiles/mem.out

pprof-cpu:
	go tool pprof -http :9000 profiles/cpu.out

# example make a iter=5 for run 1-5ths iteration
mock:
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/mock_repo.go -package=mocks
	mockgen -source=internal/client/client.go -destination=internal/client/mocks/mock_client.go -package=mocks