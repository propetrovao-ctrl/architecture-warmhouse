#!/bin/bash

# Script to run the application with database migrations
set -e

echo "Starting Smart Home API with migrations..."

# Wait for database to be ready
echo "Waiting for database to be ready..."
until pg_isready -h postgres -p 5432 -U postgres; do
  echo "Database is unavailable - sleeping"
  sleep 2
done

echo "Database is ready!"

# Run migrations
echo "Running database migrations..."
./migrate -database-url="$DATABASE_URL" -command=up

echo "Migrations completed successfully!"

# Start the application
echo "Starting Smart Home API..."
exec ./smarthome
