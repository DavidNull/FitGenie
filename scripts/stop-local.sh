#!/bin/bash

# Stop and remove PostgreSQL container
docker stop fitgenie-postgres 2>/dev/null || true
docker rm fitgenie-postgres 2>/dev/null || true 