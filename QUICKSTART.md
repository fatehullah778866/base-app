# Quick Start Guide

Get the Base App Service running locally in 5 minutes!

## Prerequisites

- Go 1.21+ installed
- Docker and Docker Compose installed
- `make` installed (optional, but recommended)

## Step 1: Clone and Setup

```bash
cd base-app
cp .env.example .env
# Edit .env if needed (defaults work for local dev)
```

## Step 2: Start Services

```bash
# Start PostgreSQL and Redis
make docker-compose-up

# Or manually:
docker-compose up -d
```

Wait a few seconds for services to be ready.

## Step 3: Run Migrations

```bash
# Run database migrations
make migrate-up

# Or manually:
./scripts/migrate.sh up
```

**Note**: If you don't have `migrate` tool installed:
```bash
# macOS
brew install golang-migrate

# Or install via Go
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Step 4: Run the Server

```bash
# Run the server
make run

# Or manually:
go run ./cmd/server/main.go
```

You should see:
```
{"level":"info","msg":"Database connection established"}
{"level":"info","msg":"Server starting","port":"8080"}
```

## Step 5: Test the API

Open a new terminal and run:

```bash
# Test health endpoint
curl http://localhost:8080/health

# Run full API test suite
./scripts/test-api.sh
```

For detailed testing instructions and methods, see [TESTING.md](TESTING.md).

## Step 6: Create a User

```bash
curl -X POST http://localhost:8080/v1/auth/signup \
  -H "Content-Type: application/json" \
  -H "X-Product-Name: test-product" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "name": "Test User",
    "terms_accepted": true,
    "terms_version": "1.0"
  }'
```

Save the `token` from the response for authenticated requests.

## Step 7: Test Authenticated Endpoints

```bash
# Replace YOUR_TOKEN with the token from signup
TOKEN="YOUR_TOKEN"

# Get current user
curl http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer $TOKEN"

# Update theme
curl -X PUT http://localhost:8080/v1/users/me/settings/theme \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "dark",
    "contrast": "high"
  }'
```

## Common Commands

```bash
# Stop services
make docker-compose-down

# View logs
docker-compose logs -f

# Reset database (WARNING: deletes all data)
make migrate-down
make migrate-up

# Clean everything
make clean
docker-compose down -v
```

## Troubleshooting

### Port Already in Use

If port 8080 is in use:
```bash
# Change PORT in .env file
PORT=8081
```

### Database Connection Failed

1. Check PostgreSQL is running: `docker-compose ps`
2. Verify connection string in `.env`
3. Check logs: `docker-compose logs postgres`

### Migrate Tool Not Found

Install it:
```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

## Next Steps

- Read [README.md](README.md) for full documentation
- Check [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment
- Review API endpoints in the codebase

## Need Help?

- Check logs: `docker-compose logs`
- Verify environment: `cat .env`
- Test database: `docker-compose exec postgres psql -U baseapp -d base_app_db`

