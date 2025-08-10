#!/bin/bash

# Docker setup script for Go-Gate reverse proxy

set -e

echo "🚀 Starting Go-Gate Docker Environment"

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose >/dev/null 2>&1; then
    echo "❌ docker-compose is not installed. Please install docker-compose first."
    exit 1
fi

# Build and start services
echo "📦 Building Go-Gate proxy server..."
docker-compose build

echo "🏃 Starting all services..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 15

# Check service health
echo "🔍 Checking service health..."

# Check proxy health
if curl -f -s http://localhost:8080/ >/dev/null; then
    echo "✅ Go-Gate proxy is healthy"
else
    echo "❌ Go-Gate proxy health check failed"
    docker-compose logs go-gate
fi

# Check backend services
for port in 3001 3002 4000; do
    if curl -f -s http://localhost:$port/ >/dev/null; then
        echo "✅ Backend service on port $port is healthy"
    else
        echo "❌ Backend service on port $port health check failed"
    fi
done

echo ""
echo "🎉 Go-Gate environment is ready!"
echo ""
echo "📋 Available services:"
echo "   • Go-Gate Proxy: http://localhost:8080"
echo "   • API Server 1:  http://localhost:3001"
echo "   • API Server 2:  http://localhost:3002"  
echo "   • Web Server:    http://localhost:4000"
echo ""
echo "🧪 To run tests:"
echo "   docker-compose run --rm test-runner"
echo ""
echo "📊 To view logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 To stop all services:"
echo "   docker-compose down"