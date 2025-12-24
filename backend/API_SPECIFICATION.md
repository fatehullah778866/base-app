# Complete API Specification

## Base URL
```
http://localhost:8080/v1
```

## Authentication

All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

## Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "Optional message"
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": { ... }
  }
}
```

## Status Codes

- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `500 Internal Server Error` - Server error

---

## Public Endpoints (No Authentication Required)

### Authentication

#### User Signup
```http
POST /auth/signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "john@example.com",
      "name": "John Doe",
      "role": "user"
    },
    "session": {
      "token": "jwt_token",
      "refresh_token": "refresh_token",
      "expires_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### User Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:** Same as signup

#### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "new_jwt_token",
    "refresh_token": "new_refresh_token",
    "expires_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Forgot Password
```http
POST /auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

#### Reset Password
```http
POST /auth/reset-password
Content-Type: application/json

{
  "token": "reset_token",
  "password": "NewSecurePass123!"
}
```

### Admin

#### Admin Login
```http
POST /admin/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "AdminPass123!"
}
```

#### Verify Admin Code
```http
POST /admin/verify-code
Content-Type: application/json

{
  "verification_code": "Kompasstech2025@"
}
```

#### Create Admin (Public)
```http
POST /admin/create
Content-Type: application/json

{
  "name": "Admin Name",
  "email": "admin@example.com",
  "password": "AdminPass123!",
  "verification_code": "Kompasstech2025@"
}
```

---

## Protected User Endpoints

### Profile & Account

#### Get Current User
```http
GET /users/me
Authorization: Bearer <token>
```

#### Update Profile
```http
PUT /users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Name",
  "email": "newemail@example.com"
}
```

#### Change Password
```http
PUT /users/me/password
Authorization: Bearer <token>
Content-Type: application/json

{
  "current_password": "OldPass123!",
  "new_password": "NewPass123!"
}
```

#### Export User Data
```http
GET /users/me/export
Authorization: Bearer <token>
```

#### Request Account Deletion
```http
POST /users/me/delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "days_until_deletion": 30
}
```

### Settings

#### Get All Settings
```http
GET /users/me/settings
Authorization: Bearer <token>
```

#### Update Profile Settings
```http
PUT /users/me/settings/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "display_name": "Display Name",
  "username": "username",
  "bio": "Bio text",
  "date_of_birth": "1990-01-01"
}
```

#### Update Security Settings
```http
PUT /users/me/settings/security
Authorization: Bearer <token>
Content-Type: application/json

{
  "two_factor_enabled": true,
  "security_questions": [...]
}
```

#### Update Privacy Settings
```http
PUT /users/me/settings/privacy
Authorization: Bearer <token>
Content-Type: application/json

{
  "profile_visibility": "public",
  "email_visibility": "private",
  "phone_visibility": "private"
}
```

#### Update Notification Settings
```http
PUT /users/me/settings/notifications
Authorization: Bearer <token>
Content-Type: application/json

{
  "email_notifications": true,
  "sms_notifications": false,
  "push_notifications": true
}
```

#### Update Account Preferences
```http
PUT /users/me/settings/preferences
Authorization: Bearer <token>
Content-Type: application/json

{
  "language": "en",
  "timezone": "UTC",
  "theme": "dark",
  "accessibility_options": {...}
}
```

#### Add Connected Account
```http
POST /users/me/settings/connected-accounts
Authorization: Bearer <token>
Content-Type: application/json

{
  "provider": "google",
  "provider_id": "google_user_id",
  "email": "user@gmail.com"
}
```

#### Remove Connected Account
```http
DELETE /users/me/settings/connected-accounts
Authorization: Bearer <token>
Content-Type: application/json

