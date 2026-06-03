<div align="center">

# Bjir E-Commerce App

**Full-stack e-commerce platform** built with Go, Gin, PostgreSQL, Redis, and JWT authentication.

[![Go CI](https://github.com/kodokbakar/bjir-ecommerce-app/actions/workflows/ci.yml/badge.svg)](https://github.com/kodokbakar/bjir-ecommerce-app/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[Features](#-features) • [Tech Stack](#-tech-stack) • [Architecture](#-architecture) • [Getting Started](#-getting-started) • [API Reference](#-api-reference) • [Database Schema](#-database-schema)

</div>

---

## 📌 Overview

Bjir E-Commerce API is a production-ready REST API backend for e-commerce platforms. Built with **clean architecture** principles, it separates concerns into handler, service, repository, and model layers for maintainability and testability.

### Key Features

| Feature | Status | Description |
|---------|--------|-------------|
| 🔐 JWT Authentication | ✅ | Access token with configurable expiry |
| 👥 Role-Based Access | ✅ | Customer & Admin role separation |
| 📦 Product Management | 🚧 | CRUD with categories, images, stock tracking |
| 🛒 Shopping Cart | 🚧 | Add, update, remove items with stock validation |
| 📋 Order Management | 🚧 | Checkout flow with status tracking |
| 💳 Payment Integration | 🚧 | Mock payment endpoint |
| 🔍 Search & Filter | 🚧 | Product search with category/price filters |
| ⚡ Redis Caching | 🚧 | Cache invalidation on writes |
| 📖 Swagger Docs | ✅ | Auto-generated API documentation |
| 🐳 Docker Ready | ✅ | Multi-stage build + docker-compose |
| 🧪 Unit Tests | ✅ | Auth, middleware, handlers, response |
| 🔄 CI/CD | ✅ | GitHub Actions (vet + test) |

> ✅ = Implemented &nbsp; 🚧 = In Progress

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Language** | [Go 1.26.3](https://go.dev/) | Primary language |
| **HTTP Router** | [Gin](https://github.com/gin-gonic/gin) | High-performance HTTP framework |
| **Database** | [PostgreSQL 17](https://www.postgresql.org/) | Primary data store |
| **DB Driver** | [pgx/v5](https://github.com/jackc/pgx) | Native PostgreSQL driver with connection pooling |
| **Cache** | [Redis 7](https://redis.io/) | Session & data caching |
| **Redis Client** | [go-redis/v9](https://github.com/redis/go-redis) | Redis client with connection pooling |
| **Auth** | [golang-jwt/v5](https://github.com/golang-jwt/jwt) | JWT token generation & validation |
| **Password** | [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | Secure password hashing |
| **Migration** | [golang-migrate](https://github.com/golang-migrate/migrate) | Database version control |
| **Validation** | [go-playground/validator](https://github.com/go-playground/validator) | Request payload validation |
| **Docs** | [Swaggo](https://github.com/swaggo/swag) | Swagger/OpenAPI auto-generation |
| **Container** | [Docker](https://www.docker.com/) | Multi-stage build, Alpine-based |
| **CI** | [GitHub Actions](https://github.com/features/actions) | Automated vet & test pipeline |

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      cmd/server/                        │
│                    main.go (entry)                      │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                    internal/                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
│  │   router/   │  │ middleware/ │  │    response/    │ │
│  │ Gin routes  │→ │ JWT + RBAC  │→ │ Standardized    │ │
│  └──────┬──────┘  └─────────────┘  │ response format │ │
│         │                           └─────────────────┘ │
│  ┌──────▼──────┐                                         │
│  │  handlers/  │  HTTP request/response                 │
│  └──────┬──────┘                                         │
│         │                                                │
│  ┌──────▼──────┐                                         │
│  │  services/  │  Business logic                        │
│  └──────┬──────┘                                         │
│         │                                                │
│  ┌──────▼──────┐                                         │
│  │ repository/ │  Database queries                      │
│  └──────┬──────┘                                         │
│         │                                                │
│  ┌──────▼──────┐  ┌─────────────┐  ┌─────────────────┐ │
│  │   models/   │  │   auth/     │  │   database/     │ │
│  │  Entities   │  │ JWT + bcrypt│  │ PG + Redis init │ │
│  └─────────────┘  └─────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

---

## Getting Started

### Prerequisites

- **Go** 1.26.3+
- **PostgreSQL** 16+
- **Redis** 7+
- **Docker** & **Docker Compose** (recommended)

### Quick Start with Docker

```bash
# Clone the repository
git clone https://github.com/kodokbakar/bjir-ecommerce-app.git
cd bjir-ecommerce-app

# Start all services (app + PostgreSQL + Redis)
docker compose up -d

# Verify
curl http://localhost:8080/health
# → {"status":"ok","message":"Go E-Commerce API is running"}
```

### Manual Setup

```bash
# 1. Clone
git clone https://github.com/kodokbakar/bjir-ecommerce-app.git
cd bjir-ecommerce-app

# 2. Copy environment file
cp .env.example .env

# 3. Start PostgreSQL and Redis (or use Docker)
docker compose up -d postgres redis

# 4. Install dependencies
go mod download

# 5. Run migrations
go run ./cmd/server

# 6. Start the server
go run ./cmd/server
```

The server starts at `http://localhost:8080`.

### Verify Installation

```bash
# Health check
curl http://localhost:8080/health

# Swagger docs
open http://localhost:8080/swagger/index.html
```

---

## API Reference

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

All protected routes require a `Bearer` token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

### Endpoints

#### System

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/health` | ❌ | Health check |
| `GET` | `/swagger/*any` | ❌ | Swagger UI |

#### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/auth/register` | ❌ | Register new user |
| `POST` | `/api/v1/auth/login` | ❌ | Login & get JWT token |

#### User

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/me` | ✅ | Get current user profile |

#### Admin

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| `GET` | `/api/v1/admin/ping` | ✅ | Admin | Admin-only endpoint |

### Request Examples

#### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Response** `201 Created`

```json
{
  "success": true,
  "message": "user registered successfully",
  "data": {
    "user_id": "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
    "email": "john@example.com",
    "role": "customer"
  }
}
```

#### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secret123"
  }'
```

**Response** `200 OK`

```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": "24h"
  }
}
```

#### Get Profile

```bash
curl http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Response** `200 OK`

```json
{
  "success": true,
  "message": "current user retrieved successfully",
  "data": {
    "user_id": "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
    "email": "john@example.com",
    "role": "customer"
  }
}
```

### Error Response Format

All errors follow a consistent format:

```json
{
  "success": false,
  "message": "invalid request body",
  "error": {
    "code": "bad_request",
    "details": [
      {
        "field": "email",
        "message": "must be a valid email"
      }
    ]
  }
}
```

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `bad_request` | 400 | Invalid request body |
| `unauthorized` | 401 | Missing or invalid token |
| `forbidden` | 403 | Insufficient permissions |
| `not_found` | 404 | Resource not found |
| `conflict` | 409 | Resource already exists |
| `internal_server_error` | 500 | Server error |

---

## Database Schema

### ER Diagram

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    users     │     │  categories  │     │   products   │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id (UUID PK) │     │ id (UUID PK) │     │ id (UUID PK) │
│ name         │     │ name         │     │ category_id  │──┐
│ email        │     │ slug (UNIQUE)│     │ name         │  │
│ password_hash│     │ description  │     │ slug (UNIQUE)│  │
│ role         │     │ is_active    │     │ description  │  │
│ is_active    │     │ created_at   │     │ price        │  │
│ created_at   │     │ updated_at   │     │ stock        │  │
│ updated_at   │     └──────────────┘     │ image_url    │  │
└──────┬───────┘                          │ is_active    │  │
       │                                  │ created_at   │  │
       │                                  │ updated_at   │  │
       │                                  └──────┬───────┘  │
       │                                         │          │
       │    ┌──────────────┐                     │          │
       │    │    carts     │                     │          │
       │    ├──────────────┤                     │          │
       │    │ id (UUID PK) │                     │          │
       ├───▶│ user_id      │                     │          │
       │    │ product_id   │◀────────────────────┘          │
       │    │ quantity     │                                │
       │    │ created_at   │                                │
       │    │ updated_at   │                                │
       │    └──────────────┘                                │
       │                                                    │
       │    ┌──────────────┐     ┌──────────────┐          │
       │    │   orders     │     │ order_items  │          │
       │    ├──────────────┤     ├──────────────┤          │
       ├───▶│ id (UUID PK) │◀───▶│ id (UUID PK) │          │
       │    │ user_id      │     │ order_id     │          │
       │    │ order_number │     │ product_id   │◀─────────┘
       │    │ status       │     │ product_name │
       │    │ total_amount │     │ quantity     │
       │    │ shipping_addr│     │ price        │
       │    │ notes        │     │ subtotal     │
       │    │ created_at   │     │ created_at   │
       │    │ updated_at   │     └──────────────┘
       │    └──────┬───────┘
       │           │
       │    ┌──────▼───────┐
       │    │   payments   │
       │    ├──────────────┤
       │    │ id (UUID PK) │
       └────│ order_id     │
            │ provider     │
            │ payment_method│
            │ transaction_id│
            │ amount       │
            │ status       │
            │ paid_at      │
            │ created_at   │
            │ updated_at   │
            └──────────────┘
```

### Tables

| Table | Description |
|-------|-------------|
| `users` | User accounts with roles (customer/admin) |
| `categories` | Product categories with slug-based lookup |
| `products` | Product catalog with pricing and stock |
| `carts` | Shopping cart items (1 row per product per user) |
| `orders` | Order headers with status tracking |
| `order_items` | Individual items within an order |
| `payments` | Payment records with status lifecycle |

### Order Status Flow

```
pending → paid → processing → shipped → delivered
    │                                        
    └──────────────────────→ cancelled
```

---

## Environment Variables

Copy `.env.example` to `.env` and configure:

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `development` | Application environment |
| `APP_PORT` | `8080` | Server port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | `postgres` | PostgreSQL password |
| `DB_NAME` | `go_ecommerce_api` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `DB_POOL_MAX_CONNS` | `10` | Max connection pool size |
| `DB_POOL_MIN_CONNS` | `0` | Min connection pool size |
| `DB_POOL_MAX_CONN_LIFETIME` | `1h` | Max connection lifetime |
| `DB_POOL_MAX_CONN_IDLE_TIME` | `30m` | Max idle connection time |
| `DB_POOL_HEALTH_CHECK_PERIOD` | `1m` | Health check interval |
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | *(empty)* | Redis password |
| `REDIS_DB` | `0` | Redis database number |
| `REDIS_POOL_SIZE` | `10` | Redis connection pool size |
| `REDIS_MIN_IDLE_CONNS` | `2` | Min idle connections |
| `REDIS_DIAL_TIMEOUT` | `5s` | Redis dial timeout |
| `REDIS_READ_TIMEOUT` | `3s` | Redis read timeout |
| `REDIS_WRITE_TIMEOUT` | `3s` | Redis write timeout |
| `JWT_SECRET` | *(required)* | JWT signing secret |
| `JWT_EXPIRES_IN` | `24h` | Token expiry duration |
| `JWT_ISSUER` | `go-ecommerce-api` | Token issuer claim |

---

## Project Structure

```
bjir-ecommerce-app/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point, dependency wiring
├── internal/
│   ├── auth/
│   │   ├── jwt.go                  # JWT manager (generate, validate)
│   │   ├── jwt_test.go             # JWT unit tests
│   │   └── password.go             # bcrypt hash & check
│   ├── config/
│   │   └── config.go               # Environment-based configuration
│   ├── database/
│   │   ├── postgres.go             # PostgreSQL connection pool
│   │   ├── redis.go                # Redis client connection
│   │   └── migration.go            # Database migration runner
│   ├── handlers/
│   │   ├── admin_handler.go        # Admin-only endpoints
│   │   ├── auth_handler.go         # Register & login handlers
│   │   ├── auth_handler_test.go    # Auth handler tests
│   │   ├── health_handler.go       # Health check endpoint
│   │   ├── me_handler.go           # Current user profile
│   │   └── response.go             # Response type definitions
│   ├── middleware/
│   │   ├── auth_middleware.go       # JWT + RBAC middleware
│   │   └── auth_middleware_test.go  # Middleware tests
│   ├── models/
│   │   ├── errors.go               # Domain error definitions
│   │   └── user.go                 # User entity
│   ├── repository/
│   │   └── user_repository.go      # User DB operations
│   ├── response/
│   │   ├── response.go             # Standardized response helpers
│   │   └── response_test.go        # Response tests
│   ├── router/
│   │   └── router.go               # Route registration
│   └── services/
│       ├── auth_service.go         # Auth business logic
│       └── auth_service_test.go    # Service tests
├── migrations/
│   ├── 000001_create_ecommerce_tables.up.sql
│   └── 000001_create_ecommerce_tables.down.sql
├── docs/
│   ├── docs.go                     # Swagger generated code
│   ├── swagger.json                # Swagger JSON spec
│   └── swagger.yaml                # Swagger YAML spec
├── .github/
│   └── workflows/
│       └── ci.yml                  # GitHub Actions CI pipeline
├── Dockerfile                      # Multi-stage Docker build
├── docker-compose.yml              # App + PostgreSQL + Redis
├── .env.example                    # Environment template
├── .gitignore                      # Git ignore rules
├── .dockerignore                   # Docker ignore rules
├── go.mod                          # Go module definition
├── go.sum                          # Go dependency checksums
└── PROJECT_BRIEF.md                # Project specification
```

---

## Testing

### Run All Tests

```bash
go test ./... -v
```

### Run Tests with Coverage

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Structure

| Package | Test File | Coverage |
|---------|-----------|----------|
| `internal/auth` | `jwt_test.go` | JWT generation & validation |
| `internal/handlers` | `auth_handler_test.go` | Register & login endpoints |
| `internal/middleware` | `auth_middleware_test.go` | Auth & RBAC middleware |
| `internal/response` | `response_test.go` | Response format helpers |
| `internal/services` | `auth_service_test.go` | Auth business logic |

---

## CI/CD

GitHub Actions pipeline runs on every push and PR to `main`:

```yaml
Pipeline: Go CI
├── Job: vet-and-test
│   ├── Checkout repository
│   ├── Setup Go 1.26.3
│   ├── Download dependencies
│   ├── Verify dependencies
│   ├── Run go vet
│   └── Run tests (verbose)
```

### Status Checks

- ✅ `go vet ./...` — No static analysis issues
- ✅ `go test ./...` — All tests pass
- ✅ `go mod verify` — Dependencies verified

---

## Docker

### Build Image

```bash
docker build -t bjir-ecommerce-app .
```

### Run Container

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e JWT_SECRET=your_secret \
  bjir-ecommerce-app
```

### Docker Compose Services

| Service | Image | Port | Purpose |
|---------|-------|------|---------|
| `app` | Built from Dockerfile | `8080` | Go API server |
| `postgres` | `postgres:17-alpine` | `5433` | PostgreSQL database |
| `redis` | `redis:7-alpine` | `6380` | Redis cache |

---

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

---

## Author

@kodokbakar](https://github.com/kodokbakar)

---

<div align="center">

**Built with ❤️**

</div>
