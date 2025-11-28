# Base-App Technical Report

**Date:** November 2025  
**Status:** Complete  
**Category:** Technical  
**Service:** all  
**Version:** 1.0

## Summary

This report provides a comprehensive technical overview of the Base-App service, including architecture, technology stack, API endpoints, database schema, and deployment architecture.

## Architecture Overview

### Layered Architecture

The Base-App service follows a clean, layered architecture pattern:

```
┌─────────────────────────────────────┐
│         HTTP Handlers               │  (Presentation Layer)
├─────────────────────────────────────┤
│         Services                    │  (Business Logic Layer)
├─────────────────────────────────────┤
│         Repositories                │  (Data Access Layer)
├─────────────────────────────────────┤
│         Database (PostgreSQL)       │  (Persistence Layer)
└─────────────────────────────────────┘
```

### Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 15 (Cloud SQL compatible)
- **Cache**: Redis (optional, for rate limiting)
- **HTTP Router**: Gorilla Mux
- **Authentication**: JWT (HS256)
- **Password Hashing**: bcrypt (cost factor 12)
- **Logging**: Zap (structured logging)
- **Validation**: go-playground/validator
- **Migrations**: golang-migrate

## Project Structure

```
base-app/
├── cmd/server/              # Application entry point
├── internal/
│   ├── config/              # Configuration management (env-based)
│   ├── database/            # Database connection & connection pooling
│   ├── handlers/            # HTTP request handlers (auth, user, theme)
│   ├── middleware/         # HTTP middleware (auth, CORS, logging, recovery)
│   ├── models/             # Domain models (User, Session, Theme, Device, Webhook)
│   ├── repositories/       # Data access layer (interfaces + implementations)
│   ├── services/           # Business logic layer
│   └── webhooks/           # Webhook emitter and dispatcher
├── pkg/
│   ├── auth/               # JWT & password utilities
│   ├── device/             # Device detection utilities
│   └── errors/             # Error handling utilities
├── migrations/             # Database migrations (up/down)
└── scripts/                # Utility scripts (testing, migration, setup)
```

## API Endpoints

### Authentication Endpoints

- `POST /v1/auth/signup` - Create new user account
- `POST /v1/auth/login` - Authenticate user
- `POST /v1/auth/refresh` - Refresh access token
- `POST /v1/auth/logout` - Revoke session

### User Endpoints

- `GET /v1/users/me` - Get current user profile
- `PUT /v1/users/me` - Update user profile

### Theme Endpoints

- `GET /v1/users/me/settings/theme` - Get theme preferences
- `PUT /v1/users/me/settings/theme` - Update theme preferences
- `POST /v1/users/me/settings/theme/sync` - Sync theme with conflict detection

### Health Check

- `GET /health` - Health check endpoint

## Database Schema

### Core Tables

**Users Table:**
- Primary key: UUID
- Email (unique, validated)
- Password hash (bcrypt)
- Profile information (name, first_name, last_name, photo_url, phone)
- Status (active, pending, suspended, deleted)
- Verification tokens (email, phone)
- Signup tracking (source, platform, campaign, referrer)
- Timestamps (created_at, updated_at, last_login_at)

**Sessions Table:**
- Primary key: UUID
- Foreign key: user_id → users(id) CASCADE)
- JWT tokens (TEXT - supports long tokens)
- Device information (device_id, device_name, device_type, os, browser)
- Location (IP address as INET, country, city)
- State (is_active, revoked_at, revoked_reason)
- Expiration timestamps

**User Devices Table:**
- Primary key: UUID
- Foreign key: → users(id) CASCADE
- Device ID (unique per user)
- Device metadata (name, type, OS, browser)
- Trust management (is_trusted, trusted_at)
- Last used timestamp

**User Settings Table:**
- Primary key: user_id (references users)
- Theme preferences (KompassUI theme settings)
- Notification preferences
- Privacy settings
- Accessibility settings

**Product Theme Overrides Table:**
- Primary key: UUID
- Unique constraint: (user_id, product_name)
- Product-specific theme settings
- Override timestamps

**Webhook Tables:**
- `webhook_subscriptions` - Webhook endpoint configuration
- `webhook_events` - Outbox pattern for reliable webhook delivery

## Configuration

### Environment Variables

**Server:**
- `PORT` - Server port (default: 8080)
- `ENV` - Environment (development, production)

**Database:**
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `DB_SSL_MODE` - SSL mode for database connection
- `DB_MAX_CONNECTIONS` - Connection pool size (default: 25)
- `DB_MAX_IDLE_CONNECTIONS` - Idle connections (default: 5)
- `DB_CONNECTION_MAX_LIFETIME` - Connection lifetime (default: 300s)

**JWT:**
- `JWT_SECRET` - Secret key for JWT signing
- `JWT_ACCESS_TOKEN_EXPIRY` - Access token expiry (default: 15m)
- `JWT_REFRESH_TOKEN_EXPIRY` - Refresh token expiry (default: 720h)

**Redis (Optional):**
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`

**Webhooks:**
- `WEBHOOK_SECRET` - Secret for webhook HMAC signing
- `WEBHOOK_MAX_RETRIES` - Maximum retry attempts (default: 3)
- `WEBHOOK_RETRY_BACKOFF_MULTIPLIER` - Exponential backoff multiplier (default: 2.0)

**Logging:**
- `LOG_LEVEL` - Log level (info, debug, warn, error)
- `LOG_FORMAT` - Log format (json, text)

## Security Features

### Authentication Security

- **JWT Tokens**: Signed with HS256, includes user and session IDs
- **Password Hashing**: bcrypt with cost factor 12
- **Token Expiry**: Short-lived access tokens (15m), long-lived refresh tokens (30d)
- **Session Validation**: Database-backed session validation

### Data Security

- **IP Address Parsing**: Removes port numbers for PostgreSQL INET type
- **Field Length Validation**: Truncates long fields to match DB constraints
- **Input Validation**: Request validation before processing
- **SQL Injection Protection**: Parameterized queries via database/sql

### Infrastructure Security

- **CORS**: Configurable cross-origin policies
- **Error Messages**: Generic error messages to prevent information leakage
- **Password Storage**: Never stored in plaintext
- **Token Storage**: Tokens stored as TEXT (supports long JWT tokens)

## Deployment Architecture

### Target Platform: Google Cloud Platform

- **Cloud Run**: Serverless container deployment
- **Cloud SQL**: Managed PostgreSQL database
- **Memorystore**: Managed Redis instance (optional)

### Containerization

- Dockerfile included for container builds
- docker-compose.yml for local development
- Cloud Build configuration for CI/CD

### Scalability

- Stateless design (sessions stored in database)
- Connection pooling for database efficiency
- Horizontal scaling via Cloud Run
- Graceful shutdown handling

## Performance Considerations

### Database

- Connection pooling (max 25 connections)
- Indexed queries (email, user_id, session tokens)
- CASCADE deletes for data consistency

### Caching

- Redis configured for rate limiting (not yet implemented)
- Future: Session caching, user profile caching

### API Performance

- Structured logging for monitoring
- Request/response logging middleware
- Error recovery middleware prevents crashes

## Related Reports

- [Implementation Summary](../implementation/implementation-summary.md)
- [Security Audit](../audits/security/initial-security-audit.md)
- [Auth Service Report](../services/auth/auth-service-complete.md)

---

**Last Updated:** November 2025

