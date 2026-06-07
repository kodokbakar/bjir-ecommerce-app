<div align="center">

# Bjir E-Commerce API

A production-ready REST API backend for e-commerce platforms, built with clean architecture principles.

[![Go CI](https://github.com/kodokbakar/bjir-ecommerce-app/actions/workflows/ci.yml/badge.svg)](https://github.com/kodokbakar/bjir-ecommerce-app/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## Features

| Feature | Description |
|---------|-------------|
| JWT Authentication | Access token with configurable expiry |
| Role-Based Access Control | Customer & Admin role separation |
| Category Management | CRUD with slug-based lookup |
| Product Management | CRUD with categories, images, stock tracking |
| Shopping Cart | Add, update, remove items with stock validation |
| Order Management | Checkout flow with status tracking |
| Mock Payment | Payment simulation with status lifecycle |
| Search & Filter | Product search with category filters + pagination |
| Redis Caching | Product listing cache with invalidation on writes |
| Request Logging | Structured logging with slog, auto log-level |
| Panic Recovery | Global error handler with JSON responses |
| Input Validation | Request payload validation on all endpoints |
| Swagger Docs | Auto-generated OpenAPI documentation |
| Docker Ready | Multi-stage build + docker-compose |
| CI/CD | GitHub Actions (go vet + go test) |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | [Go 1.26.3](https://go.dev/) |
| HTTP Router | [Gin](https://github.com/gin-gonic/gin) |
| Database | [PostgreSQL 17](https://www.postgresql.org/) |
| DB Driver | [pgx/v5](https://github.com/jackc/pgx) |
| Cache | [Redis 7](https://redis.io/) via [go-redis/v9](https://github.com/redis/go-redis) |
| Auth | [golang-jwt/v5](https://github.com/golang-jwt/jwt) + bcrypt |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Docs | [Swaggo](https://github.com/swaggo/swag) (Swagger/OpenAPI) |
| Migration | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Container | [Docker](https://www.docker.com/) (multi-stage, Alpine) |
| CI | [GitHub Actions](https://github.com/features/actions) |

---

## Architecture

```
cmd/server/main.go          Entry point, dependency wiring
        │
internal/
├── router/                 Gin route registration
├── middleware/              JWT, RBAC, recovery, request logger
├── handlers/               HTTP request/response handlers
├── services/               Business logic layer
├── repository/             Database queries (pgx)
├── models/                 Domain entities
├── auth/                   JWT manager + password hashing
├── config/                 Environment-based configuration
├── database/               PostgreSQL + Redis connections
├── response/               Standardized response helpers
└── docs/                   Swagger generated code
```

Clean architecture: **Handler → Service → Repository → Database**. Each layer depends only on the layer below it.

---

## Getting Started

### Prerequisites

- Go 1.26.3+
- PostgreSQL 16+
- Redis 7+
- Docker & Docker Compose (recommended)

### Quick Start (Docker)

```bash
git clone https://github.com/kodokbakar/bjir-ecommerce-app.git
cd bjir-ecommerce-app
docker compose up -d
curl http://localhost:8080/health
```

### Manual Setup

```bash
git clone https://github.com/kodokbakar/bjir-ecommerce-app.git
cd bjir-ecommerce-app
cp .env.example .env
docker compose up -d postgres redis
go mod download
go run ./cmd/server
```

Server starts at `http://localhost:8080`. Swagger UI at `http://localhost:8080/swagger/index.html`.

---

## API Reference

### Endpoints

#### System

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/health` | No | Health check |
| `HEAD` | `/health` | No | Health check (HEAD) |
| `GET` | `/swagger/*any` | No | Swagger UI |

#### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/auth/register` | No | Register new user |
| `POST` | `/api/v1/auth/login` | No | Login & get JWT token |

#### User

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/me` | Yes | Get current user profile |

#### Categories

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| `GET` | `/api/v1/categories` | No | — | List all categories |
| `GET` | `/api/v1/categories/:id` | No | — | Get category by ID |
| `GET` | `/api/v1/categories/slug/:slug` | No | — | Get category by slug |
| `POST` | `/api/v1/categories` | Yes | Admin | Create category |
| `PUT` | `/api/v1/categories/:id` | Yes | Admin | Update category |
| `DELETE` | `/api/v1/categories/:id` | Yes | Admin | Delete category |

#### Products

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| `GET` | `/api/v1/products` | No | — | List products (search, filter, paginate) |
| `GET` | `/api/v1/products/:id` | No | — | Get product by ID |
| `GET` | `/api/v1/products/slug/:slug` | No | — | Get product by slug |
| `POST` | `/api/v1/products` | Yes | Admin | Create product |
| `POST` | `/api/v1/products/:id/image` | Yes | Admin | Upload product image |
| `PUT` | `/api/v1/products/:id` | Yes | Admin | Update product |
| `DELETE` | `/api/v1/products/:id` | Yes | Admin | Delete product |

#### Cart

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/cart` | Yes | View cart with total |
| `POST` | `/api/v1/cart/items` | Yes | Add item to cart |
| `PUT` | `/api/v1/cart/items/:id` | Yes | Update cart item quantity |
| `DELETE` | `/api/v1/cart/items/:id` | Yes | Remove cart item |

#### Orders

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/orders/checkout` | Yes | Checkout cart → order |
| `GET` | `/api/v1/orders` | Yes | List my orders |
| `GET` | `/api/v1/orders/:id` | Yes | Get order detail |
| `PATCH` | `/api/v1/admin/orders/:id/status` | Yes (Admin) | Update order status |

#### Payments

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/payments/pay` | Yes | Mock payment |

### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "not_found",
    "message": "resource not found"
  }
}
```

| Code | Status | Description |
|------|--------|-------------|
| `bad_request` | 400 | Invalid request body |
| `unauthorized` | 401 | Missing or invalid token |
| `forbidden` | 403 | Insufficient permissions |
| `not_found` | 404 | Resource not found |
| `conflict` | 409 | Resource already exists |
| `internal_server_error` | 500 | Server error |

---

## Database

### ER Diagram

```
users ─────┬──────── categories
           │              │
           │           products
           │              │
         carts ───────────┘
           │
         orders ──── order_items
           │
         payments
```

### Tables

| Table | Description |
|-------|-------------|
| `users` | User accounts with roles (customer/admin) |
| `categories` | Product categories with slug-based lookup |
| `products` | Product catalog with pricing, stock, and image |
| `carts` | Shopping cart items (1 row per product per user) |
| `orders` | Order headers with status tracking |
| `order_items` | Individual items within an order (price snapshot) |
| `payments` | Payment records with status lifecycle |

### Order Status Flow

```
pending → paid → shipped → delivered
   │
   └→ canceled
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
| `DB_PASSWORD` | *(required)* | PostgreSQL password |
| `DB_NAME` | `go_ecommerce_api` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | *(empty)* | Redis password |
| `JWT_SECRET` | *(required)* | JWT signing secret |
| `JWT_EXPIRES_IN` | `24h` | Token expiry duration |

> Full list of connection pool and timeout variables available in `.env.example`.

---

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Test Coverage

| Package | Tests |
|---------|-------|
| `internal/auth` | JWT generation & validation |
| `internal/config` | Config loading |
| `internal/database` | Connection setup |
| `internal/handlers` | Auth, category, product, cart, order, payment handlers |
| `internal/middleware` | Auth, RBAC, recovery, request logger |
| `internal/repository` | User, category, product, cart, order, payment repositories |
| `internal/response` | Response format helpers |
| `internal/router` | Route registration + unknown routes |
| `internal/services` | Auth, category, product, cart, order, payment services |

---

## CI/CD

GitHub Actions pipeline runs on every push and pull request to `main`:

```
Pipeline: Go CI
├── Setup Go 1.26.3
├── Download dependencies
├── go mod verify
├── go vet ./...
└── go test ./... -v
```

---

## Docker

### Build & Run

```bash
docker build -t bjir-ecommerce-app .
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e JWT_SECRET=your_secret \
  bjir-ecommerce-app
```

### Docker Compose

| Service | Image | Port |
|---------|-------|------|
| `app` | Built from Dockerfile | `8080` |
| `postgres` | `postgres:17-alpine` | `5433` |
| `redis` | `redis:7-alpine` | `6380` |

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

[@kodokbakar](https://github.com/kodokbakar)

</div>
