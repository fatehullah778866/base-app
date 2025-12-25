# Backend Independence Documentation

## Overview

The Base App backend is designed to be **100% frontend-independent**. It works with any frontend framework (React, Vue, Angular, Next.js, Svelte, etc.) without modification.

## Core Principles

### 1. API-First Design
- Pure REST API
- JSON-based communication
- Stateless requests
- Standard HTTP methods
- Resource-based URLs

### 2. No Frontend Dependencies
- ✅ No frontend code in backend
- ✅ No frontend-specific logic
- ✅ No frontend build artifacts
- ✅ Pure API endpoints
- ✅ Frontend serving is optional

### 3. Standard Communication
- Consistent JSON response format
- Standard error responses
- RESTful conventions
- HTTP status codes

## Architecture

### API-Only Mode (Default)
The backend runs in **API-only mode** by default:
- No frontend serving
- Pure REST API
- Frontend can be served separately
- Better for production

### Optional Frontend Serving
Frontend serving can be enabled via environment variable:
```bash
FRONTEND_DIR=../frontend
```

This is optional and not required for API functionality.

## CORS Configuration

### Universal Access
- ✅ Accepts requests from ANY origin
- ✅ CORS enabled for all origins
- ✅ Works with any frontend
- ✅ Standard HTTP methods supported

### Headers
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## Frontend Compatibility

### Supported Frameworks
✅ **React** - Use fetch/axios  
✅ **Vue** - Use axios/vue-resource  
✅ **Angular** - Use HttpClient  
✅ **Next.js** - Use fetch/axios  
✅ **Svelte** - Use fetch/axios  
✅ **Vanilla JS** - Use fetch  
✅ **Mobile Apps** - Use HTTP clients  
✅ **Desktop Apps** - Use HTTP libraries  

### Integration Examples

#### Vanilla JavaScript
```javascript
const API_BASE_URL = 'http://localhost:8080/v1';

async function login(email, password) {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });
  
  const data = await response.json();
  if (data.success) {
    localStorage.setItem('access_token', data.data.session.token);
    localStorage.setItem('user', JSON.stringify(data.data.user));
    return data.data;
  }
  throw new Error(data.error.message);
}

async function getUsers(token) {
  const response = await fetch(`${API_BASE_URL}/admin/users`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
  
  const data = await response.json();
  return data.data || [];
}
```

#### React Example
```javascript
import { useState, useEffect } from 'react';

const API_BASE_URL = 'http://localhost:8080/v1';

function useUsers() {
  const [users, setUsers] = useState([]);
  const token = localStorage.getItem('access_token');
  
  useEffect(() => {
    async function fetchUsers() {
      const response = await fetch(`${API_BASE_URL}/admin/users`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });
      const data = await response.json();
      setUsers(data.data || []);
    }
    fetchUsers();
  }, [token]);
  
  return users;
}
```

#### Vue.js Example
```javascript
import { ref, onMounted } from 'vue';

const API_BASE_URL = 'http://localhost:8080/v1';

export function useUsers() {
  const users = ref([]);
  const token = localStorage.getItem('access_token');
  
  onMounted(async () => {
    const response = await fetch(`${API_BASE_URL}/admin/users`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    const data = await response.json();
    users.value = data.data || [];
  });
  
  return { users };
}
```

#### Advanced Search Example
```javascript
async function searchNearMe(lat, lng, radius = 5) {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`${API_BASE_URL}/search`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: '',
      type: 'all',
      latitude: lat,
      longitude: lng,
      radius: radius, // in kilometers
      limit: 50
    }),
  });
  
  const data = await response.json();
  return data.data?.data?.results || [];
}
```

## API Endpoints

### Base URL
```
http://localhost:8080/v1
```

### Authentication
All protected endpoints require:
```
Authorization: Bearer <access_token>
```

### Response Format
**Success:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Optional message"
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

## Complete API Coverage

### Public Endpoints
- `POST /v1/auth/signup` - User registration
- `POST /v1/auth/login` - User login
- `POST /v1/auth/refresh` - Refresh access token
- `POST /v1/auth/forgot-password` - Request password reset
- `POST /v1/auth/reset-password` - Reset password with token
- `POST /v1/admin/login` - Admin login
- `POST /v1/admin/verify-code` - Verify admin verification code
- `POST /v1/admin/create` - Create admin account (requires verification)
- `GET /v1/cruds/templates/active` - Get active CRUD templates (for users)

### Protected User Endpoints

#### Profile & Settings
- `GET /v1/users/me` - Get current user
- `PUT /v1/users/me` - Update profile
- `PUT /v1/users/me/password` - Change password
- `GET /v1/users/me/export` - Export user data
- `POST /v1/users/me/delete` - Request account deletion
- `GET /v1/users/me/settings` - Get all settings
- `PUT /v1/users/me/settings/*` - Update specific settings category

