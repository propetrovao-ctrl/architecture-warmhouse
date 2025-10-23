#!/bin/bash

# Database initialization script
set -e

echo "Initializing Smart Home database..."

# Wait for database to be ready
echo "Waiting for database to be ready..."
until pg_isready -h postgres -p 5432 -U postgres; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Run migrations
echo "Running database migrations..."
cd smart_home
go run cmd/migrate/main.go -database-url="postgres://postgres:postgres@postgres:5432/smarthome" -command=up

echo "Database initialization completed successfully!"
