# 测试环境部署指南

## 快速部署

### 1. 一键启动所有服务

```bash
# 复制环境变量文件
cp .env.example .env

# 启动所有服务（自动重建镜像）
docker-compose up -d --build
```

### 2. 查看服务状态

```bash
# 查看日志
docker-compose logs -f

# 查看服务状态
docker-compose ps
```

### 3. 默认管理员账号

系统启动后会自动创建默认管理员账号：

- **用户名**: `root`
- **密码**: `admin123`
- **角色**: 管理员 (admin)

> ⚠️ **重要提示**: 首次登录后请立即修改默认密码！

## API 测试

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

### 2. 登录测试

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "root", "password": "admin123"}'
```

### 3. 测试管理接口

```bash
# 获取用户列表（需要管理员权限）
curl -X GET http://localhost:8080/api/admin/user \
  -H "Authorization: Bearer <your_token>"
```

## 配置说明

### 环境变量 (.env 文件)

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DB_PASSWORD` | 数据库密码 | `123456` |
| `AUTO_CREATE_ROOT` | 是否自动创建 root 用户 | `true` |
| `PYTHON_SERVICE_URL` | Python 推理服务地址 | `http://inference:5001` |

### 禁用自动创建 root 用户

如果不需要自动创建 root 用户，可以在 `.env` 文件中设置：

```bash
AUTO_CREATE_ROOT=false
```

然后重启服务：

```bash
docker-compose down
docker-compose up -d
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| Go API | 8080 | REST API 服务 |
| Python 推理 | 5001 | 模型识别服务 |
| MySQL | 3306 | 数据库服务 |

## 常见问题

### 1. root 用户创建失败

查看 API 服务日志：

```bash
docker-compose logs api
```

### 2. 需要重置 root 用户

删除数据库卷并重启：

```bash
docker-compose down -v
docker-compose up -d --build
```

### 3. 修改默认密码

登录后通过管理接口修改密码：

```bash
curl -X POST http://localhost:8080/api/auth/change-password \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{"old_password": "admin123", "new_password": "your_new_password"}'
```

## 生产环境建议

1. **修改默认密码**: 部署后立即修改 root 用户密码
2. **使用环境变量**: 不要将敏感信息硬编码在代码中
3. **HTTPS**: 配置 HTTPS 加密通信
4. **备份数据库**: 定期备份数据库数据
5. **日志监控**: 配置日志收集和监控告警
