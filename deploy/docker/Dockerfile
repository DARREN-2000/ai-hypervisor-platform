# Multi-stage Dockerfile for AI Hypervisor Platform services

# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make build-base

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build API Server
ARG SERVICE=api-server
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always)" \
    -o /app/bin/${SERVICE} \
    ./cmd/${SERVICE}

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates libvirt

# Create non-root user
RUN addgroup -g 1000 aihypervisor && \
    adduser -D -u 1000 -G aihypervisor aihypervisor

# Copy binary from builder
ARG SERVICE=api-server
COPY --from=builder /app/bin/${SERVICE} /usr/local/bin/aihypervisor-service

# Set user
USER aihypervisor:aihypervisor

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run service
ENTRYPOINT ["/usr/local/bin/aihypervisor-service"]
CMD ["--config", "/etc/aihypervisor/config.yaml"]
