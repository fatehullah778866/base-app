# Base App Backend

A modern, production-ready Go backend API that is **100% frontend-independent**. Works with any frontend framework (React, Vue, Angular, Next.js, etc.) without modification.

## 🚀 Quick Start

```bash
cd backend
go run ./cmd/server/main.go
```

Backend runs on: `http://localhost:8080`

## ✨ Features

### ✅ Complete Authentication System
- User signup/login
- Admin authentication
- Token refresh
- Password reset
- Session management

### ✅ User Management
- Profile management
- Settings (8 categories)
- Account control
- Data export

### ✅ Dashboard System
- CRUD operations
- Item management

### ✅ Notification System
- Real-time notifications
- Read/unread tracking

### ✅ Messaging System
- Send/receive messages
- Conversations

### ✅ Search System
- Global search
- Type filtering

### ✅ Admin Features
- User management
- Admin settings
- Custom CRUDs with templates

### ✅ Modern CRUD System
- Pre-built templates (Portfolio, Visa, Products, Blog, Events, Contacts)
- Custom entity creation
- Schema validation

## 📚 Documentation

- **[API_SPECIFICATION.md](./API_SPECIFICATION.md)** - Complete API reference with examples
- **[API_ENDPOINTS.md](./API_ENDPOINTS.md)** - List of all endpoints
- **[BACKEND_INDEPENDENCE.md](./BACKEND_INDEPENDENCE.md)** - Backend independence guide
- **[FRONTEND_MIGRATION_GUIDE.md](./FRONTEND_MIGRATION_GUIDE.md)** - Guide for changing frontend
- **[CRUD_SYSTEM_GUIDE.md](./CRUD_SYSTEM_GUIDE.md)** - CRUD system documentation
- **[docs/CRUD_TEMPLATES.md](./docs/CRUD_TEMPLATES.md)** - Template documentation

## 🔌 API Base URL

```
http://localhost:8080/v1
```

## 🔒 Authentication

All protected endpoints require:
```
Authorization: Bearer <access_token>
```

## 🌐 CORS Configuration

✅ **Backend accepts requests from ANY origin**
✅ **CORS enabled for all origins**
✅ **Works with any frontend**

## 🎯 Frontend Independence

✅ **100% Frontend-Agnostic**
- No frontend dependencies
- Works with React, Vue, Angular, Next.js, or any framework
- RESTful API - standard HTTP methods
- JSON-based communication
- Complete API documentation

✅ **Change Frontend Anytime**
- Backend doesn't need changes
- Same API endpoints
- Same response format
- Works with any HTTP client

## 🏗️ Architecture

```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── repositories/    # Data access
│   ├── models/          # Data models
│   └── middleware/      # HTTP middleware
├── pkg/                 # Shared packages
├── migrations/          # Database migrations
└── docs/                # Documentation
```

## 🔧 Configuration

### Environment Variables

```bash
# Server
PORT=8080

# Database
DB_PATH=./app.db

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Optional: Frontend serving (leave empty for API-only mode)
FRONTEND_DIR=../frontend
```

### API-Only Mode (Recommended)

By default, backend runs in **API-only mode**:
- No frontend serving
- Pure REST API
- Frontend can be served separately
- Better for production

To enable frontend serving:
```bash
export FRONTEND_DIR=../frontend
go run ./cmd/server/main.go
```

## 📡 API Endpoints

### Public Endpoints
- `POST /auth/signup` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh token
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/reset-password` - Reset password
- `POST /admin/login` - Admin login
- `POST /admin/verify-code` - Verify admin code
- `POST /admin/create` - Create admin account

### Protected User Endpoints
- `GET /users/me` - Get current user
- `PUT /users/me` - Update profile
- `GET /users/me/settings` - Get settings
- `PUT /users/me/settings/*` - Update settings
- `GET /dashboard/items` - List dashboard items
- `POST /dashboard/items` - Create item
- `GET /notifications` - Get notifications
- `POST /messages` - Send message
- `GET /search` - Search

### Protected Admin Endpoints
- `GET /admin/users` - List users
- `POST /admin/users` - Create user
- `GET /admin/cruds/templates` - Get CRUD templates
- `POST /admin/cruds/templates/{name}/create` - Create from template

See [API_SPECIFICATION.md](./API_SPECIFICATION.md) for complete documentation.

## 🧪 Testing

### Test with cURL
```bash
# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Get profile
curl -X GET http://localhost:8080/v1/users/me \
  -H "Authorization: Bearer <token>"
```

### Test with Postman
1. Import API specification
2. Test all endpoints
3. No frontend needed

## 🔐 Security Features

- ✅ JWT authentication
- ✅ Password hashing (bcrypt)
- ✅ CORS protection
- ✅ Rate limiting
- ✅ Security headers
- ✅ Input validation
- ✅ SQL injection protection
- ✅ XSS protection

## 📊 Health Checks

- `GET /health` - Health check
- `GET /health/ready` - Readiness check
- `GET /health/live` - Liveness check
- `GET /metrics` - Prometheus metrics

## 🎨 CRUD Templates

Pre-built templates for common use cases:
- **Portfolio** - Project showcase
- **Visa** - Visa management
- **Products** - E-commerce catalog
- **Blog Posts** - Content management
- **Events** - Event management
- **Contacts** - CRM system

See [CRUD_SYSTEM_GUIDE.md](./CRUD_SYSTEM_GUIDE.md) for details.

## 🚀 Production Deployment

### Recommended Setup
1. **Backend:** Deploy as API-only service
2. **Frontend:** Deploy separately (CDN, separate server, etc.)
3. **Database:** Use production SQLite or migrate to PostgreSQL/MySQL
4. **Environment:** Set proper environment variables

### Docker (Optional)
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server/main.go
CMD ["./server"]
```

## 📝 License

[Your License Here]

## 🤝 Support

- **API Documentation:** See `API_SPECIFICATION.md`
- **Migration Guide:** See `FRONTEND_MIGRATION_GUIDE.md`
- **Backend Independence:** See `BACKEND_INDEPENDENCE.md`

---

**✅ Backend is 100% independent and ready for any frontend!**
