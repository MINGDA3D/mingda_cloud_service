# 3D打印机数据采集系统 API文档 v1.0
## 基础信息
| 项目 | 说明 |
| --- | --- |
| 基础URL | [https://sptdata.3dmingda.com/api/v1](https://sptdata.3dmingda.com/api/v1) |
| 传输协议 | HTTP |
| 请求方式 | POST |
| 数据格式 | JSON |
| 加密方式 | AES-256 |
| 字符编码 | UTF-8 |


## 通用请求头
| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| Content-Type | string | 是 | 固定值: application/json |
| Authorization | string | 是 | 格式: Bearer {token} |
| Device-SN | string | 是 | 设备序列号 |
| Timestamp | number | 是 | 请求时间戳(毫秒) |
| Sign | string | 是 | 请求签名 |


## API接口列表
### 1. 设备基础信息接口
#### 1.1 设备信息上报
> 上报设备基础信息和软件版本信息
>

**请求URL**

```plain
POST /device/info
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| device_info | object | 是 | 设备基础信息 |
| └─ device_sn | string | 是 | 设备SN码 |
| └─ device_model | string | 是 | 机型 |
| └─ hardware_version | string | 是 | 硬件版本号 |
| software_versions | object | 是 | 软件版本信息 |
| └─ klipper | string | 是 | Klipper版本 |
| └─ klipper_screen | string | 是 | KlipperScreen版本 |
| └─ moonraker | string | 否 | moonraker版本 |
| └─ mainsail | string | 否 | mainsail版本 |
| └─ crowsnest | string | 否 | crowsnest版本 |
| └─ firmware | object | 是 | 固件版本信息 |
|     └─ mainboard | string | 是 | 主板固件版本 |
|     └─ printhead | string | 是 | 打印头板固件版本 |
|     └─ leveling | string | 否 | 快速调平板固件版本 |


**响应参数**

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| code | number | 状态码 |
| message | string | 状态信息 |
| data | object | 响应数据 |
| └─ success | boolean | 是否成功 |


**请求示例**

```json
{
    "device_info": {
        "device_sn": "MD123456789",
        "device_model": "MD-400D",
        "hardware_version": "V1.0"
    },
    "software_versions": {
        "klipper": "v0.11.0",
        "klipper_screen": "v1.0.0",
        "moonraker": "v0.8.0",
        "mainsail": "v2.5.0",
        "crowsnest": "v1.0.0",
        "firmware": {
            "mainboard": "v1.2.3",
            "printhead": "v1.0.1",
            "leveling": "v1.1.0"
        }
    }
}
```

**响应示例**

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "success": true
    }
}
```

### 2. 设备状态接口
#### 2.1 网络状态上报
> 上报设备网络状态信息
>

**请求URL**

```plain
POST /device/network
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| device_sn | string | 是 | 设备SN码 |
| ip_address | string | 是 | IP地址 |
| mac_address | string | 是 | MAC地址 |
| network_type | string | 否 | 网络类型(wifi/ethernet) |
| signal_strength | number | 否 | 信号强度(wifi时必填) |


#### 2.2 系统状态上报
> 上报设备系统资源使用状态
>

**请求URL**

```plain
POST /device/status
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| device_sn | string | 是 | 设备SN码 |
| storage | object | 是 | 存储信息 |
| └─ total | number | 是 | 总存储空间(MB) |
| └─ used | number | 是 | 已使用空间(MB) |
| └─ free | number | 是 | 剩余空间(MB) |
| cpu | object | 是 | CPU信息 |
| └─ usage | number | 是 | CPU使用率(%) |
| └─ temperature | number | 是 | CPU温度(℃) |
| memory | object | 是 | 内存信息 |
| └─ total | number | 是 | 总内存(MB) |
| └─ used | number | 是 | 已使用内存(MB) |
| └─ free | number | 是 | 剩余内存(MB) |


**响应参数**

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| code | number | 状态码 |
| message | string | 状态信息 |
| data | object | 响应数据 |
| └─ success | boolean | 是否成功 |


**请求示例**

```json
{
    "device_sn": "MD123456789",
    "storage": {
        "total": 32768,
        "used": 15360,
        "free": 17408
    },
    "cpu": {
        "usage": 35.6,
        "temperature": 45.2
    },
    "memory": {
        "total": 2048,
        "used": 1024,
        "free": 1024
    }
}
```

**响应示例**

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "success": true
    }
}
```

### 3. 告警管理接口
#### 3.1 告警数据上报
> 上报设备告警信息
>

**请求URL**

```plain
POST /device/alarm
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| device_sn | string | 是 | 设备SN码 |
| error_code | string | 是 | Klipper错误码 |
| error_message | string | 是 | 错误描述 |
| severity | string | 是 | 严重程度(error/warning/info) |
| alarm_time | string | 是 | 告警发生时间(ISO8601格式) |
| details | object | 否 | 详细信息(JSON格式) |


