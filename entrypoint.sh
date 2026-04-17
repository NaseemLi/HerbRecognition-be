#!/bin/sh
set -e

echo "=== 检查 ONNX Runtime 库 ==="
echo "库文件是否存在:"
ls -la /usr/local/lib/libonnxruntime.so || echo "NOT FOUND"

echo ""
echo "库依赖检查:"
ldd /usr/local/lib/libonnxruntime.so 2>/dev/null || echo "ldd failed"

echo ""
echo "动态链接库缓存:"
ldconfig -p | grep onnx || echo "not in cache"

echo ""
echo "LD_LIBRARY_PATH: $LD_LIBRARY_PATH"

echo ""
echo "=== 启动服务 ==="
exec ./main
