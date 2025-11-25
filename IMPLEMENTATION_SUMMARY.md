# Base App API - Detailed Implementation Summary

## Overview

The Base App Service is a production-ready Go backend API designed as a shared authentication and user management service for multiple products. It provides JWT-based authentication, user profile management, theme preferences, session management, and device tracking.

## Architecture

### Layered Architecture

The application follows a clean, layered architecture:

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

## Core Components

### 1. Configuration (`internal/config/`)

**Environment-based configuration** with sensible defaults:

- **Server**: Port, environment (dev/prod)
- **Database**: Connection pool settings (max connections, idle connections, lifetime)
- **JWT**: Secret, access token expiry (15m), refresh token expiry (720h)
- **Redis**: Optional cache configuration
- **Webhooks**: Retry logic, HMAC signing
- **Logging**: Level (info/debug), format (json/text)

### 2. Database Layer (`internal/database/`)

**Connection Management**:
- PostgreSQL connection pooling
- Configurable max connections (default: 25)
- Idle connection management (default: 5)
- Connection lifetime management (default: 300s)
- Transaction helper for atomic operations

**Database Schema**:
- **Users**: Email, password hash, profile info, status, signup tracking
- **Sessions**: JWT tokens, device info, IP address, expiration
- **User Devices**: Device tracking, trust management
- **User Settings**: Theme preferences, notifications, privacy
- **Product Theme Overrides**: Product-specific theme customization
- **Webhook Subscriptions**: Webhook endpoint configuration
- **Webhook Events**: Outbox pattern for reliable webhook delivery

### 3. Authentication System

#### JWT Token Generation (`pkg/auth/jwt.go`)

- **Algorithm**: HS256 (HMAC SHA-256)
- **Token Claims**:
  - `user_id`: UUID of the user
  - `session_id`: UUID of the session
  - Standard JWT claims (exp, iat, nbf, jti)
- **Token Pair**: Access token + Refresh token
- **Expiry**: Configurable (default: 15m access, 30 days refresh)

#### Password Security (`pkg/auth/password.go`)

- **Hashing**: bcrypt with cost factor 12
- **Verification**: Secure password comparison
- **No plaintext storage**: Passwords never stored in plaintext

#### Authentication Flow

1. **Signup**:
   - Validates email uniqueness
   - Hashes password with bcrypt
   - Creates user with "pending" status
   - Creates session with JWT tokens
   - Tracks device and IP address
   - Returns user info + tokens

2. **Login**:
   - Validates email/password
   - Checks user status (active/pending)
   - Creates/updates device record
   - Creates new session
   - Updates last login timestamp
   - Returns user info + tokens + device status

3. **Token Refresh**:
   - Validates refresh token
   - Checks session validity
   - Generates new token pair
   - Updates session in database

4. **Logout**:
   - Revokes current session
   - Option to revoke all user sessions
   - Marks session as inactive

### 4. Middleware Stack

#### Authentication Middleware (`internal/middleware/auth.go`)

- Validates Bearer token from Authorization header
- Extracts user ID and session ID from JWT claims
- Injects user context into request
- Returns 401 for invalid/missing tokens

#### CORS Middleware (`internal/middleware/cors.go`)

- Handles cross-origin requests
- Configurable allowed origins, methods, headers

#### Logging Middleware (`internal/middleware/logging.go`)

- Structured request/response logging
- Logs method, path, status, duration, user agent

#### Error Recovery (`internal/middleware/recovery.go`)

- Panic recovery
- Returns 500 error for unhandled panics
- Prevents server crashes

### 5. API Endpoints

#### Public Endpoints (No Authentication)

**POST `/v1/auth/signup`**
- Creates new user account
- Requires: email, password (min 8 chars), name, terms acceptance
- Optional: first_name, last_name, phone, marketing_consent
- Headers: `X-Product-Name`, `X-Device-ID`, `X-Device-Name`
- Returns: User info + access token + refresh token

**POST `/v1/auth/login`**
- Authenticates user
- Requires: email, password
- Optional: remember_me
- Headers: `X-Device-ID`, `X-Device-Name`
- Returns: User info + tokens + device status

**POST `/v1/auth/refresh`**
- Refreshes access token
- Requires: refresh_token
- Returns: New access token + expiry

**GET `/health`**
- Health check endpoint
- Returns: `{"status":"healthy"}`

#### Protected Endpoints (Require Authentication)

**POST `/v1/auth/logout`**
- Revokes current session
- Optional: `revoke_all_sessions` flag
- Returns: Success message

**GET `/v1/users/me`**
- Gets current user profile
- Returns: Full user details (email, name, status, timestamps)

