# Base App API Documentation

**Version:** 1.0  
**Base URL:** `https://api.example.com/v1` (Production) | `http://localhost:8080/v1` (Development)  
**Last Updated:** November 2025

## Table of Contents

1. [Getting Started](#getting-started)
2. [Authentication](#authentication)
3. [API Endpoints](#api-endpoints)
4. [Error Handling](#error-handling)
5. [Rate Limits](#rate-limits)
6. [Webhooks](#webhooks)
7. [Theme Sync Guide](#theme-sync-guide)
8. [Integration Guides](#integration-guides)
9. [cURL Examples](#curl-examples)

---

## Getting Started

### Base URL

```
Production: https://api.example.com/v1
Development: http://localhost:8080/v1
```

### Content Type

All requests must include:
```
Content-Type: application/json
```

### Response Format

All responses follow this structure:

**Success Response:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": { ... }  // Optional, for validation errors
  }
}
```

### Quick Start

1. **Signup** to create an account
2. **Login** to get access tokens
3. Use the **access token** in `Authorization: Bearer <token>` header
4. **Refresh** tokens before they expire

---

## Authentication

### Overview

Base App uses **JWT (JSON Web Tokens)** for authentication. The API uses a token pair system:

- **Access Token**: Short-lived (15 minutes), used for API requests
- **Refresh Token**: Long-lived (30 days), used to get new access tokens

### Authentication Flow

```
1. Signup/Login → Receive access_token + refresh_token
2. Use access_token in Authorization header
3. When access_token expires → Use refresh_token to get new access_token
4. Logout → Revoke tokens
```

### Headers

**For Authenticated Requests:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Optional Headers:**
```
X-Product-Name: your-product-name      # For signup tracking
X-Device-ID: unique-device-id          # For device tracking
X-Device-Name: My Device               # Device name
```

### Token Structure

JWT tokens contain:
- `user_id`: User UUID
- `session_id`: Session UUID
- `exp`: Expiration timestamp
- `iat`: Issued at timestamp

### Token Refresh

Access tokens expire after **15 minutes**. Use the refresh token to get a new access token:

```bash
curl -X POST https://api.example.com/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "your_refresh_token"}'
```

---

## API Endpoints

### Health Check

#### GET `/health`

Check if the API is running.

**No authentication required**

**Response:**
```json
{
  "status": "healthy"
}
```

**cURL Example:**
```bash
curl https://api.example.com/health
```

---

### Authentication Endpoints

#### POST `/v1/auth/signup`

Create a new user account.

**No authentication required**

**Request Headers:**
```
Content-Type: application/json
X-Product-Name: your-product-name (optional)
X-Device-ID: device-id (optional)
X-Device-Name: Device Name (optional)
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "name": "John Doe",
  "first_name": "John",              // Optional
  "last_name": "Doe",                 // Optional
  "phone": "+1234567890",             // Optional
  "marketing_consent": false,         // Optional
  "terms_accepted": true,              // Required
  "terms_version": "1.0"              // Required
}
```

**Validation Rules:**
- `email`: Valid email format, unique
- `password`: Minimum 8 characters
- `name`: Required
- `terms_accepted`: Must be `true`
- `terms_version`: Required

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "f3627e1d-3c21-4164-b679-b092f295868e",
      "email": "user@example.com",
      "name": "John Doe",
      "email_verified": false,
      "status": "pending"
    },
    "session": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_at": "2025-11-25T15:54:28+05:00"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Validation error
- `409 Conflict`: Email already exists
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl -X POST https://api.example.com/v1/auth/signup \
  -H "Content-Type: application/json" \
  -H "X-Product-Name: my-product" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "name": "John Doe",
    "terms_accepted": true,
    "terms_version": "1.0"
  }'
```

---

#### POST `/v1/auth/login`

Authenticate and receive access tokens.

**No authentication required**

**Request Headers:**
```
Content-Type: application/json
X-Device-ID: device-id (optional)
X-Device-Name: Device Name (optional)
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "remember_me": false  // Optional
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "f3627e1d-3c21-4164-b679-b092f295868e",
      "email": "user@example.com",
      "name": "John Doe",
      "email_verified": false,
      "status": "active"
    },
    "session": {
      "id": "af2cd970-21f4-4174-8e32-fd5456070e47",
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_at": "2025-11-25T15:54:28+05:00"
    },
    "device": {
      "id": "device-123",
      "is_new_device": false
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Validation error
- `401 Unauthorized`: Invalid credentials
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl -X POST https://api.example.com/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!"
  }'
```

---

#### POST `/v1/auth/refresh`

Refresh access token using refresh token.

**No authentication required**

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-11-25T16:09:28+05:00"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Validation error
- `401 Unauthorized`: Invalid or expired refresh token
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl -X POST https://api.example.com/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token"
  }'
```

---

#### POST `/v1/auth/logout`

Logout and revoke current session.

**Authentication required**

**Request Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "revoke_all_sessions": false  // Optional, default: false
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "Logged out successfully",
    "sessions_revoked": 1
  }
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl -X POST https://api.example.com/v1/auth/logout \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{"revoke_all_sessions": false}'
```

---

### User Endpoints

#### GET `/v1/users/me`

Get current authenticated user's profile.

**Authentication required**

**Request Headers:**
```
Authorization: Bearer <access_token>
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "f3627e1d-3c21-4164-b679-b092f295868e",
    "email": "user@example.com",
    "email_verified": false,
    "name": "John Doe",
    "first_name": "John",
    "last_name": "Doe",
    "photo_url": null,
    "phone": null,
    "phone_verified": false,
    "status": "active",
    "created_at": "2025-11-25T15:39:28.65086Z",
    "updated_at": "2025-11-25T15:39:28.65086Z"
  }
}
```

**Error Responses:**
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl https://api.example.com/v1/users/me \
  -H "Authorization: Bearer your_access_token"
```

---

#### PUT `/v1/users/me`

Update current user's profile.

**Authentication required**

**Request Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Jane Doe",           // Optional
  "first_name": "Jane",         // Optional
  "last_name": "Doe",           // Optional
  "phone": "+1234567890",       // Optional
  "photo_url": "https://..."    // Optional
}
```

**Note:** Password updates are not supported via this endpoint.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "f3627e1d-3c21-4164-b679-b092f295868e",
    "email": "user@example.com",
    "name": "Jane Doe",
    "updated_at": "2025-11-25T15:45:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl -X PUT https://api.example.com/v1/users/me \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "first_name": "Jane"
  }'
```

---

### Theme Endpoints

#### GET `/v1/users/me/settings/theme`

Get user's theme preferences.

**Authentication required**

**Request Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `product` (optional): Product name for product-specific theme override

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "theme": "dark",
    "contrast": "high",
    "text_direction": "ltr",
    "brand": null,
    "source": "global",
    "product_override": null,
    "synced_at": "2025-11-25T15:39:30Z",
    "localStorage_keys": {
      "theme": "kompassui-theme",
      "contrast": "kompassui-contrast",
      "text_direction": "kompassui-text-direction"
    }
  }
}
```

**Response Fields:**
- `theme`: `"auto"`, `"light"`, or `"dark"`
- `contrast`: `"standard"`, `"high"`, or `"low"`
- `text_direction`: `"auto"`, `"ltr"`, or `"rtl"`
- `brand`: Optional brand identifier
- `source`: `"global"` or `"product_override"`
- `product_override`: Product name if using product-specific theme

**Error Responses:**
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
# Get global theme
curl https://api.example.com/v1/users/me/settings/theme \
  -H "Authorization: Bearer your_access_token"

# Get product-specific theme
curl "https://api.example.com/v1/users/me/settings/theme?product=my-product" \
  -H "Authorization: Bearer your_access_token"
```

---

#### PUT `/v1/users/me/settings/theme`

Update user's theme preferences.

**Authentication required**

**Request Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "theme": "dark",              // Optional: "auto", "light", "dark"
  "contrast": "high",            // Optional: "standard", "high", "low"
  "text_direction": "ltr",       // Optional: "auto", "ltr", "rtl"
  "brand": "my-brand"            // Optional: Brand identifier
}
```

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "theme": "dark",
    "contrast": "high",
    "text_direction": "ltr",
    "brand": null,
    "synced_at": "2025-11-25T15:45:00Z",
    "message": "Theme preferences updated successfully"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl -X PUT https://api.example.com/v1/users/me/settings/theme \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "dark",
    "contrast": "high"
  }'
```

---

#### POST `/v1/users/me/settings/theme/sync`

Sync client theme with server (conflict detection).

**Authentication required**

**Request Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "theme": "dark",
  "contrast": "high",
  "text_direction": "ltr",
  "brand": "my-brand",           // Optional
  "client_timestamp": "2025-11-25T15:40:00Z"  // Optional
}
```

**Response:** `200 OK`

**No Conflicts (Synced):**
```json
{
  "success": true,
  "data": {
    "synced": true,
    "server_theme": {
      "theme": "dark",
      "contrast": "high",
      "text_direction": "ltr",
      "brand": null,
      "synced_at": "2025-11-25T15:45:00Z"
    },
    "conflicts": []
  }
}
```

**With Conflicts (Not Synced):**
```json
{
  "success": true,
  "data": {
    "synced": false,
    "server_theme": {
      "theme": "light",
      "contrast": "standard",
      "text_direction": "auto",
      "brand": null,
      "synced_at": "2025-11-25T15:50:00Z"
    },
    "conflicts": ["theme", "contrast"]
  }
}
```

**Conflict Detection:**
- Server timestamp is compared with client timestamp
- If server theme was updated more recently, conflicts are returned
- Client should use server theme when conflicts exist

**Error Responses:**
- `400 Bad Request`: Invalid request
- `401 Unauthorized`: Invalid or missing token
- `500 Internal Server Error`: Server error

**cURL Example:**
```bash
curl -X POST https://api.example.com/v1/users/me/settings/theme/sync \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "dark",
    "contrast": "high",
    "text_direction": "ltr"
  }'
```

---

## Error Handling

### Error Response Format

All errors follow this structure:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {  // Optional, for validation errors
      "fields": {
        "email": ["Invalid email format"],
        "password": ["Password must be at least 8 characters"]
      }
    }
  }
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_REQUEST` | 400 | Malformed request or invalid JSON |
| `VALIDATION_ERROR` | 422 | Request validation failed |
| `UNAUTHORIZED` | 401 | Missing or invalid authentication token |
| `CONFLICT` | 409 | Resource conflict (e.g., email already exists) |
| `NOT_FOUND` | 404 | Resource not found |
| `INTERNAL_ERROR` | 500 | Internal server error |

### Validation Errors

When validation fails, the `details` field contains field-specific errors:

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "fields": {
        "email": ["Invalid email format"],
        "password": ["Password must be at least 8 characters"],
        "terms_accepted": ["terms_accepted is required"]
      }
    }
  }
}
```

### Common Error Scenarios

**1. Missing Authorization Header:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing authorization header"
  }
}
```

**2. Invalid Token:**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid token"
  }
}
```

