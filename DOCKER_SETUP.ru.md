## 🐳 Docker Compose Setup - Подробное руководство

### Что входит в Docker Compose

**PostgreSQL 16 Service:**
- Container: `gophkeeper-postgres`
- Image: `postgres:16-alpine`
- Port: `5432` (локальный) → `5432` (контейнер)
- Database: `gophkeeper`
- User: `gophkeeper`
- Password: `gophkeeper_password`
- Health check: `pg_isready` каждые 10 сек
- Volume: `postgres_data:/var/lib/postgresql/data`

**GophKeeper Server Service:**
- Container: `gophkeeper-server`
- Build: `Dockerfile` (Debian Bookworm-based)
- Port: `8080` (локальный) → `8080` (контейнер)
- Depends on: `postgres` (service_healthy)
- Health check: `curl http://localhost:8080/health` каждые 30 сек

### Makefile Commands

```bash
make docker-up      # Запуск PostgreSQL + сервер
make docker-down    # Остановка контейнеров
make docker-logs    # Просмотр логов (live)
make docker-build   # Пересборка Docker образа
make docker-clean   # Удаление контейнеров + volumes
```

### Переменные окружения

Сервер получает при запуске:

```yaml
DATABASE_URL: postgres://gophkeeper:gophkeeper_password@postgres:5432/gophkeeper?sslmode=disable
SERVER_ADDRESS: 0.0.0.0:8080
SECRET_KEY: your-secret-key-here
LOG_LEVEL: info
```

### Статус контейнеров

```bash
# Показать статус всех контейнеров (с health check)
docker-compose ps

# Ожидаемый вывод:
# gophkeeper-postgres  postgres:16-alpine  Up (healthy)
# gophkeeper-server    gophkeep-*          Up (health: starting/healthy)
```
Управление данными

**Персистентность:**
```bash
# Остановить контейнеры (данные сохранятся)
docker-compose down

# Перезапустить (данные восстановятся)
docker-compose up -d
```

**Удаление данных:**
```bash
# ВАЖНО: Это удалит все данные БД!
docker-compose down -v
```

### Подключение к БД

**Из контейнера:**
```bash
docker-compose exec postgres psql -U gophkeeper -d gophkeeper
```

**С хоста (если установлен psql):**
```bash
psql -h localhost -p 5432 -U gophkeeper -d gophkeeper
# Password: gophkeeper_password
```

**ТиHealth Checks

```bash
# Server health check
curl http://localhost:8080/health
# Response: {"status":"ok"}

# Регистрация пользователя
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# Логин
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
curl http://localhost:8080/health

# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"password123"}'
```

### Логирование и debug

```bash
# Просмотр логов сервера
docker-compose logs gophkeeper-server

# Просмотр логов PostgreSQL
docker-compose logs postgres

# Live logs (следить за изменениями)
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f gophkeeper-server
```
что занимает порт
lsof -i :5432

# Вариант 1: Остановить локальный PostgreSQL
sudo systemctl stop postgresql  # Linux
brew services stop postgresql   # macOS

# Вариант 2: Использовать другой порт
# Отредактируйте docker-compose.yml, измените:
# ports:
#   - "5433:5432"  # используем 5433 вместо 5432
```

**Порт 8080 занят:**
```bash
# Проверить что занимает порт
lsof -i :8080

# Отредактируйте docker-compose.yml, измените:
# ports:
#   - "8081:8080"  # используем 8081 вместо 8080
```

**Контейнер PostgreSQL не становится healthy:**
```bash
# Просмотр логов
docker-compose logs postgres

# Пересборить контейнер
docker-compose down -v
docker-compose up -d
```

**Контейнер сервера падает:**
```bash
# Просмотр логов ошибок
docker-compose logs gophkeeper-server

# Проверить DATABASE_URL
docker-compose exec gophkeeper-server env | grep DATABASE_URL

# Пересобрать образ
make docker-build
docker-compose up -d
```

**Ошибка подключения к БД:**
```bash
# Проверить что PostgreSQL healthy
docker-compose ps postgres

# Проверить переменные окружения
docker-compose config | grep -A 10 gophkeeper-server

# Проверить что они правильно переданы
docker-compose exec gophkeeper-server env | grep DATABASE
```bash
# Проверить какой процесс занимает порт
lsof -i :5432

# Использовать другой порт в docker-compose.yml
# Измените строку: "5432:5432" на "5433:5432"
```

**Ошибка подключения сервера к БД:**
- Убедитесь, что PostgreSQL контейнер `healthy`
- Проверьте логи: `docker-compose logs postgres`
- Проверьте переменную `DATABASE_URL` в compose файле

**Хочу использовать свой пароль для PostgreSQL:**
Отредактируйте `docker-compose.yml` или используйте переменные окружения:

```bash
POSTGRES_PASSWORD=my-secure-password docker-compose up -d
```
