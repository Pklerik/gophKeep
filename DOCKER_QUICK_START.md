# Docker Deployment - Quick Reference

## 📋 Создаваемые файлы

### 1. `docker-compose.yml`
```yaml
✅ PostgreSQL 16 сервис
✅ GophKeeper Server сервис  
✅ Health checks для обоих сервисов
✅ Volume для персистентности данных БД
✅ Автоматическое ожидание readiness БД перед стартом сервера
```

### 2. `Dockerfile`
```dockerfile
✅ Multi-stage build (builder + final)
✅ Go 1.25 Alpine для компиляции
✅ Финальный образ на Alpine 3.19
✅ Health check endpoint
✅ Минимальный размер образа
```

### 3. `Makefile` - новые команды
```makefile
✅ make docker-up      - Start PostgreSQL и сервер
✅ make docker-down    - Stop контейнеры  
✅ make docker-logs    - Просмотр логов в реальном времени
✅ make docker-build   - Пересобрать Docker образ
✅ make docker-clean   - Удалить контейнеры и volumes
```

### 4. `internal/config/server/config.go`
- ✅ Поддержка DATABASE_URL (для PostgreSQL)
- ✅ Backward compatible с DATABASE_PATH (для SQLite)
- ✅ Чтение переменных из .env

## 🚀 Быстрый старт

### Запуск
```bash
make docker-up
```

### Проверка
```bash
docker-compose ps
curl http://localhost:8080/health
```

### Логи
```bash
make docker-logs
```

### Остановка
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

## 📚 Доступные документы

- **DOCKER.md** - Основная документация по Docker
- **DOCKER_SETUP.md** - Детальное руководство по настройке
- **.env.example** - Пример переменных окружения

## ⚙️ Структура

```
gophKeep/
├── docker-compose.yml       # ← Docker Compose конфиг
├── Dockerfile              # ← Dockerfile для сервера
├── .env.example            # ← Пример env переменных
├── DOCKER.md              # ← Документация
├── DOCKER_SETUP.md        # ← Гайд по настройке
├── Makefile               # ← Обновлен с Docker командами
├── go.mod
├── go.sum
├── cmd/
│   └── gophkeep/
│       ├── main.go
│       ├── server/
│       └── client/
└── internal/
    ├── config/
    │   └── server/config.go  # ← Обновлена для DATABASE_URL
    └── ...
```

## 💡 Быстрые команды

```bash
# Полный цикл: запуск, логи, остановка
make docker-up && make docker-logs &
# ... сделайте тесты ...
make docker-down

# Перестартовать
make docker-down && make docker-up

# Очистить и заново
make docker-clean && make docker-up
```
