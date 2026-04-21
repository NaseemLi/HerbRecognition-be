# API 接口文档

**基础信息**
- Base URL: `http://localhost:8080`
- 认证方式：JWT Token，Header: `Authorization: Bearer <token>`
- 响应格式：统一 JSON 格式

## 响应格式

```json
{
  "code": 200,
  "message": "成功",
  "data": {}
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 状态码（200=成功，400=参数错误，401=未授权，500=服务器错误） |
| message | string | 响应消息 |
| data | object/array | 响应数据 |

---

## 认证模块 `/api/auth`

### 1. 用户注册
- **接口**: `POST /api/auth/register`
- **权限**: 公开
- **请求参数**:
```json
{
  "username": "string (3-32 字符，必填)",
  "password": "string (最少 6 字符，必填)"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "注册成功",
  "data": null
}
```

### 2. 用户登录
- **接口**: `POST /api/auth/login`
- **权限**: 公开
- **请求参数**:
```json
{
  "username": "string (必填)",
  "password": "string (必填)"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "role": "user",
      "avatar": "",
      "created_at": "2026-02-23T06:01:45.365Z"
    }
  }
}
```

### 3. 修改密码
- **接口**: `POST /api/auth/change-password`
- **权限**: 需登录
- **请求参数**:
```json
{
  "old_password": "string (必填)",
  "new_password": "string (最少 6 字符，必填)"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "密码修改成功",
  "data": null
}
```

---

## 用户模块 `/api/user`

### 1. 获取用户资料
- **接口**: `GET /api/user/profile`
- **权限**: 需登录
- **响应**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "role": "user",
      "avatar": "/uploads/avatars/avatar_xxx.jpg",
      "created_at": "2026-02-23T06:01:45.365Z"
    }
  }
}
```

### 2. 更新用户资料
- **接口**: `PUT /api/user/profile`
- **权限**: 需登录
- **请求参数**: `application/json`
```json
{
  "username": "string (3-32 字符，可选)",
  "avatar": "string (头像URL，可选)"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "资料更新成功",
  "data": {
    "user": {
      "id": 1,
      "username": "newusername",
      "role": "user",
      "avatar": "/uploads/avatars/avatar_xxx.jpg",
      "created_at": "2026-02-23T06:01:45.365Z"
    }
  }
}
```

### 3. 上传头像
- **接口**: `POST /api/user/avatar`
- **权限**: 需登录
- **请求**: `multipart/form-data`
  - `avatar`: 图片文件（必填）
  - 支持格式：JPG, PNG, GIF, WEBP
  - 最大大小：2MB
- **响应**:
```json
{
  "code": 200,
  "message": "头像上传成功",
  "data": {
    "avatar_url": "/uploads/avatars/avatar_xxx.jpg"
  }
}
```

---

## 识别模块 `/api/recognize`

### 1. 上传图片并识别
- **接口**: `POST /api/recognize/upload`
- **权限**: 需登录
- **请求参数**: `multipart/form-data`
  - `image`: 图片文件（必填）
  - 支持格式：JPG, PNG, GIF, WEBP
  - 最大大小：5MB
- **响应**:
```json
{
  "code": 200,
  "message": "识别成功",
  "data": {
    "record_id": 1,
    "herb_id": 6,
    "herb_name": "鸡内金",
    "confidence": 98.5,
    "image_url": "/uploads/images/xxx.jpg"
  }
}
```

### 2. 获取识别历史
- **接口**: `GET /api/recognize/history`
- **权限**: 需登录
- **Query 参数**:
  - `page`: 页码（默认 1）
  - `page_size`: 每页数量（默认 10，最大 50）
- **响应**:
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "id": 1,
        "image_url": "/uploads/images/xxx.jpg",
        "herb_id": 6,
        "herb_name": "鸡内金",
        "confidence": 98.5,
        "created_at": "2026-02-23T06:01:45.365Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 10
  }
}
```

### 3. 删除识别记录
- **接口**: `DELETE /api/recognize/history`
- **权限**: 需登录
- **请求参数**: `application/json`
```json
{
  "id": 1
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

---

## 药材查询模块 `/api/herb`

### 1. 获取药材列表
- **接口**: `GET /api/herb`
- **权限**: 需登录
- **Query 参数**:
  - `category`: 分类筛选（可选）
  - `page`: 页码（默认 1）
  - `page_size`: 每页数量（默认 10，最大 50）
- **响应**:
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "鸡内金",
        "scientific": "Endothelium Corneum Gigeriae Galli",
        "alias": "鸡肫皮，鸡黄皮",
        "category": "消食药",
        "description": "鸡的干燥沙囊内壁",
        "effects": "健胃消食，涩精止遗，通淋化石",
        "usage": "内服：煎汤，3-10g",
        "image_url": "/uploads/herbs/xxx.jpg"
      }
    ],
    "total": 20,
    "page": 1,
    "page_size": 10
  }
}
```

### 2. 搜索药材
- **接口**: `GET /api/herb/search`
- **权限**: 需登录
- **Query 参数**:
  - `keyword`: 关键词（必填，支持名称/别名/学名模糊搜索）
  - `page`: 页码（默认 1）
  - `page_size`: 每页数量（默认 10，最大 50）
- **响应**: 同 `获取药材列表`

### 3. 获取药材详情
- **接口**: `GET /api/herb/:id`
- **权限**: 需登录
- **路径参数**: `id` - 药材 ID
- **响应**:
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": 1,
    "name": "鸡内金",
    "scientific": "Endothelium Corneum Gigeriae Galli",
    "alias": "鸡肫皮，鸡黄皮",
    "category": "消食药",
    "description": "鸡的干燥沙囊内壁",
    "effects": "健胃消食，涩精止遗，通淋化石",
    "usage": "内服：煎汤，3-10g",
    "image_url": "/uploads/herbs/xxx.jpg"
  }
}
```

