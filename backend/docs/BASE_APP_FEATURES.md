# Base App Features Documentation

## Overview

This document provides a comprehensive list of all features and functionality available in the Base App.

## Core Features

### 1. Authentication System
- ✅ User Registration (Signup)
- ✅ User Login
- ✅ Admin Login
- ✅ Admin Account Creation (with verification code)
- ✅ Password Reset (Forgot/Reset Password)
- ✅ Token Refresh
- ✅ Session Management
- ✅ Logout
- ✅ Logout from All Devices

### 2. User Management
- ✅ User Profile Management
- ✅ Profile Picture/Avatar Upload (File upload with preview)
- ✅ Profile Picture Formats (JPG, PNG, GIF)
- ✅ File Size Validation (5MB max)
- ✅ Profile Card Display (Modal with user info)
- ✅ User Data Export
- ✅ Account Deletion Request
- ✅ Account Deactivation/Reactivation
- ✅ User Status Management (Active, Disabled, Pending)
- ✅ User Role Management (Admin, User)
- ✅ User Session Management
- ✅ Last Login Tracking

### 3. Settings System (8 Categories)

#### 3.1 Profile Settings
- ✅ Full Name / Display Name
- ✅ Username
- ✅ Profile Picture / Avatar
- ✅ Email Address
- ✅ Phone Number
- ✅ Bio / About Me
- ✅ Date of Birth

#### 3.2 Security Settings
- ✅ Change Password
- ✅ Forgot / Reset Password
- ✅ Password Rules Display
- ✅ Two-Factor Authentication (2FA) On/Off
- ✅ Security Questions
- ✅ Active Sessions / Logged-in Devices
- ✅ Log Out from All Devices

#### 3.3 Privacy Settings
- ✅ Profile Visibility (public / private)
- ✅ Email Visibility Control
- ✅ Phone Visibility Control
- ✅ Who Can Message
- ✅ Search Visibility
- ✅ Data Sharing Preferences

#### 3.4 Notification Settings
- ✅ Email Notifications On/Off
- ✅ SMS Notifications On/Off
- ✅ Push Notifications On/Off
- ✅ Message Notifications
- ✅ Alert Notifications
- ✅ Promotion Notifications
- ✅ Security Notifications

#### 3.5 Account Preferences
- ✅ Language Selection
- ✅ Time Zone
- ✅ Theme (Light / Dark / Auto)
- ✅ Font Size
- ✅ High Contrast Mode
- ✅ Reduced Motion
- ✅ Screen Reader Support

#### 3.6 Connected Accounts
- ✅ Google Login Integration
- ✅ Facebook Login Integration
- ✅ Apple Login Integration
- ✅ Add Connected Account
- ✅ Remove Connected Account
- ✅ View Connected Accounts

#### 3.7 Data & Account Control
- ✅ Download Personal Data
- ✅ Delete Account
- ✅ Deactivate Account (Temporary)
- ✅ Reactivate Account
- ✅ Account Deletion Scheduling

#### 3.8 Help & Support
- ✅ Change Email Help
- ✅ Password Help
- ✅ Contact Support
- ✅ Report a Problem

### 4. Dashboard System
- ✅ Create Dashboard Items
- ✅ Read Dashboard Items
- ✅ Update Dashboard Items
- ✅ Delete Dashboard Items
- ✅ Archive Dashboard Items
- ✅ Item Categories
- ✅ Item Status Management
- ✅ Item Priority
- ✅ Item Metadata (JSON)

### 5. Notification System
- ✅ Real-time Notifications
- ✅ Notification Types (Message, Alert, Promotion, Security, System)
- ✅ Read/Unread Status
- ✅ Mark as Read
- ✅ Mark All as Read
- ✅ Delete Notifications
- ✅ Unread Count Badge (Real-time updates)
- ✅ Notification History
- ✅ Real-time Notification Polling
- ✅ Notification Settings Control (Enable/Disable polling)
- ✅ Notification Title and Message
- ✅ Notification Timestamps

