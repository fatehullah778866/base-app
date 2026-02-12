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

FROM gcr.io/distroless/static-debian11
WORKDIR /app

# Copy built binary
COPY --from=builder /out/server /app/server

# Copy migrations, frontend, and uploads if present (preserve structure)
COPY --from=builder /src/migrations /app/migrations
COPY --from=builder /src/frontend /app/frontend
COPY --from=builder /src/uploads /app/uploads

# Default listen port and environment
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/server"]
