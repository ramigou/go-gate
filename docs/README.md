# Go-Gate L7 Reverse Proxy

A lightweight Layer 7 reverse proxy server built in Go, supporting host-based and path-based routing with load balancing.

## Features

- **L7 Routing**: Route requests based on host, path, or both
- **Load Balancing**: Weighted round-robin load balancing across multiple upstreams
- **Health Checks**: Basic health check configuration (structure ready for implementation)
- **Request Logging**: HTTP access logs with response time and status
- **Flexible Configuration**: YAML-based configuration with validation

## Project Structure

```
go-gate/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── config/config.go        # Configuration management
│   ├── proxy/server.go         # Main proxy server logic
│   ├── router/router.go        # Request routing and upstream selection
│   └── middleware/logging.go   # HTTP middleware (logging)
├── configs/config.yaml         # Example configuration
└── docs/README.md             # This documentation
```

## Getting Started

### Prerequisites

- Go 1.24.6 or later
- Basic understanding of reverse proxies and HTTP

### Running the Server

1. **Build the application:**
   ```bash
   go build -o go-gate cmd/server/main.go
   ```

2. **Run with default config:**
   ```bash
   ./go-gate
   ```

3. **Run with custom config:**
   ```bash
   ./go-gate -config /path/to/your/config.yaml
   ```

### Configuration

The server uses YAML configuration files. See `configs/config.yaml` for an example.

#### Configuration Structure:

- **server**: Server settings (port)
- **upstreams**: Backend server definitions with weights and health check settings
- **routes**: Routing rules with host/path matching

#### Route Matching:

- **Host matching**: Supports exact match and wildcard (`*.example.com`)
- **Path matching**: Supports exact match and prefix matching (`/api/*`)
- **Load balancing**: Weighted distribution among multiple upstreams

### Example Usage

With the provided config, the proxy will:

1. Route `/api/*` requests to API servers (api_server_1, api_server_2)
2. Route `admin.example.com` requests to api_server_1
3. Route `*.example.com` requests to web_server
4. Fall back to all servers for other requests

### Development Commands

- **Build**: `go build cmd/server/main.go`
- **Run**: `go run cmd/server/main.go`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`
- **Tidy modules**: `go mod tidy`

### Next Steps for Enhancement

1. **Health Checking**: Implement actual health check monitoring
2. **Metrics**: Add Prometheus metrics for monitoring
3. **TLS Support**: Add HTTPS/TLS termination
4. **Rate Limiting**: Implement request rate limiting
5. **Circuit Breaker**: Add circuit breaker pattern for upstream failures
6. **WebSocket Support**: Add WebSocket proxying capabilities