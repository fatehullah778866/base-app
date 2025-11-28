# Base App Monorepo

This workspace now has a clear split between backend (Go API) and frontend (static HTML/JS).

- `backend/` – Go service with auth, user management, admin dashboard APIs. Run from inside this folder.
- `frontend/` – Landing page, user auth/dashboard, and admin dashboard UI that talks to the backend APIs.

## Quick start

```bash
# backend
cd backend
go run cmd/server/main.go

# frontend (served by backend automatically from ../frontend)
# open http://localhost:8080/ in your browser
```

Default admin: `admin@gmail.com` / `admin123`.
