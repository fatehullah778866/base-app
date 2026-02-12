# Multi-stage Dockerfile for Base App
# Builds the Go backend and includes frontend/static assets and migrations.

FROM golang:1.24-bullseye AS builder
WORKDIR /src
ENV CGO_ENABLED=0

# Copy entire project to allow the build to find modules and frontend/migrations
COPY . .

# Ensure expected directories exist so later COPY steps don't fail in Cloud Build
RUN mkdir -p /src/migrations /src/frontend /src/uploads || true

WORKDIR /src/backend

# Ensure modules are downloaded and build the server
RUN go env -w GOPROXY=https://proxy.golang.org,direct \
    && go mod download \
    && go build -ldflags "-s -w" -o /out/server ./cmd/server

# Create runtime upload dir in builder (will be copied into final image)
RUN mkdir -p /tmp/uploads || true

FROM gcr.io/distroless/static-debian11
WORKDIR /app

# Copy built binary
COPY --from=builder /out/server /app/server

# Copy migrations, frontend, and uploads if present (preserve structure)
COPY --from=builder /src/migrations /app/migrations
COPY --from=builder /src/frontend /app/frontend
COPY --from=builder /src/uploads /app/uploads

# Copy empty/writeable runtime upload dir created in builder into final image
COPY --from=builder /tmp/uploads /tmp/uploads

# Ensure runtime writable locations and sensible defaults for Cloud Run
# Use SQLite database in /tmp (writable in containers) to allow startup without external DB
ENV DB_DRIVER=sqlite
ENV DB_SQLITE_PATH=file:/tmp/app.db?_pragma=foreign_keys(ON)
# Default listen port and upload dir
ENV PORT=8080
ENV UPLOAD_DIR=/tmp/uploads

EXPOSE 8080

ENTRYPOINT ["/app/server"]
