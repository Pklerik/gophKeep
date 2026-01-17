# GophKeeper Docker & PostgreSQL Deployment

## 📋 What's Included

Full Docker Compose support with PostgreSQL for running the server in a container.

### Files
- **docker-compose.yml** - Config for running PostgreSQL and GophKeeper Server
- **Dockerfile** - Multi-stage build (Debian-based) for the server
- **.env.example** - Example environment variables
- **Makefile** - Added `docker-*` commands

## 🚀 Quick Start

### 1️⃣ Ensure Docker and Docker Compose are installed
```bash
docker --version      # >= 20.10
docker-compose --version  # >= 1.29
```

### 2️⃣ Start the service
```bash
make docker-up
```

This will start:
- **PostgreSQL 16** on local port **5432**
- **GophKeeper Server** on local port **8080**

### 3️⃣ Check the status
```bash
docker-compose ps
```

Both services should be in the `Up` status (PostgreSQL - `healthy`)

### 4️⃣ Test
```bash
# Check the server
curl http://localhost:8080/health

# Connect to the database
psql -h localhost -U gophkeeper -d gophkeeper
# Password: gophkeeper_password
```

## 🐳 Docker Services

### PostgreSQL Service
```yaml
container_name: gophkeeper-postgres
image: postgres:16-alpine
ports: 5432:5432
environment:
  POSTGRES_DB: gophkeeper
  POSTGRES_USER: gophkeeper
  POSTGRES_PASSWORD: gophkeeper_password
healthcheck:
  - pg_isready every 10 seconds
volumes:
  - postgres_data:/var/lib/postgresql/data
```

### GophKeeper Server Service
```yaml
container_name: gophkeeper-server
build: Dockerfile (Debian-based, multi-stage)
ports: 8080:8080
depends_on: postgres (service_healthy)
environment:
  SERVER_ADDRESS: 0.0.0.0:8080
  DATABASE_URL: postgres://gophkeeper:gophkeeper_password@postgres:5432/gophkeeper
  SECRET_KEY: ${SECRET_KEY:-your-secret-key-here}
  LOG_LEVEL: info
healthcheck:
  - curl http://localhost:8080/health every 30 seconds
```

## 📝 Makefile Commands

```bash
make docker-up       # Start PostgreSQL + server
make docker-down     # Stop containers
make docker-logs     # View logs in real-time
make docker-build    # Rebuild Docker image
make docker-clean    # Remove containers and volumes
```

## ⚙️ Production Configuration

1. **Copy .env**:
   ```bash
   cp .env.example .env
   ```

2. **Update credentials in .env**:
   ```env
   SECRET_KEY=<generate-with: openssl rand -base64 32>
   POSTGRES_PASSWORD=<strong-password>
   ```

3. **Start with .env**:
   ```bash
   docker-compose --env-file .env up -d
   ```

## 🔍 Troubleshooting

### Server cannot connect to DB
```bash
# Check logs
docker-compose logs gophkeeper-server

# Ensure PostgreSQL is healthy
docker-compose logs postgres | grep "ready"
```

### Port is in use
```bash
# Example: change PostgreSQL port to 5433
docker-compose.yml:
  postgres:
    ports:
      - "5433:5432"
```

### Clean everything and restart
```bash
make docker-clean
make docker-up
```

## 📊 Monitoring

```bash
# Live logs of all services
make docker-logs

# Logs of a specific service
docker-compose logs -f gophkeeper-server
docker-compose logs -f postgres

# Container status
docker-compose ps
```

## 🔗 Testing Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

## 📚 Documentation

- **DOCKER_QUICK_START.md** - Start here
- **DOCKER_SETUP.md** - Detailed instructions
- **DOCKER.md** - General information

## 💾 Data

PostgreSQL data is stored in the Docker volume `postgres_data`. Running `docker-compose down -v` will delete the volume.

To persist data between restarts, use `docker-compose down` (without -v):
```bash
docker-compose down  # data will be preserved
docker-compose up -d # data will be restored
```

---

**Ready to launch:**
```bash
make docker-up && make docker-logs &
```