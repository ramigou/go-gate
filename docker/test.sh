#!/bin/bash

# Docker test script for Go-Gate reverse proxy

set -e

echo "🧪 Running Go-Gate Docker Tests"

# Check if services are running
if ! docker-compose ps | grep -q "Up"; then
    echo "❌ Services are not running. Please run './docker/start.sh' first."
    exit 1
fi

echo "⏳ Waiting for services to be ready..."
sleep 5

echo "🔄 Running comprehensive tests..."

# Run the test container
docker-compose run --rm test-runner

echo ""
echo "📊 Service Status:"
docker-compose ps

echo ""
echo "📋 Test completed!"
echo ""
echo "💡 Tips:"
echo "   • View real-time logs: docker-compose logs -f go-gate"
echo "   • Test manually: curl http://localhost:8080/api/test"
echo "   • Stop services: docker-compose down"