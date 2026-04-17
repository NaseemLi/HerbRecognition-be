# 变更日志

## [2026-03-01] - 安全加固与代码重构

### 新增

#### 配置项

- `jwt.secret` - JWT 密钥配置，支持从配置文件读取
- `jwt.expire_hours` - Token 过期时间配置（默认 168 小时/7 天）
- `cors.allow_origins` - CORS 白名单配置，支持多域名
- `admin.username` - 默认管理员用户名配置
- `admin.password` - 默认管理员密码配置（留空自动生成随机密码）
- `database.max_idle_conns` - 数据库最大空闲连接数
- `database.max_open_conns` - 数据库最大打开连接数
- `database.conn_max_lifetime` - 数据库连接最大生命周期（秒）

#### 工具包

- `pkg/jwtutil/jwt.go` - JWT 统一工具包
  - `GenerateToken()` - 生成 JWT Token
  - `ParseToken()` - 解析 JWT Token
- `pkg/upload/upload.go` - 文件上传统一工具包
  - `UploadFile()` - 通用文件上传（支持配置化）
  - `DeleteFile()` - 删除文件
  - `DefaultImageConfig` - 默认图片上传配置
  - `AvatarConfig` - 头像上传配置

### 变更

#### 安全性改进

- **JWT 密钥**: 从硬编码改为配置文件读取，不再在代码中暴露密钥
- **CORS 配置**: 从 `Access-Control-Allow-Origin: *` 改为白名单机制
- **默认密码**:
  - 支持从配置文件设置默认管理员密码
  - 密码为空时自动生成 12 位随机密码
  - 日志中仅在首次创建时显示随机密码，后续不再打印明文

#### 架构改进

- **分层架构修复**: `handler/app/user.go` 中的 `GetProfile` 方法不再直接访问数据库，改为调用 Service 层
- **数据库连接池**: 添加连接池配置支持

#### 代码重构

- **文件上传逻辑**: 抽取公共工具函数，消除三处重复代码
  - `service/auth_service.go` - UploadAvatar
  - `service/recognize_service.go` - UploadImage
  - `service/admin_service.go` - UploadAndSetImage

### 移除

- `middleware/auth.go` 中硬编码的 `jwtSecret` 变量
- `service/auth_service.go` 中硬编码的 `jwtSecret` 变量和 `generateToken` 函数
- 各 Service 中重复的文件上传验证逻辑

### 配置迁移指南

旧配置 (`config.yaml`):

```yaml
server:
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: 123456
  dbname: herb_recognition
  charset: utf8mb4
  parseTime: true
  loc: Local

model_service:
  url: "http://127.0.0.1:5000"
```

新配置 (`config.yaml`):

```yaml
server:
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: 123456
  dbname: herb_recognition
  charset: utf8mb4
  parseTime: true
  loc: Local
  max_idle_conns: 10 # 新增
  max_open_conns: 100 # 新增
  conn_max_lifetime: 3600 # 新增

model_service:
  url: "http://127.0.0.1:5000"

# 新增 JWT 配置
jwt:
  secret: "your-super-secret-key-change-in-production"
  expire_hours: 168

# 新增 CORS 配置
cors:
  allow_origins:
    - "http://localhost:3000"
    - "http://localhost:5173"
    - "http://127.0.0.1:3000"
    - "http://127.0.0.1:5173"

# 新增管理员配置
admin:
  username: "root"
  password: "" # 留空自动生成随机密码
```

### 生产环境部署注意事项

1. **JWT 密钥**: 务必修改 `jwt.secret` 为强随机字符串（建议 32 位以上）
2. **CORS 白名单**: 配置实际的前端域名，不要使用开发环境的 localhost
3. **管理员密码**: 建议设置固定密码或记录自动生成的随机密码
4. **数据库连接池**: 根据服务器配置调整连接池参数

### API 兼容性

本次更新**不影响 API 接口设计**，所有接口路径、参数、响应格式保持不变。

---

## 历史版本

(后续版本更新记录将追加在此处)
