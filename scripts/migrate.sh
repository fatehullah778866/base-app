#!/bin/bash

# Migration script for Base App Service
# Usage: ./scripts/migrate.sh [up|down|create] [migration_name]

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-baseapp}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-base_app_db}
DB_SSL_MODE=${DB_SSL_MODE:-disable}

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# Check if migrate tool is installed
if ! command -v migrate &> /dev/null; then
    echo "Error: migrate tool not found. Install it with:"
    echo "  brew install golang-migrate"
    echo "  or"
    echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

ACTION=${1:-up}
MIGRATION_NAME=${2:-}

case $ACTION in
    up)
        echo "Running migrations UP..."
        migrate -path migrations -database "$DATABASE_URL" up
        ;;
    down)
        echo "Running migrations DOWN..."
        migrate -path migrations -database "$DATABASE_URL" down
        ;;
    create)
        if [ -z "$MIGRATION_NAME" ]; then
            echo "Error: Migration name required for 'create' action"
            echo "Usage: ./scripts/migrate.sh create migration_name"
            exit 1
        fi
        echo "Creating migration: $MIGRATION_NAME"
        migrate create -ext sql -dir migrations -seq "$MIGRATION_NAME"
        ;;
    *)
        echo "Usage: $0 [up|down|create] [migration_name]"
        exit 1
        ;;
esac

echo "Migration completed successfully!"