{
  "provider": "google"
}
```

#### Deactivate Account
```http
POST /users/me/settings/account/deactivate
Authorization: Bearer <token>
```

#### Reactivate Account
```http
POST /users/me/settings/account/reactivate
Authorization: Bearer <token>
```

#### Get Active Sessions
```http
GET /users/me/settings/sessions
Authorization: Bearer <token>
```

#### Logout All Devices
```http
POST /users/me/settings/sessions/logout-all
Authorization: Bearer <token>
```

### Dashboard

#### Create Dashboard Item
```http
POST /dashboard/items
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Item Title",
  "description": "Item description"
}
```

#### List Dashboard Items
```http
GET /dashboard/items?status=active
Authorization: Bearer <token>
```

#### Get Dashboard Item
```http
GET /dashboard/items/{id}
Authorization: Bearer <token>
```

#### Update Dashboard Item
```http
PUT /dashboard/items/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "description": "Updated description"
}
```

#### Delete Dashboard Item
```http
DELETE /dashboard/items/{id}
Authorization: Bearer <token>
```

#### Archive Dashboard Item
```http
POST /dashboard/items/{id}/archive
Authorization: Bearer <token>
```

### Notifications

#### Get Notifications
```http
GET /notifications?unread_only=true&limit=50
Authorization: Bearer <token>
```

#### Get Unread Count
```http
GET /notifications/unread-count
Authorization: Bearer <token>
```

#### Mark as Read
```http
POST /notifications/read
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": "notification_id"
}
```

#### Mark All as Read
```http
POST /notifications/read-all
Authorization: Bearer <token>
```

#### Delete Notification
```http
DELETE /notifications
Authorization: Bearer <token>
Content-Type: application/json

{
  "id": "notification_id"
}
```

### Messaging

#### Send Message
```http
POST /messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "recipient_id": "user_uuid",
  "content": "Message content"
}
```

#### Get Conversations
```http
GET /messages/conversations
Authorization: Bearer <token>
```

#### Get Messages
```http
GET /messages?conversation_id={id}&limit=50
Authorization: Bearer <token>
```

#### Mark Message as Read
```http
POST /messages/read
Authorization: Bearer <token>
Content-Type: application/json

{
  "message_id": "message_uuid"
}
```

#### Get Unread Message Count
```http
GET /messages/unread-count
Authorization: Bearer <token>
```

### Search

#### Search
```http
GET /search?q=query&type=all&limit=20
Authorization: Bearer <token>
```

**Query Parameters:**
- `q` - Search query (required)
- `type` - Type: `all`, `users`, `dashboard_items`, `messages`
- `limit` - Results limit (default: 20)

### Account Switching

#### Switch Account
```http
POST /account/switch
Authorization: Bearer <token>
Content-Type: application/json