**3. Expired Token:**
- Token expires after 15 minutes
- Use refresh token to get new access token

**4. Email Already Exists:**
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "email already exists"
  }
}
```

---

## Rate Limits

### Current Limits

**Note:** Rate limiting is configured but not yet enforced. Limits will be implemented in future releases.

**Planned Limits:**
- **Authentication endpoints**: 10 requests/minute per IP
- **API endpoints**: 100 requests/minute per access token
- **Webhook delivery**: 60 events/minute per subscription (configurable)

### Rate Limit Headers

When rate limiting is enforced, responses will include:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1638360000
```

### Rate Limit Exceeded Response

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please try again later."
  }
}
```

---

## Webhooks

### Overview

Base App supports webhooks for event-driven integrations. Webhooks are delivered via HTTP POST requests to your configured endpoints.

### Webhook Events

#### Available Events

| Event Type | Description | Triggered When |
|------------|-------------|----------------|
| `user.created` | New user account created | User signs up |
| `user.updated` | User profile updated | User updates profile |
| `user.status.changed` | User status changed | User status changes |
| `session.created` | New session created | User logs in or signs up |
| `session.revoked` | Session revoked | User logs out |
| `theme.updated` | Theme preferences updated | User updates theme |
| `device.trusted` | Device marked as trusted | Device trust status changes |

### Webhook Payload Format

```json
{
  "event_id": "uuid",
  "event_type": "user.created",
  "event_version": "1.0",
  "event_source": "base_app",
  "timestamp": "2025-11-25T15:39:28Z",
  "user_id": "f3627e1d-3c21-4164-b679-b092f295868e",
  "payload": {
    // Event-specific payload data
  }
}
```

### Webhook Headers

Each webhook request includes:

```
Content-Type: application/json
X-Webhook-Signature: sha256=<hmac_signature>
X-Webhook-Timestamp: 1638360000
X-Webhook-Event-ID: <event_uuid>
X-Webhook-Event-Type: user.created
```

### Webhook Signature Verification

Webhooks are signed using HMAC SHA-256. Verify signatures to ensure requests are from Base App:

**Signature Format:**
```
sha256=<hmac_sha256(timestamp + "." + payload_json, webhook_secret)>
```

**Verification Steps:**

1. Extract `X-Webhook-Timestamp` and `X-Webhook-Signature` headers
2. Get the raw request body (JSON string)
3. Create message: `timestamp + "." + body`
4. Compute HMAC SHA-256: `hmac_sha256(message, webhook_secret)`
5. Compare with signature from header

**Example (Node.js):**
```javascript
const crypto = require('crypto');

