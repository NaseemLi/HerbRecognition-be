# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Copy source code (needed for go build to download correct deps)
COPY . .

# Set Go proxy and build in one step
RUN GOPROXY="https://goproxy.cn,https://goproxy.io,direct" GOSUMDB=off \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go


# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/configs/config.docker.yaml ./configs/config.yaml

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
