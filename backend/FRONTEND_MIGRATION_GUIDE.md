# Frontend Migration Guide

## Overview

This guide helps you migrate to a new frontend framework while keeping the backend unchanged. The backend is **100% frontend-agnostic** and works with any frontend.

## Quick Start

### Step 1: Understand the API

Read `API_SPECIFICATION.md` to understand all available endpoints.

### Step 2: Set Up Your New Frontend

Choose your framework:
- React
- Vue
- Angular
- Next.js
- Svelte
- Or any other framework

### Step 3: Create API Client

Create an API client for your framework:

#### React Example
```typescript
// api/client.ts
const API_BASE_URL = 'http://localhost:8080/v1';

export const apiClient = {
  async request(endpoint: string, options: RequestInit = {}) {
    const token = localStorage.getItem('access_token');
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
        ...options.headers,
      },
    });
    return response.json();
  },
  
  async login(email: string, password: string) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  },
};
```

#### Vue Example
```typescript
// api/client.ts
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/v1',
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
```

### Step 4: Implement Authentication

```typescript
// auth.ts
export async function login(email: string, password: string) {
  const response = await apiClient.login(email, password);
  if (response.success) {
    localStorage.setItem('access_token', response.data.session.token);
    localStorage.setItem('refresh_token', response.data.session.refresh_token);
    return response.data;
  }
  throw new Error(response.error.message);
}
```

### Step 5: Implement Token Refresh

```typescript
// auth.ts
export async function refreshToken() {
  const refreshToken = localStorage.getItem('refresh_token');
  const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });
  const data = await response.json();
  if (data.success) {
    localStorage.setItem('access_token', data.data.token);
    return data.data.token;
  }
  // Redirect to login
  window.location.href = '/login';
}
```

## API Endpoints Reference

### Authentication
- `POST /auth/signup` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh token
- `POST /auth/logout` - Logout
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password

### User Profile
- `GET /users/me` - Get current user
- `PUT /users/me` - Update profile
- `PUT /users/me/password` - Change password

### Settings
- `GET /users/me/settings` - Get all settings
- `PUT /users/me/settings/profile` - Update profile settings
- `PUT /users/me/settings/security` - Update security settings
- `PUT /users/me/settings/privacy` - Update privacy settings
- `PUT /users/me/settings/notifications` - Update notifications
- `PUT /users/me/settings/preferences` - Update preferences
- `POST /users/me/settings/connected-accounts` - Add connected account
- `DELETE /users/me/settings/connected-accounts` - Remove connected account

### Dashboard
- `GET /dashboard/items` - List items
- `POST /dashboard/items` - Create item
- `PUT /dashboard/items/{id}` - Update item
- `DELETE /dashboard/items/{id}` - Delete item

### Notifications
- `GET /notifications` - Get notifications
- `POST /notifications/read` - Mark as read
- `POST /notifications/read-all` - Mark all as read

### Messaging
- `GET /messages/conversations` - Get conversations
- `POST /messages` - Send message
- `GET /messages` - Get messages

### Admin
- `GET /admin/users` - List users
- `POST /admin/users` - Create user
- `PUT /admin/users/{id}` - Update user
- `DELETE /admin/users/{id}` - Delete user
- `GET /admin/cruds/templates` - Get CRUD templates
- `POST /admin/cruds/templates/{name}/create` - Create from template

## Response Format

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

## Error Handling

```typescript
try {
  const response = await apiClient.request('/endpoint');
  if (response.success) {
    // Handle success
  } else {
    // Handle error
    console.error(response.error.message);
  }
} catch (error) {
  // Handle network error
  console.error('Network error:', error);
}
```

## Token Management

### Store Tokens
```typescript
// After login
localStorage.setItem('access_token', token);
localStorage.setItem('refresh_token', refreshToken);
```

### Use Tokens
```typescript
// In API requests
headers: {
  'Authorization': `Bearer ${localStorage.getItem('access_token')}`
}
```

### Refresh Tokens
```typescript
// Before token expires
if (isTokenExpiringSoon()) {
  await refreshToken();
}
```

## Role-Based Routing

```typescript
// After login
const role = response.data.user.role;

if (role === 'admin') {
  router.push('/admin-dashboard');
} else {
  router.push('/dashboard');
}
```

## Testing

### Test Backend Independently
```bash
# Using cURL
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

### Test with Postman
1. Import API specification
2. Test all endpoints
3. No frontend needed

## Common Patterns

### Form Submission
```typescript
async function handleSubmit(formData) {
  try {
    const response = await apiClient.request('/endpoint', {
      method: 'POST',
      body: JSON.stringify(formData),
    });
    
    if (response.success) {
      // Success handling
    }
  } catch (error) {
    // Error handling
  }
}
```

### File Upload
```typescript
async function uploadFile(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch(`${API_BASE_URL}/files/upload/image`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: formData,
  });
  
  return response.json();
}
```

## Backend Configuration

No changes needed! Backend works with any frontend.

### CORS
Already configured to accept requests from any origin.

### API Base URL
```
http://localhost:8080/v1
```

## Migration Checklist

- [ ] Read API documentation
- [ ] Set up new frontend project
- [ ] Create API client
- [ ] Implement authentication
- [ ] Implement token refresh
- [ ] Create login/signup pages
- [ ] Create dashboard pages
- [ ] Create settings pages
- [ ] Implement role-based routing
- [ ] Test all endpoints
- [ ] Deploy frontend

## Support

- **API Documentation:** `API_SPECIFICATION.md`
- **Endpoint List:** `API_ENDPOINTS.md`
- **Backend Independence:** `BACKEND_INDEPENDENCE.md`

## Conclusion

✅ **Backend doesn't need changes**
✅ **Use any frontend framework**
✅ **Same API endpoints**
✅ **Same response format**
✅ **Complete documentation**

**Your backend is ready for any frontend!**

