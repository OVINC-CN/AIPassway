# AIPassway

A lightweight and flexible dynamic reverse proxy service written in Go.

## Overview

AIPassway is a dynamic reverse proxy that routes requests to different backend services based on URL path prefixes. It extracts service keys from incoming request paths and forwards them to corresponding backend hosts configured via environment variables.

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

## Configuration

Configure backend services using environment variables with the following pattern:

```bash
APP_REAL_HOST_<SERVICE_KEY>=<BACKEND_URL>
```

### Environment Variables

| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `APP_REAL_HOST_<KEY>` | Backend host URL for service key | Yes | `APP_REAL_HOST_API=https://api.example.com` |
| `APP_FORWARD_PROXY_URL` | Forward proxy URL for outbound requests | No | `http://proxy.example.com:8080` |

## Error Handling

- **501 Not Implemented**: Returned when service key is missing or backend host is not configured
- **500 Internal Server Error**: Returned when backend host URL parsing fails

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository.