# 表情包生成器 API 文档

## 1. 用户认证接口

### 1.1 用户注册

**接口地址**: `POST /api/user/register`

**请求参数**:
```json
{
  "username": "string",      // 用户名
  "password": "string",      // 密码
  "email": "string"          // 邮箱（可选）
}
```

**响应结果**:
```json
{
  "code": 0,                 // 状态码，0-成功，其他-失败
  "message": "string",       // 返回信息
  "data": {
    "userId": "string",      // 用户ID
    "username": "string"     // 用户名
  }
}
```

### 1.2 用户登录

**接口地址**: `POST /api/user/login`

**请求参数**:
```json
{
  "username": "string",      // 用户名
  "password": "string"       // 密码
}
```

**响应结果**:
```json
{
  "code": 0,                 // 状态码，0-成功，其他-失败
  "message": "string",       // 返回信息
  "data": {
    "userId": "string",      // 用户ID
    "username": "string",    // 用户名
    "token": "string"        // 登录凭证
  }
}
```

### 1.3 修改密码

**接口地址**: `POST /api/user/change-password`

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "oldPassword": "string",   // 原密码
  "newPassword": "string"    // 新密码
}
```

**响应结果**:
```json
{
  "code": 0,                 // 状态码，0-成功，其他-失败
  "message": "string"        // 返回信息
}
```

## 2. 表情包生成接口

### 2.1 上传参考图生成表情包

**接口地址**: `POST /api/meme/generate-from-upload`

**请求头**:
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**请求参数**:
```
file: 文件                // 参考图文件
topText: string           // 顶部文字（可选）
bottomText: string        // 底部文字（可选）
fontSize: number          // 字体大小（可选）
fontColor: string         // 字体颜色（可选）
```

**响应结果**:
```json
{
  "code": 0,                 // 状态码，0-成功，其他-失败
  "message": "string",       // 返回信息
  "data": {
    "memeId": "string",      // 表情包ID
    "imageUrl": "string"     // 生成的表情包URL
  }
}
```

### 2.2 联网生成表情包

**接口地址**: `POST /api/meme/generate-from-network`

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "keywords": "string",      // 关键词
  "topText": "string",       // 顶部文字（可选）
  "bottomText": "string",    // 底部文字（可选）
  "fontSize": 20,            // 字体大小（可选）
  "fontColor": "#FFFFFF"     // 字体颜色（可选）
}
```

**响应结果**:
```json
{
  "code": 0,                 // 状态码，0-成功，其他-失败
  "message": "string",       // 返回信息
  "data": {
    "memeId": "string",      // 表情包ID
    "imageUrl": "string"     // 生成的表情包URL
  }
}
```

### 2.3 不联网生成表情包

**接口地址**: `POST /api/meme/generate-from-template`

**请求头**:
```
Authorization: Bearer <token>
```

**请求参数**:
```json
{
  "templateId": "string",    // 模板ID
  "topText": "string",       // 顶部文字（可选）
  "bottomText": "string",    // 底部文字（可选）
  "fontSize": 20,            // 字体大小（可选）
  "fontColor": "#FFFFFF"     // 字体颜色（可选）
}
```

**响应结果**:
```json
{
  "code": 0,                 // 状态码，0-成功，其他-失败
  "message": "string",       // 返回信息
  "data": {
    "memeId": "string",      // 表情包ID
    "imageUrl": "string"     // 生成的表情包URL
  }
}
```