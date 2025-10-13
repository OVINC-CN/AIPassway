# AIPassway

A lightweight and flexible dynamic reverse proxy service written in Go.

## Overview

AIPassway is a dynamic reverse proxy that routes requests to different backend services based on URL path prefixes. It
extracts service keys from incoming request paths and forwards them to corresponding backend hosts configured via
environment variables.

## Features

- **Dynamic Routing**: Routes requests based on the first path segment (service key)
- **Environment-based Configuration**: Backend hosts configured through environment variables
- **Forward Proxy Support**: Optional forward proxy support for outbound requests
- **Request Logging**: Comprehensive request/response logging with timing information
- **JSON Structured Logging**: Uses structured JSON logging for better observability
- **Health Monitoring**: Built-in HTTP server with middleware support

## How It Works

1. Extracts service key from the first segment of the request path
2. Looks up the corresponding backend host from environment variables
3. Forwards the request to the backend service
4. Returns the response to the client
5. Logs all request/response details

### Example Request Flow

```
Client Request: GET /api/v1/users
                     ↓
Service Key: "api" → Environment Variable: APP_REAL_HOST_API=https://api.example.com
                     ↓
Backend Request: GET https://api.example.com/v1/users
                     ↓
Response forwarded back to client
```

## Environment Variables

| Variable                   | Description                                                        | Required                                                                  | Example                                           |
|----------------------------|--------------------------------------------------------------------|---------------------------------------------------------------------------|---------------------------------------------------|
| `APP_REAL_HOST_<KEY>`      | Backend host URL for service key                                   | Yes                                                                       | `APP_REAL_HOST_API=https://api.example.com`       |
| `APP_FORWARD_PROXY_URL`    | Forward proxy URL for outbound requests                            | No                                                                        | `http://proxy.example.com:8080`                   |
| `APP_IDLE_TIMEOUT`         | Idle connection timeout in seconds                                 | No <br/>(default: 600)                                                    | `APP_IDLE_TIMEOUT=600`                            |
| `APP_HEADER_TIMEOUT`       | Response header timeout in seconds                                 | No <br/>(default: 60)                                                     | `APP_HEADER_TIMEOUT=60`                           |
| `APP_INTERNAL_NETWORKS`    | comma-separated list of internal network CIDRs                     | No <br/>(default: 127.0.0.0/8, 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) | `APP_INTERNAL_NETWORKS=192.168.1.0/24,10.0.0.0/8` |
| `APP_PUBLIC_AUTH_TOKEN`    | token for public authentication (used in X-AI-Passway-Auth header) | No <br/>(default: random uuid4 each reboot)                               | `APP_PUBLIC_AUTH_TOKEN=your_token_here`           |
| `APP_ENABLE_TRACE`         | enable opentelemetry tracing                                       | No <br/>(default: disabled)                                               | `APP_ENABLE_TRACE=true`                           |
| `APP_TRACE_ENDPOINT`       | otlp trace endpoint                                                | No <br/>(default: 127.0.0.1:4317)                                         | `APP_TRACE_ENDPOINT=jaeger:4317`                  |
| `OTEL_SERVICE_NAME`        | service name for tracing                                           | No <br/>(default: ai-passway)                                             | `APP_SERVICE_NAME=my-service`                     |
| `OTEL_RESOURCE_ATTRIBUTES` | resource attributes for tracing, comma-separated key-value pairs   | No <br/>(default: none)                                                   | `OTEL_RESOURCE_ATTRIBUTES=key1=val1,key2=val2`    |

## Error Handling

- **401 Unauthorized**: Returned when the `X-AI-Passway-Auth` header is missing or invalid and visiting through public
  network
- **501 Not Implemented**: Returned when service key is missing or backend host is not configured
- **500 Internal Server Error**: Returned when backend host URL parsing fails

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.