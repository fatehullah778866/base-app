# Base App Service

A production-ready Go backend service for user authentication, account management, and theme preferences. Designed to be used as a shared service across multiple products. This code now lives under `backend/` (run commands from this folder). Static assets are served from `../frontend` by default; override with `FRONTEND_DIR`.

## Features

- **Authentication**: Signup, login, refresh token, logout
- **User Management**: Profile management, account settings
- **Theme Management**: Global and product-specific theme preferences (KompassUI integration)
- **Session Management**: Multi-device session tracking and revocation
- **Device Management**: Device tracking and trust management
- **Webhooks**: Event-driven webhook system with retry logic and HMAC signing
- **Security**: JWT-based authentication, password hashing (bcrypt), rate limiting
- **Admin Mode**: Built-in admin account with dashboard APIs for user management, access requests, and activity logs

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL (Cloud SQL compatible)
- **Cache**: Redis (for rate limiting and sessions)
- **HTTP Router**: Gorilla Mux
- **Logging**: Zap
- **Validation**: go-playground/validator

## Project Structure

```
base-app/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and helpers
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware (auth, logging, CORS, recovery)
│   ├── models/          # Data models
│   ├── repositories/    # Data access layer
│   ├── services/        # Business logic layer
│   └── webhooks/        # Webhook emitter and dispatcher
├── pkg/
│   ├── auth/            # Authentication utilities (JWT, password hashing)
│   ├── device/          # Device detection utilities
│   └── errors/          # Error handling utilities
├── migrations/          # Database migrations
└── tests/               # Test files
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL 12+
- Redis (optional, for rate limiting)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/kompass-tech/base-app.git
cd base-app
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run database migrations:
```bash
# Using your preferred migration tool (e.g., migrate, golang-migrate)
migrate -path migrations -database "postgres://user:password@localhost/base_app_db?sslmode=disable" up
```

5. Run the server:
```bash
go run cmd/server/main.go
```

The server will start on port 8080 (or the port specified in `PORT` environment variable).

## Configuration

Environment variables:

```bash
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=baseapp
DB_PASSWORD=password
DB_NAME=base_app_db
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONNECTION_MAX_LIFETIME=300s

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=change-me-in-production
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=720h

# Webhooks
WEBHOOK_SECRET=change-me-in-production
WEBHOOK_MAX_RETRIES=3
WEBHOOK_RETRY_BACKOFF_MULTIPLIER=2.0

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REDIS_KEY_PREFIX=ratelimit:

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## API Endpoints

### Authentication

- `POST /v1/auth/signup` - Create a new user account
- `POST /v1/auth/login` - Login with email and password
- `POST /v1/auth/refresh` - Refresh access token
- `POST /v1/auth/logout` - Logout (revoke session)

### Users

- `GET /v1/users/me` - Get current user profile
- `PUT /v1/users/me` - Update user profile
- `POST /v1/requests` - Create a user access/request ticket
- `GET /v1/requests` - List access requests for current user

### Theme

- `GET /v1/users/me/settings/theme` - Get theme preferences
- `PUT /v1/users/me/settings/theme` - Update theme preferences
- `POST /v1/users/me/settings/theme/sync` - Sync theme with server

### Admin

- `POST /v1/admin/login` - Admin-only login (`admin@gmail.com` / `admin123` seeded)
- `GET /v1/admin/users` - List/search users
- `POST /v1/admin/users/{id}/status` - Enable/disable a user
- `GET /v1/admin/logs` - View recent activity logs
- `POST /v1/admin/admins` - Create another admin account
- `GET /v1/admin/requests` - List all user requests
- `POST /v1/admin/requests/{id}/status` - Approve/reject/pending with feedback

## Docker

Build the Docker image:

```bash
docker build -t base-app-service .
```

Run the container:

```bash
docker run -p 8080:8080 --env-file .env base-app-service
```

## Development

### Running Tests

```bash
go test ./...
```

### API Testing

For testing the API endpoints, see [TESTING.md](TESTING.md) for comprehensive testing instructions.

**Check prerequisites:**
```bash
./scripts/check-prerequisites.sh
```

**Quick test:**
```bash
# Run automated test suite
./scripts/test-api.sh
```

The test script will automatically check if the server is running before testing.

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## Deployment

The service is designed to be deployed on Google Cloud Platform:

- **Cloud Run**: Serverless container deployment
- **Cloud SQL**: Managed PostgreSQL database
- **Memorystore**: Managed Redis instance

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed deployment instructions and `cloudbuild.yaml` for CI/CD pipeline configuration.

## Documentation

### Quick Start Guides

- [QUICKSTART.md](QUICKSTART.md) - Get started in 5 minutes
- [TESTING.md](TESTING.md) - Comprehensive testing guide
- [DEPLOYMENT.md](DEPLOYMENT.md) - Production deployment guide
- [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - Complete API reference

### Technical Reports

For detailed technical documentation, implementation reports, security audits, and service-specific reports, see the [Reports Center](docs/reports/README.md):

- **Technical Reports**: Architecture, API specs, database schema
- **Implementation Reports**: Status, milestones, feature completion
- **Security Audits**: Security reviews and recommendations
- **Service Reports**: Auth, Theme, and Webhook service documentation
- **Analysis Reports**: Performance, code coverage, optimization

## License

[Your License Here]
