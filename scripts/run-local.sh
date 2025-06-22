#!/bin/bash

echo " Starting FitGenie locally..."

if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo " Starting PostgreSQL database..."
docker run -d \
    --name fitgenie-postgres \
    -e POSTGRES_DB=fitgenie \
    -e POSTGRES_USER=fitgenie \
    -e POSTGRES_PASSWORD=fitgenie \
    -p 5432:5432 \
    postgres:15 2>/dev/null || echo "📦 PostgreSQL container already running"

echo " Waiting for database to be ready..."
sleep 5

until docker exec fitgenie-postgres pg_isready -U fitgenie -d fitgenie > /dev/null 2>&1; do
    echo " Still waiting for database..."
    sleep 2
done

echo "✅ Database is ready!"

echo " Starting FitGenie API server..."
(cd ../ && go run main.go)