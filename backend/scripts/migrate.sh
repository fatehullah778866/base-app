#!/bin/bash

# Migration script for SQLite-based Base App Service
# Usage: ./scripts/migrate.sh [up|down]

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

DB_SQLITE_PATH=${DB_SQLITE_PATH:-file:app.db?_pragma=foreign_keys(ON)}

# Normalize to a filesystem path for the sqlite3 CLI
SQLITE_FILE=${DB_SQLITE_PATH#file:}
SQLITE_FILE=${SQLITE_FILE%%\?*}
if [ -z "$SQLITE_FILE" ]; then
    SQLITE_FILE="app.db"
fi

ACTION=${1:-up}

if ! command -v sqlite3 >/dev/null 2>&1; then
    echo "Error: sqlite3 command not found. Please install SQLite CLI to run migrations."
    exit 1
fi

case $ACTION in
    up)
        echo "Running migrations UP against $SQLITE_FILE ..."
        for file in $(ls migrations/*.sql | sort); do
            echo "Applying $file"
            sqlite3 "$SQLITE_FILE" < "$file"
        done
        ;;
    down)
        echo "Running migrations DOWN against $SQLITE_FILE ..."
        for file in $(ls migrations/*down.sql | sort); do
            echo "Applying $file"
            sqlite3 "$SQLITE_FILE" < "$file"
        done
        ;;
    *)
        echo "Usage: $0 [up|down]"
        exit 1
        ;;
esac

echo "Migration completed successfully!"

