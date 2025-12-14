# Multi-stage Dockerfile for Health Monitor

# Build stage
FROM golang:1.24.11-alpine AS builder

# Install build dependencies including gcc for CGO
RUN apk add --no-cache git make ca-certificates tzdata gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
ARG VERSION=dev

RUN BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') && \
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o health-monitor \
    ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 healthmon && \
    adduser -D -u 1000 -G healthmon healthmon

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/health-monitor .

# Copy configuration examples
COPY --from=builder /build/configs /configs

# Copy web static files
COPY --from=builder /build/web /app/web

# Create data directory
RUN mkdir -p /data && chown -R healthmon:healthmon /data /app

# Switch to non-root user
USER healthmon

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/health-monitor", "--version"]

# Run the application
ENTRYPOINT ["/app/health-monitor"]
CMD ["--config", "/configs/example.yaml"]
