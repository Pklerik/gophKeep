FROM golang:1.25 AS builder

WORKDIR /build

# Install build dependencies for SQLite3
RUN apt-get update && apt-get install -y --no-install-recommends \
    git gcc libc-dev libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.Version=v1.0.0 -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
    -o gophkeeper-server ./cmd/gophkeep

# Final stage - use debian slim for runtime
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates libsqlite3-0 curl \
    && rm -rf /var/lib/apt/lists/*

# Copy binary from builder
COPY --from=builder /build/gophkeeper-server .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run server
CMD ["./gophkeeper-server", "server"]