**PUT `/v1/users/me`**
- Updates user profile
- Optional fields: name, first_name, last_name, phone, photo_url
- Returns: Updated user info

**GET `/v1/users/me/settings/theme`**
- Gets theme preferences
- Query param: `?product=product-name` (for product-specific theme)
- Returns: Theme, contrast, text_direction, brand, sync info

**PUT `/v1/users/me/settings/theme`**
- Updates theme preferences
- Optional fields: theme, contrast, text_direction, brand
- Returns: Updated theme preferences

**POST `/v1/users/me/settings/theme/sync`**
- Syncs client theme with server
- Detects conflicts based on timestamps
- Returns: Server theme + conflict list

### 6. Theme Management System

#### Features

- **Global Theme**: User-wide theme preferences
- **Product Overrides**: Product-specific theme customization
- **Conflict Detection**: Timestamp-based sync conflict detection
- **KompassUI Integration**: localStorage key mapping for frontend

#### Theme Properties

- **Theme**: `auto`, `light`, `dark`
- **Contrast**: `standard`, `high`, `low`
- **Text Direction**: `auto`, `ltr`, `rtl`
- **Brand**: Optional brand identifier

#### Sync Logic

- Compares client and server timestamps
- Server wins if server timestamp is newer
- Returns conflicts list when conflicts detected
- Updates server theme if no conflicts

### 7. Session Management

#### Features

- **Multi-device Support**: Track multiple devices per user
- **Device Information**: IP address, user agent, device ID, device name
- **Session Expiration**: Configurable token expiry
- **Session Revocation**: Single or all sessions
- **Active Session Tracking**: `is_active` flag for session state

#### Session Lifecycle

1. Created on signup/login
2. Updated on token refresh
3. Tracked via `last_used_at` timestamp
4. Revoked on logout or expiration
5. Cleaned up via CASCADE on user deletion

### 8. Device Management

#### Features

- **Device Tracking**: Unique device ID per user
- **Device Metadata**: Name, type, OS, browser
- **Trust Management**: `is_trusted` flag for trusted devices
- **Location Tracking**: Country and city (optional)
- **Last Used**: Timestamp tracking

#### Device Lifecycle

- Created on first login with device ID
- Updated on subsequent logins
- Linked to sessions via device_id
- Can be marked as trusted

### 9. Error Handling (`pkg/errors/`)

#### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": { ... }  // Optional
  }
}
```

#### Error Codes

- `INVALID_REQUEST`: Bad request format
- `VALIDATION_ERROR`: Request validation failed
- `UNAUTHORIZED`: Authentication/authorization failed
- `CONFLICT`: Resource conflict (e.g., email exists)
- `NOT_FOUND`: Resource not found
- `INTERNAL_ERROR`: Server error

#### Validation Errors

- Field-level validation messages
- Detailed error per field
- Uses go-playground/validator tags

### 10. Security Features

#### Authentication Security

- **JWT Tokens**: Signed with HS256, includes user and session IDs
- **Password Hashing**: bcrypt with cost factor 12
- **Token Expiry**: Short-lived access tokens (15m), long-lived refresh tokens (30d)
- **Session Validation**: Database-backed session validation

#### Data Security

- **IP Address Parsing**: Removes port numbers for PostgreSQL INET type
- **Field Length Validation**: Truncates long fields to match DB constraints
- **Input Validation**: Request validation before processing
- **SQL Injection Protection**: Parameterized queries via database/sql

#### Infrastructure Security

- **CORS**: Configurable cross-origin policies
- **Error Messages**: Generic error messages to prevent information leakage
- **Password Storage**: Never stored in plaintext
- **Token Storage**: Tokens stored as TEXT (supports long JWT tokens)

## Database Schema Details

### Users Table

- **Primary Key**: UUID
- **Email**: Unique, validated format
- **Status**: `active`, `pending`, `suspended`, `deleted`
- **Password**: Hash + `password_changed_at`
- **Verification**: Email and phone verification tokens
- **Signup Tracking**: Source, platform, campaign, referrer
- **Indexes**: Email, status, created_at, verification tokens

### Sessions Table

- **Primary Key**: UUID
- **Foreign Key**: `user_id` → users(id) CASCADE
- **Tokens**: TEXT (supports long JWT tokens)
- **Device Info**: Device ID, name, type, OS, browser
- **Location**: IP address (INET), country, city
- **State**: `is_active`, `revoked_at`, `revoked_reason`
- **Indexes**: User ID, token, refresh_token, active sessions, expiration

### User Devices Table

- **Primary Key**: UUID
- **Foreign Key**: `user_id` → users(id) CASCADE
- **Device ID**: Unique per user
- **Trust Management**: `is_trusted`, `trusted_at`
- **Indexes**: User ID, device ID, trusted devices

### User Settings Table

- **Primary Key**: `user_id` (references users)
- **Theme Preferences**: KompassUI theme settings
- **Notifications**: Email, push, SMS preferences
- **Privacy**: Visibility settings
- **Accessibility**: High contrast, reduced motion, screen reader

### Product Theme Overrides Table

- **Primary Key**: UUID
- **Unique Constraint**: (user_id, product_name)
- **Theme Overrides**: Product-specific theme settings
- **Indexes**: User ID, product name

## Request/Response Examples

### Signup Request

```json
POST /v1/auth/signup
Headers:
  Content-Type: application/json
  X-Product-Name: my-product
  X-Device-ID: device-123
