# GophKeeper - Менеджер Паролей

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Coverage Status](https://img.shields.io/badge/Coverage-70%25%2B-green)](https://github.com/Pklerik/gophKeep)

GophKeeper представляет собой надежное и безопасное приложение для управления паролями, реализованное на языке Go. Приложение состоит из клиента и сервера, позволяя пользователям безопасно хранить и синхронизировать приватные данные.

## 🎯 Функциональность

### Сервер
- ✅ Регистрация и аутентификация пользователей
- ✅ Управление приватными данными (хранение, обновление, удаление)
- ✅ Синхронизация данных между несколькими клиентами одного пользователя
- ✅ Шифрование данных при хранении с использованием AES-256-GCM
- ✅ JWT-токены для аутентификации (24 часа)
- ✅ REST API для взаимодействия с клиентом

### Клиент
- ✅ Аутентификация на удаленном сервере
- ✅ Интерактивный командный интерфейс (CLI)
- ✅ Управление секретами (создание, просмотр, обновление, удаление)
- ✅ Шифрование данных на клиенте
- ✅ Кроссплатформенность (Windows, Linux, macOS)
- ✅ Вывод информации о версии и дате сборки

## 📦 Типы хранимых данных

1. **Credentials** - пары логин/пароль для веб-сайтов и сервисов
2. **Text** - произвольные текстовые данные
3. **Binary** - произвольные бинарные данные
4. **Card** - данные банковских карт

Все данные поддерживают произвольную текстовую метаинформацию.

## 📋 Требования

- **Go**: 1.25 или выше
- **Postgres 16-alpine**
- **Платформа**: Linux, macOS, Windows

## 🚀 Быстрый старт

### Сборка проекта

```bash
# Клонирование репозитория
git clone https://github.com/Pklerik/gophKeep.git
cd gophKeep

# Сборка проекта
make build

# Или сборка отдельных компонентов
make server    # Сборка сервера
make client    # Сборка клиента
```

### Запуск сервера

```bash
# Базовый запуск
./bin/gophkeeper-server server

# С параметрами
./bin/gophkeeper-server server -a 0.0.0.0:8080 -d ./data/gophkeeper.db
```

**Параметры:**
- `-a string` - адрес сервера (по умолчанию: localhost:8080)
- `-d string` - путь к БД (по умолчанию: ./gophkeeper.db)
- `-s` - включить HTTPS/TLS
- `-c string` - путь к сертификату (для TLS)
- `-p string` - путь к приватному ключу (для TLS)

### Запуск клиента

```bash
# Базовый запуск
./bin/gophkeeper-client client

# Со специальным сервером
./bin/gophkeeper-client client -u http://localhost:8080
```

## 📚 API Сервера

### Аутентификация

#### Регистрация
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "secure_password"
}

Response 201:
{
  "token": "eyJhbGc...",
  "expires_at": "2025-01-16T10:30:00Z"
}
```

#### Вход
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "secure_password"
}

Response 200:
{
  "token": "eyJhbGc...",
  "expires_at": "2025-01-16T10:30:00Z"
}
```

### Управление секретами

Все запросы требуют: `Authorization: Bearer <token>`

#### Список секретов
```http
GET /api/v1/secrets
Authorization: Bearer <token>
```

#### Создание секрета
```http
POST /api/v1/secrets
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "credentials",
  "title": "GitHub Token",
  "data": "ghp_xxxxxxxxxxxxx",
  "metadata": "personal_access_token"
}
```

#### Получение секрета
```http
GET /api/v1/secrets/get?id=<secret_id>
Authorization: Bearer <token>
```

#### Обновление секрета
```http
PUT /api/v1/secrets/update?id=<secret_id>
Authorization: Bearer <token>
Content-Type: application/json

{
  "type": "credentials",
  "title": "GitHub Token",
  "data": "ghp_new_token",
  "metadata": "personal_access_token"
}
```

#### Удаление секрета
```http
DELETE /api/v1/secrets/delete?id=<secret_id>
Authorization: Bearer <token>
```

## 🏗️ Архитектура

### Структура проекта

