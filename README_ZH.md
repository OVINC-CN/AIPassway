# AIPassway

中文 | [English](README.md)

一个用 Go 编写的轻量级、灵活的动态反向代理服务。

## ✨ 特性

- 🚀 **动态路由**: 基于 URL 路径前缀（服务密钥）进行请求路由
- ⚙️ **环境变量配置**: 通过环境变量进行简单配置
- 🌐 **转发代理支持**: 可选的出站请求转发代理
- 📝 **结构化日志**: JSON 格式的日志，包含完整的请求/响应详情
- 🔐 **公网认证**: 内置的公网访问令牌认证机制
- 📊 **OpenTelemetry 追踪**: 支持 OTLP 分布式追踪
- 🛡️ **网络安全**: 自动检测内网/公网访问

## 📖 工作原理

AIPassway 从传入请求路径中提取服务密钥，并将其转发到相应的后端服务：

```
客户端请求: GET /api/v1/users
                 ↓
服务密钥: "api" → 环境变量: APP_REAL_HOST_API=https://api.example.com
                 ↓
后端请求: GET https://api.example.com/v1/users
                 ↓
响应转发回客户端
```

### 请求流程

1. 从请求路径的第一段提取服务密钥（例如 `/api/...` → `api`）
2. 从环境变量 `APP_REAL_HOST_<KEY>` 中查找对应的后端主机
3. 将请求转发到后端服务
4. 将响应返回给客户端
5. 记录所有请求/响应详情及耗时信息

## 🚀 快速开始

### 使用 Docker

```bash
docker run -d \
  -p 8000:8000 \
  -e APP_REAL_HOST_API=https://api.example.com \
  -e APP_PUBLIC_AUTH_TOKEN=your_secret_token \
  ghcr.io/ovinc-cn/aipassway:latest
```

### 使用 Go

```bash
# 克隆仓库
git clone https://github.com/OVINC-CN/AIPassway.git
cd AIPassway

# 构建
go build -o ai-passway .

# 配置环境变量
export APP_REAL_HOST_API=https://api.example.com
export APP_PUBLIC_AUTH_TOKEN=your_secret_token

# 运行
./ai-passway
```

### 使用 Docker Compose

```yaml
version: '3.8'
services:
  aipassway:
    image: ghcr.io/ovinc-cn/aipassway:latest
    ports:
      - "8000:8000"
    environment:
      - APP_REAL_HOST_API=https://api.example.com
      - APP_REAL_HOST_WEB=https://web.example.com
      - APP_PUBLIC_AUTH_TOKEN=your_secret_token
      - APP_ENABLE_TRACE=true
      - APP_TRACE_ENDPOINT=jaeger:4317
```

## ⚙️ 配置

### 环境变量

| 变量 | 描述 | 必需 | 默认值 | 示例 |
|------|------|------|--------|------|
| `APP_REAL_HOST_<KEY>` | 服务密钥对应的后端主机 URL | 是 | - | `APP_REAL_HOST_API=https://api.example.com` |
| `APP_FORWARD_PROXY_URL` | 出站请求的转发代理 URL | 否 | - | `http://proxy.example.com:8080` |
| `APP_IDLE_TIMEOUT` | 空闲连接超时时间（秒） | 否 | 600 | `600` |
| `APP_HEADER_TIMEOUT` | 响应头超时时间（秒） | 否 | 60 | `60` |
| `APP_INTERNAL_NETWORKS` | 内网 CIDR 列表，逗号分隔 | 否 | `127.0.0.0/8, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16` | `192.168.1.0/24,10.0.0.0/8` |
| `APP_PUBLIC_AUTH_TOKEN` | 公网访问认证令牌（X-AI-Passway-Auth 请求头） | 否 | 每次重启随机生成 UUID v4 | `your_secret_token` |
| `APP_ENABLE_TRACE` | 启用 OpenTelemetry 追踪 | 否 | false | `true` |
| `APP_TRACE_ENDPOINT` | OTLP 追踪端点 | 否 | `127.0.0.1:4317` | `jaeger:4317` |
| `OTEL_SERVICE_NAME` | 追踪服务名称 | 否 | `ai-passway` | `my-service` |
| `OTEL_RESOURCE_ATTRIBUTES` | 追踪资源属性（逗号分隔的 key=value 对） | 否 | - | `key1=val1,key2=val2` |

### 服务密钥映射

服务密钥通过环境变量映射到后端主机：

```bash
# 将 "api" 密钥映射到 https://api.example.com
export APP_REAL_HOST_API=https://api.example.com

# 将 "web" 密钥映射到 https://web.example.com
export APP_REAL_HOST_WEB=https://web.example.com

# 现在对 /api/* 的请求将转发到 https://api.example.com/*
# 对 /web/* 的请求将转发到 https://web.example.com/*
```

## 🔐 认证

### 公网访问

从公网（非内网 IP）访问时，请求必须包含 `X-AI-Passway-Auth` 请求头：

```bash
curl -H "X-AI-Passway-Auth: your_secret_token" \
  http://your-proxy.com/api/v1/users
```

### 内网访问

来自内网（通过 `APP_INTERNAL_NETWORKS` 配置）的请求会自动绕过认证。

## 📊 可观测性

### 日志

AIPassway 使用结构化 JSON 日志，包含以下字段：

```json
{
  "level": "info",
  "msg": "request completed",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "duration": "123ms",
  "remote_addr": "192.168.1.100",
  "service_key": "api",
  "backend_host": "https://api.example.com"
}
```

### 追踪

启用 OpenTelemetry 分布式追踪：

```bash
export APP_ENABLE_TRACE=true
export APP_TRACE_ENDPOINT=jaeger:4317
export OTEL_SERVICE_NAME=ai-passway
```

## 🚨 错误处理

| 状态码 | 描述 | 原因 |
|--------|------|------|
| 401 Unauthorized | 认证缺失或无效 | 公网请求缺少或错误的 `X-AI-Passway-Auth` 请求头 |
| 501 Not Implemented | 服务未配置 | 服务密钥缺失或未配置对应的后端主机 |
| 500 Internal Server Error | 后端错误 | 后端主机 URL 解析失败或内部错误 |

## 🏗️ 架构

```
┌─────────┐         ┌──────────────┐         ┌─────────────┐
│  客户端  │────────▶│  AIPassway   │────────▶│   后端服务   │
│         │         │  (反向代理)   │         │             │
│         │◀────────│              │◀────────│             │
└─────────┘         └──────────────┘         └─────────────┘
                           │
                           │ 可选
                           ▼
                    ┌──────────────┐
                    │  转发代理    │
                    └──────────────┘
```

## 🛠️ 开发

### 前置要求

- Go 1.24 或更高版本
- Docker（可选）

### 构建

```bash
go build -o ai-passway .
```

### 运行测试

```bash
go test ./...
```

### 构建 Docker 镜像

```bash
docker build -t aipassway:latest .
```

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。

## 📧 联系方式

- GitHub: [OVINC-CN/AIPassway](https://github.com/OVINC-CN/AIPassway)
- Issues: [GitHub Issues](https://github.com/OVINC-CN/AIPassway/issues)
