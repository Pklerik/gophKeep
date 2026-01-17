# GophKeeper - Password Manager

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Coverage Status](https://img.shields.io/badge/Coverage-70%25%2B-green)](https://github.com/Pklerik/gophKeep)

GophKeeper is a reliable and secure password management application implemented in Go. The application consists of a client and a server, allowing users to securely store and synchronize private data.

## 🎯 Features

### Server
- ✅ User registration and authentication
- ✅ Management of private data (storage, update, deletion)
- ✅ Data synchronization between multiple clients of the same user
- ✅ Data encryption at rest using AES-256-GCM
- ✅ JWT tokens for authentication (24 hours)
- ✅ REST API for client interaction

### Client
- ✅ Authentication with the remote server
- ✅ Interactive Command Line Interface (CLI)
- ✅ Secret management (create, view, update, delete)
- ✅ Client-side data encryption
- ✅ Cross-platform (Windows, Linux, macOS)
- ✅ Displays version and build date information

## 📦 Supported Data Types

1. **Credentials** - login/password pairs for websites and services
2. **Text** - arbitrary text data
3. **Binary** - arbitrary binary data
4. **Card** - bank card details

All data supports arbitrary text metadata.

## 📋 Requirements

- **Go**: 1.25 or higher
- **SQLite3**
- **Platform**: Linux, macOS, Windows

## 🚀 Quick Start

### Build the Project

```bash
# Clone the repository
git clone https://github.com/Pklerik/gophKeep.git
cd gophKeep

# Build the project
make build

# Or build individual components
make server    # Build the server
make client    # Build the client
```

### Start the Server

```bash
# Basic start
./bin/gophkeeper-server server

# With parameters
./bin/gophkeeper-server server -a 0.0.0.0:8080 -d ./data/gophkeeper.db
```

**Parameters:**
- `-a string` - server address (default: localhost:8080)
- `-d string` - path to the database (default: ./gophkeeper.db)
- `-s` - enable HTTPS/TLS
- `-c string` - path to the certificate (for TLS)
- `-p string` - path to the private key (for TLS)

### Start the Client

```bash
# Basic start
./bin/gophkeeper-client client

# With a specific server
./bin/gophkeeper-client client -u http://localhost:8080
```

## 📚 Server API

### Authentication

#### Register
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

#### Login
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

### Secret Management

All requests require: `Authorization: Bearer <token>`

#### List Secrets
```http
GET /api/v1/secrets
Authorization: Bearer <token>
```

#### Create Secret
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

#### Get Secret
```http
GET /api/v1/secrets/get?id=<secret_id>
Authorization: Bearer <token>
```

#### Update Secret
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

#### Delete Secret
```http
DELETE /api/v1/secrets/delete?id=<secret_id>
Authorization: Bearer <token>
```

## 🏗️ Architecture

### Project Structure

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
```