### 6. Messaging System
- ✅ Send Messages
- ✅ Receive Messages
- ✅ Conversations (One-on-one)
- ✅ Message Threads
- ✅ Mark Messages as Read
- ✅ Archive Messages
- ✅ Unread Message Count (Real-time badge)
- ✅ Message Search
- ✅ User Search within Messages (Find users to message)
- ✅ Conversation List with Latest Message
- ✅ Real-time Message Polling
- ✅ Message Settings Control (Enable/Disable polling)
- ✅ Message Subject and Content
- ✅ Message Timestamps

### 7. Search System
- ✅ Global Search (All Entities)
- ✅ Dashboard Items Search
- ✅ Messages Search
- ✅ Users Search
- ✅ Notifications Search
- ✅ Custom CRUD Search
- ✅ Location-based Search
- ✅ Map-based Search (Interactive Map)
- ✅ Search Near Me (Geolocation-based with radius)
- ✅ Advanced Filters:
  - Date Range
  - Location (Country, City)
  - Latitude/Longitude coordinates
  - Search Radius (1km - 100km)
  - Category
  - Status
  - Entity Type
- ✅ Live Search (Real-time dropdown results)
- ✅ Search History
- ✅ Full-Text Search (FTS5)
- ✅ Reverse Geocoding (Coordinates to Address)

### 8. Admin Features

#### 8.1 User Management
- ✅ List Users
- ✅ Get User Details
- ✅ Create User
- ✅ Update User
- ✅ Delete User
- ✅ Update User Status
- ✅ View User Sessions
- ✅ Revoke User Sessions

#### 8.2 Admin Settings
- ✅ Get Admin Settings
- ✅ Update Admin Settings
- ✅ Admin Verification Code Management
- ✅ Change Verification Code (With confirmation)
- ✅ Default Verification Code: `Kompasstech2025@`
- ✅ Verification Code Validation
- ✅ Settings Persistence in Database

#### 8.3 Custom CRUD System
- ✅ Create Custom CRUD Entities
- ✅ List CRUD Entities
- ✅ Get CRUD Entity
- ✅ Update CRUD Entity
- ✅ Delete CRUD Entity
- ✅ Create CRUD Data
- ✅ List CRUD Data
- ✅ Update CRUD Data
- ✅ Delete CRUD Data
- ✅ Schema Validation

#### 8.4 CRUD Templates
- ✅ Dynamic Template Creation (Admin)
- ✅ Template Schema Builder (JSON-based)
- ✅ Template Field Management (Add/Remove fields)
- ✅ Template Categories
- ✅ Template Icons
- ✅ Active/Inactive Template Status
- ✅ System Templates (Pre-built)
- ✅ Custom Templates (Admin-created)
- ✅ Create Entity from Template
- ✅ List Available Templates (Active only for users)
- ✅ Template Management (Create, Read, Update, Delete)
- ✅ Template Display Name and Description

### 9. File Management
- ✅ Image Upload
- ✅ Document Upload
- ✅ File Download
- ✅ File Deletion
- ✅ File Size Limits
- ✅ File Type Validation

### 10. Account Switching
- ✅ Switch Between Accounts
- ✅ Switch History
- ✅ Account Context Management

## Technical Features

### Security
- ✅ JWT Authentication
- ✅ Refresh Tokens
- ✅ Password Hashing (bcrypt)
- ✅ Rate Limiting
- ✅ CORS Protection
- ✅ Security Headers
- ✅ Input Validation
- ✅ SQL Injection Protection
- ✅ XSS Protection
- ✅ CSRF Protection

### Performance
- ✅ Database Indexing
- ✅ Efficient Queries
- ✅ Connection Pooling
- ✅ Caching Support (In-Memory)
- ✅ Full-Text Search (FTS5)
- ✅ Pagination Support

### Monitoring
- ✅ Health Checks
- ✅ Readiness Checks
- ✅ Liveness Checks
- ✅ Prometheus Metrics
- ✅ Request Logging
- ✅ Error Logging

### Database
- ✅ SQLite Database
- ✅ Database Migrations
- ✅ Schema Evolution
- ✅ Foreign Key Constraints
- ✅ Soft Deletes
- ✅ Full-Text Search Indexes