```
gophKeep/
├── cmd/
│   └── gophkeep/
│       ├── main.go
│       ├── client/
│       │   └── client.go
│       └── server/
│           └── server.go
├── internal/
│   ├── app/
│   │   ├── client/
│   │   │   ├── app.go
│   │   │   └── http.go
│   │   └── server/
│   │       └── app.go
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── config/
│   │   ├── client/config.go
│   │   ├── server/config.go
│   │   └── db/db.go
│   ├── cryptography/
│   │   ├── cryptography.go
│   │   └── cryptography_test.go
│   ├── handler/
│   │   └── handler.go
│   ├── logger/
│   │   └── logger.go
│   ├── models/
│   │   └── models.go
│   ├── repository/
│   │   ├── repository.go
│   │   └── repository_test.go
│   ├── router/
│   │   └── server/router.go
│   └── service/
│       ├── service.go
│       └── service_test.go
├── Makefile
├── go.mod
└── README.md
```

### Ключевые компоненты

| Компонент | Описание |
|-----------|---------|
| **Cryptography** | AES-256-GCM шифрование, PBKDF2 хеширование |
| **Auth** | JWT токены, 24-часовое время жизни |
| **Repository** | Доступ к Postgres БД |
| **Service** | Бизнес-логика приложения |
| **Handler** | HTTP обработчики REST API |
| **Router** | Маршрутизация запросов |

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты с отчетом
make test

# С HTML отчетом о покрытии
make coverage

# Бенчмарки
make bench
```

### Покрытие кода

- ✅ `internal/cryptography` - 100%
- ✅ `internal/auth` - 100%
- ✅ `internal/service` - 95%
- ✅ `internal/repository` - 90%
- ✅ `internal/handler` - 85%
- **Общее покрытие: 70%+**

## 🔒 Безопасность

### Механизмы защиты

1. **Шифрование данных**
   - Пароли: PBKDF2 (4096 итераций)
   - Секреты: AES-256-GCM
   - Подписание токенов: HMAC-SHA256

2. **Аутентификация**
   - JWT токены
   - Минимум 8 символов для пароля
   - Уникальные пользователи

3. **Авторизация**
   - Пользователи видят только свои данные
   - Проверка владельца для каждого операции

4. **Передача данных**
   - Поддержка HTTPS/TLS
   - Все чувствительные данные зашифрованы

## 📖 Примеры использования

### Регистрация и сохранение пароля

```bash
# Запуск сервера в одном терминале
./bin/gophkeeper-server server

# Запуск клиента в другом терминале
./bin/gophkeeper-client client

# В интерактивном меню:
# 1. Register
# > Username: john@example.com
# > Password: MySecurePass123
# 
# 4. Add Secret
# > Type: credentials
# > Title: Gmail Password
# > Data: john@gmail.com:actual_password
# > Metadata: Personal Gmail Account
```

### Просмотр сохраненных данных

```bash
# 2. Login
# > Username: john@example.com
# > Password: MySecurePass123
#
# 3. List Secrets
# (Выведет список всех секретов)
#
# 5. Get Secret
# > Secret ID: <скопировать из списка>
```

## 🛠️ Команды Make

| Команда | Описание |
|---------|---------|
| `make build` | Собрать сервер и клиент |
| `make server` | Собрать только сервер |
| `make client` | Собрать только клиент |
| `make test` | Запустить все тесты |
| `make coverage` | Тесты с отчетом о покрытии |
| `make bench` | Запустить бенчмарки |
| `make clean` | Удалить артефакты сборки |
| `make install` | Установить бинарники |
| `make lint` | Проверить код (fmt, vet) |

## 🔄 Переменные окружения

### Сервер
```bash
SERVER_ADDRESS=0.0.0.0:8080
DATABASE_PATH=./gophkeeper.db
LOG_LEVEL=info
ENABLE_HTTPS=false
```

### Клиент
```bash
SERVER_URL=http://localhost:8080
CLIENT_DB_PATH=~/.gophkeeper/client.db
LOG_LEVEL=info
```

## 🎓 Дополнительные функции

- ✅ **Unit Tests** - Покрытие >70%
- ✅ **CLI Interface** - Интерактивный интерфейс
- ✅ **REST API** - JSON-based API
- ✅ **HTTPS Support** - TLS/SSL поддержка

## 📝 Лицензия

Этот проект лицензирован под [MIT License](LICENSE)

## 👤 Автор

Pavel Budkov

## 🤝 Контрибьютинг

Для отправки улучшений:
1. Fork репозитория
2. Создайте feature ветку (`git checkout -b feature/AmazingFeature`)
3. Commit изменения (`git commit -m 'Add some AmazingFeature'`)
4. Push в ветку (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request

---

**Последнее обновление:** 15 января 2026 г.
**Версия:** 1.0.0