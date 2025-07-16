#!/bin/bash

if ! docker info > /dev/null 2>&1; then
    exit 1
fi

docker run -d \
    --name fitgenie-postgres \
    -e POSTGRES_DB=fitgenie \
    -e POSTGRES_USER=fitgenie \
    -e POSTGRES_PASSWORD=fitgenie \
    -p 5432:5432 \
    postgres:15 2>/dev/null || true

sleep 5

until docker exec fitgenie-postgres pg_isready -U fitgenie -d fitgenie > /dev/null 2>&1; do
    sleep 2
done

(cd ../ && go run main.go)