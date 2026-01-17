# GophKeeper Docker & PostgreSQL Deployment

## 📋 Что было добавлено

Полная поддержка Docker Compose с PostgreSQL для запуска сервера в контейнере.

### Файлы
- **docker-compose.yml** - Конфиг для запуска PostgreSQL и GophKeeper Server
- **Dockerfile** - Multi-stage build (Debian-based) для сервера
- **.env.example** - Пример переменных окружения
- **Makefile** - добавлены команды `docker-*`

## 🚀 Быстрый старт

### 1️⃣ Убедитесь что установлены Docker и Docker Compose
```bash
docker --version      # >= 20.10
docker-compose --version  # >= 1.29
```

### 2️⃣ Запустите сервис
```bash
make docker-up
```

Это запустит:
- **PostgreSQL 16** на локальном порту **5432**
- **GophKeeper Server** на локальном порту **8080**

### 3️⃣ Проверьте статус
```bash
docker-compose ps
```

Оба сервиса должны быть в статусе `Up` (PostgreSQL - `healthy`)

### 4️⃣ Протестируйте
```bash
# Проверить сервер
curl http://localhost:8080/health

# Подключиться к БД
psql -h localhost -U gophkeeper -d gophkeeper
# Password: gophkeeper_password
```

## � Docker Services

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
  - pg_isready каждые 10 секунд
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
  - curl http://localhost:8080/health каждые 30 секунд
```

## 📝 Команды Makefile

```bash
make docker-up       # Запустить PostgreSQL + сервер
make docker-down     # Остановить контейнеры
make docker-logs     # Просмотр логов в реальном времени
make docker-build    # Пересобрать Docker образ
make docker-clean    # Удалить контейнеры и volumes
```

## ⚙️ Конфигурация для Production

1. **Скопируйте .env**:
   ```bash
   cp .env.example .env
   ```

2. **Обновите credentials в .env**:
   ```env
   SECRET_KEY=<generate-with: openssl rand -base64 32>
   POSTGRES_PASSWORD=<strong-password>
   ```

3. **Запустите с .env**:
   ```bash
   docker-compose --env-file .env up -d
   ```

## 🔍 Troubleshooting

### Сервер не подключается к БД
```bash
# Проверьте логи
docker-compose logs gophkeeper-server

# Убедитесь PostgreSQL здоров
docker-compose logs postgres | grep "ready"
```

### Порт занят
```bash
# Пример: меняем порт PostgreSQL на 5433
docker-compose.yml:
  postgres:
    ports:
      - "5433:5432"
```

### Хочу очистить всё и заново
```bash
make docker-clean
make docker-up
```

## 📊 Мониторинг

```bash
# Живые логи всех сервисов
make docker-logs

# Логи конкретного сервиса
docker-compose logs -f gophkeeper-server
docker-compose logs -f postgres

# Статус контейнеров
docker-compose ps
```

## 🔗 Endpoints для тестирования

```bash
# Health check
curl http://localhost:8080/health

# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Логин
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

## 📚 Документация

- **DOCKER_QUICK_START.md** - Начните отсюда
- **DOCKER_SETUP.md** - Детальные инструкции
- **DOCKER.md** - Общая информация

## 💾 Данные

PostgreSQL данные сохраняются в Docker volume `postgres_data`. При `docker-compose down -v` volume будет удален.

Для персистентности между перезапусками используйте `docker-compose down` (без -v):
```bash
docker-compose down  # данные сохранятся
docker-compose up -d # восстановятся
```

---

**Готово к запуску:**
```bash
make docker-up && make docker-logs &
```
