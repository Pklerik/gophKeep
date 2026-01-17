# GophKeeper Docker Setup

## Requirements

- **Docker** (>= 20.10)
- **Docker Compose** (>= 1.29)
- Free ports: 5432 (PostgreSQL), 8080 (Server)

## 🚀 Quick Start

### 1. Start all services
```bash
make docker-up
```

This will start:
- PostgreSQL 16 on local port 5432
- GophKeeper Server on local port 8080
- Both services with built-in health checks

### 2. Check status
```bash
docker-compose ps
```

Expected output:
```
NAME                 IMAGE                        STATUS
gophkeeper-postgres  postgres:16-alpine           Up (healthy)
gophkeeper-server    gophkeep-gophkeeper-server   Up (healthy)
```

### 3. Verify functionality
```bash
# Server health check
curl http://localhost:8080/health
# Response: {"status":"ok"}

# Connect to the database
psql -h localhost -p 5432 -U gophkeeper -d gophkeeper
# Password: gophkeeper_password
```

### 4. View logs
```bash
make docker-logs           # Live logs of all services
docker-compose logs postgres           # Only PostgreSQL
docker-compose logs gophkeeper-server  # Only the server
```

### 5. Stop services
```bash
make docker-down    # Stop (data will be preserved)
make docker-clean   # Full cleanup (including data)
```

## API Examples

### Health Check
```bash
curl http://localhost:8080/health
```

### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## 🔧 Troubleshooting

### Port already in use
```bash
# Find the process on the port
lsof -i :5432   # PostgreSQL
lsof -i :8080   # Server

# Use a different port in docker-compose.yml
# Change the line: "5432:5432" to "5433:5432"
```

### Container won't start
```bash
# View errors
docker-compose logs postgres
docker-compose logs gophkeeper-server

# Rebuild the image
make docker-build
docker-compose up -d
```

### Connection to DB refused
```bash
# Check that PostgreSQL is healthy
docker-compose ps postgres

# Check DATABASE_URL
docker-compose exec gophkeeper-server env | grep DATABASE_URL

# Check DB logs
docker-compose logs postgres | tail -50
```

## 📚 Documentation

- **DOCKER_README.md** - Introduction & Quick Start
- **DOCKER_QUICK_START.md** - Command Checklist
- **DOCKER_SETUP.md** - Detailed Guide
- **.env.example** - Environment Variables Template

Run:
```bash
docker-compose up -d
```

## Manual Database Connection

```bash
psql -h localhost -U gophkeeper -d gophkeeper
```

Password: `gophkeeper_password`

## Health Check

The server has a health check; verify the status:

```bash
docker-compose ps
```

## Data

PostgreSQL data is stored in the Docker volume `postgres_data`. Running `docker-compose down -v` will delete the volume.

## Local Testing

After starting the containers, the server is available at:
- API: http://localhost:8080
- Health check: http://localhost:8080/health