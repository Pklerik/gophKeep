# Docker Deployment - Quick Reference

## 📋 Created Files

### 1. `docker-compose.yml`
```yaml
✅ PostgreSQL 16 service
✅ GophKeeper Server service  
✅ Health checks for both services
✅ Volume for database data persistence
✅ Automatic readiness wait for DB before server starts
```

### 2. `Dockerfile`
```dockerfile
✅ Multi-stage build (builder + final)
✅ Go 1.25 Alpine for compilation
✅ Final image on Alpine 3.19
✅ Health check endpoint
✅ Minimal image size
```

### 3. `Makefile` - new commands
```makefile
✅ make docker-up      - Start PostgreSQL and server
✅ make docker-down    - Stop containers  
✅ make docker-logs    - View logs in real-time
✅ make docker-build   - Rebuild Docker image
✅ make docker-clean   - Remove containers and volumes
```

### 4. `internal/config/server/config.go`
- ✅ Support for DATABASE_URL (for PostgreSQL)
- ✅ Backward compatible with DATABASE_PATH (for SQLite)
- ✅ Reads variables from .env

## 🚀 Quick Start

### Start
```bash
make docker-up
```

### Verify
```bash
docker-compose ps
curl http://localhost:8080/health
```

### Logs
```bash
make docker-logs
```

### Stop
```bash
make docker-down
```

## 🔑 Default Credentials

| Service | User | Password | Port | URL |
|---------|------|----------|------|-----|
| PostgreSQL | gophkeeper | gophkeeper_password | 5432 | localhost:5432 |
| Server | - | - | 8080 | http://localhost:8080 |

## 📝 Makefile Commands

```bash
# Docker commands
make docker-up      # Start PostgreSQL + Server
make docker-down    # Stop containers
make docker-logs    # View logs (follow)
make docker-build   # Rebuild image
make docker-clean   # Remove containers + volumes

# Other commands
make build          # Build binaries
make test           # Run tests
make coverage       # Run tests with coverage
```

## 🐳 Service Details

### PostgreSQL
- Image: `postgres:16-alpine`
- Container: `gophkeeper-postgres`
- Port: `5432`
- Health check: `pg_isready` every 10s
- Volume: `postgres_data:/var/lib/postgresql/data`
- Database: `gophkeeper`

### GophKeeper Server
- Build: Multi-stage Dockerfile (Debian Bookworm)
- Container: `gophkeeper-server`
- Port: `8080`
- Health check: `curl http://localhost:8080/health` every 30s
- Depends on: `postgres` (service_healthy)

## 🎯 Common Tasks

### Backup Database
```bash
docker-compose exec postgres pg_dump -U gophkeeper -d gophkeeper > backup.sql
```

### Restore Database
```bash
docker-compose exec -T postgres psql -U gophkeeper -d gophkeeper < backup.sql
```

### Connect to Database
```bash
# From container
docker-compose exec postgres psql -U gophkeeper -d gophkeeper

# From host (if psql installed)
psql -h localhost -p 5432 -U gophkeeper -d gophkeeper
```

### View Container Logs
```bash
docker-compose logs              # All logs once
docker-compose logs -f           # All logs follow
docker-compose logs gophkeeper-server  # Only server
docker-compose logs postgres     # Only database
```

## ✅ Health Checks

### Server Health
```bash
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

### Container Status
```bash
docker-compose ps
# Both should show (healthy) or (starting)
```

## 📚 Available Documents

- **DOCKER.md** - Main Docker documentation
- **DOCKER_SETUP.md** - Detailed setup guide
- **.env.example** - Example environment variables

## ⚙️ Structure

```
gophKeep/
├── docker-compose.yml       # ← Docker Compose config
├── Dockerfile              # ← Dockerfile for the server
├── .env.example            # ← Example env variables
├── DOCKER.md              # ← Documentation
├── DOCKER_SETUP.md        # ← Setup guide
├── Makefile               # ← Updated with Docker commands
├── go.mod
├── go.sum
├── cmd/
│   └── gophkeep/
│       ├── main.go
│       ├── server/
│       └── client/
└── internal/
    ├── config/
    │   └── server/config.go  # ← Updated for DATABASE_URL
    └── ...
```

## 💡 Quick Commands

```bash
# Full cycle: start, logs, stop
make docker-up && make docker-logs &
# ... run tests ...
make docker-down

# Restart
make docker-down && make docker-up

# Clean and restart
make docker-clean && make docker-up
```