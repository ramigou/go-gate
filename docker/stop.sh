#!/bin/bash

# Docker cleanup script for Go-Gate reverse proxy

set -e

echo "🛑 Stopping Go-Gate Docker Environment"

# Stop and remove containers
echo "📦 Stopping all services..."
docker-compose down

# Optional: Remove images (uncomment if needed)
# echo "🗑️  Removing Go-Gate images..."
# docker-compose down --rmi local

# Optional: Remove volumes (uncomment if needed)
# echo "💾 Removing volumes..."
# docker-compose down -v

echo "✅ Go-Gate environment stopped successfully"
echo ""
echo "💡 To start again:"
echo "   ./docker/start.sh"