function verifyWebhookSignature(timestamp, body, signature, secret) {
  const message = `${timestamp}.${body}`;
  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(message)
    .digest('hex');
  
  return signature === `sha256=${expectedSignature}`;
}
```

### Webhook Retry Logic

- **Max Attempts**: 3 (configurable per subscription)
- **Retry Backoff**: Exponential backoff (default: 2x multiplier)
- **Retry Schedule**: 
  - 1st retry: 1 minute
  - 2nd retry: 2 minutes
  - 3rd retry: 4 minutes

### Webhook Response

Your endpoint should return:
- **200-299**: Success (webhook delivered)
- **400-499**: Client error (will retry)
- **500-599**: Server error (will retry)

**Recommended Response:**
```json
{
  "received": true
}
```

### Webhook Subscription Management

**Note:** Webhook subscription endpoints are not yet implemented. This will be available in a future release.

Planned endpoints:
- `POST /v1/webhooks/subscriptions` - Create subscription
- `GET /v1/webhooks/subscriptions` - List subscriptions
- `PUT /v1/webhooks/subscriptions/:id` - Update subscription
- `DELETE /v1/webhooks/subscriptions/:id` - Delete subscription

---

## Theme Sync Guide

### Overview

Theme sync allows clients to synchronize theme preferences with the server while detecting conflicts.

### Sync Flow

```
1. Client loads theme from localStorage
2. Client calls GET /v1/users/me/settings/theme
3. Compare client timestamp with server synced_at
4. If conflicts → Use server theme
5. If no conflicts → Sync client theme to server
```

### Conflict Detection

Conflicts occur when:
- Server theme was updated more recently than client theme
- Different values exist for theme, contrast, or text_direction

### Implementation Example

**JavaScript/TypeScript:**
```javascript
async function syncTheme(clientTheme) {
  // Get server theme
  const serverResponse = await fetch('/v1/users/me/settings/theme', {
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  });
  
  const serverTheme = await serverResponse.json();
  
  // Sync with conflict detection
  const syncResponse = await fetch('/v1/users/me/settings/theme/sync', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      theme: clientTheme.theme,
      contrast: clientTheme.contrast,
      text_direction: clientTheme.text_direction,
      client_timestamp: clientTheme.synced_at
    })
  });
  
  const syncResult = await syncResponse.json();
  
  if (syncResult.data.synced) {
    // Successfully synced
    updateLocalStorage(syncResult.data.server_theme);
  } else {
    // Conflicts detected - use server theme
    updateLocalStorage(syncResult.data.server_theme);
    showConflictMessage(syncResult.data.conflicts);
  }
}
```

### localStorage Keys

For KompassUI integration, use these localStorage keys:

```javascript
localStorage.setItem('kompassui-theme', theme);
localStorage.setItem('kompassui-contrast', contrast);
localStorage.setItem('kompassui-text-direction', text_direction);
```

### Product-Specific Themes

To get product-specific theme override:

```javascript
const response = await fetch(
  '/v1/users/me/settings/theme?product=my-product',
  {
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  }
);
```

If `source` is `"product_override"`, the theme is product-specific.

---

## Integration Guides

### Next.js Integration

#### 1. Install Dependencies

```bash
npm install axios  # or fetch API
```

#### 2. Create API Client

**`lib/api-client.ts`:**
```typescript
import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/v1';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token refresh on 401
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          });
          localStorage.setItem('access_token', data.data.token);
          // Retry original request
          return apiClient.request(error.config);
        } catch (refreshError) {
          // Refresh failed, redirect to login
          window.location.href = '/login';
        }
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

