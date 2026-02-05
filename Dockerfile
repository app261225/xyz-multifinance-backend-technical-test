# Multi-stage build for optimized final image

# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o xyz-multifinance .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and TLS
RUN apk --no-cache add ca-certificates

# Copy certificate for database connection
COPY ca.pem .

# Copy .env file (should be mounted or set via environment variables in production)
COPY .env .

# Copy the binary from builder
COPY --from=builder /app/xyz-multifinance .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./xyz-multifinance"]
