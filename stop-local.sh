#!/bin/bash

echo "🛑 Stopping FitGenie local environment..."

# Stop and remove PostgreSQL container
echo "📦 Stopping PostgreSQL database..."
docker stop fitgenie-postgres 2>/dev/null || echo "📦 PostgreSQL container was not running"
docker rm fitgenie-postgres 2>/dev/null || echo "📦 PostgreSQL container was already removed"

echo "✅ Local environment stopped!" 