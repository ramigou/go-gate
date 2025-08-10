# Docker Setup for Go-Gate

*Read this in other languages: [한국어](DOCKER.ko.md), [日本語](DOCKER.ja.md)*

This guide explains how to run Go-Gate reverse proxy server and testing environment using Docker.

## Quick Start

### Prerequisites

- Docker Engine 20.10 or later
- Docker Compose 2.0 or later

### One-Command Setup

```bash
# Start all services
./docker/start.sh
```

This will:
- Build the Go-Gate proxy server
- Start 3 mock backend servers (API Server 1, API Server 2, Web Server)
- Start the reverse proxy on port 8080
- Perform health checks
- Display service status

## Docker Components

### Services

| Service | Description | Port | URL |
|---------|-------------|------|-----|
| `go-gate` | Reverse proxy server | 8080 | http://localhost:8080 |
| `api-server-1` | Mock API backend (weight: 2) | 3001 | http://localhost:3001 |
| `api-server-2` | Mock API backend (weight: 1) | 3002 | http://localhost:3002 |
| `web-server` | Mock web backend | 4000 | http://localhost:4000 |
| `test-runner` | Automated test container | - | - |

### Network

All services run on a custom bridge network `go-gate-network` for isolated communication.

## Usage

### Start Environment

```bash
# Start all services in background
./docker/start.sh

# Or manually with docker-compose
docker-compose up -d
```

### Run Tests

```bash
# Run automated tests
./docker/test.sh

# Or manually
docker-compose run --rm test-runner
```

### Stop Environment

```bash
# Stop all services
./docker/stop.sh

# Or manually
docker-compose down
```

## Testing

### Automated Tests

The test runner performs comprehensive testing:

```bash
docker-compose run --rm test-runner
```

Tests include:
- API route load balancing
- Default route distribution
- Host-based routing
- Service health checks

### Manual Testing

```bash
# Test API routing (load balanced)
curl http://localhost:8080/api/users

# Test default routing
curl http://localhost:8080/default

# Test host-based routing
curl -H "Host: admin.example.com" http://localhost:8080/dashboard
curl -H "Host: www.example.com" http://localhost:8080/home

# Test POST requests
curl -X POST -H "Content-Type: application/json" \
     -d '{"test": "data"}' \
     http://localhost:8080/api/submit
```

## Monitoring

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f go-gate
docker-compose logs -f api-server-1

# Live proxy logs only
docker-compose logs -f go-gate | grep -E "(api_server|web_server)"
```

### Service Status

```bash
# Check running services
docker-compose ps

# Check service health
docker-compose exec go-gate wget --spider -q http://localhost:8080/
```

## Configuration

### Docker Configuration

The Docker setup uses `configs/docker-config.yaml` which is optimized for container networking:

- Uses container hostnames (`api-server-1`, `api-server-2`, `web-server`)
- Configured for Docker network communication
- Same routing rules as local development

### Environment Variables

You can override configuration using environment variables:

```bash
# Custom port
PROXY_PORT=9090 docker-compose up -d

# Custom config file
CONFIG_FILE=configs/production-config.yaml docker-compose up -d
```

## Development

### Building Custom Image

```bash
# Build Go-Gate image
docker-compose build go-gate

# Or manually
docker build -t go-gate:latest .
```

### Development Mode

For development with live reloading:

```bash
# Mount source code
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## Troubleshooting

### Common Issues

**Services won't start:**
```bash
# Check Docker status
docker info

# Check port conflicts
lsof -i :8080 -i :3001 -i :3002 -i :4000

# View service logs
docker-compose logs
```

**Proxy returns 502 errors:**
```bash
# Check backend services
docker-compose ps
docker-compose logs api-server-1 api-server-2 web-server

# Test backend services directly
curl http://localhost:3001/
curl http://localhost:3002/
curl http://localhost:4000/
```

**Tests fail:**
```bash
# Wait for services to be ready
sleep 10

# Check proxy health
curl -f http://localhost:8080/

# Run tests with verbose output
docker-compose run --rm test-runner
```

### Cleanup

```bash
# Remove all containers and networks
docker-compose down

# Remove containers, networks, and images
docker-compose down --rmi local

# Remove everything including volumes
docker-compose down -v --rmi all
```

## Production Deployment

For production deployment:

1. Use production configuration
2. Set up proper logging
3. Configure health monitoring
4. Use container orchestration (Kubernetes, Docker Swarm)
5. Set up SSL/TLS termination
6. Configure persistent volumes if needed

```bash
# Production example
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```