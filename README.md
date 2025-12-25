# Base App

A comprehensive full-stack application with user management, admin dashboard, messaging, notifications, advanced search, CRUD templates, and extensive settings management.

## 🚀 Features

### Authentication & User Management
- ✅ User registration and login
- ✅ Admin login with verification code
- ✅ Password reset (forgot/reset password)
- ✅ JWT-based authentication
- ✅ Session management
- ✅ Profile management with file upload
- ✅ Account deactivation/reactivation

### User Dashboard
- ✅ Dashboard items management (CRUD)
- ✅ Custom CRUD entities creation
- ✅ CRUD templates browsing and usage
- ✅ Real-time notifications with badge counts
- ✅ Messaging system with conversations
- ✅ Advanced search with map integration
- ✅ Search Near Me (geolocation-based)
- ✅ Profile card display
- ✅ Comprehensive settings (8 categories)

### Admin Dashboard
- ✅ User management (View, Edit, Delete, Toggle Status)
- ✅ CRUD templates management
- ✅ Custom CRUDs management
- ✅ Admin settings (verification code management)
- ✅ All user dashboard features
- ✅ Enhanced user cards with action buttons

### Search System
- ✅ Global search across all entities
- ✅ Live search with dropdown results
- ✅ Advanced search modal with filters
- ✅ Map-based location search (Leaflet.js)
- ✅ Search Near Me with radius selection (1-100km)
- ✅ Reverse geocoding (coordinates to address)
- ✅ Location-based filtering
- ✅ Search history

### Messaging System
- ✅ One-on-one conversations
- ✅ User search within messages
- ✅ Real-time message polling
- ✅ Unread message count badges
- ✅ Message threading
- ✅ Settings control for polling

### Notification System
- ✅ Real-time notifications
- ✅ Multiple notification types
- ✅ Unread count badges
- ✅ Mark as read (individual/bulk)
- ✅ Real-time polling
- ✅ Settings control

### Settings (8 Categories)
1. **Profile Settings** - Name, email, phone, bio, profile picture
2. **Security Settings** - Password change, 2FA, active sessions
3. **Privacy Settings** - Visibility controls, messaging permissions
4. **Notification Settings** - Email, SMS, push notifications
5. **Account Preferences** - Language, timezone, theme, accessibility
6. **Connected Accounts** - Google, Facebook, Apple integration
7. **Data & Account Control** - Export data, delete/deactivate account
8. **Help & Support** - Support resources

### CRUD Templates System
- ✅ Dynamic template creation (Admin)
- ✅ Template schema builder (JSON-based)
- ✅ Field management (Add/Remove fields)
- ✅ Template categories and icons
- ✅ Active/Inactive template status
- ✅ User access to active templates
- ✅ One-click entity creation from templates

### File Management
- ✅ Image upload (JPG, PNG, GIF)
- ✅ File size validation (5MB max)
- ✅ Profile picture upload with preview

## 🛠️ Technology Stack

### Backend
- **Go 1.21+** - Programming language
- **SQLite** - Database (modernc.org/sqlite)
- **Gorilla Mux** - HTTP router
- **JWT** - Authentication
- **Zap Logger** - Structured logging
- **bcrypt** - Password hashing

### Frontend
- **HTML5/CSS3/JavaScript** - Core technologies
- **Leaflet.js** - Interactive maps
- **OpenStreetMap** - Map tiles
- **Nominatim API** - Geocoding

## 📋 Prerequisites

- Go 1.21 or higher
- Git
- Modern web browser

## 🔧 Installation

### 1. Clone the Repository
```bash
git clone https://github.com/kompass-tech/base-app.git
cd base-app
```

### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Set environment variables (create .env file)
PORT=8080
ENV=development
DB_PATH=./app.db
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Run database migrations
go run cmd/server/main.go migrate

# Start the server
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

### 3. Frontend Setup

The frontend is already included and will be served by the backend. No additional setup required.

## 🚀 Quick Start

1. **Start the Backend**
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **Access the Application**
   - User Login: `http://localhost:8080/`
   - User Dashboard: `http://localhost:8080/dashboard`
   - Admin Dashboard: `http://localhost:8080/admin-dashboard`
   - Settings: `http://localhost:8080/settings`

3. **Create Admin Account**
   - Click "Create Admin" on login page
   - Enter verification code: `Kompasstech2025@`
   - Fill in admin details and create account

## 📚 API Documentation

### Base URL
```
http://localhost:8080/v1
```

### Authentication
All protected endpoints require:
```
Authorization: Bearer <access_token>
```

### Key Endpoints

#### Public
- `POST /v1/auth/signup` - User registration
- `POST /v1/auth/login` - User login
- `POST /v1/auth/forgot-password` - Request password reset
- `POST /v1/auth/reset-password` - Reset password
- `POST /v1/admin/login` - Admin login
- `POST /v1/admin/verify-code` - Verify admin code
- `POST /v1/admin/create` - Create admin account