**响应参数**

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| code | number | 状态码 |
| message | string | 状态信息 |
| data | object | 响应数据 |
| └─ success | boolean | 是否成功 |
| └─ alarm_id | string | 告警ID |


**请求示例**

```json
{
    "device_sn": "MD123456789",
    "error_code": "E1001",
    "error_message": "打印头温度异常",
    "severity": "error",
    "alarm_time": "2024-03-14T10:30:00Z",
    "details": {
        "temperature": 280,
        "target": 200,
        "component": "extruder"
    }
}
```

**响应示例**

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "success": true,
        "alarm_id": "ALM202403141030123456"
    }
}
```

### 4. 打印任务接口
#### 4.1 打印任务状态上报
> 上报打印任务状态信息
>

**请求URL**

```plain
POST /device/print/status
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| device_sn | string | 是 | 设备SN码 |
| task_id | string | 是 | 打印任务ID |
| file_name | string | 是 | 打印文件名 |
| status | string | 是 | 任务状态 |
| start_time | string | 否 | 开始时间(ISO8601格式) |
| end_time | string | 否 | 结束时间(ISO8601格式) |
| progress | number | 否 | 打印进度(%) |
| duration | number | 是 | 打印时长(秒) |
| filament_used | number | 否 | 耗材使用量(mm) |
| layers_completed | number | 否 | 已完成层数 |
| error_code | string | 否 | 错误码(如果因错误中断) |
| cancel_reason | string | 否 | 取消原因 |


**状态说明**

| 状态值 | 说明 |
| --- | --- |
| idle | 空闲状态 |
| printing | 打印中 |
| paused | 已暂停 |
| resumed | 已恢复 |
| completed | 已完成 |
| cancelled | 已取消 |
| error | 错误中断 |


**响应参数**

| 参数 | 类型 | 说明 |
| --- | --- | --- |
| code | number | 状态码 |
| message | string | 状态信息 |
| data | object | 响应数据 |
| └─ success | boolean | 是否成功 |


**请求示例**

```json
{
    "device_sn": "MD123456789",
    "task_id": "PT123456789123456",
    "file_name": "test_model.gcode",
    "status": "printing",
    "start_time": "2024-03-14T10:00:00Z",
    "progress": 45.5,
    "duration": 1800,
    "filament_used": 1250.5,
    "layers_completed": 150
}
```

**响应示例**

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "success": true
    }
}
```

### 1. 定时上报策略
| 数据类型 | 上报间隔 | 说明 |
| --- | --- | --- |
| 网络状态 | 10分钟 | 包含IP、网络类型、信号强度等 |
| 系统状态 | 5分钟 | 包含CPU、内存、存储等资源使用情况 |
| 软件版本 | 24小时 | 或版本发生变更时立即上报 |


### 2. 事件触发上报
| 事件类型 | 触发条件 | 说明 |
| --- | --- | --- |
| 告警信息 | 实时 | 错误发生时立即上报 |
| 打印状态 | 状态变化 | 状态发生改变时上报 |
| 版本变更 | 变更时 | 软件或固件版本更新后上报 |


### 3. 数据优先级
| 优先级 | 数据类型 | 说明 |
| --- | --- | --- |
| 高 | 告警数据 | 网络恢复后优先上报 |
| 中 | 打印状态 | 次优先级上报 |
| 低 | 系统状态 | 最后上报 |


## 安全认证机制
### 1. 设备认证
**认证流程**

1. 设备首次连接时进行认证
2. 使用设备SN码和密钥生成认证请求
3. 服务器验证通过后返回token
4. 后续请求使用token进行认证

**Token格式**

```json
{
    "device_sn": "设备序列号",
    "iat": "签发时间",
    "exp": "过期时间",
    "type": "token类型"
}
```

### 2. 数据加密
**AES-256加密流程**

```typescript
// 1. 生成随机AES密钥
const aesKey = crypto.randomBytes(32);