#### 3. Authentication Hook

**`hooks/useAuth.ts`:**
```typescript
import { useState, useEffect } from 'react';
import apiClient from '@/lib/api-client';

export function useAuth() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('access_token');
    if (token) {
      apiClient.get('/users/me')
        .then(({ data }) => setUser(data.data))
        .catch(() => {
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
        })
        .finally(() => setLoading(false));
    } else {
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const { data } = await apiClient.post('/auth/login', { email, password });
    localStorage.setItem('access_token', data.data.session.token);
    localStorage.setItem('refresh_token', data.data.session.refresh_token);
    setUser(data.data.user);
    return data;
  };

  const logout = async () => {
    await apiClient.post('/auth/logout');
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    setUser(null);
  };

  return { user, loading, login, logout };
}
```

#### 4. Theme Sync Hook

**`hooks/useTheme.ts`:**
```typescript
import { useState, useEffect } from 'react';
import apiClient from '@/lib/api-client';

export function useTheme() {
  const [theme, setTheme] = useState({
    theme: 'auto',
    contrast: 'standard',
    text_direction: 'auto',
  });

  useEffect(() => {
    // Load from localStorage first
    const localTheme = {
      theme: localStorage.getItem('kompassui-theme') || 'auto',
      contrast: localStorage.getItem('kompassui-contrast') || 'standard',
      text_direction: localStorage.getItem('kompassui-text-direction') || 'auto',
    };
    setTheme(localTheme);

    // Sync with server
    syncTheme(localTheme);
  }, []);

  const syncTheme = async (clientTheme: any) => {
    try {
      const { data } = await apiClient.post('/users/me/settings/theme/sync', {
        ...clientTheme,
        client_timestamp: new Date().toISOString(),
      });

      if (data.data.synced) {
        // Update localStorage with server theme
        localStorage.setItem('kompassui-theme', data.data.server_theme.theme);
        localStorage.setItem('kompassui-contrast', data.data.server_theme.contrast);
        localStorage.setItem('kompassui-text-direction', data.data.server_theme.text_direction);
        setTheme(data.data.server_theme);
      } else {
        // Conflicts - use server theme
        setTheme(data.data.server_theme);
        console.warn('Theme conflicts:', data.data.conflicts);
      }
    } catch (error) {
      console.error('Theme sync failed:', error);
    }
  };

  const updateTheme = async (updates: any) => {
    try {
      const { data } = await apiClient.put('/users/me/settings/theme', updates);
      setTheme(data.data);
      // Update localStorage
      localStorage.setItem('kompassui-theme', data.data.theme);
      localStorage.setItem('kompassui-contrast', data.data.contrast);
      localStorage.setItem('kompassui-text-direction', data.data.text_direction);
    } catch (error) {
      console.error('Theme update failed:', error);
    }
  };

  return { theme, updateTheme, syncTheme };
}
```

