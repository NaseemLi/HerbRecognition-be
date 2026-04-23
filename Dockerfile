# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25
ARG ONNXRUNTIME_VERSION=1.24.1

FROM --platform=$TARGETPLATFORM golang:${GO_VERSION}-bookworm AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

ENV GOPROXY=https://goproxy.cn,direct

RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -o /out/main ./cmd/server/main.go

FROM --platform=$TARGETPLATFORM debian:bookworm-slim AS onnxruntime

ARG TARGETARCH
ARG ONNXRUNTIME_VERSION

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates wget tar && \
    rm -rf /var/lib/apt/lists/*

RUN case "${TARGETARCH}" in \
        amd64) ort_arch="x64" ;; \
        arm64) ort_arch="aarch64" ;; \
        *) echo "Unsupported TARGETARCH: ${TARGETARCH}" >&2; exit 1 ;; \
    esac && \
    wget -q "https://github.com/microsoft/onnxruntime/releases/download/v${ONNXRUNTIME_VERSION}/onnxruntime-linux-${ort_arch}-${ONNXRUNTIME_VERSION}.tgz" && \
    tar -xzf "onnxruntime-linux-${ort_arch}-${ONNXRUNTIME_VERSION}.tgz" && \
    install -Dm644 \
        "onnxruntime-linux-${ort_arch}-${ONNXRUNTIME_VERSION}/lib/libonnxruntime.so.${ONNXRUNTIME_VERSION}" \
        /out/usr/local/lib/libonnxruntime.so

FROM --platform=$TARGETPLATFORM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata libstdc++6 && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /app/models/onnx /app/uploads/images /app/uploads/herbs

COPY --from=builder /out/main ./main
COPY --from=onnxruntime /out/usr/local/lib/libonnxruntime.so /usr/local/lib/libonnxruntime.so
COPY models/onnx/herb.onnx ./models/onnx/herb.onnx
COPY models/onnx/classes.txt ./models/onnx/classes.txt
COPY configs/config.docker.yaml ./configs/config.yaml

ENV LD_LIBRARY_PATH=/usr/local/lib
ENV ONNX_RUNTIME_LIB=/usr/local/lib/libonnxruntime.so

EXPOSE 8080

CMD ["./main"]