## API Endpoints Summary

### Public Endpoints
- `POST /v1/auth/signup` - User registration
- `POST /v1/auth/login` - User login
- `POST /v1/auth/refresh` - Refresh token
- `POST /v1/auth/forgot-password` - Request password reset
- `POST /v1/auth/reset-password` - Reset password with token
- `POST /v1/admin/login` - Admin login
- `POST /v1/admin/verify-code` - Verify admin verification code
- `POST /v1/admin/create` - Create admin account (with verification)

### Protected User Endpoints

#### User Profile
- `GET /v1/users/me` - Get current user profile
- `PUT /v1/users/me` - Update user profile
- `PUT /v1/users/me/password` - Change password
- `GET /v1/users/me/export` - Export user data
- `POST /v1/users/me/delete` - Request account deletion

#### Settings
- `GET /v1/users/me/settings` - Get all settings
- `PUT /v1/users/me/settings/profile` - Update profile settings
- `PUT /v1/users/me/settings/security` - Update security settings
- `PUT /v1/users/me/settings/privacy` - Update privacy settings
- `PUT /v1/users/me/settings/notifications` - Update notification settings
- `PUT /v1/users/me/settings/preferences` - Update account preferences
- `GET /v1/users/me/settings/sessions` - Get active sessions
- `POST /v1/users/me/settings/sessions/logout-all` - Logout all devices
- `POST /v1/users/me/settings/connected-accounts` - Add connected account
- `DELETE /v1/users/me/settings/connected-accounts` - Remove connected account
- `POST /v1/users/me/settings/account/deactivate` - Deactivate account
- `POST /v1/users/me/settings/account/reactivate` - Reactivate account
- `POST /v1/users/me/settings/account/delete` - Request account deletion

#### Dashboard
- `GET /v1/dashboard/items` - List dashboard items
- `POST /v1/dashboard/items` - Create dashboard item
- `GET /v1/dashboard/items/{id}` - Get dashboard item
- `PUT /v1/dashboard/items/{id}` - Update dashboard item
- `DELETE /v1/dashboard/items/{id}` - Delete dashboard item
- `POST /v1/dashboard/items/{id}/archive` - Archive dashboard item

#### Notifications
- `GET /v1/notifications` - Get notifications (with pagination)
- `GET /v1/notifications/unread-count` - Get unread notification count
- `POST /v1/notifications/read` - Mark notification as read
- `POST /v1/notifications/read-all` - Mark all notifications as read

#### Messaging
- `GET /v1/messages/conversations` - Get all conversations
- `GET /v1/messages/conversations/{id}` - Get messages in conversation
- `POST /v1/messages` - Send message
- `POST /v1/messages/{id}/read` - Mark message as read

#### Search
- `POST /v1/search` - Advanced search (supports query params and JSON body)
  - Query parameters: `q`, `type`, `location`, `country`, `city`, `latitude`, `longitude`, `radius`, `category`, `status`, `date_from`, `date_to`
  - JSON body: Full search request with all filters
- `GET /v1/search/history` - Get search history
- `DELETE /v1/search/history` - Clear search history

#### Custom CRUDs (User)
- `GET /v1/cruds/entities` - List user's CRUD entities
- `POST /v1/cruds/entities` - Create CRUD entity
- `GET /v1/cruds/entities/{id}` - Get CRUD entity
- `POST /v1/cruds/entities/{id}/data` - Create CRUD data
- `GET /v1/cruds/entities/{id}/data` - List CRUD data
- `GET /v1/cruds/templates/active` - Get active templates (public endpoint for users)
- `GET /v1/cruds/templates/{name}` - Get template details
- `POST /v1/cruds/templates/{name}/create` - Create entity from template

#### File Upload
- `POST /v1/files/upload/image` - Upload image file (multipart/form-data)
  - Supports: JPG, PNG, GIF
  - Max size: 5MB
  - Returns: File URL and metadata

#### Theme
- `GET /v1/users/me/settings/theme` - Get user theme
- `PUT /v1/users/me/settings/theme` - Update user theme
- `POST /v1/users/me/settings/theme/sync` - Sync theme across devices

