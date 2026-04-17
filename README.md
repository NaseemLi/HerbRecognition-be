# 中草药识别系统后端

基于 Go + Gin + GORM + ONNX Runtime 的中草药图像识别后端服务。

当前识别模块已经迁移到 Go 进程内执行，直接加载 ONNX 模型推理，不再依赖独立的 Python 推理服务，也不再使用 `.pth` 模型部署链路。

## 技术栈

| 组件 | 技术 |
| --- | --- |
| Web 框架 | Go 1.25 + Gin |
| 数据库 | MySQL 8.0 + GORM |
| 认证 | JWT |
| 配置管理 | Viper |
| 日志 | Zap |
| 模型推理 | ONNX Runtime + ONNX |
| 容器化 | Docker + Docker Compose |

## 当前架构

- API、鉴权、上传、识别、查询都在同一个 Go 服务中完成。
- 模型文件位于 `models/onnx/herb.onnx`。
- 分类标签文件位于 `models/onnx/classes.txt`。
- Docker 构建时会按目标架构自动下载匹配的 ONNX Runtime 动态库。
- 当前 Dockerfile 支持 `linux/amd64` 和 `linux/arm64` 构建。

## 快速开始

### 方式一：Docker Compose

```bash
docker compose up -d --build
docker compose logs -f api
```

启动后可用以下命令检查服务：

```bash
curl http://localhost:8080/health
```

### 服务端口

| 服务 | 端口 | 说明 |
| --- | --- | --- |
| API | 8080 | REST API + ONNX 推理 |
| MySQL | 3306 | 数据库 |

## 项目结构

```text
.
├── cmd/server/main.go           # 服务入口
├── configs/                     # 配置文件
├── docs/                        # 部署、接口等文档
├── internal/                    # 业务代码
│   ├── config/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── routes/
│   └── service/
├── models/onnx/                 # ONNX 模型与标签
│   ├── herb.onnx
│   └── classes.txt
├── pkg/
│   ├── jwtutil/
│   ├── logger/
│   ├── onnx/                    # ONNX Runtime 封装
│   ├── response/
│   └── upload/
├── scripts/init.sql             # 数据库初始化脚本
├── Dockerfile                   # API 镜像
└── docker-compose.yml           # API + MySQL 编排
```

## 本地开发

### 1. 启动数据库

```bash
docker compose up -d db
```

### 2. 准备 ONNX Runtime

本地直接运行 Go 服务时，需要系统中存在 ONNX Runtime 动态库，或者通过环境变量显式指定路径：

```bash
export ONNX_RUNTIME_LIB=/path/to/libonnxruntime.so
```

Linux 安装细节见 [docs/LINUX_DEPLOY.md](docs/LINUX_DEPLOY.md)。

### 3. 启动服务

```bash
go run cmd/server/main.go
```

## 配置说明

默认配置文件：

- 本地开发：`configs/config.yaml`
- Docker 部署：`configs/config.docker.yaml`

当前服务支持通过环境变量覆盖关键配置：

| 变量名 | 说明 | 默认值 |
| --- | --- | --- |
| `SERVER_PORT` | 服务端口 | `8080` |
| `SERVER_MODE` | Gin 运行模式 | `debug` 或 `release` |
| `DB_HOST` | 数据库地址 | 本地默认 `127.0.0.1`，Docker 默认 `db` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASSWORD` | 数据库密码 | `123456` |
| `DB_NAME` | 数据库名 | `herb_recognition` |
| `ONNX_MODEL_PATH` | ONNX 模型路径 | `./models/onnx/herb.onnx` |
| `CLASSES_PATH` | 类别文件路径 | `./models/onnx/classes.txt` |
| `ONNX_RUNTIME_LIB` | ONNX Runtime 动态库路径 | 自动探测 |
| `ADMIN_USERNAME` | 默认管理员用户名 | `root` |
| `ADMIN_PASSWORD` | 默认管理员密码 | 由配置文件决定 |

## Docker 部署

### 构建与启动

```bash
docker compose build api
docker compose up -d api
```

### 查看日志

```bash
docker compose logs -f api
docker compose logs -f db
```

### 无缓存重建

```bash
docker compose build --no-cache api
docker compose up -d api
```

## 多架构镜像

如果需要同时发布 `amd64` 和 `arm64` 镜像，可以使用 `buildx`：

```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t <registry>/herb-api:latest \
  --push .
```

如果只需要导出当前架构镜像：

```bash
docker save -o herb-api.tar herbrecognition-be-api:latest
```

## 常用接口

| 模块 | 路径 | 说明 |
| --- | --- | --- |
| 健康检查 | `/health` | 服务状态 |
| 认证 | `/api/auth/*` | 登录、注册、改密 |
| 识别 | `/api/recognize/*` | 上传图片、识别、历史记录 |
| 药材 | `/api/herb/*` | 药材查询、搜索 |
| 管理端 | `/api/admin/*` | 后台管理 |

请求示例：

```bash
curl -X POST http://localhost:8080/api/recognize/upload \
  -H "Authorization: Bearer <token>" \
  -F "image=@photo.jpg"
```

## 常见问题

### 1. ONNX Runtime 动态库找不到

优先检查以下几项：

- `ONNX_RUNTIME_LIB` 是否指向正确的动态库文件
- 动态库是否与当前系统架构一致
- Docker 镜像是否使用了最新构建结果

### 2. ONNX Runtime 版本不匹配

如果出现类似以下错误：

```text
The requested API version [...] is not available
```

说明 Go binding 和动态库版本不匹配。需要保证 `github.com/yalue/onnxruntime_go` 使用的头文件版本与镜像内下载的 ONNX Runtime 版本一致。

### 3. 识别接口返回 500

优先看 `docker compose logs -f api`，常见原因有：

- ONNX Runtime 未成功初始化
- 模型文件路径错误
- 类别文件缺失
- 上传文件无法解码

## 相关文档

- [API 接口文档](docs/API.md)
- [Linux 部署说明](docs/LINUX_DEPLOY.md)
