# 明达3D打印云服务平台 - 数据采集服务

## 项目概述
本项目是明达3D打印云服务平台的数据采集服务模块，负责采集和管理3D打印设备的各类数据，包括设备状态、打印任务、告警信息等。

### 主要功能
- 设备基础信息采集
- 设备状态监控
- 打印任务跟踪
- 告警信息管理
- 数据安全传输
- 离线数据缓存

## 技术栈
- 语言：Go 1.21+
- Web框架：Gin
- ORM：GORM
- 数据库：MySQL 8.0+
- 缓存：Redis 6.0+
- 消息队列：RabbitMQ
- API文档：Swagger

## 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- RabbitMQ 3.8+

### 安装步骤

1. 克隆项目
```bash
git clone https://github.com/yourusername/mingda-cloud-service.git
cd mingda-cloud-service
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境变量
```bash
cp configs/config.example.yaml configs/config.yaml
# 编辑 config.yaml 配置文件
```

4. 初始化数据库
```bash
# 使用 sql 文件初始化数据库结构
mysql -u your_username -p your_database < scripts/init.sql
```

5. 运行服务
```bash
go run cmd/server/main.go
```

### 配置说明
配置文件位于 `configs/config.yaml`，主要包含以下配置项：

```yaml
server:
  port: 8080
  mode: debug  # debug/release

database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  dbname: md_device_db

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

rabbitmq:
  host: localhost
  port: 5672
  username: guest
  password: guest
```

## API文档
API文档使用Swagger生成，服务启动后可访问：
```
http://localhost:8080/swagger/index.html
```

## 项目结构
```
mingda_cloud_service/
├── api/                    # API接口定义
├── cmd/                    # 主程序入口
├── configs/                # 配置文件
├── internal/               # 内部代码
├── pkg/                    # 可对外公开的包
├── scripts/                # 构建、部署脚本
├── test/                   # 测试文件
└── docs/                   # 文档
```

## 开发规范
- 代码风格遵循Go官方规范
- 使用golangci-lint进行代码检查
- 提交前运行单元测试
- 遵循RESTful API设计规范
- 使用统一的错误处理机制

## 测试
运行单元测试：
```bash
go test ./...
```

运行带覆盖率的测试：
```bash
go test -cover ./...
```

## 部署
### Docker部署
```bash
# 构建镜像
docker build -t mingda-cloud-service .

# 运行容器
docker run -d -p 8080:8080 mingda-cloud-service
```

### 手动部署
1. 编译
```bash
go build -o mingda-cloud-service cmd/server/main.go
```

2. 运行
```bash
./mingda-cloud-service
```

## 监控
服务提供了以下监控指标：
- HTTP请求统计
- 系统资源使用情况
- 数据采集任务状态
- 告警信息统计

可通过Prometheus + Grafana进行监控。

## 贡献指南
1. Fork 项目
2. 创建特性分支
3. 提交变更
4. 推送到分支
5. 创建Pull Request

## 许可证
[MIT License](LICENSE)

## 联系方式
- 项目维护者：[维护者姓名]
- 邮箱：[联系邮箱]
