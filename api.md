# 后端 API 文档

本文档详细描述了后端服务的 API 接口，用于视频内容的生成和查询。后端服务封装了阿里云的 DashScope AI 视频生成模型。


## 2. API 端点

### 2.1 用户认证

#### 2.1.1 用户注册

此端点用于创建新用户账户。

- **URL**: `/api/v1/users/register`
- **方法**: `POST`
- **Content-Type**: `multipart/form-data`

##### 请求体

| 字段 | 类型 | 是否必须 | 描述 |
| :--- | :--- | :--- | :--- |
| `phone` | string | 是 | 用户的手机号码 (必须为中国大陆11位手机号)。 |
| `password` | string | 是 | 用户密码 (长度至少为8位)。 |

##### 响应体

**成功响应 (HTTP 200)**:
```json
{
  "code": 0,
  "message": "Registration successful"
}
```

**失败响应 (HTTP 400/500)**:
```json
{
  "code": 1,
  "message": "Invalid phone number format"
}
```
```json
{
  "code": 1,
  "message": "user with this phone number already exists"
}
```

---

#### 2.1.2 用户登录

此端点用于用户登录，成功后会返回一个 JWT (JSON Web Token)。

- **URL**: `/api/v1/users/login`
- **方法**: `POST`
- **Content-Type**: `multipart/form-data`

##### 请求体

| 字段 | 类型 | 是否必须 | 描述 |
| :--- | :--- | :--- | :--- |
| `phone` | string | 是 | 用户的手机号码。 |
| `password` | string | 是 | 用户密码。 |

##### 响应体

**成功响应 (HTTP 200)**:
```json
{
  "code": 0,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicGhvbmUiOiIxMzgwMDAwMDAwMCIsImV4cCI6MTcxODQzODQyMywiaXNzIjoiZW1v..."
  }
}
```

**失败响应 (HTTP 401)**:
```json
{
  "code": 1,
  "message": "invalid phone or password"
}
```

---

#### 2.1.3 修改密码

此端点用于已登录用户修改自己的密码。**此接口需要认证**。

- **URL**: `/api/v1/users/change-password`
- **方法**: `POST`
- **Content-Type**: `multipart/form-data`
- **认证**: `Authorization: Bearer <token>`

##### 请求头

| Key | Value |
| :--- | :--- |
| `Authorization` | `Bearer ` + 从登录接口获取的 JWT |

##### 请求体

| 字段 | 类型 | 是否必须 | 描述 |
| :--- | :--- | :--- | :--- |
| `newPassword` | string | 是 | 用户的新密码 (长度至少为8位)。 |

##### 响应体

**成功响应 (HTTP 200)**:
```json
{
  "code": 0,
  "message": "Password changed successfully"
}
```

**失败响应 (HTTP 401/500)**:
```json
{
  "code": 1,
  "message": "Invalid or expired JWT"
}
```

### 2.2 创建视频生成任务

- **认证**: `Authorization: Bearer <token>`

此端点用于创建一个新的视频生成任务。后端会根据 `type` 字段调用不同的 AI 模型：
- `text_to_video`: 调用 `wanx2.1-t2v-turbo` 模型。
- `image_to_video`: 调用 `wan2.2-i2v-flash` 模型。

- **URL**: `/api/v1/video/create`
- **方法**: `POST`
- **Content-Type**: `multipart/form-data`
- **认证**: `Authorization: Bearer <token>`

#### 请求体 (form-data)

| 字段 | 类型 | 是否必须 | 描述 |
| :--- | :--- | :--- | :--- |
| `type` | string | 是 | 生成类型，必须为 `text_to_video` 或 `image_to_video`。 |
| `prompt` | string | 是 | 描述视频内容的核心文本。 |
| `negative_prompt` | string | 否 | 反向提示词，用于排除不希望出现的内容。 |
| `size` | string | `type`为`text_to_video`时是 | 视频分辨率，格式为 "宽*高"。**可用值参考附录A**。 |
| `resolution` | string | `type`为`image_to_video`时是 | 视频分辨率档位。**可用值参考附录B**。 |
| `img_base64` | string | `type`为`image_to_video`时是 | 输入图片的 Base64 编码字符串，**必须为完整的 Data URI 格式**。例如: `data:image/png;base64,iVBORw0KGgo...` |

#### 响应体 (`CreateTaskResponse`)

**成功响应 (HTTP 200)**:

```json
{
  "code": 200,
  "message": "任务创建成功",
  "data": {
    "job_id": "job_xxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

**失败响应 (HTTP 400/500)**:

```json
{
  "error": "错误描述信息"
}
```

### 2.2 创建视频生成任务（高级）文生表情包（联网）

- **认证**: `Authorization: Bearer <token>`

**目前联网生成是根据角色，通过联网，用文本描述样貌之后进行生成。效果远不及图生表情。建议通过图生表情**

把上面这句话展示给用户

- **URL**: `/api/v1/video/create_with_prompt`
- **方法**: `POST`
- **Content-Type**: `multipart/form-data`
- **认证**: `Authorization: Bearer <token>`

#### 请求体 (form-data)

| 字段 | 类型 | 是否必须 | 描述 |
| :--- | :--- | :--- | :--- |
| `role` | string | 是 | 视频的核心角色名称。例如："孙悟空"、"马里奥"。 |
| `source` | string | 否 | 角色的来源，用于帮助 AI 更精确地识别。例如："七龙珠"、"任天堂游戏"。 |
| `action` | string | 是 | 角色执行的核心动作。例如："正在跳舞"、"正在奔跑"。 |
| `size` | string | 是 | 视频分辨率，格式为 "宽*高"。**可用值参考附录A**。 |

#### 响应体 (`CreateTaskResponse`)

响应结构与标准创建接口完全相同。

**成功响应 (HTTP 200)**:

```json
{
  "code": 200,
  "message": "任务创建成功",
  "data": {
    "job_id": "job_xxxxxxxxxxxxxxxxxxxxxxxx"
  }
}
```

**失败响应 (HTTP 400/500)**:

```json
{
  "error": "错误描述信息"
}
```

### 2.3 查询任务结果

- **认证**: `Authorization: Bearer <token>`

此端点用于查询指定任务的当前状态和结果。小程序端应通过轮询此接口来获取最终的视频 URL。

- **URL**: `/api/v1/video/query/:job_id`
- **方法**: `GET`
- **认证**: `Authorization: Bearer <token>`

#### URL 参数

| 参数 | 类型 | 描述 |
| :--- | :--- | :--- |
| `job_id` | string | 要查询的任务 ID，由创建任务接口返回。 |

#### 响应体 (`QueryTaskResponse`)

响应体中的 `status` 字段表示任务的当前状态。

**任务成功 (SUCCEEDED)**:
当任务成功后，后端会将生成的 `.mp4` 视频转换为 `.gif` 格式，并返回 GIF 的 URL。

```json
{
  "code": 200,
  "data": {
    "job_id": "job_xxxxxxxxxxxxxxxxxxxxxxxx",
    "status": "SUCCEEDED",
    "video_url": "https://host:port/tasks/job_xxxxxxxxxxxxxxxxxxxxxxxx.gif"
  }
}
```

**任务进行中 (RUNNING / PENDING)**:

```json
{
  "code": 200,
  "data": {
    "job_id": "job_xxxxxxxxxxxxxxxxxxxxxxxx",
    "status": "RUNNING"
  }
}
```

**任务失败 (FAILED)**:

```json
{
  "code": 200,
  "data": {
    "job_id": "job_xxxxxxxxxxxxxxxxxxxxxxxx",
    "status": "FAILED",
    "error_message": "视频生成失败的具体原因"
  }
}
```

**任务不存在 (UNKNOWN)**:

```json
{
    "code": 500,
    "message": "任务不存在",
    "data": {
        "job_id": "job_invalid_id",
        "status": "UNKNOWN"
    }
}
```

## 3. 任务状态 (Status)

| 状态 | 描述 |
| :--- | :--- |
| `PENDING` | 任务已提交，正在排队等待 AI 模型处理。 |
| `RUNNING` | 任务正在由 AI 模型处理中。 |
| `SUCCEEDED` | 任务成功完成，`video_url` 字段会包含生成的 GIF 链接。 |
| `FAILED` | 任务处理失败，`error_message` 字段会包含失败原因。 |
| `UNKNOWN` | 任务不存在或状态未知。 |

## 4. 调用流程示例

1.  **发起请求创建任务**
    小程序向 `/api/v1/video/create` 发送 `POST` 请求。

    **请求示例 (文生视频)**:
    ```json
    {
      "type": "text_to_video",
      "prompt": "一只可爱的猫在打篮球",
      "size": "960*960"
    }
    ```

2.  **获取 `job_id`**
    后端返回 `job_id`。
    ```json
    {
      "code": 200,
      "message": "任务创建成功",
      "data": {
        "job_id": "job_1234567890abcdef"
      }
    }
    ```

3.  **轮询查询任务结果**
    小程序使用获取到的 `job_id`，定期 (例如每 5-10 秒) 向 `/api/v1/video/query/job_1234567890abcdef` 发送 `GET` 请求。

4.  **处理轮询响应**
    -   如果 `status` 是 `PENDING` 或 `RUNNING`，则继续轮询。
    -   如果 `status` 是 `SUCCEEDED`，则获取 `video_url` 并展示给用户，停止轮询。
    -   如果 `status` 是 `FAILED`，则向用户显示 `error_message`，停止轮询。

---

## 附录

## 注意：由于表情包不需要太高画质，所以直接不需要用户选择，传数据时直接624*624（文生视频），480P（图生视频）

### 附录A: `size` 参数可用值 (文生视频)

`size` 参数用于 `text_to_video` 类型，基于 `wanx2.1-t2v-turbo` 模型，支持 480P 和 720P 档位。

**480P 档位**:
- `"832*480"` (16:9)
- `"480*832"` (9:16)
- `"624*624"` (1:1)

**720P 档位 (默认)**:
- `"1280*720"` (16:9)
- `"720*1280"` (9:16)
- `"960*960"` (1:1)
- `"1088*832"` (4:3)
- `"832*1088"` (3:4)

### 附录B: `resolution` 参数可用值 (图生视频)

`resolution` 参数用于 `image_to_video` 类型，基于 `wan2.2-i2v-flash` 模型。

- `"480P"`
- `"720P"` (默认)