#### User Endpoints
- `GET /v1/users/me` - Get current user
- `PUT /v1/users/me` - Update profile
- `GET /v1/users/me/settings` - Get all settings
- `PUT /v1/users/me/settings/*` - Update settings
- `GET /v1/dashboard/items` - List dashboard items
- `POST /v1/dashboard/items` - Create dashboard item
- `GET /v1/notifications` - Get notifications
- `GET /v1/messages/conversations` - Get conversations
- `POST /v1/messages` - Send message
- `POST /v1/search` - Advanced search
- `POST /v1/files/upload/image` - Upload image

#### Admin Endpoints
- `GET /v1/admin/users` - List all users
- `POST /v1/admin/users` - Create user
- `PUT /v1/admin/users/{id}` - Update user
- `DELETE /v1/admin/users/{id}` - Delete user
- `GET /v1/admin/cruds/templates` - Get templates
- `POST /v1/admin/cruds/templates` - Create template
- `GET /v1/admin/cruds/entities` - List CRUD entities
- `GET /v1/admin/settings` - Get admin settings
- `PUT /v1/admin/settings` - Update admin settings

For complete API documentation, see [backend/docs/BASE_APP_FEATURES.md](backend/docs/BASE_APP_FEATURES.md)

## 📁 Project Structure

```
BASEAPP/
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── handlers/        # HTTP handlers
│   │   ├── services/        # Business logic
│   │   ├── repositories/    # Data access
│   │   ├── models/          # Domain models
│   │   ├── middleware/      # HTTP middleware
│   │   └── database/        # Database connection
│   ├── migrations/          # Database migrations
│   ├── docs/                # Documentation
│   └── pkg/                 # Shared packages
├── frontend/
│   ├── index.html           # Login/Signup page
│   ├── dashboard.html       # User dashboard
│   ├── admin-dashboard.html # Admin dashboard
│   ├── settings.html        # Settings page
│   ├── css/                 # Stylesheets
│   └── js/                  # JavaScript files
└── README.md                # This file
```

## 🔐 Default Credentials

### Admin Verification Code
```
Kompasstech2025@
```

**Note:** Admins can change this code from the admin dashboard settings.

## 🌟 Key Features Details

### Advanced Search
- **Live Search**: Real-time results as you type
- **Map Search**: Click on map to set location
- **Search Near Me**: Uses browser geolocation with configurable radius
- **Filters**: Location, date range, category, status, entity type
- **Distance Calculation**: Haversine formula for accurate results

### Messaging
- Search for users within messages modal
- One-on-one conversations
- Real-time message updates
- Unread count badges
- Message threading

### Notifications
- Real-time notification updates
- Multiple notification types
- Unread count badges
- Mark as read functionality
- Settings control

### CRUD Templates
- Admins create templates with custom schemas
- Users browse and use active templates
- Dynamic field management
- JSON-based schema definition
- One-click entity creation

## 📖 Documentation

- **[BASE_APP_FEATURES.md](backend/docs/BASE_APP_FEATURES.md)** - Complete feature list
- **[BACKEND_INDEPENDENCE.md](backend/docs/BACKEND_INDEPENDENCE.md)** - Backend independence guide
- **[CODE_QUALITY.md](backend/docs/CODE_QUALITY.md)** - Code quality standards

## 🔒 Security Features

- JWT authentication with refresh tokens
- Password hashing with bcrypt
- Rate limiting
- CORS protection
- Security headers
- Input validation
- SQL injection protection
- XSS protection
- CSRF protection

## 🧪 Testing

### Health Check
```bash
curl http://localhost:8080/health
```

### Test API
```bash
# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

## 🚀 Deployment

### Production Build
```bash
cd backend
go build -o server cmd/server/main.go
./server
```

### Environment Variables (Production)
```bash
ENV=production
JWT_SECRET=<strong-secret-key>
PORT=8080
DB_PATH=/data/app.db
```

## 📝 Development

### Running in Development
```bash
cd backend
go run cmd/server/main.go
```

### Database Migrations
```bash
# Run migrations
go run cmd/server/main.go migrate
```

### Project Structure
The backend follows Clean Architecture principles:
- **Handlers**: HTTP request/response handling
- **Services**: Business logic
- **Repositories**: Data access
- **Models**: Domain entities

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Commit with descriptive messages
5. Push to your branch
6. Create a Pull Request

## 📄 License

This project is proprietary software.

## 🆘 Support

For issues or questions:
1. Check the documentation in `backend/docs/`
2. Review error logs
3. Check GitHub issues
4. Contact support

## 🎯 Roadmap

- [ ] WebSocket support for real-time updates
- [ ] Email notifications
- [ ] SMS notifications
- [ ] Mobile app support
- [ ] Enhanced analytics
- [ ] Multi-language support
- [ ] Advanced reporting

## 📞 Contact

- **Repository**: https://github.com/kompass-tech/base-app
- **Organization**: Kompass Tech

---

**Built with ❤️ using Go and modern web technologies**

