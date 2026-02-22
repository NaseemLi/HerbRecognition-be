#!/bin/bash

echo "=========================================="
echo "   中草药识别系统"
echo "=========================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

CONDA_ENV_NAME="herb-inference"
GO_PORT=8080
PYTHON_PORT=5001

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查端口占用
check_port() {
    if lsof -i :$1 -sTCP:LISTEN -t >/dev/null ; then
        return 0
    fi
    return 1
}

# 清理函数
cleanup() {
    echo ""
    log_info "正在停止所有服务..."

    if [ -n "$GO_PID" ] && kill -0 $GO_PID 2>/dev/null; then
        kill $GO_PID
        log_info "Go 服务已停止"
    fi

    if [ -n "$PYTHON_PID" ] && kill -0 $PYTHON_PID 2>/dev/null; then
        kill $PYTHON_PID
        log_info "Python 推理服务已停止"
    fi

    exit 0
}

trap cleanup SIGINT SIGTERM

# 检查 Go 端口
if check_port $GO_PORT; then
    log_error "端口 $GO_PORT 已被占用"
    exit 1
fi

# 检查 Python 端口
if check_port $PYTHON_PORT; then
    log_error "端口 $PYTHON_PORT 已被占用"
    exit 1
fi

# 启动 Go 服务
log_info "正在启动 Go API 服务..."
go run cmd/server/main.go &
GO_PID=$!
sleep 2

# 启动 Python 推理服务
log_info "正在启动 Python 推理服务..."

# 加载 conda
if [ -f "$(conda info --base)/etc/profile.d/conda.sh" ]; then
    source "$(conda info --base)/etc/profile.d/conda.sh"
    conda activate $CONDA_ENV_NAME
else
    log_error "Conda 未安装或配置不正确"
    exit 1
fi

cd "$SCRIPT_DIR/services/inference-service"
python app.py &
PYTHON_PID=$!

log_info "=========================================="
log_info "服务启动完成!"
log_info "------------------------------------------"
log_info "Go API 服务：http://localhost:$GO_PORT"
log_info "Python 推理：http://localhost:$PYTHON_PORT"
log_info "=========================================="
log_info "按 Ctrl+C 停止所有服务"

# 等待子进程
wait
