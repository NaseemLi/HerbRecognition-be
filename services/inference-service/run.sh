#!/bin/bash

echo "=========================================="
echo "启动中草药识别推理服务"
echo "=========================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

CONDA_ENV_NAME="herb-inference"
PORT=5001

# 加载 conda
source "$(conda info --base)/etc/profile.d/conda.sh"
conda activate $CONDA_ENV_NAME

# 检查端口
if lsof -i :$PORT -sTCP:LISTEN -t >/dev/null ; then
    echo "端口 $PORT 已被占用"
    exit 1
fi

echo "启动 Python 推理服务..."

python app.py &
PID=$!

trap "echo ''; echo '正在停止服务...'; kill $PID; exit 0" SIGINT

# 等待子进程
wait $PID