#### Dashboard & CRUDs
- `GET /v1/dashboard/items` - List dashboard items
- `POST /v1/dashboard/items` - Create dashboard item
- `PUT /v1/dashboard/items/{id}` - Update dashboard item
- `DELETE /v1/dashboard/items/{id}` - Delete dashboard item
- `GET /v1/cruds/entities` - List user's CRUD entities
- `POST /v1/cruds/entities` - Create CRUD entity
- `GET /v1/cruds/templates/{name}` - Get template details
- `POST /v1/cruds/templates/{name}/create` - Create entity from template

#### Notifications & Messaging
- `GET /v1/notifications` - Get notifications
- `GET /v1/notifications/unread-count` - Get unread count
- `POST /v1/notifications/read` - Mark as read
- `GET /v1/messages/conversations` - Get conversations
- `POST /v1/messages` - Send message
- `POST /v1/messages/{id}/read` - Mark message as read

#### Search & Files
- `POST /v1/search` - Advanced search (supports location, radius, filters)
- `GET /v1/search/history` - Get search history
- `POST /v1/files/upload/image` - Upload image file

### Protected Admin Endpoints

#### User Management
- `GET /v1/admin/users` - List all users
- `GET /v1/admin/users/{id}` - Get user details
- `POST /v1/admin/users` - Create user
- `PUT /v1/admin/users/{id}` - Update user
- `DELETE /v1/admin/users/{id}` - Delete user
- `POST /v1/admin/users/{id}/status` - Update user status

#### Admin Settings
- `GET /v1/admin/settings` - Get admin settings
- `PUT /v1/admin/settings` - Update admin settings (verification code, etc.)

#### CRUD Management
- `GET /v1/admin/cruds/entities` - List all CRUD entities
- `POST /v1/admin/cruds/entities` - Create CRUD entity
- `GET /v1/admin/cruds/templates` - Get all templates
- `POST /v1/admin/cruds/templates` - Create template
- `PUT /v1/admin/cruds/templates/id/{id}` - Update template
- `DELETE /v1/admin/cruds/templates/id/{id}` - Delete template

## Migration Guide

### Changing Frontend
When you change your frontend:

1. **No Backend Changes Needed** ✅
2. **Update API Client** - Use your framework's HTTP client
3. **Update Token Storage** - Use framework's storage solution
4. **Update Routing** - Handle redirects in your framework
5. **Test Endpoints** - Use same API endpoints

### Example: React Migration
**Old (HTML/JS):**
```javascript
fetch('/v1/auth/login', { ... })
```

**New (React):**
```javascript
// Same endpoint, just use React hooks
const login = async (email, password) => {
  const response = await fetch('/v1/auth/login', { ... });
  // Same response handling
}
```

**Backend:** No changes needed! ✅

## Testing Backend Independently

### Using cURL
```bash
# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Get user profile
curl -X GET http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer <token>"
```

### Using Postman/Insomnia
- Import API specification
- Test all endpoints
- No frontend needed

### Using API Clients
- Any HTTP client works
- Standard REST API
- JSON responses

## Configuration

### Environment Variables
```bash
# Server Configuration
PORT=8080                    # Server port (default: 8080)
ENV=development              # Environment: development or production
FRONTEND_DIR=../frontend     # Optional - only for serving static files

# Database Configuration
DB_PATH=./app.db             # SQLite database file path

# JWT Configuration
JWT_SECRET=your-secret-key   # JWT signing secret (REQUIRED)
JWT_ACCESS_EXPIRY=15m        # Access token expiry (default: 15 minutes)
JWT_REFRESH_EXPIRY=7d        # Refresh token expiry (default: 7 days)

# Email Configuration (Optional)
SMTP_HOST=smtp.gmail.com     # SMTP server host
SMTP_PORT=587                # SMTP server port
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourapp.com
```

### Response Format Standards

#### Success Response
```json
{
  "success": true,
  "data": {
    // Response data here
  },
  "message": "Optional success message"
}
```

#### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message"
  }
}
```

#### Paginated Response
```json
{
  "success": true,
  "data": {
    "results": [...],
    "count": 100,
    "limit": 20,
    "offset": 0
  }
}
```

### No Frontend Required
- ✅ Backend runs independently
- ✅ Can be tested without frontend
- ✅ API-first design
- ✅ Frontend is optional
- ✅ Works with any HTTP client
- ✅ Standard REST API conventions

## Benefits

### 1. Flexibility
- Use any frontend framework
- Change frontend anytime
- No backend modifications needed

### 2. Scalability
- Frontend and backend can scale independently
- Deploy separately
- Use CDN for frontend

### 3. Development
- Frontend and backend teams work independently
- API contracts define integration
- Easy testing

### 4. Production
- API-only mode for microservices
- Frontend on CDN
- Better performance
- Easier deployment

## Conclusion

✅ **Backend is 100% independent**  
✅ **Works with any frontend**  
✅ **No frontend dependencies**  
✅ **Complete API documentation**  
✅ **Production ready**  

**You can change your frontend anytime without touching the backend!**

