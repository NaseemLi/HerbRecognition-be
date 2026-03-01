# 中草药识别系统后端

基于 Go + Gin + GORM + PyTorch 的中草药图像识别后端服务

## 技术栈

| 组件         | 技术                    |
| ------------ | ----------------------- |
| **Web 框架** | Go 1.25 + Gin           |
| **数据库**   | MySQL 8.0 + GORM        |
| **认证**     | JWT (golang-jwt)        |
| **配置管理** | Viper                   |
| **日志**     | Zap                     |
| **AI 推理**  | Python + PyTorch        |
| **容器化**   | Docker + Docker Compose |

## 快速开始

### 环境要求

- Docker & Docker Compose
- Go 1.25+ (本地开发)
- Python 3.10+ (本地开发)

### 一键启动

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 服务端口

| 服务        | 端口 | 说明         |
| ----------- | ---- | ------------ |
| Go API      | 8080 | REST API     |
| Python 推理 | 5001 | 模型推理服务 |
| MySQL       | 3306 | 数据库       |

## 项目结构

```
.
├── cmd/                      # 应用入口
│   └── server/
│       └── main.go
├── internal/                 # 内部业务逻辑
│   ├── handler/             # HTTP 处理器
│   ├── service/             # 业务服务层
│   ├── repository/          # 数据访问层
│   ├── model/               # 数据模型
│   ├── routes/              # 路由配置
│   ├── middleware/          # 中间件
│   ├── client/              # 外部服务客户端
│   └── config/              # 配置加载
├── pkg/                      # 公共包
│   ├── response/            # 统一响应格式
│   ├── jwtutil/             # JWT 工具包
│   ├── upload/              # 文件上传工具包
│   ├── logger/              # 日志工具
│   └── errors/              # 错误定义
├── models/                   # AI 模型文件
│   └── best_herb_model.pth
├── services/                 # Python 推理服务
│   └── inference-service/
├── scripts/                  # 脚本文件
│   └── init.sql             # 数据库初始化脚本
├── configs/                  # 配置文件
├── uploads/                  # 上传文件存储
├── Dockerfile                # Go API 镜像
├── Dockerfile.python         # Python 推理镜像
└── docker-compose.yml        # 服务编排
```

## API 文档

详细接口文档请参考 [API.md](API.md)

### 主要接口

| 模块 | 路径               | 说明                      |
| ---- | ------------------ | ------------------------- |
| 认证 | `/api/auth/*`      | 注册、登录、修改密码      |
| 识别 | `/api/recognize/*` | 上传图片识别、历史记录    |
| 药材 | `/api/herb/*`      | 药材查询、搜索            |
| 管理 | `/api/admin/*`     | 后台管理（需 admin 权限） |

### 请求示例

```bash
# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}'

# 识别图片
curl -X POST http://localhost:8080/api/recognize/upload \
  -H "Authorization: Bearer <token>" \
  -F "image=@photo.jpg"

# 搜索药材
curl -X GET "http://localhost:8080/api/herb/search?keyword=枸杞" \
  -H "Authorization: Bearer <token>"
```

## 开发指南

### 本地开发

```bash
# 启动数据库和推理服务
docker-compose up -d db inference

# 本地运行 Go 服务
go run cmd/server/main.go
```

### 数据库配置

默认数据库配置：

- Host: `localhost` (本地) / `db` (Docker)
- Port: `3306`
- User: `root`
- Password: `123456`
- Database: `herb_recognition`

### 环境变量

| 变量名               | 说明                | 默认值                  |
| -------------------- | ------------------- | ----------------------- |
| `PYTHON_SERVICE_URL` | Python 推理服务地址 | `http://localhost:5001` |
| `DB_HOST`            | 数据库地址          | `localhost`             |
| `DB_PASSWORD`        | 数据库密码          | `123456`                |

## Docker 部署

### 构建镜像

```bash
# 构建所有服务
docker-compose build

# 单独构建
docker-compose build api      # Go API
docker-compose build inference # Python 推理
```

### 生产环境建议

1. 修改默认数据库密码
2. 使用 `.env` 文件管理敏感配置
3. 配置 HTTPS
4. 设置日志持久化

## 常见问题

### 服务启动失败

```bash
# 查看日志
docker-compose logs api
docker-compose logs inference

# 检查端口占用
lsof -i :8080
lsof -i :5001
lsof -i :3306

# 重启服务
docker-compose restart
```

### 重新构建

```bash
# 强制重新构建（不使用缓存）
docker-compose build --no-cache
docker-compose up -d
```

## 相关文档

- [API 接口文档](docs/API.md)
