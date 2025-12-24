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

### Integration Example
```javascript
// Works in any framework
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
    return data.data;
  }
  throw new Error(data.error.message);
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
- User signup/login
- Admin login
- Password reset
- Admin verification

### Protected User Endpoints
- Profile management
- Settings (8 categories)
- Dashboard CRUD
- Notifications
- Messaging
- Search
- File uploads

### Protected Admin Endpoints
- User management
- Admin settings
- Custom CRUDs
- CRUD templates

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
# Server
PORT=8080
FRONTEND_DIR=../frontend  # Optional - only for serving static files

# Database
DB_PATH=./app.db

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
```

### No Frontend Required
- Backend runs independently
- Can be tested without frontend
- API-first design
- Frontend is optional

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

