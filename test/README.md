# Testing Go-Gate Reverse Proxy

*Read this in other languages: [한국어](README.ko.md), [日本語](README.ja.md)*

This directory contains testing tools and scripts for the Go-Gate L7 reverse proxy server.

## Quick Start

### 1. Start Mock Backend Servers

```bash
# Terminal 1: Start mock backend servers
cd test/
python3 mock-servers.py
```

This will start three mock HTTP servers:
- `localhost:3001` - API Server 1 (weight: 2)
- `localhost:3002` - API Server 2 (weight: 1)  
- `localhost:4000` - Web Server (weight: 1)

### 2. Start the Proxy Server

```bash
# Terminal 2: Start the reverse proxy
./go-gate -config configs/test-config.yaml
```

The proxy will start on `localhost:8080` and route requests to the mock servers.

### 3. Run Tests

```bash
# Terminal 3: Run automated tests
cd test/
./test-requests.sh
```

## Test Files

### `mock-servers.py`
Python script that creates three HTTP servers for testing:
- Returns JSON responses with server info, request details, and timestamps
- Handles both GET and POST requests
- Includes `/health` endpoint for health checks
- Logs all incoming requests

### `test-requests.sh`
Shell script that tests various proxy scenarios:
- Path-based routing (`/api/*`)
- Host-based routing (`admin.example.com`, `*.example.com`)
- Load balancing verification
- Default routing fallback
- POST request handling
- Health check endpoints

### `test-config.yaml`
Simplified configuration for testing (without health check settings).

## Manual Testing

### Basic Requests

```bash
# Test default routing (load balanced across all servers)
curl http://localhost:8080/

# Test API routing (load balanced between API servers)
curl http://localhost:8080/api/users

# Test POST request
curl -X POST -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/api/submit
```

### Host-based Routing

First, add test domains to your `/etc/hosts` file:
```bash
# Add these lines to /etc/hosts
127.0.0.1 admin.example.com
127.0.0.1 www.example.com
```

Then test:
```bash
# Should route to API Server 1 only
curl -H "Host: admin.example.com" http://localhost:8080/dashboard

# Should route to Web Server only
curl -H "Host: www.example.com" http://localhost:8080/home
```

## Expected Behavior

### Path-based Routing
- Requests to `/api/*` are load balanced between API Server 1 (weight 2) and API Server 2 (weight 1)
- API Server 1 should receive ~67% of requests, API Server 2 ~33%

### Host-based Routing
- `admin.example.com` requests go to API Server 1 only
- `*.example.com` requests (like `www.example.com`) go to Web Server only

### Default Routing
- All other requests are load balanced across all three servers
- Distribution: API Server 1 (50%), API Server 2 (25%), Web Server (25%)

## Troubleshooting

### Proxy won't start
- Check if port 8080 is available: `lsof -i :8080`
- Verify configuration syntax: `./go-gate -config configs/test-config.yaml`

### Mock servers won't start
- Check if ports 3001, 3002, 4000 are available
- Ensure Python 3 is installed: `python3 --version`

### Requests fail
- Verify all servers are running with `ps aux | grep -E "(go-gate|python)"`
- Check proxy logs for error messages
- Test mock servers directly: `curl http://localhost:3001/`

### Host-based routing not working
- Ensure `/etc/hosts` entries are added correctly
- Use `curl -H "Host: domain.com"` for testing
- Check DNS resolution: `nslookup admin.example.com`

## Advanced Testing

### Load Testing
```bash
# Test load distribution (requires 'ab' - Apache Bench)
ab -n 100 -c 10 http://localhost:8080/api/test

# Or with curl in a loop
for i in {1..20}; do
  curl -s http://localhost:8080/api/test | grep '"server"'
done | sort | uniq -c
```

### Error Scenarios
```bash
# Stop one backend server and test failover
# Kill API Server 2 (Ctrl+C in mock-servers.py terminal, then restart without port 3002)
curl http://localhost:8080/api/test  # Should still work with API Server 1

# Test with all backends down
# Stop mock-servers.py completely
curl http://localhost:8080/  # Should return 502 Bad Gateway
```

## Monitoring

### Real-time Logs
Watch proxy logs to see request routing:
```bash
./go-gate -config configs/test-config.yaml | grep -E "(api_server|web_server)"
```

### Request Distribution Analysis
```bash
# Analyze which servers are handling requests
./test-requests.sh 2>&1 | grep "Server=" | sort | uniq -c
```