---

### Mobile Integration (React Native / iOS / Android)

#### React Native Example

**`services/api.ts`:**
```typescript
import AsyncStorage from '@react-native-async-storage/async-storage';
import axios from 'axios';

const API_BASE_URL = 'https://api.example.com/v1';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token
apiClient.interceptors.request.use(async (config) => {
  const token = await AsyncStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = await AsyncStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const { data } = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          });
          await AsyncStorage.setItem('access_token', data.data.token);
          return apiClient.request(error.config);
        } catch {
          // Redirect to login
        }
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

**Usage:**
```typescript
import apiClient from './services/api';

// Login
const login = async (email: string, password: string) => {
  const { data } = await apiClient.post('/auth/login', { email, password });
  await AsyncStorage.setItem('access_token', data.data.session.token);
  await AsyncStorage.setItem('refresh_token', data.data.session.refresh_token);
  return data.data.user;
};

// Get user profile
const getProfile = async () => {
  const { data } = await apiClient.get('/users/me');
  return data.data;
};

// Update theme
const updateTheme = async (theme: string, contrast: string) => {
  const { data } = await apiClient.put('/users/me/settings/theme', {
    theme,
    contrast,
  });
  return data.data;
};
```

---

## cURL Examples

### Complete Authentication Flow

```bash
# 1. Signup
curl -X POST http://localhost:8080/v1/auth/signup \
  -H "Content-Type: application/json" \
  -H "X-Product-Name: my-product" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "name": "John Doe",
    "terms_accepted": true,
    "terms_version": "1.0"
  }'

