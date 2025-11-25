#!/bin/bash

# Local PostgreSQL Setup Helper Script
# This script helps set up PostgreSQL locally (macOS with Homebrew)

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Setting up local PostgreSQL database...${NC}"
echo "=================================="
echo ""

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo -e "${RED}✗${NC} Homebrew not found"
    echo "Install from: https://brew.sh"
    exit 1
fi

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo -e "${YELLOW}PostgreSQL not found. Installing...${NC}"
    brew install postgresql@14
    brew services start postgresql@14
    sleep 3
else
    echo -e "${GREEN}✓${NC} PostgreSQL is installed"
fi

# Start PostgreSQL service
echo -e "${YELLOW}Starting PostgreSQL service...${NC}"
brew services start postgresql@14 || brew services restart postgresql@14
sleep 2

# Check if PostgreSQL is running
if ! pg_isready -h localhost -p 5432 &> /dev/null; then
    echo -e "${RED}✗${NC} PostgreSQL is not running"
    echo "Try manually: brew services start postgresql@14"
    exit 1
fi

echo -e "${GREEN}✓${NC} PostgreSQL is running"

# Create database and user
echo ""
echo -e "${YELLOW}Creating database and user...${NC}"

# Default values
DB_NAME=${DB_NAME:-base_app_db}
DB_USER=${DB_USER:-baseapp}
DB_PASSWORD=${DB_PASSWORD:-password}

# Create user (ignore error if exists)
psql -h localhost -U $(whoami) -d postgres -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || echo "User already exists or using default user"

# Create database (ignore error if exists)
psql -h localhost -U $(whoami) -d postgres -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;" 2>/dev/null || echo "Database already exists"

echo -e "${GREEN}✓${NC} Database '$DB_NAME' and user '$DB_USER' are ready"
echo ""
echo "Database connection details:"
echo "  Host: localhost"
echo "  Port: 5432"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Password: $DB_PASSWORD"
echo ""
echo "Update your .env file with these values, then run:"
echo "  make migrate-up"
echo "  make run"