### Protected Admin Endpoints

#### User Management
- `GET /v1/admin/users` - List all users (with search)
- `GET /v1/admin/users/{id}` - Get user details
- `POST /v1/admin/users` - Create new user
- `PUT /v1/admin/users/{id}` - Update user
- `DELETE /v1/admin/users/{id}` - Delete user
- `POST /v1/admin/users/{id}/status` - Update user status
- `GET /v1/admin/users/{id}/sessions` - Get user sessions
- `DELETE /v1/admin/users/{id}/sessions` - Revoke user sessions

#### Admin Settings
- `GET /v1/admin/settings` - Get admin settings
- `PUT /v1/admin/settings` - Update admin settings (including verification code)

#### Custom CRUDs (Admin)
- `GET /v1/admin/cruds/entities` - List all CRUD entities
- `POST /v1/admin/cruds/entities` - Create CRUD entity
- `GET /v1/admin/cruds/entities/{id}` - Get CRUD entity
- `PUT /v1/admin/cruds/entities/{id}` - Update CRUD entity
- `DELETE /v1/admin/cruds/entities/{id}` - Delete CRUD entity
- `POST /v1/admin/cruds/entities/{id}/data` - Create CRUD data
- `GET /v1/admin/cruds/entities/{id}/data` - List CRUD data
- `GET /v1/admin/cruds/data/{id}` - Get CRUD data
- `PUT /v1/admin/cruds/data/{id}` - Update CRUD data
- `DELETE /v1/admin/cruds/data/{id}` - Delete CRUD data

#### CRUD Templates (Admin)
- `GET /v1/admin/cruds/templates` - Get all templates (with filters)
- `POST /v1/admin/cruds/templates` - Create new template
- `GET /v1/admin/cruds/templates/{name}` - Get template by name
- `PUT /v1/admin/cruds/templates/id/{id}` - Update template
- `DELETE /v1/admin/cruds/templates/id/{id}` - Delete template

#### Admin Management
- `GET /v1/admin/admins` - List all admins
- `POST /v1/admin/admins` - Add new admin

#### Access Requests
- `GET /v1/admin/requests` - List access requests
- `POST /v1/admin/requests/{id}/status` - Update request status

#### Activity Logs
- `GET /v1/admin/logs` - Get activity logs

## Database Schema

### Core Tables
- `users` - User accounts (id, email, name, password_hash, role, status, created_at, updated_at, deleted_at)
- `sessions` - User sessions (id, user_id, token, refresh_token, expires_at, ip_address, user_agent, device_id, created_at, last_used_at)
- `devices` - User devices (id, user_id, device_id, device_name, ip_address, location_country, location_city, last_seen_at, created_at)
- `user_settings_comprehensive` - All user settings (user_id, profile_visibility, email_visibility, phone_visibility, allow_messaging, search_visibility, email_notifications, push_notifications, language, timezone, theme, font_size, etc.)
- `dashboard_items` - Dashboard items (id, user_id, title, description, category, status, priority, metadata, created_at, updated_at, deleted_at)
- `messages` - Messages (id, conversation_id, sender_id, recipient_id, subject, content, is_read, created_at)
- `conversations` - Message conversations (id, user1_id, user2_id, last_message_at, created_at)
- `notifications` - Notifications (id, user_id, type, title, message, is_read, created_at)
- `search_history` - Search history (id, user_id, query, search_type, results_count, created_at)
- `admin_settings` - Admin settings (admin_id, dashboard_layout, default_permissions, notification_preferences, theme_preferences, admin_verification_code, created_at, updated_at)
- `custom_crud_entities` - Custom CRUD entities (id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at)
- `custom_crud_data` - Custom CRUD data (id, entity_id, data, created_at, updated_at)
- `crud_templates` - CRUD templates (id, name, display_name, description, category, icon, schema, created_by, is_active, is_system, created_at, updated_at)
- `password_resets` - Password reset tokens (id, user_id, token, expires_at, used_at, created_at)
- `access_requests` - Access requests (id, user_id, title, details, status, feedback, created_at, updated_at)
- `activity_logs` - Activity logs (id, user_id, action, resource_type, resource_id, ip_address, user_agent, created_at)