// 2. 使用RSA加密AES密钥
const encryptedKey = crypto.publicEncrypt(
    SERVER_PUBLIC_KEY,
    aesKey
);

// 3. 使用AES加密数据
const encryptedData = crypto.createCipheriv(
    'aes-256-gcm',
    aesKey,
    iv
).update(JSON.stringify(data));

// 4. 组装加密请求
const request = {
    key: encryptedKey.toString('base64'),
    data: encryptedData.toString('base64'),
    iv: iv.toString('base64'),
    tag: tag.toString('base64')
};
```

### 3. 签名验证
**签名生成规则**

```typescript
// 1. 组装签名字符串
const signString = `${timestamp}${deviceSN}${JSON.stringify(sortedData)}`;

// 2. 使用HMAC-SHA256生成签名
const sign = crypto.createHmac('sha256', DEVICE_SECRET)
    .update(signString)
    .digest('hex');
```

## 错误码完整说明
### 1. 系统级错误 (1000-1999)
| 错误码 | 说明 | 处理建议 |
| --- | --- | --- |
| 1000 | 请求参数错误 | 检查请求参数格式和必填项 |
| 1001 | 认证失败 | 检查token是否有效或重新认证 |
| 1002 | 签名无效 | 检查签名算法和密钥 |
| 1003 | 数据上报过期 | 检查数据时间戳是否在有效期内 |
| 1004 | 服务器内部错误 | 联系服务器管理员 |
| 1005 | 请求频率超限 | 降低请求频率 |
| 1006 | 服务暂时不可用 | 稍后重试 |


### 2. 设备相关错误 (2000-2999)
| 错误码 | 说明 | 处理建议 |
| --- | --- | --- |
| 2000 | 设备未注册 | 先进行设备注册 |
| 2001 | 设备长时间未上报 | 检查网络连接 |
| 2002 | 设备类型不支持 | 检查设备型号是否正确 |
| 2003 | 设备SN码无效 | 检查SN码格式 |
| 2004 | 设备数据缓存已满 | 清理本地缓存数据 |
| 2005 | 设备固件版本过低 | 更新设备固件 |


### 3. 数据相关错误 (3000-3999)
| 错误码 | 说明 | 处理建议 |
| --- | --- | --- |
| 3000 | 数据格式错误 | 检查数据格式是否符合规范 |
| 3001 | 数据解密失败 | 检查加密参数和密钥 |
| 3002 | 数据校验失败 | 检查数据完整性 |
| 3003 | 数据保存失败 | 重试或联系管理员 |
| 3004 | 数据过期 | 检查数据时间是否有效 |
| 3005 | 批量数据格式错误 | 检查批量数据格式 |


### 4. 打印任务相关错误 (4000-4999)
| 错误码 | 说明 | 处理建议 |
| --- | --- | --- |
| 4000 | 任务ID不存在 | 检查任务ID是否正确 |
| 4001 | 任务状态无效 | 检查状态值是否在允许范围 |
| 4002 | 文件名无效 | 检查文件名格式 |
| 4003 | 打印参数错误 | 检查打印参数设置 |
| 4004 | 任务已结束 | 不能对已结束任务操作 |
| 4005 | 重复的任务ID | 使用新的任务ID |


## 附录
### 1. 数据格式规范
**时间格式**

+ 使用ISO8601格式
+ 示例：2024-03-14T10:30:00Z

**数值类型**

+ 温度：保留1位小数
+ 百分比：保留2位小数
+ 内存/存储：整数，单位MB

### 2. 最佳实践
**网络处理**

+ 使用指数退避算法进行重试
+ 设置合理的超时时间
+ 维护请求队列

**错误处理**

+ 实现错误重试机制
+ 记录详细错误日志
+ 关键错误告警通知