{
  "target_account_id": "account_uuid"
}
```

#### Get Switch History
```http
GET /account/switch/history
Authorization: Bearer <token>
```

### Logout

#### Logout
```http
POST /auth/logout
Authorization: Bearer <token>
```

---

## Protected Admin Endpoints

### User Management

#### List Users
```http
GET /admin/users?search=query
Authorization: Bearer <admin_token>
```

#### Get User
```http
GET /admin/users/{id}
Authorization: Bearer <admin_token>
```

#### Create User
```http
POST /admin/users
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "User Name",
  "email": "user@example.com",
  "password": "Password123!",
  "role": "user",
  "status": "active"
}
```

#### Update User
```http
PUT /admin/users/{id}
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Updated Name",
  "email": "newemail@example.com",
  "status": "active"
}
```

#### Delete User
```http
DELETE /admin/users/{id}
Authorization: Bearer <admin_token>
```

#### Update User Status
```http
POST /admin/users/{id}/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "status": "active|pending|disabled|deleted"
}
```

#### Get User Sessions
```http
GET /admin/users/{id}/sessions
Authorization: Bearer <admin_token>
```

#### Revoke User Sessions
```http
DELETE /admin/users/{id}/sessions
Authorization: Bearer <admin_token>
```

### Admin Settings

#### Get Admin Settings
```http
GET /admin/settings
Authorization: Bearer <admin_token>
```

#### Update Admin Settings
```http
PUT /admin/settings
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "admin_verification_code": "NewCode123@"
}
```

### Custom CRUDs

#### Get All Templates
```http
GET /admin/cruds/templates
Authorization: Bearer <admin_token>
```

#### Get Template
```http
GET /admin/cruds/templates/{name}
Authorization: Bearer <admin_token>
```

#### Create Entity from Template
```http
POST /admin/cruds/templates/{name}/create
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "display_name": "Custom Name",
  "description": "Optional description"
}
```

#### Create Custom Entity
```http
POST /admin/cruds/entities
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "entity_name": "custom_entity",
  "display_name": "Custom Entity",
  "description": "Description",
  "schema": { ... }
}
```

#### List Entities
```http
GET /admin/cruds/entities?active_only=true
Authorization: Bearer <admin_token>
```

#### Get Entity
```http
GET /admin/cruds/entities/{id}
Authorization: Bearer <admin_token>
```

#### Update Entity
```http
PUT /admin/cruds/entities/{id}
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "display_name": "Updated Name",
  "schema": { ... }
}
```

#### Delete Entity
```http
DELETE /admin/cruds/entities/{id}
Authorization: Bearer <admin_token>
```

#### Create Data Entry
```http
POST /admin/cruds/entities/{entity_id}/data
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "field1": "value1",
  "field2": "value2"
}
```

#### List Data Entries
```http
GET /admin/cruds/entities/{entity_id}/data?limit=50&offset=0
Authorization: Bearer <admin_token>
```

#### Get Data Entry
```http
GET /admin/cruds/data/{id}
Authorization: Bearer <admin_token>
```

#### Update Data Entry
```http
PUT /admin/cruds/data/{id}
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "field1": "updated_value"
}
```

#### Delete Data Entry
```http
DELETE /admin/cruds/data/{id}
Authorization: Bearer <admin_token>
```

---

## Error Codes

| Code | Description |
|------|-------------|
| `INVALID_REQUEST` | Invalid request format or missing required fields |
| `UNAUTHORIZED` | Missing or invalid authentication token |
| `FORBIDDEN` | Insufficient permissions |
| `NOT_FOUND` | Resource not found |
| `CONFLICT` | Resource already exists |
| `BAD_REQUEST` | Invalid data or validation failed |
| `INTERNAL_ERROR` | Server error |

---

## CORS Configuration

The backend is configured to accept requests from any origin. CORS headers are automatically added to all responses.

**Allowed Methods:** GET, POST, PUT, DELETE, OPTIONS
**Allowed Headers:** Content-Type, Authorization
**Allowed Origins:** * (all origins)

---

## Rate Limiting

- **Rate Limit:** 100 requests per minute per IP/user
- **Response Headers:** Rate limit information included in response headers

---

## File Uploads

### Upload Image
```http
POST /files/upload/image
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <image_file>
```

### Upload Document
```http
POST /files/upload/document
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <document_file>
```

### Download File
```http
GET /files/download?file_id={id}
Authorization: Bearer <token>
```

### Delete File
```http
DELETE /files/delete?file_id={id}
Authorization: Bearer <token>
```

---

## Health Checks

### Health Check
```http
GET /health
```

### Readiness Check
```http
GET /health/ready
```

### Liveness Check
```http
GET /health/live
```

### Metrics
```http
GET /metrics
```

---

## Notes for Frontend Developers

1. **Token Management:** Store `access_token` and `refresh_token` securely
2. **Token Refresh:** Implement automatic token refresh before expiry
3. **Error Handling:** Always check `success` field in responses
4. **CORS:** Backend supports CORS from any origin
5. **Content-Type:** Always set `Content-Type: application/json` for JSON requests
6. **Authorization:** Include `Authorization: Bearer <token>` header for protected endpoints
7. **Pagination:** Use `limit` and `offset` query parameters for paginated endpoints
8. **Validation:** Backend validates all inputs - handle validation errors appropriately

---

## Frontend Independence

✅ **Backend is completely frontend-agnostic**
✅ **Works with React, Vue, Angular, Next.js, or any framework**
✅ **RESTful API - standard HTTP methods**
✅ **JSON-based communication**
✅ **CORS enabled for all origins**
✅ **No frontend dependencies**
✅ **Complete API documentation**
✅ **Standard error responses**

The backend can be used with **any frontend framework** without modification!

