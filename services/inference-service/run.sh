#!/bin/bash

echo "=========================================="
echo "启动中草药识别推理服务"
echo "=========================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

CONDA_ENV_NAME="herb-inference"

echo ""
echo "检查 Conda 环境：$CONDA_ENV_NAME"
if ! conda env list | grep -q "^$CONDA_ENV_NAME "; then
    echo "错误：Conda 环境 '$CONDA_ENV_NAME' 不存在"
    echo "请先运行：conda create -n $CONDA_ENV_NAME python=3.10 -y"
    echo "然后运行：pip install -r requirements.txt"
    exit 1
fi

source "$(conda info --base)/etc/profile.d/conda.sh"
conda activate $CONDA_ENV_NAME

echo ""
echo "启动 Python 推理服务..."
echo "服务地址：http://0.0.0.0:5001"
echo "日志文件：/tmp/inference_service.log"
echo ""

nohup python app.py > /tmp/inference_service.log 2>&1 &
PID=$!
echo "服务已启动 (PID: $PID)"
echo ""
echo "查看日志：tail -f /tmp/inference_service.log"
echo "停止服务：kill $PID"
echo "检查状态：curl http://localhost:5001/health"
