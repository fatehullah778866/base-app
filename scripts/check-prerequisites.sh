#!/bin/bash

# Prerequisites Check Script for Base App Service
# Usage: ./scripts/check-prerequisites.sh

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Checking prerequisites for Base App Service...${NC}"
echo "=================================="
echo ""

ALL_OK=true

# Check Go
echo -n "Checking Go installation... "
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} Found: $GO_VERSION"
else
    echo -e "${RED}✗${NC} Go not found"
    echo "  Install from: https://go.dev/dl/"
    ALL_OK=false
fi

# Check PostgreSQL
echo -n "Checking PostgreSQL... "
if command -v psql &> /dev/null; then
    PSQL_VERSION=$(psql --version | awk '{print $3}')
    echo -e "${GREEN}✓${NC} Found: $PSQL_VERSION"
    
    # Check if PostgreSQL is running
    if command -v pg_isready &> /dev/null; then
        if pg_isready -h localhost -p 5432 &> /dev/null; then
            echo -e "  ${GREEN}✓${NC} PostgreSQL is running on localhost:5432"
        else
            echo -e "  ${YELLOW}⚠${NC} PostgreSQL is installed but not running on localhost:5432"
            echo "    Start with: brew services start postgresql@14 (macOS)"
            echo "    Or check: sudo systemctl status postgresql (Linux)"
        fi
    fi
else
    echo -e "${YELLOW}⚠${NC} PostgreSQL client not found"
    echo "  Install from: https://www.postgresql.org/download/"
    echo "  Or use Docker: docker compose up -d"
fi

# Check Docker (optional but recommended)
echo -n "Checking Docker... "
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | tr -d ',')
    echo -e "${GREEN}✓${NC} Found: $DOCKER_VERSION"
    
    # Check if Docker is running
    if docker info &> /dev/null; then
        echo -e "  ${GREEN}✓${NC} Docker daemon is running"
        
        # Check docker-compose
        if command -v docker-compose &> /dev/null || docker compose version &> /dev/null; then
            echo -e "  ${GREEN}✓${NC} Docker Compose is available"
        else
            echo -e "  ${YELLOW}⚠${NC} Docker Compose not found"
        fi
    else
        echo -e "  ${YELLOW}⚠${NC} Docker daemon is not running"
        echo "    Start Docker Desktop or Docker daemon"
    fi
else
    echo -e "${YELLOW}⚠${NC} Docker not found (optional, but recommended)"
    echo "  Install from: https://www.docker.com/products/docker-desktop"
fi

# Check migrate tool
echo -n "Checking golang-migrate... "
if command -v migrate &> /dev/null; then
    MIGRATE_VERSION=$(migrate -version 2>&1 | head -1 | awk '{print $2}')
    echo -e "${GREEN}✓${NC} Found: $MIGRATE_VERSION"
else
    echo -e "${YELLOW}⚠${NC} golang-migrate not found"
    echo "  Install with: brew install golang-migrate (macOS)"
    echo "  Or: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
fi

# Check jq (optional, for better test output)
echo -n "Checking jq... "
if command -v jq &> /dev/null; then
    JQ_VERSION=$(jq --version)
    echo -e "${GREEN}✓${NC} Found: $JQ_VERSION"
else
    echo -e "${YELLOW}⚠${NC} jq not found (optional, for better JSON output)"
    echo "  Install with: brew install jq (macOS)"
fi

# Check if server is running
echo ""
echo -n "Checking if server is running... "
if curl -s http://localhost:8080/health &> /dev/null; then
    echo -e "${GREEN}✓${NC} Server is running on http://localhost:8080"
    SERVER_RUNNING=true
elif curl -s http://localhost:8081/health &> /dev/null; then
    echo -e "${GREEN}✓${NC} Server is running on http://localhost:8081"
    SERVER_RUNNING=true
else
    echo -e "${YELLOW}⚠${NC} Server is not running"
    SERVER_RUNNING=false
fi

# Check port availability
echo ""
echo -n "Checking port 8080... "
if lsof -ti:8080 &> /dev/null; then
    PORT_PID=$(lsof -ti:8080 | head -1)
    PORT_PROCESS=$(ps -p $PORT_PID -o comm= 2>/dev/null || echo "unknown")
    echo -e "${YELLOW}⚠${NC} Port 8080 is in use by PID $PORT_PID ($PORT_PROCESS)"
    echo "    You may need to use a different port: PORT=8081 make run"
else
    echo -e "${GREEN}✓${NC} Port 8080 is available"
fi

# Summary
echo ""
echo "=================================="
if [ "$ALL_OK" = true ] && [ "$SERVER_RUNNING" = true ]; then
    echo -e "${GREEN}All checks passed! Ready to test.${NC}"
    exit 0
elif [ "$SERVER_RUNNING" = true ]; then
    echo -e "${GREEN}Server is running. You can run tests.${NC}"
    exit 0
else
    echo -e "${YELLOW}Some prerequisites are missing or server is not running.${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Start database:"
    echo "   - With Docker: make docker-compose-up"
    echo "   - Or local PostgreSQL: brew services start postgresql@14"
    echo ""
    echo "2. Run migrations:"
    echo "   make migrate-up"
    echo ""
    echo "3. Start server:"
    echo "   make dev  # (sets up + runs)"
    echo "   # OR"
    echo "   make run  # (if already set up)"
    echo ""
    echo "4. Run tests:"
    echo "   ./scripts/test-api.sh"
    exit 1
fi

