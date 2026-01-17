## 🐳 Docker Compose Setup - Detailed Guide

### What is included in Docker Compose

**PostgreSQL 16 Service:**
- Container: `gophkeeper-postgres`
- Image: `postgres:16-alpine`
- Port: `5432` (host) → `5432` (container)
- Database: `gophkeeper`
- User: `gophkeeper`
- Password: `gophkeeper_password`
- Health check: `pg_isready` every 10s
- Volume: `postgres_data:/var/lib/postgresql/data`

**GophKeeper Server Service:**
- Container: `gophkeeper-server`
- Build: `Dockerfile` (Debian Bookworm-based)
- Port: `8080` (host) → `8080` (container)
- Depends on: `postgres` (service_healthy)
- Health check: `curl http://localhost:8080/health` every 30s

### Makefile Commands

```bash
make docker-up      # Start PostgreSQL + server
make docker-down    # Stop containers
make docker-logs    # View logs (live)
make docker-build   # Rebuild Docker image
make docker-clean   # Remove containers + volumes
```

### Environment variables

The server receives the following on startup:

```yaml
DATABASE_URL: postgres://gophkeeper:gophkeeper_password@postgres:5432/gophkeeper?sslmode=disable
SERVER_ADDRESS: 0.0.0.0:8080
SECRET_KEY: your-secret-key-here
LOG_LEVEL: info
```

### Container status

```bash
# Show status of all containers (with health check)
docker-compose ps

# Expected output:
# gophkeeper-postgres  postgres:16-alpine  Up (healthy)
# gophkeeper-server    gophkeeper-*        Up (health: starting/healthy)
```

Data management

**Persistence:**
```bash
# Stop containers (data will be preserved)
docker-compose down

# Restart (data will be restored)
docker-compose up -d
```

**Remove data:**
```bash
# IMPORTANT: This will delete all DB data!
docker-compose down -v
```

### Connect to the DB

**From the container:**
```bash
docker-compose exec postgres psql -U gophkeeper -d gophkeeper
```

**From the host (if psql is installed):**
```bash
psql -h localhost -p 5432 -U gophkeeper -d gophkeeper
# Password: gophkeeper_password
```

### Health checks

```bash
# Server health check
curl http://localhost:8080/health
# Response: {"status":"ok"}

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Logging and debug

```bash
# View server logs
docker-compose logs gophkeeper-server

# View PostgreSQL logs
docker-compose logs postgres

# Live logs (follow)
docker-compose logs -f

# Logs for specific service
docker-compose logs -f gophkeeper-server
```

Check what is occupying a port:
```bash
lsof -i :5432
```

# Option 1: Stop local PostgreSQL
```bash
sudo systemctl stop postgresql  # Linux
brew services stop postgresql   # macOS
```

# Option 2: Use another port
Edit docker-compose.yml and change:
```yaml
ports:
  - "5433:5432"  # use 5433 instead of 5432
```

**Port 8080 is occupied:**
```bash
# Check what occupies the port
lsof -i :8080

# Edit docker-compose.yml and change:
# ports:
#   - "8081:8080"  # use 8081 instead of 8080
```

**PostgreSQL container not becoming healthy:**
```bash
# View logs
docker-compose logs postgres

# Recreate the container
docker-compose down -v
docker-compose up -d
```

**Server container crashes:**
```bash
# Check error logs
docker-compose logs gophkeeper-server

# Check DATABASE_URL
docker-compose exec gophkeeper-server env | grep DATABASE_URL

# Rebuild the image
make docker-build
docker-compose up -d
```

**Database connection error:**
```bash
# Check that PostgreSQL is healthy
docker-compose ps postgres

# Check environment variables
docker-compose config | grep -A 10 gophkeeper-server

# Verify they are passed correctly
docker-compose exec gophkeeper-server env | grep DATABASE
```

```bash
# Check which process uses the port
lsof -i :5432

# Use a different port in docker-compose.yml:
# Change "5432:5432" to "5433:5432"
```

**Server cannot connect to DB:**
- Ensure the PostgreSQL container is `healthy`
- Check logs: `docker-compose logs postgres`
- Verify `DATABASE_URL` in the compose file

**I want to use my own PostgreSQL password:**
Edit `docker-compose.yml` or use environment variables:

```bash
POSTGRES_PASSWORD=my-secure-password docker-compose up -d
```
