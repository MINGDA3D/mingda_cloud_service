## 更新用户信息
允许已授权的用户通过此接口更新自己的信息。

:::tips
+ **URL**：`/api/user`
+ **Method**：`PUT`
+ **需要登录**：<font style="background:#F6E1AC;color:#664900">是</font>
+ **需要鉴权**：<font style="background:#F6E1AC;color:#664900">是</font>

:::

### 请求参数
| 参数 | 类型 | 约束 |
| --- | --- | --- |
| `first_name` | String | 1 到 30 个字符 |
| `last_name` | String | 1 到 30 个字符 |


:::info
💡 注意，id 和 email 字段目前是只读属性，不允许通过此接口进行修改。

:::

### 请求示例
可以仅传递部分请求参数。

```json
{
    "first_name": "John"
}
```

可以通过传递空字符串来清除数据。

```json
{
    "last_name": ""
}
```

### 成功响应
:::tips
**条件**：请求参数合法，并且用户身份校验通过。

**状态码：**`200 OK`

**响应示例**：响应会将修改后的用户信息数据返回，一个`id`为 1234 的用户设置他们的姓名后将会返回：

:::

```json
{
    "id": 1234,
    "first_name": "Joe",
    "last_name": "Bloggs",
    "email": "joe25@example.com"
}
```

### 错误响应
:::tips
**条件**：请求数据非法，例如 fist_name 的长度过长。

**状态码**：`400 BAD REQUEST`

响应示例:

:::

```json
{
    "first_name": "Please provide maximum 30 character or empty string"
}
```

### 注意事项
:::info
💡 如果用户的用户信息不存在，将会使用请求的数据创建一个新的用户信息。

:::





### 接口定义
#### 1. 初始化云服务连接
**接口名称**: `InitCloudService`

+ **描述**: 首次开机验证云服务权限。
+ **请求地址**: `/api/initCloudService`
+ **请求方法**: POST
+ **请求参数**:
    - `sn`: 设备序列号（String）
    - `model`: 机型信息（String）
    - `version`: 客户端版本号（String）
+ **响应格式**:
    - 成功: `{ "status": "success", "message": "Verification successful." }`
    - 失败: `{ "status": "error", "message": "Invalid SN or model." }`

#### 2. 图像预处理与上传
**接口名称**: `PreprocessAndUploadImage`

+ **描述**: 对采集到的图像进行初步处理，并将其上传至云端。
+ **请求地址**: `/api/preprocessAndUploadImage`
+ **请求方法**: POST
+ **请求头**:
    - `Authorization`: Bearer 
    - `Content-Type`: multipart/form-data
+ **请求参数**:
    - `imageFile`: 图像文件（File）
    - `userId`: 用户ID（String）
    - `deviceId`: 设备ID（String）
+ **响应格式**:
    - 成功: `{ "status": "success", "taskId": "<unique_task_id>" }`
    - 失败: `{ "status": "error", "message": "Upload failed." }`

#### 3. 云端接收并确认
**接口名称**: `ConfirmCloudReceipt`

+ **描述**: 云端接收到上传的数据包后，对其进行完整性检查和初步验证。
+ **请求地址**: `/api/confirmCloudReceipt`
+ **请求方法**: GET
+ **请求参数**:
    - `taskId`: 唯一任务ID（String）
+ **响应格式**:
    - 成功: `{ "status": "success", "message": "Data received and verified." }`
    - 失败: `{ "status": "error", "message": "Data verification failed." }`

#### 4. 推入Redis推理请求队列
**接口名称**: `PushToInferenceQueue`

+ **描述**: 将经过验证的任务正式推送到名为`inference_queue`的列表中。
+ **请求地址**: `/api/pushToInferenceQueue`
+ **请求方法**: POST
+ **请求参数**:
    - `taskId`: 唯一任务ID（String）
    - `data`: JSON格式的任务数据（Object）
+ **响应格式**:
    - 成功: `{ "status": "success", "message": "Task pushed to queue." }`
    - 失败: `{ "status": "error", "message": "Failed to push task." }`

#### 5. 设置超时时间 & 启动定时器
**接口名称**: `SetTimeoutAndStartTimer`

+ **描述**: 设置合理的超时时间和启动定时器以监测云端推理过程中的超时情况。
+ **请求地址**: `/api/setTimeoutAndStartTimer`
+ **请求方法**: POST
+ **请求参数**:
    - `taskId`: 唯一任务ID（String）
    - `timeoutSeconds`: 超时秒数（Integer）
+ **响应格式**:
    - 成功: `{ "status": "success", "message": "Timer started." }`
    - 失败: `{ "status": "error", "message": "Failed to start timer." }`

#### 6. 终端定时查询
**接口名称**: `QueryInferenceResult`

+ **描述**: 终端设备根据自身的业务逻辑设定一定的周期主动向Redis发起查询请求。
+ **请求地址**: `/api/queryInferenceResult`
+ **请求方法**: GET
+ **请求参数**:
    - `taskId`: 唯一任务ID（String）
+ **响应格式**:
    - 成功: `{ "status": "success", "result": "<inference_result>", "timestamp": "<timestamp>" }`
    - 失败: `{ "status": "error", "message": "No result found." }`

#### 7.终端推理
**接口名称**:predict

+ **描述**: 终端推理接口。
+ **请求地址**: `localhost:5000/predict`
+ **请求方法**: PUT
+ **请求参数**:
    - `file`: 需要推理的图片（File）
+ **响应格式**:
    - 成功: `{"results":[{"bbox":[2202.1162109375,1116.47900390625,3997.7099609375,3059.36376953125],"class":"spaghetti","confidence":0.25059378147125244}]}`
    - 失败: `{"error":"Invalid image file"}`

### 编写接口文档的基本结构
在编写具体的接口文档时，应该遵循以下基本结构：

1. **接口简介**：简要介绍该接口的功能及作用。
2. **请求地址**：提供API的具体URL路径。
3. **请求方法**：指明HTTP请求类型（GET, POST等）。
4. **请求头**（如果适用）：列出所有必要的HTTP头部信息。
5. **请求参数**：详细说明每一个输入参数的意义、类型以及是否必需。
6. **响应格式**：定义返回值的格式，包括成功和失败两种情形下的预期输出。
7. **错误码**：列出可能遇到的各种错误代码及其含义。
8. **示例**：给出完整的请求示例和相应的响应示例，以便开发者能够快速理解和测试。

此外，为了确保接口文档的质量，还需要注意以下几点：

+ **清晰准确**：避免使用模糊不清的语言，确保每个字段都有明确的定义。
+ **易于查找**：组织好文档的内容，使用户可以轻松定位所需信息。
+ **持续更新**：随着项目的迭代发展，及时修正和完善文档内容，保持其时效性。
+ **工具辅助**：利用Swagger、Apifox等工具自动生成部分文档内容，减少手动编写的负担。







