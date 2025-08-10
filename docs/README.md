# Go-Gate L7 Reverse Proxy

*Read this in other languages: [한국어](README.ko.md), [日本語](README.ja.md)*

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

You can run Go-Gate in two ways: **Docker** (recommended) or **Local Development**.

### 🐳 Option 1: Docker (Recommended)

**Prerequisites:**
- Docker Desktop installed and running

**Quick Start:**
```bash
# Start everything (proxy + mock servers)
./docker/start.sh

# Run tests
./docker/test.sh

# Stop services
./docker/stop.sh
```

**Pros:** ✅ No dependencies to install ✅ Consistent environment ✅ One-command setup  
**Cons:** ❌ Requires Docker ❌ Slower startup

See [Docker Setup Guide](DOCKER.md) for detailed instructions.

### 💻 Option 2: Local Development

**Prerequisites:**
- Go 1.24.6 or later
- Python 3.7+ (for mock servers)
- Basic understanding of reverse proxies and HTTP

**Quick Start:**

```bash
# Terminal 1: Start mock backend servers
cd test/
python3 mock-servers.py

# Terminal 2: Start proxy server
go build -o go-gate cmd/server/main.go
./go-gate -config configs/test-config.yaml

# Terminal 3: Run tests
cd test/
./test-requests.sh
```

**Pros:** ✅ Faster startup ✅ Direct debugging ✅ No Docker needed  
**Cons:** ❌ Need Go/Python installed ❌ Manual setup ❌ Environment differences

### 🔄 Switching Between Methods

Both approaches work independently:
- **Docker files don't affect local development**
- **Local setup doesn't interfere with Docker**
- **Same test scripts work for both methods**

## 📊 Docker vs Local Comparison

| Feature | 🐳 Docker | 💻 Local |
|---------|-----------|---------|
| **Setup Time** | ~2-3 minutes | ~1 minute |
| **Dependencies** | Only Docker Desktop | Go 1.24.6+, Python 3.7+ |
| **Environment** | Consistent across systems | Depends on local setup |
| **Debugging** | Container logs | Direct IDE debugging |
| **Resource Usage** | Higher (containers) | Lower (native processes) |
| **Isolation** | Complete isolation | Uses system resources |
| **CI/CD Ready** | ✅ Perfect | ❌ Needs setup |
| **Beginner Friendly** | ✅ One command | ❌ Multiple steps |
| **Development Speed** | Slower builds | Faster iteration |

### 💡 When to Use Each

**Use Docker when:**
- 🚀 Getting started quickly
- 🔄 Testing in CI/CD pipelines  
- 👥 Ensuring team consistency
- 🐛 Isolating from system dependencies
- 📦 Simulating production environment

**Use Local when:**
- ⚡ Rapid development iteration
- 🐞 Debugging with IDE breakpoints
- 🔧 Customizing configurations frequently
- 💾 Conserving system resources
- 📚 Learning Go development patterns

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

## Docker Deployment

Go-Gate includes full Docker support for easy deployment and testing:

- **Complete Environment**: Proxy + mock backend servers
- **One-Command Setup**: `./docker/start.sh`
- **Automated Testing**: `./docker/test.sh`
- **Production Ready**: Optimized Dockerfile and docker-compose

See [Docker Documentation](DOCKER.md) for complete setup instructions.

### Next Steps for Enhancement

1. **Health Checking**: Implement actual health check monitoring
2. **Metrics**: Add Prometheus metrics for monitoring
3. **TLS Support**: Add HTTPS/TLS termination
4. **Rate Limiting**: Implement request rate limiting
5. **Circuit Breaker**: Add circuit breaker pattern for upstream failures
6. **WebSocket Support**: Add WebSocket proxying capabilities