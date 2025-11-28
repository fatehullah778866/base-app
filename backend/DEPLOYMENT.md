# Deployment Guide

This guide covers deploying the Base App Service to Google Cloud Platform.

## Prerequisites

1. Google Cloud Project with billing enabled
2. `gcloud` CLI installed and configured
3. Docker installed locally (for local testing)
4. Service account with necessary permissions

## Local Development Setup

### 1. Start Local Services

```bash
# Start PostgreSQL and Redis using Docker Compose
make docker-compose-up

# Or manually:
docker-compose up -d
```

### 2. Run Database Migrations

```bash
# Run migrations
make migrate-up

# Or manually:
./scripts/migrate.sh up
```

### 3. Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your local settings
# For local development, defaults should work:
# DB_HOST=localhost
# REDIS_HOST=localhost
```

### 4. Run the Server

```bash
# Run the server
make run

# Or manually:
go run ./cmd/server/main.go
```

The server will start on `http://localhost:8080`

### 5. Test the API

```bash
# Run API tests
./scripts/test-api.sh

# Or test manually with curl:
curl http://localhost:8080/health
```

## GCP Deployment

### 1. Enable Required APIs

```bash
gcloud services enable \
    cloudbuild.googleapis.com \
    run.googleapis.com \
    sqladmin.googleapis.com \
    redis.googleapis.com \
    secretmanager.googleapis.com
```

### 2. Create Cloud SQL Instance

```bash
gcloud sql instances create base-app-db \
    --database-version=POSTGRES_15 \
    --tier=db-f1-micro \
    --region=us-central1 \
    --root-password=CHANGE_ME
```

### 3. Create Database

```bash
gcloud sql databases create base_app_db \
    --instance=base-app-db
```

### 4. Create Redis Instance (Memorystore)

```bash
gcloud redis instances create base-app-redis \
    --size=1 \
    --region=us-central1 \
    --redis-version=redis_7_0
```

### 5. Store Secrets in Secret Manager

```bash
# JWT Secret
echo -n "your-super-secret-jwt-key" | gcloud secrets create jwt-secret --data-file=-

# Webhook Secret
echo -n "your-super-secret-webhook-key" | gcloud secrets create webhook-secret --data-file=-

# Database Password
echo -n "your-db-password" | gcloud secrets create db-password --data-file=-
```

### 6. Create VPC Connector (for Cloud SQL and Redis access)

```bash
gcloud compute networks vpc-access connectors create base-app-vpc-connector \
    --region=us-central1 \
    --network=default \
    --range=10.8.0.0/28
```

### 7. Grant Cloud Run Service Account Access

```bash
PROJECT_ID=$(gcloud config get-value project)
SERVICE_ACCOUNT="${PROJECT_ID}@appspot.gserviceaccount.com"

# Grant Secret Manager access
gcloud secrets add-iam-policy-binding jwt-secret \
    --member="serviceAccount:${SERVICE_ACCOUNT}" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding webhook-secret \
    --member="serviceAccount:${SERVICE_ACCOUNT}" \
    --role="roles/secretmanager.secretAccessor"

gcloud secrets add-iam-policy-binding db-password \
    --member="serviceAccount:${SERVICE_ACCOUNT}" \
    --role="roles/secretmanager.secretAccessor"

# Grant Cloud SQL access
gcloud projects add-iam-policy-binding ${PROJECT_ID} \
    --member="serviceAccount:${SERVICE_ACCOUNT}" \
    --role="roles/cloudsql.client"
```

### 8. Update cloudbuild.yaml

Edit `cloudbuild.yaml` with your project-specific values:
- Project ID
- Region
- Cloud SQL instance connection name
- VPC connector name

### 9. Deploy Using Cloud Build

```bash
# Submit build
gcloud builds submit --config cloudbuild.yaml

# Or trigger from GitHub (if configured)
# Push to main branch triggers automatic deployment
```

### 10. Run Migrations on Cloud SQL

```bash
# Get Cloud SQL connection name
INSTANCE_CONNECTION_NAME=$(gcloud sql instances describe base-app-db --format="value(connectionName)")

# Run migrations using Cloud SQL Proxy or Cloud Build
gcloud builds submit --config=migration-build.yaml
```

### 11. Set Environment Variables

After deployment, update Cloud Run service with environment variables:

```bash
gcloud run services update base-app-service \
    --region=us-central1 \
    --update-env-vars="DB_HOST=/cloudsql/${INSTANCE_CONNECTION_NAME},DB_NAME=base_app_db,DB_USER=postgres"
```

## Environment Variables for Production

Set these in Cloud Run:

```bash
# Server
PORT=8080
ENV=production

# Database (Cloud SQL)
DB_HOST=/cloudsql/PROJECT_ID:REGION:INSTANCE_NAME
DB_PORT=5432
DB_USER=postgres
DB_NAME=base_app_db
DB_SSL_MODE=disable

# Redis (Memorystore)
REDIS_HOST=10.x.x.x  # Internal IP from Memorystore
REDIS_PORT=6379

# Secrets (from Secret Manager)
JWT_SECRET (from secret manager)
WEBHOOK_SECRET (from secret manager)
DB_PASSWORD (from secret manager)
```

## Monitoring

### View Logs

```bash
gcloud run services logs read base-app-service --region=us-central1
```

### View Metrics

Visit Cloud Console → Cloud Run → base-app-service → Metrics

## Troubleshooting

### Database Connection Issues

1. Verify VPC connector is active
2. Check Cloud SQL instance is running
3. Verify service account has `cloudsql.client` role
4. Check connection name format: `/cloudsql/PROJECT:REGION:INSTANCE`

### Secret Access Issues

1. Verify service account has `secretmanager.secretAccessor` role
2. Check secret names match in Cloud Run configuration
3. Verify secrets exist in Secret Manager

### Build Failures

1. Check Cloud Build logs in Console
2. Verify Dockerfile builds locally: `docker build -t test .`
3. Check image permissions in Container Registry

## Rollback

```bash
# List revisions
gcloud run revisions list --service=base-app-service --region=us-central1

# Rollback to previous revision
gcloud run services update-traffic base-app-service \
    --region=us-central1 \
    --to-revisions=REVISION_NAME=100
```