## Technology Stack

### Backend
- **Language**: Go (Golang) 1.21+
- **Database**: SQLite (modernc.org/sqlite)
- **Router**: Gorilla Mux
- **Authentication**: JWT (JSON Web Tokens)
- **Logging**: Zap Logger (Structured logging)
- **Validation**: go-playground/validator
- **Password Hashing**: bcrypt
- **UUID**: google/uuid
- **HTTP Client**: Standard library

### Frontend (Current Implementation)
- **HTML5** - Structure
- **CSS3** - Styling (Custom CSS with CSS Variables)
- **JavaScript (ES6+)** - Functionality
- **Leaflet.js** - Interactive maps for location search
- **OpenStreetMap** - Map tiles
- **Nominatim API** - Reverse geocoding

### Third-Party Services
- **Nominatim** - Geocoding and reverse geocoding
- **OpenStreetMap** - Map tiles (via Leaflet)

## Frontend Features

### User Dashboard
- ✅ Responsive design with cards layout
- ✅ Navbar with search, notifications, messages, profile
- ✅ Avatar with dropdown menu (Profile, Settings, Logout)
- ✅ Live search with dropdown results
- ✅ Advanced search modal with map integration
- ✅ Search Near Me with radius selection
- ✅ Notifications modal with unread count badge
- ✅ Messages modal with user search and conversations
- ✅ Profile card modal
- ✅ Dashboard items management
- ✅ CRUD templates browsing
- ✅ Custom CRUD creation and management
- ✅ Real-time polling for messages and notifications
- ✅ Settings page with 8 categories

### Admin Dashboard
- ✅ User management with action buttons (View, Edit, Toggle Status, Delete)
- ✅ CRUD Templates management (Create, Edit, Delete)
- ✅ Custom CRUDs management
- ✅ Admin Settings (Verification code management)
- ✅ All user dashboard features (search, messaging, notifications)
- ✅ Enhanced user cards with detailed information
- ✅ Template creation with dynamic field builder
- ✅ Schema JSON editor with example loader

## Key Features Details

### Search System
- **Live Search**: Real-time search results in dropdown as you type
- **Advanced Search**: Modal with comprehensive filters
- **Map Search**: Interactive map to select location by clicking
- **Search Near Me**: Uses browser geolocation with configurable radius (1-100km)
- **Reverse Geocoding**: Converts coordinates to readable addresses
- **Distance Calculation**: Haversine formula for accurate distance filtering
- **Multi-Entity Search**: Search across users, dashboard items, messages, notifications, CRUDs

### Messaging System
- **User Search**: Search for users within messages modal
- **Conversations**: One-on-one conversation threads
- **Real-time Updates**: Polling for new messages
- **Unread Counts**: Badge showing unread message count
- **Message Threading**: Organized conversation view
- **Settings Control**: Enable/disable message polling in settings

### Notification System
- **Real-time Updates**: Polling for new notifications
- **Unread Counts**: Badge showing unread notification count
- **Notification Types**: Message, Alert, Promotion, Security, System
- **Mark as Read**: Individual and bulk read operations
- **Settings Control**: Enable/disable notification polling in settings

### CRUD Templates System
- **Dynamic Creation**: Admin can create templates with custom fields
- **Field Types**: Support for various field types (string, number, date, etc.)
- **Schema Builder**: JSON-based schema definition
- **Template Categories**: Organize templates by category
- **Active/Inactive**: Control template availability
- **User Access**: Users can browse and use active templates
- **Template to Entity**: One-click entity creation from template

### Admin Verification System
- **Default Code**: `Kompasstech2025@`
- **Code Verification**: Required before admin account creation
- **Code Management**: Admins can change verification code
- **Code Persistence**: Stored in database per admin
- **System-wide Code**: First admin's code used system-wide

## Conclusion

The Base App provides a comprehensive set of features for user management, settings, dashboard, messaging, notifications, search, and admin functionality. All features are production-ready and fully integrated.

