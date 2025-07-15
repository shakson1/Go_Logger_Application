# syntax=docker/dockerfile:1
FROM golang:1.21 as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install build dependencies for SQLite
RUN apt-get update && apt-get install -y gcc libsqlite3-dev

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main .

# Final stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/*.db .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/metrics || exit 1

# Run the binary
CMD ["./main"] 