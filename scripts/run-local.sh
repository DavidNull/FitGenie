#!/bin/bash

echo "🚀 Starting FitGenie locally..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Start PostgreSQL in Docker
echo "📦 Starting PostgreSQL database..."
docker run -d \
    --name fitgenie-postgres \
    -e POSTGRES_DB=fitgenie \
    -e POSTGRES_USER=fitgenie \
    -e POSTGRES_PASSWORD=fitgenie \
    -p 5432:5432 \
    postgres:15 2>/dev/null || echo "📦 PostgreSQL container already running"

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for database to be ready..."
sleep 5

# Check if database is ready
until docker exec fitgenie-postgres pg_isready -U fitgenie -d fitgenie > /dev/null 2>&1; do
    echo "⏳ Still waiting for database..."
    sleep 2
done

echo "✅ Database is ready!"

# Run the Go application from the parent directory
echo "🏃 Starting FitGenie API server..."
(cd ../ && go run main.go)