# Save the token from response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 2. Get user profile
curl http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer $TOKEN"

# 3. Update profile
curl -X PUT http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "first_name": "Jane"
  }'

# 4. Update theme
curl -X PUT http://localhost:8080/v1/users/me/settings/theme \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "dark",
    "contrast": "high"
  }'

# 5. Get theme
curl http://localhost:8080/v1/users/me/settings/theme \
  -H "Authorization: Bearer $TOKEN"

# 6. Refresh token
REFRESH_TOKEN="your_refresh_token"
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}"

# 7. Logout
curl -X POST http://localhost:8080/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"revoke_all_sessions": false}'
```

### Error Handling Examples

```bash
# Invalid email format
curl -X POST http://localhost:8080/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "SecurePass123!",
    "name": "John Doe",
    "terms_accepted": true,
    "terms_version": "1.0"
  }'

# Missing authorization
curl http://localhost:8080/v1/users/me

# Invalid token
curl http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer invalid_token"
```

---

## Best Practices

### Security

1. **Never expose tokens** in client-side code or logs
2. **Use HTTPS** in production
3. **Store tokens securely**: Use secure storage (Keychain on iOS, Keystore on Android)
4. **Validate webhook signatures** before processing
5. **Rotate refresh tokens** periodically

### Performance

1. **Cache user profile** to reduce API calls
2. **Batch theme updates** when possible
3. **Use token refresh** before expiration
4. **Implement retry logic** with exponential backoff

### Error Handling

1. **Handle 401 errors** by refreshing tokens
2. **Show user-friendly messages** for validation errors
3. **Log errors** for debugging
4. **Implement offline support** where possible

### Theme Sync

1. **Sync on app startup** to get latest preferences
2. **Handle conflicts gracefully** by using server theme
3. **Update localStorage** after successful sync
4. **Sync after local changes** to keep server updated

---

## Support & Resources

- **API Base URL**: `https://api.example.com/v1`
- **Documentation**: This file
- **Status Page**: `https://status.example.com`
- **Support Email**: support@example.com

---

## Changelog

### Version 1.0 (November 2025)
- Initial API release
- Authentication endpoints
- User management
- Theme preferences
- Session management
- Device tracking

---

**Last Updated:** November 2025  
**API Version:** 1.0

