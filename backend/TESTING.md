# Testing Base App API with Cursor AI

## Prerequisites

1. **Check prerequisites**: Run `./scripts/check-prerequisites.sh` to verify all requirements

2. **Start the server**: 
   ```bash
   make dev  # (sets up database + runs server)
   # OR
   make setup  # (sets up database only)
   make run     # (runs server)
   ```

3. Server running at `http://localhost:8080`

### Quick Setup

If you don't have Docker, you can set up PostgreSQL locally:
```bash
# macOS with Homebrew
./scripts/setup-local-db.sh

# Then run migrations and start server
make migrate-up
make run
```

## Testing methods

### Method 1: Using Cursor AI Chat

Ask Cursor AI to:

- "Test the signup endpoint with email test@example.com"

- "Create a curl command to login"

- "Get the current user profile with token X"

### Method 2: Using the test script

```bash
# Run automated test suite
./scripts/test-api.sh
```

### Method 3: Manual curl commands

**1. Health Check:**

```bash
curl http://localhost:8080/health
```

**2. Signup:**

```bash
curl -X POST http://localhost:8080/v1/auth/signup \
  -H "Content-Type: application/json" \
  -H "X-Product-Name: test-product" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!",
    "name": "Test User",
    "terms_accepted": true,
    "terms_version": "1.0"
  }'
```

**3. Login:**

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPass123!"
  }'
```

**4. Get Current User (replace TOKEN):**

```bash
curl http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**5. Update Theme:**

```bash
curl -X PUT http://localhost:8080/v1/users/me/settings/theme \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"theme": "dark", "contrast": "high"}'
```

## Quick test flow

1. Signup → get token
2. Use token for authenticated requests
3. Test theme endpoints
4. Test refresh token
5. Test logout

## Using Cursor AI terminal

- Open Cursor AI terminal (`` Ctrl+` ``)
- Run commands directly
- Copy/paste curl commands
- View responses in terminal

## Troubleshooting

- Server not running? Check: `curl http://localhost:8080/health`
- Database error? Run: `make migrate-up`
- Port conflict? Change `PORT` in `.env`

**Tip:** Ask Cursor AI: "Help me test the Base App API endpoints" for interactive testing.

