# Linux 部署指南

## 依赖要求

### 系统依赖
- Go 1.25+
- MySQL 8.0
- ONNX Runtime 库 (libonnxruntime.so)

### 安装 ONNX Runtime (Linux)

#### 方法 1: 手动安装

```bash
# amd64 示例
wget https://github.com/microsoft/onnxruntime/releases/download/v1.24.1/onnxruntime-linux-x64-1.24.1.tgz
tar -xzf onnxruntime-linux-x64-1.24.1.tgz

# 安装到系统目录
sudo cp onnxruntime-linux-x64-1.24.1/lib/libonnxruntime.so.1.24.1 /usr/local/lib/libonnxruntime.so
sudo ldconfig

# 清理
rm -rf onnxruntime-linux-x64-1.24.1*
```

如果是 `arm64` 机器，将下载文件名中的 `x64` 改为 `aarch64`。

#### 方法 2: 手动指定库路径

将 `libonnxruntime.so` 放在以下任一位置：
- `./models/onnx/libonnxruntime.so`
- `./lib/libonnxruntime.so`
- `/usr/local/lib/libonnxruntime.so`

或者直接设置环境变量：

```bash
export ONNX_RUNTIME_LIB=/path/to/libonnxruntime.so
```

## 编译

```bash
# 启用 CGO (ONNX Runtime 需要)
CGO_ENABLED=1 go build -o herb-server ./cmd/server/main.go
```

## 运行

### 方式 1: 直接运行

```bash
# 确保库在系统路径中
export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

# 运行
./herb-server
```

### 方式 2: Docker (推荐)

```bash
docker compose up --build
```

## 常见问题

### 问题 1: 找不到 libonnxruntime.so

**错误信息：**
```text
未找到 ONNX Runtime 动态库
```

**解决：**
```bash
# 检查库是否存在
ls -la /usr/local/lib/libonnxruntime.so

# 更新动态链接库缓存
sudo ldconfig

# 或设置环境变量
export LD_LIBRARY_PATH=/path/to/lib:$LD_LIBRARY_PATH
```

### 问题 2: 缺少 GLIBC 依赖

**错误信息：**
```
Error loading shared library libgcompat.so.0
```

**解决：**
```bash
# Alpine Linux
sudo apk add gcompat

# Ubuntu/Debian
sudo apt-get install libc6
```

### 问题 3: 库版本不匹配

**错误信息：**
```text
version 'GLIBCXX_3.4.29' not found
```

**解决：**
```bash
# 更新 libstdc++
# Ubuntu/Debian
sudo apt-get install libstdc++6

# Alpine
sudo apk add libstdc++ libgcc
```

## 平台支持

| 平台 | 库文件名 | 测试状态 |
|------|---------|---------|
| macOS | libonnxruntime.dylib | ✅ 已测试 |
| Linux x64 | libonnxruntime.so | ✅ 支持 |
| Linux ARM | libonnxruntime.so | ✅ 支持 |
| Windows | onnxruntime.dll | ⚠️ 未测试 |

## 验证安装

```bash
# 检查依赖
ldd herb-server

# 检查是否能找到库
ldconfig -p | grep onnxruntime

# 运行测试
curl -X POST http://localhost:8080/api/recognize/upload \
  -H "Authorization: Bearer <token>" \
  -F "image=@test.jpg"
```
