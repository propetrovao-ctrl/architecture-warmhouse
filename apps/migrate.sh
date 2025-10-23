#!/bin/bash

# Migration management script for Smart Home API

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
DATABASE_URL="postgres://postgres:postgres@localhost:5432/smarthome"
COMMAND="up"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--database-url)
            DATABASE_URL="$2"
            shift 2
            ;;
        -c|--command)
            COMMAND="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -d, --database-url URL    Database connection URL"
            echo "  -c, --command COMMAND     Migration command (up, status)"
            echo "  -h, --help               Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Run migrations with default settings"
            echo "  $0 -c status                          # Show migration status"
            echo "  $0 -d postgres://user:pass@host/db    # Use custom database URL"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

echo -e "${YELLOW}Smart Home API Migration Tool${NC}"
echo "Database URL: $DATABASE_URL"
echo "Command: $COMMAND"
echo ""

# Check if we're running in Docker
if [ -f /.dockerenv ]; then
    echo -e "${GREEN}Running in Docker container${NC}"
    ./migrate -database-url="$DATABASE_URL" -command="$COMMAND"
else
    echo -e "${GREEN}Running locally${NC}"
    cd smart_home
    go run cmd/migrate/main.go -database-url="$DATABASE_URL" -command="$COMMAND"
fi

echo -e "${GREEN}Migration completed successfully!${NC}"
