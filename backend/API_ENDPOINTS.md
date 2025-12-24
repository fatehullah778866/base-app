# Base App API Endpoints Documentation

## Base URL
`http://localhost:8080/v1`

## Authentication
All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

---

## Public Endpoints

### Authentication
- `POST /auth/signup` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh access token
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password with token
- `POST /admin/login` - Admin login

---

## User Endpoints (Protected)

### Profile & Account
- `GET /users/me` - Get current user profile
- `PUT /users/me` - Update user profile
- `GET /users/me/export` - Export user data
- `POST /users/me/delete` - Request account deletion

### Settings (8 Categories)
- `GET /users/me/settings` - Get all settings
- `PUT /users/me/settings/profile` - Update profile settings (username, display name, bio, date of birth)
- `PUT /users/me/settings/security` - Update security settings (2FA, security questions)
- `PUT /users/me/settings/privacy` - Update privacy settings (visibility, messaging, search)
- `PUT /users/me/settings/notifications` - Update notification preferences
- `PUT /users/me/settings/preferences` - Update account preferences (language, timezone, theme, accessibility)
- `POST /users/me/settings/connected-accounts` - Add connected account (Google, Facebook, etc.)
- `DELETE /users/me/settings/connected-accounts` - Remove connected account
- `POST /users/me/settings/delete-account` - Schedule account deletion
- `POST /users/me/settings/deactivate` - Deactivate account
- `POST /users/me/settings/reactivate` - Reactivate account

### Theme
- `GET /users/me/settings/theme` - Get theme preferences
- `PUT /users/me/settings/theme` - Update theme preferences
- `POST /users/me/settings/theme/sync` - Sync theme with server

### Dashboard (CRUD)
- `POST /dashboard/items` - Create dashboard item
- `GET /dashboard/items` - List dashboard items (optional: ?status=active|archived)
- `GET /dashboard/items/{id}` - Get dashboard item by ID
- `PUT /dashboard/items/{id}` - Update dashboard item
- `DELETE /dashboard/items/{id}` - Delete dashboard item permanently
- `POST /dashboard/items/{id}/archive` - Archive dashboard item (soft delete)

### Notifications
- `GET /notifications` - Get notifications (?unread_only=true&limit=50)
- `GET /notifications/unread-count` - Get unread notification count
- `POST /notifications/read` - Mark notification as read
- `POST /notifications/read-all` - Mark all notifications as read
- `DELETE /notifications` - Delete notification

### Messaging
- `POST /messages` - Send message to user
- `GET /messages/conversations` - Get all conversations
- `GET /messages` - Get messages (?conversation_id={id}&limit=50)
- `POST /messages/read` - Mark message as read
- `GET /messages/unread-count` - Get unread message count

### Account Switching
- `POST /account/switch` - Switch account context
- `GET /account/switch/history` - Get account switch history

### Search
- `GET /search` - Search (?q=query&type=all|users|dashboard_items|messages&limit=20)

### Access Requests
- `POST /requests` - Create access request
- `GET /requests` - List user's access requests

---

## Admin Endpoints (Protected + Admin Role Required)

### User Management (CRUD)
- `GET /admin/users` - List all users (?search=query)
- `GET /admin/users/{id}` - Get user by ID
- `POST /admin/users` - Create new user
- `PUT /admin/users/{id}` - Update user
- `DELETE /admin/users/{id}` - Delete user
- `POST /admin/users/{id}/status` - Update user status (active/pending/disabled/deleted)
- `GET /admin/users/{id}/sessions` - Get user sessions
- `DELETE /admin/users/{id}/sessions` - Revoke all user sessions

### Admin Management
- `GET /admin/admins` - List all admins (?search=query)
- `POST /admin/admins` - Create new admin account

### Admin Settings
- `GET /admin/settings` - Get admin settings
- `PUT /admin/settings` - Update admin settings

### Custom CRUD Entities (Admin can create multiple CRUDs)
- `POST /admin/cruds/entities` - Create custom CRUD entity
- `GET /admin/cruds/entities` - List all CRUD entities (?active_only=true)
- `GET /admin/cruds/entities/{id}` - Get CRUD entity by ID
- `PUT /admin/cruds/entities/{id}` - Update CRUD entity
- `DELETE /admin/cruds/entities/{id}` - Delete CRUD entity

### Custom CRUD Data
- `POST /admin/cruds/entities/{id}/data` - Create data for entity
- `GET /admin/cruds/entities/{id}/data` - List data for entity (?limit=50&offset=0)
- `GET /admin/cruds/data/{id}` - Get data by ID
- `PUT /admin/cruds/data/{id}` - Update data
- `DELETE /admin/cruds/data/{id}` - Delete data

### Activity Logs
- `GET /admin/logs` - Get activity logs (?limit=200)

### Access Requests Management
- `GET /admin/requests` - List all access requests (?status=pending|approved|rejected)
- `POST /admin/requests/{id}/status` - Update request status

---

## Example Requests

### User Signup
```json
POST /v1/auth/signup
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890",
  "terms_accepted": true,
  "terms_version": "1.0"
}
```

### Create Dashboard Item
```json
POST /v1/dashboard/items
Authorization: Bearer <token>
{
  "title": "My Task",
  "description": "Complete this task",
  "category": "work",
  "priority": 1
}
```

### Create Custom CRUD Entity (Admin)
```json
POST /v1/admin/cruds/entities
Authorization: Bearer <admin_token>
{
  "entity_name": "products",
  "display_name": "Products",
  "description": "Product catalog",
  "schema": {
    "fields": [
      {"name": "name", "type": "string", "required": true},
      {"name": "price", "type": "number", "required": true},
      {"name": "description", "type": "string"}
    ]
  }
}
```

### Create Data for Custom Entity (Admin)
```json
POST /v1/admin/cruds/entities/{entity_id}/data
Authorization: Bearer <admin_token>
{
  "name": "Product Name",
  "price": 99.99,
  "description": "Product description"
}
```

---

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
    "message": "Error message"
  }
}
```

---

## Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

---

## Notes
- All timestamps are in RFC3339 format
- All UUIDs are in standard UUID format
- Pagination uses `limit` and `offset` query parameters
- Search uses `q` query parameter
- Filters use query parameters (e.g., `?status=active`)

