# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Copy source code
COPY . .

# Download ONNX Runtime library (Linux version)
RUN wget -q https://github.com/microsoft/onnxruntime/releases/download/v1.21.0/onnxruntime-linux-x64-1.21.0.tgz && \
    tar -xzf onnxruntime-linux-x64-1.21.0.tgz && \
    mkdir -p models/onnx && \
    mv onnxruntime-linux-x64-1.21.0/lib/libonnxruntime.so.1.21.0 ./models/onnx/libonnxruntime.so && \
    rm -rf onnxruntime-linux-x64-1.21.0*

# Build
RUN GOPROXY="https://goproxy.cn,https://goproxy.io,direct" GOSUMDB=off \
    CGO_ENABLED=1 GOOS=linux go build -a -o main ./cmd/server/main.go


# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca certificates and dependencies
RUN apk --no-cache add ca-certificates tzdata libstdc++ libgcc

# Copy binary and required files
COPY --from=builder /app/main .
COPY --from=builder /app/models/onnx/libonnxruntime.so /usr/local/lib/
COPY --from=builder /app/models/onnx/herb.onnx ./models/onnx/
COPY --from=builder /app/models/onnx/classes.txt ./models/onnx/
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/configs/config.docker.yaml ./configs/config.yaml

# Update library path
ENV LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
