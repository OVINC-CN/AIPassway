# AIPassway

[中文](README_ZH.md) | English

A lightweight and flexible dynamic reverse proxy service written in Go.

## ✨ Features

- 🚀 **Dynamic Routing**: Routes requests based on URL path prefixes (service keys)
- ⚙️ **Environment-based Configuration**: Simple configuration through environment variables
- 🌐 **Forward Proxy Support**: Optional forward proxy for outbound requests
- 📝 **Structured Logging**: JSON-formatted logging with comprehensive request/response details
- 🔐 **Public Network Authentication**: Built-in token authentication for public network access
- 📊 **OpenTelemetry Tracing**: Distributed tracing support with OTLP
- 🛡️ **Network Security**: Automatic internal/public network detection

## 📖 How It Works

AIPassway extracts service keys from incoming request paths and forwards them to corresponding backend hosts:

```
Client Request: GET /api/v1/users
                     ↓
Service Key: "api" → Environment Variable: APP_REAL_HOST_API=https://api.example.com
                     ↓
Backend Request: GET https://api.example.com/v1/users
                     ↓
Response forwarded back to client
```

### Request Flow

1. Extract service key from the first segment of the request path (e.g., `/api/...` → `api`)
2. Look up the corresponding backend host from environment variable `APP_REAL_HOST_<KEY>`
3. Forward the request to the backend service
4. Return the response to the client
5. Log all request/response details with timing information

## 🚀 Quick Start

### Using Docker

```bash
docker run -d \
  -p 8000:8000 \
  -e APP_REAL_HOST_API=https://api.example.com \
  -e APP_PUBLIC_AUTH_TOKEN=your_secret_token \
  ghcr.io/ovinc-cn/aipassway:latest
```

### Using Go

```bash
# Clone the repository
git clone https://github.com/OVINC-CN/AIPassway.git
cd AIPassway

# Build
go build -o ai-passway .

# Configure environment variables
export APP_REAL_HOST_API=https://api.example.com
export APP_PUBLIC_AUTH_TOKEN=your_secret_token

# Run
./ai-passway
```

### Using Docker Compose

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

## ⚙️ Configuration

### Environment Variables

| Variable                   | Description                                                        | Required | Default                                                  | Example                                     |
|----------------------------|--------------------------------------------------------------------|----------|----------------------------------------------------------|---------------------------------------------|
| `APP_REAL_HOST_<KEY>`      | Backend host URL for service key                                   | Yes      | -                                                        | `APP_REAL_HOST_API=https://api.example.com` |
| `APP_FORWARD_PROXY_URL`    | Forward proxy URL for outbound requests                            | No       | -                                                        | `http://proxy.example.com:8080`             |
| `APP_IDLE_TIMEOUT`         | Idle connection timeout (seconds)                                  | No       | 600                                                      | `600`                                       |
| `APP_HEADER_TIMEOUT`       | Response header timeout (seconds)                                  | No       | 60                                                       | `60`                                        |
| `APP_INTERNAL_NETWORKS`    | Comma-separated list of internal network CIDRs                     | No       | `127.0.0.0/8, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16` | `192.168.1.0/24,10.0.0.0/8`                 |
| `APP_PUBLIC_AUTH_TOKEN`    | Token for public network authentication (X-AI-Passway-Auth header) | No       | Random UUID v4 on each restart                           | `your_secret_token`                         |
| `APP_ENABLE_TRACE`         | Enable OpenTelemetry tracing                                       | No       | false                                                    | `true`                                      |
| `APP_TRACE_ENDPOINT`       | OTLP trace endpoint                                                | No       | `127.0.0.1:4317`                                         | `jaeger:4317`                               |
| `OTEL_SERVICE_NAME`        | Service name for tracing                                           | No       | `ai-passway`                                             | `my-service`                                |
| `OTEL_RESOURCE_ATTRIBUTES` | Resource attributes for tracing (comma-separated key=value pairs)  | No       | -                                                        | `key1=val1,key2=val2`                       |

### Service Key Mapping

Service keys are mapped to backend hosts using environment variables:

```bash
# Map "api" key to https://api.example.com
export APP_REAL_HOST_API=https://api.example.com

# Map "web" key to https://web.example.com
export APP_REAL_HOST_WEB=https://web.example.com

# Now requests to /api/* will be forwarded to https://api.example.com/*
# and requests to /web/* will be forwarded to https://web.example.com/*
```

## 🔐 Authentication

### Public Network Access

When accessing from public networks (non-internal IPs), requests must include the `X-AI-Passway-Auth` header:

```bash
curl -H "X-AI-Passway-Auth: your_secret_token" \
  http://your-proxy.com/api/v1/users
```

### Internal Network Access

Requests from internal networks (configured via `APP_INTERNAL_NETWORKS`) bypass authentication automatically.

## 📊 Observability

### Logging

AIPassway uses structured JSON logging with the following fields:

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

### Tracing

Enable distributed tracing with OpenTelemetry:

```bash
export APP_ENABLE_TRACE=true
export APP_TRACE_ENDPOINT=jaeger:4317
export OTEL_SERVICE_NAME=ai-passway
```

## 🚨 Error Handling

| Status Code               | Description                       | Cause                                                                          |
|---------------------------|-----------------------------------|--------------------------------------------------------------------------------|
| 401 Unauthorized          | Missing or invalid authentication | `X-AI-Passway-Auth` header is missing or incorrect for public network requests |
| 501 Not Implemented       | Service not configured            | Service key is missing or no backend host configured for the service key       |
| 500 Internal Server Error | Backend error                     | Backend host URL parsing fails or internal errors                              |

## 🏗️ Architecture

```
┌─────────┐         ┌──────────────┐         ┌─────────────┐
│ Client  │────────▶│  AIPassway   │────────▶│  Backend    │
│         │         │  (Reverse    │         │  Services   │
│         │◀────────│   Proxy)     │◀────────│             │
└─────────┘         └──────────────┘         └─────────────┘
                           │
                           │ Optional
                           ▼
                    ┌──────────────┐
                    │Forward Proxy │
                    └──────────────┘
```

## 🛠️ Development

### Prerequisites

- Go 1.24 or later
- Docker (optional)

### Build

```bash
go build -o ai-passway .
```

### Run Tests

```bash
go test ./...
```

### Build Docker Image

```bash
docker build -t aipassway:latest .
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Contact

- GitHub: [OVINC-CN/AIPassway](https://github.com/OVINC-CN/AIPassway)
- Issues: [GitHub Issues](https://github.com/OVINC-CN/AIPassway/issues)