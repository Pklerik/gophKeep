# GophKeeper Docker Setup

## Требования

- **Docker** (>= 20.10)
- **Docker Compose** (>= 1.29)
- Свободные порты: 5432 (PostgreSQL), 8080 (Server)

## 🚀 Быстрый старт

### 1. Запуск всех сервисов
```bash
make docker-up
```

Это запустит:
- PostgreSQL 16 на локальном порту 5432
- GophKeeper Server на локальном порту 8080
- Оба сервиса с встроенными health checks

### 2. Проверка статуса
```bash
docker-compose ps
```

Ожидаемый вывод:
```
NAME                 IMAGE                        STATUS
gophkeeper-postgres  postgres:16-alpine           Up (healthy)
gophkeeper-server    gophkeep-gophkeeper-server   Up (healthy)
```

### 3. Проверка работоспособности
```bash
# Health check сервера
curl http://localhost:8080/health
# Response: {"status":"ok"}

# Подключение к БД
psql -h localhost -p 5432 -U gophkeeper -d gophkeeper
# Password: gophkeeper_password
```

### 4. Просмотр логов
```bash
make docker-logs           # Live logs всех сервисов
docker-compose logs postgres           # Только PostgreSQL
docker-compose logs gophkeeper-server  # Только сервер
```

### 5. Остановка сервисов
```bash
make docker-down    # Остановка (данные сохранятся)
make docker-clean   # Полная очистка (включая данные)
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
# Найти процесс на порту
lsof -i :5432   # PostgreSQL
lsof -i :8080   # Server

# Использовать другой порт в docker-compose.yml
# Измените строку: "5432:5432" на "5433:5432"
```

### Container won't start
```bash
# Просмотр ошибок
docker-compose logs postgres
docker-compose logs gophkeeper-server

# Пересборить образ
make docker-build
docker-compose up -d
```

### Connection to DB refused
```bash
# Проверить что PostgreSQL healthy
docker-compose ps postgres

# Проверить DATABASE_URL
docker-compose exec gophkeeper-server env | grep DATABASE_URL

# Проверить логи БД
docker-compose logs postgres | tail -50
```

## 📚 Documentation

- **DOCKER_README.md** - Introduction & Quick Start
- **DOCKER_QUICK_START.md** - Command Checklist
- **DOCKER_SETUP.md** - Detailed Guide
- **.env.example** - Environment Variables Template

И запустите:
```bash
docker-compose up -d
```

## Ручное подключение к БД

```bash
psql -h localhost -U gophkeeper -d gophkeeper
```

Пароль: `gophkeeper_password`

## Проверка здоровья

Сервер имеет health check, проверьте статус:

```bash
docker-compose ps
```

## Данные

PostgreSQL данные сохраняются в Docker volume `postgres_data`. При `docker-compose down -v` volume будет удален.

## Локальное тестирование

После запуска контейнеров, сервер доступен по адресу:
- API: http://localhost:8080
- Health check: http://localhost:8080/health
