# Hospital EMR System - Deployment Guide

## Overview

This guide provides step-by-step instructions for deploying the Hospital EMR System in various environments.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development](#local-development)
3. [Docker Deployment](#docker-deployment)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [Production Deployment](#production-deployment)
6. [Database Setup](#database-setup)
7. [Security Configuration](#security-configuration)
8. [Monitoring and Logging](#monitoring-and-logging)

---

## Prerequisites

### Required Software

- **Go**: 1.21 or higher
- **PostgreSQL**: 15 or higher (or Neon account)
- **Redis**: 7 or higher
- **NATS**: 2.10 or higher
- **Docker**: 20.10 or higher (optional)
- **Kubernetes**: 1.25 or higher (for K8s deployment)

### Hardware Requirements

**Minimum** (Development):
- CPU: 2 cores
- RAM: 4 GB
- Storage: 20 GB

**Recommended** (Production):
- CPU: 8 cores
- RAM: 16 GB
- Storage: 100 GB SSD

---

## Local Development

### 1. Clone Repository

```bash
git clone https://github.com/hospital-emr/backend.git
cd backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` and configure:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=emr_database
JWT_SECRET=your_jwt_secret_key_at_least_32_chars
ENCRYPTION_KEY=your_32_byte_encryption_key_here
```

### 4. Setup Database

```bash
# Start PostgreSQL (or use Neon)
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=postgres123 \
  -e POSTGRES_DB=emr_database \
  -p 5432:5432 \
  postgres:15-alpine

# Run migrations
make migrate-up

# Seed initial data
make seed
```

### 5. Start Services

```bash
# Start Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# Start NATS
docker run -d --name nats -p 4222:4222 nats:latest

# Run application
make run
```

The API will be available at `http://localhost:8080`

### 6. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@hospital-emr.com","password":"admin123"}'
```

---

## Docker Deployment

### Using Docker Compose

The easiest way to deploy all services:

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

### Build Custom Image

```bash
# Build image
docker build -t hospital-emr-backend:latest .

# Run container
docker run -d \
  --name emr-api \
  -p 8080:8080 \
  --env-file .env \
  hospital-emr-backend:latest
```

---

## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (EKS, GKE, AKS, or local minikube)
- kubectl configured
- Helm 3 (optional)

### 1. Create Namespace

```bash
kubectl create namespace emr-system
```

### 2. Create Secrets

```bash
# Database credentials
kubectl create secret generic db-credentials \
  --from-literal=username=postgres \
  --from-literal=password=your_secure_password \
  -n emr-system

# JWT secret
kubectl create secret generic jwt-secret \
  --from-literal=secret=your_jwt_secret_key \
  -n emr-system

# Encryption key
kubectl create secret generic encryption-key \
  --from-literal=key=your_32_byte_encryption_key \
  -n emr-system
```

### 3. Deploy PostgreSQL

```bash
kubectl apply -f deployments/kubernetes/postgres-statefulset.yaml -n emr-system
```

### 4. Deploy Redis

```bash
kubectl apply -f deployments/kubernetes/redis-deployment.yaml -n emr-system
```

### 5. Deploy NATS

```bash
kubectl apply -f deployments/kubernetes/nats-deployment.yaml -n emr-system
```

### 6. Deploy Application

```bash
kubectl apply -f deployments/kubernetes/api-deployment.yaml -n emr-system
kubectl apply -f deployments/kubernetes/api-service.yaml -n emr-system
```

### 7. Deploy Ingress

```bash
kubectl apply -f deployments/kubernetes/ingress.yaml -n emr-system
```

### 8. Verify Deployment

```bash
kubectl get pods -n emr-system
kubectl get services -n emr-system
kubectl logs -f deployment/emr-api -n emr-system
```

---

## Production Deployment

### Using Terraform (AWS)

```bash
cd deployments/terraform/aws

# Initialize
terraform init

# Plan
terraform plan

# Apply
terraform apply
```

### Neon PostgreSQL Setup

1. Create Neon account: https://neon.tech
2. Create new project
3. Get connection string
4. Update `.env`:
   ```env
   DB_HOST=your-project.neon.tech
   DB_PORT=5432
   DB_USER=your_user
   DB_PASSWORD=your_password
   DB_NAME=your_database
   DB_SSLMODE=require
   ```

### Environment Variables for Production

```env
APP_ENV=production
APP_PORT=8080

# Database
DB_HOST=your-neon-host.neon.tech
DB_SSLMODE=require
DB_MAX_CONNECTIONS=100

# JWT (use strong random keys)
JWT_SECRET=<generate-with-openssl-rand-base64-32>
JWT_EXPIRATION_HOURS=24

# Encryption
ENCRYPTION_KEY=<generate-with-openssl-rand-bytes-32>
DATA_ENCRYPTION_ENABLED=true

# Security
MFA_ISSUER=Hospital-EMR
SESSION_TIMEOUT_MINUTES=30
RATE_LIMIT_REQUESTS_PER_MINUTE=100

# CORS (update with your frontend URLs)
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com,https://admin.yourdomain.com

# Monitoring
ENABLE_METRICS=true
LOG_LEVEL=info
LOG_FORMAT=json
```

### Generate Secure Keys

```bash
# JWT Secret
openssl rand -base64 32

# Encryption Key (32 bytes)
openssl rand -base64 32
```

---

## Database Setup

### Migrations

```bash
# Run all migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create name=add_new_table
```

### Seed Data

```bash
# Seed initial data (roles, permissions, admin user)
make seed
```

Default credentials after seeding:
- Admin: `admin@hospital-emr.com` / `admin123`
- Doctor: `doctor@hospital-emr.com` / `doctor123`

**⚠️ Change these passwords immediately in production!**

### Backup and Restore

```bash
# Backup
pg_dump -h localhost -U postgres -d emr_database > backup.sql

# Restore
psql -h localhost -U postgres -d emr_database < backup.sql
```

---

## Security Configuration

### SSL/TLS Certificates

```bash
# Generate self-signed certificate (development only)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout server.key -out server.crt

# Use Let's Encrypt (production)
certbot certonly --standalone -d api.yourdomain.com
```

### Firewall Rules

Allow only necessary ports:
```bash
# API
ufw allow 8080/tcp

# PostgreSQL (only from application)
ufw allow from <app-ip> to any port 5432

# Redis (only from application)
ufw allow from <app-ip> to any port 6379
```

### Environment Isolation

- Development: `APP_ENV=development`
- Staging: `APP_ENV=staging`
- Production: `APP_ENV=production`

---

## Monitoring and Logging

### Prometheus Metrics

Access metrics at: `http://localhost:9090/metrics`

### Grafana Dashboards

Access Grafana at: `http://localhost:3001`
- Username: `admin`
- Password: `admin123`

### Log Aggregation

Logs are structured in JSON format and can be sent to:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- CloudWatch (AWS)
- Stackdriver (GCP)

### Health Checks

- Health: `GET /health`
- Readiness: `GET /ready`

---

## Scaling

### Horizontal Scaling

```bash
# Scale API pods
kubectl scale deployment emr-api --replicas=5 -n emr-system
```

### Database Scaling

For Neon:
- Compute scales automatically
- Read replicas for read-heavy workloads

### Caching Strategy

- Session data: Redis
- API responses: Redis with TTL
- Static data: Application memory cache

---

## Troubleshooting

### Application Won't Start

1. Check logs: `docker-compose logs api`
2. Verify database connection: `psql -h localhost -U postgres`
3. Check environment variables: `env | grep DB_`

### Database Connection Issues

```bash
# Test connection
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME

# Check if database exists
psql -h $DB_HOST -U $DB_USER -l
```

### High Memory Usage

1. Check connection pool settings
2. Monitor with: `kubectl top pods -n emr-system`
3. Adjust resource limits in deployment

---

## Rollback Procedure

### Docker

```bash
# Rollback to previous image
docker-compose down
docker-compose up -d --force-recreate
```

### Kubernetes

```bash
# Rollback deployment
kubectl rollout undo deployment/emr-api -n emr-system

# Check rollout status
kubectl rollout status deployment/emr-api -n emr-system
```

---

## Support

For deployment support:
- Documentation: https://docs.hospital-emr.com
- Email: devops@hospital-emr.com
- Slack: #emr-deployment

---

## Checklist

Before going to production:

- [ ] Change all default passwords
- [ ] Generate secure JWT secret (32+ chars)
- [ ] Generate secure encryption key (32 bytes)
- [ ] Configure SSL/TLS certificates
- [ ] Setup database backups
- [ ] Configure monitoring and alerting
- [ ] Test disaster recovery plan
- [ ] Review security policies
- [ ] Configure log retention
- [ ] Setup rate limiting
- [ ] Enable MFA for admin accounts
- [ ] Document runbook procedures
- [ ] Conduct security audit
- [ ] Load testing completed
- [ ] Backup and restore tested