---

## 管理后台模块 `/api/admin`

**所有管理接口需要 `admin` 角色权限**

### 药材管理

#### 1. 获取药材列表
- **接口**: `GET /api/admin/herb`
- **Query 参数**:
  - `category`: 分类筛选（可选）
  - `page`: 页码（默认 1）
  - `page_size`: 每页数量（默认 10，最大 50）
- **响应**: 同 `获取药材列表`

#### 2. 创建药材
- **接口**: `POST /api/admin/herb`
- **请求参数**: `application/json`
```json
{
  "name": "string (必填，1-64 字符)",
  "scientific": "string (可选)",
  "alias": "string (可选)",
  "category": "string (可选)",
  "description": "string (可选)",
  "effects": "string (可选)",
  "usage": "string (可选)",
  "image_url": "string (可选，图片 URL)"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "鸡内金",
    "scientific": "...",
    "alias": "...",
    "category": "...",
    "description": "...",
    "effects": "...",
    "usage": "...",
    "image_url": ""
  }
}
```

#### 3. 更新药材
- **接口**: `PUT /api/admin/herb`
- **请求参数**: `application/json`
```json
{
  "id": 1,
  "name": "string (必填)",
  "scientific": "string",
  "alias": "string",
  "category": "string",
  "description": "string",
  "effects": "string",
  "usage": "string",
  "image_url": "string (可选，图片 URL)"
}
```
- **响应**: 同 `创建药材`

#### 4. 删除药材
- **接口**: `DELETE /api/admin/herb`
- **请求参数**: `application/json`
```json
{
  "id": 1
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```

#### 5. 批量删除药材
- **接口**: `DELETE /api/admin/herb/batch`
- **请求参数**: `application/json`
```json
{
  "ids": [1, 2, 3]
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "批量删除成功",
  "data": null
}
```

#### 6. 上传药材图片
- **接口**: `POST /api/admin/herb/upload-image`
- **请求**: `multipart/form-data`
  - `image`: 图片文件
- **响应**:
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "image_url": "/uploads/herbs/xxx.jpg"
  }
}
```

---

### 用户管理

#### 1. 获取用户列表
- **接口**: `GET /api/admin/user`
- **Query 参数**:
  - `role`: 角色筛选（`user` 或 `admin`）
  - `page`: 页码（默认 1）
  - `page_size`: 每页数量（默认 10，最大 50）
- **响应**:
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "testuser",
        "role": "user",
        "avatar": ""
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 10
  }
}
```

#### 2. 修改用户角色
- **接口**: `POST /api/admin/user/role`
- **请求参数**: `application/json`
```json
{
  "user_id": 1,
  "role": "admin"
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "修改成功",
  "data": null
}
```

#### 3. 删除用户
- **接口**: `DELETE /api/admin/user`
- **请求参数**: `application/json`
```json
{
  "user_id": 1
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```
- **错误说明**:
  - `不能删除自己` — 管理员不能删除自己的账号
  - `不能删除管理员` — 不能删除其他管理员账号
  - `用户不存在` — 目标用户不存在

#### 4. 批量删除用户
- **接口**: `DELETE /api/admin/user/batch`
- **请求参数**: `application/json`
```json
{
  "user_ids": [1, 2, 3]
}
```
- **响应**:
```json
{
  "code": 200,
  "message": "删除成功",
  "data": null
}
```
- **错误说明**:
  - `不能删除自己` — 列表中包含管理员自己的 ID
  - `不能删除管理员` — 列表中包含其他管理员账号
  - `部分用户不存在` — 列表中有不存在的用户 ID

---

## 健康检查

### 健康检查
- **接口**: `GET /health`
- **权限**: 公开
- **响应**:
```json
{
  "code": 200,
  "message": "健康",
  "data": null
}
```

---

## 静态资源

### 访问上传的图片
- **路径**: `/uploads/*filepath`
- **示例**: `http://localhost:8080/uploads/images/xxx.jpg`

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权（未登录或 token 无效） |
| 403 | 禁止访问（权限不足） |
| 500 | 服务器内部错误 |

---

## 前端开发提示

1. **Token 存储**: 登录成功后将 `token` 存入 localStorage/sessionStorage
2. **请求拦截**: 所有需要认证的请求添加 Header: `Authorization: Bearer <token>`
3. **图片上传**: 使用 `FormData`，字段名为 `image`
4. **分页**: 所有列表接口支持分页，默认每页 10 条
5. **角色判断**: `user.role === 'admin'` 判断是否为管理员