Body:
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John Doe",
  "terms_accepted": true,
  "terms_version": "1.0"
}
```

### Signup Response

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe",
      "email_verified": false,
      "status": "pending"
    },
    "session": {
      "token": "eyJhbGci...",
      "refresh_token": "eyJhbGci...",
      "expires_at": "2025-11-25T15:54:28+05:00"
    }
  }
}
```

### Login Request

```json
POST /v1/auth/login
Body:
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

### Get User Profile

```json
GET /v1/users/me
Headers:
  Authorization: Bearer eyJhbGci...
```

### Update Theme

```json
PUT /v1/users/me/settings/theme
Headers:
  Authorization: Bearer eyJhbGci...
Body:
{
  "theme": "dark",
  "contrast": "high"
}
```

## Testing Infrastructure

### Test Scripts

1. **`scripts/test-api.sh`**: Comprehensive API test suite
   - Tests all endpoints sequentially
   - Extracts and uses tokens automatically
   - Color-coded output
   - macOS compatible

2. **`scripts/check-prerequisites.sh`**: Prerequisites checker
   - Verifies Go, PostgreSQL, Docker, migrate tool
   - Checks server status
   - Provides setup instructions

3. **`scripts/setup-local-db.sh`**: Local database setup helper
   - Automates PostgreSQL setup on macOS
   - Creates database and user

### Test Coverage

- Health check
- User signup
- User login
- Get current user
- Update theme
- Get theme
- Refresh token
- Logout

## Deployment

### Server Configuration

- **Graceful Shutdown**: Handles SIGINT/SIGTERM
- **Timeouts**: Read (15s), Write (15s), Idle (60s)
- **Connection Pooling**: Configurable database connections
- **Logging**: Structured JSON logging for production

### Environment Variables

All configuration via environment variables with sensible defaults:
- Database connection settings
- JWT secrets and expiry
- Server port and environment
- Redis configuration (optional)
- Webhook settings

### Docker Support

- Dockerfile included
- docker-compose.yml for local development
- Cloud Build configuration for GCP deployment

## Key Implementation Details

### IP Address Handling

- Extracts IP from `X-Forwarded-For`, `X-Real-IP`, or `RemoteAddr`
- Removes port numbers (PostgreSQL INET type requirement)
- Handles both IPv4 and IPv6 addresses

### Field Length Validation

- Truncates long fields to match database constraints:
  - `signup_source`: VARCHAR(100)
  - `device_name`: VARCHAR(255)
  - `device_id`: VARCHAR(255)

### Token Storage

- Migrated from VARCHAR(255) to TEXT
- Supports JWT tokens of any length (typically 200-400+ characters)
- Migration: `002_fix_token_lengths.up.sql`

### Error Recovery

- Panic recovery middleware prevents server crashes
- Returns proper HTTP error responses
- Logs errors for debugging

## API Response Format

### Success Response

```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

## Future Enhancements (Infrastructure Ready)

- **Webhooks**: Event-driven webhook system with retry logic
- **Rate Limiting**: Redis-based rate limiting (configured but not implemented)
- **Email Verification**: Email verification token system
- **Phone Verification**: Phone verification code system
- **Multi-product Support**: Product-specific theme overrides

## Summary

The Base App API is a well-architected, production-ready service providing:

✅ **Secure Authentication**: JWT-based with refresh tokens  
✅ **User Management**: Complete user lifecycle management  
✅ **Theme System**: Global and product-specific theme preferences  
✅ **Session Management**: Multi-device session tracking  
✅ **Device Tracking**: Device management and trust  
✅ **Error Handling**: Comprehensive error responses  
✅ **Security**: Password hashing, input validation, SQL injection protection  
✅ **Testing**: Automated test scripts and infrastructure  
✅ **Documentation**: Comprehensive guides and API documentation  
✅ **Deployment Ready**: Docker, migrations, configuration management  

The implementation follows Go best practices with clean architecture, separation of concerns, and production-ready features.

