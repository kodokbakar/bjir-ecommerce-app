<div align="center">

# Bjir E-Commerce

**Sharp catalog. Hard shadows. No generic shelf.**

A full-stack e-commerce platform with Go API backend and React frontend.  
Brutalist design, clean architecture, production-ready.

[![Go CI](https://github.com/kodokbakar/bjir-ecommerce-app/actions/workflows/ci.yml/badge.svg)](https://github.com/kodokbakar/bjir-ecommerce-app/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.26.3-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=white)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-6-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite&logoColor=white)](https://vite.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

---

## Overview

Monorepo containing a Go REST API backend (`internal/`) and a React SPA frontend (`frontend/`). The backend handles auth, products, categories, cart, orders, payments, and contact — all with RBAC. The frontend ships as a Vite SPA with landing page, product catalog, shopping cart, checkout flow, and admin dashboard.

| Layer | Stack | Lines |
|-------|-------|-------|
| Backend (Go) | Gin + pgx + go-redis + golang-jwt | 22,801 |
| Frontend (React) | React 19 + TypeScript 6 + Vite 8 | 13,140 |

---

## Features

### Backend

| Feature | Description |
|---------|-------------|
| JWT Auth | Access token with configurable expiry |
| RBAC | Customer & Admin role separation |
| Categories | CRUD with slug-based lookup |
| Products | CRUD with categories, images, stock tracking |
| Shopping Cart | Add, update, remove items with stock validation |
| Orders | Checkout flow with status tracking |
| Payments | Mock payment with status lifecycle |
| Contact | Store contact message submission |
| Search & Filter | Product search with category filters + pagination |
| Redis Caching | Product listing cache with invalidation on writes |
| Structured Logging | slog with auto log-level |
| Validation | Request payload validation on all endpoints |
| Swagger Docs | Auto-generated OpenAPI documentation |
| Docker | Multi-stage build + docker-compose |
| CI/CD | GitHub Actions (go vet + go test) |

### Frontend

| Feature | Description |
|---------|-------------|
| Landing Page | Hero, featured products (8), about, contact form, footer |
| Auth Pages | Login, register with role-aware redirect |
| Product Catalog | Grid/list, search, category filter, pagination |
| Product Detail | Image, price, stock state, add to cart |
| Shopping Cart | Quantity update, remove, checkout |
| Checkout | Address form, order summary, payment |
| Order Management | Order list, detail, status tracking |
| Dashboard | Role-aware: customer orders / admin panel |
| Admin Panel | Products, categories, orders CRUD |
| Responsive | 320px–1440px, hamburger menu, 44px touch targets |
| SEO | Meta tags, OG/Twitter cards, JSON-LD, sitemap |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Language** | [Go 1.26.3](https://go.dev/) |
| **HTTP Router** | [Gin](https://github.com/gin-gonic/gin) |
| **Database** | [PostgreSQL 17](https://www.postgresql.org/) |
| **DB Driver** | [pgx/v5](https://github.com/jackc/pgx) |
| **Cache** | [Redis 7](https://redis.io/) via [go-redis/v9](https://github.com/redis/go-redis) |
| **Auth** | [golang-jwt/v5](https://github.com/golang-jwt/jwt) + bcrypt |
| **Validation** | [go-playground/validator](https://github.com/go-playground/validator) |
| **Docs** | [Swaggo](https://github.com/swaggo/swag) |
| **Migration** | [golang-migrate](https://github.com/golang-migrate/migrate) |
| **Container** | Docker (multi-stage, Alpine) |
| | |
| **Language** | [TypeScript 6](https://www.typescriptlang.org/) |
| **UI Library** | [React 19](https://react.dev/) |
| **Bundler** | [Vite 8](https://vite.dev/) |
| **Styling** | Tailwind CSS 4 + custom brutalist CSS |
| **Icons** | [Lucide React](https://lucide.dev/) |
| **HTTP Client** | [Axios](https://axios-http.com/) |
| **Routing** | [React Router 7](https://reactrouter.com/) |
| **Testing** | [Vitest 4](https://vitest.dev/) + [jsdom](https://github.com/jsdom/jsdom) |
| **CI** | [GitHub Actions](https://github.com/features/actions) |

---

## Architecture

```
bjir-ecommerce-app/
│
├── cmd/server/main.go          Entry point, dependency wiring
│
├── internal/
│   ├── router/                 Gin route registration
│   ├── middleware/             JWT, RBAC, recovery, request logger
│   ├── handlers/               HTTP request/response handlers
│   ├── services/               Business logic layer
│   ├── repository/             Database queries (pgx)
│   ├── models/                 Domain entities
│   ├── auth/                   JWT manager + password hashing
│   ├── config/                 Environment-based configuration
│   ├── database/               PostgreSQL + Redis connections
│   ├── response/               Standardized response helpers
│   └── docs/                   Swagger generated code
│
└── frontend/
    ├── src/
    │   ├── pages/              17 route pages (LandingPage, Products, Cart, etc.)
    │   ├── components/         Reusable UI (ProductImage, Toast, Pagination, etc.)
    │   ├── services/           API client layer per domain (auth, product, cart, etc.)
    │   ├── hooks/              useAuth, useCartCount
    │   ├── context/            AuthContext, ToastContext
    │   └── utils/              formatRupiah, getStockState, authRouting, etc.
    ├── public/                 Static assets: OG image, sitemap.xml, robots.txt
    └── test/                   Test setup + vitest config
```

Clean architecture — **Handler → Service → Repository → Database**.  
Frontend follows — **Page → Component → Service → API**.

---

## Landing Page

The landing page ships as part of the SPA and covers the full customer journey:

| Section | Notes |
|---------|-------|
| **Navbar** | Auth-aware desktop + mobile hamburger, Escape key |
| **Hero** | Auth-aware CTAs, loading placeholders, entrance animation |
| **Featured Products** | 8 newest from API, stock states (in/low/out), skeleton loading |
| **About** | 4 value propositions (Pengiriman Cepat, Pembayaran Aman, etc.) |
| **Contact Form** | Name/email/message, validation, toast feedback, POST `/v1/contact` |
| **Footer** | Quick links, social links (GitHub, Email), copyright auto-year |
| **SEO** | Meta tags, OG/Twitter cards, JSON-LD Structured Data, sitemap.xml |

---

## Getting Started

### Prerequisites

- Go 1.26.3+
- PostgreSQL 17+
- Redis 7+
- Node.js 22+
- Docker & Docker Compose (recommended)

### Quick Start (Docker — Full Stack)

```bash
git clone https://github.com/kodokbakar/bjir-ecommerce-app.git
cd bjir-ecommerce-app
docker compose up -d
curl http://localhost:8080/health
```

### Frontend Dev

```bash
cd frontend
npm install
npm run dev          # → http://localhost:5173
npm run test:run     # 42 tests, 7 files
```

### Backend Dev

```bash
cp .env.example .env
docker compose up -d postgres redis
go mod download
go run ./cmd/server  # → http://localhost:8080
```

### Seeder

```bash
go run cmd/seed/main.go
```

### Production Build

```bash
cd frontend
npm run build          # → dist/
npm run start          # serve on PORT
```

### Deployment

Live API: [bjir-ecommerce-app-production.up.railway.app](https://bjir-ecommerce-app-production.up.railway.app)  
Swagger UI: [bjir-ecommerce-app-production.up.railway.app/swagger/index.html](https://bjir-ecommerce-app-production.up.railway.app/swagger/index.html)  
Frontend: [frontend-bjir-ecommerce-production.up.railway.app](https://frontend-bjir-ecommerce-production.up.railway.app) (root)

---

## Testing

### Backend (Go)

```bash
go test ./...
go test ./... -v
go test ./... -coverprofile=coverage.out
```

10 packages — auth, config, database, handlers, middleware, repository, response, router, server, services.

### Frontend (React)

```bash
cd frontend
npm run test:run
```

7 test files — 42 tests total (unit + integration).

---

## API Reference

### System

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/health` | No | Health check |
| `HEAD` | `/health` | No | Health check (HEAD) |
| `GET` | `/swagger/*any` | No | Swagger UI |

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/auth/register` | No | Register new user |
| `POST` | `/api/v1/auth/login` | No | Login & get JWT token |

### User

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/me` | Yes | Get current user profile |

### Categories

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| `GET` | `/api/v1/categories` | No | — | List all categories |
| `GET` | `/api/v1/categories/:id` | No | — | Get category by ID |
| `GET` | `/api/v1/categories/slug/:slug` | No | — | Get category by slug |
| `POST` | `/api/v1/categories` | Yes | Admin | Create category |
| `PUT` | `/api/v1/categories/:id` | Yes | Admin | Update category |
| `DELETE` | `/api/v1/categories/:id` | Yes | Admin | Delete category |

### Products

| Method | Endpoint | Auth | Role | Description |
|--------|----------|------|------|-------------|
| `GET` | `/api/v1/products` | No | — | List products (search, filter, paginate) |
| `GET` | `/api/v1/products/:id` | No | — | Get product by ID |
| `GET` | `/api/v1/products/slug/:slug` | No | — | Get product by slug |
| `POST` | `/api/v1/products` | Yes | Admin | Create product |
| `POST` | `/api/v1/products/:id/image` | Yes | Admin | Upload product image |
| `PUT` | `/api/v1/products/:id` | Yes | Admin | Update product |
| `DELETE` | `/api/v1/products/:id` | Yes | Admin | Delete product |

### Cart

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/cart` | Yes | View cart with total |
| `POST` | `/api/v1/cart/items` | Yes | Add item to cart |
| `PUT` | `/api/v1/cart/items/:id` | Yes | Update cart item quantity |
| `DELETE` | `/api/v1/cart/items/:id` | Yes | Remove cart item |

### Orders

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/orders/checkout` | Yes | Checkout cart → order |
| `GET` | `/api/v1/orders` | Yes | List my orders |
| `GET` | `/api/v1/orders/:id` | Yes | Get order detail |
| `PATCH` | `/api/v1/admin/orders/:id/status` | Yes (Admin) | Update order status |

### Payments

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/payments/pay` | Yes | Mock payment |

### Contact

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/v1/contact` | No | Submit contact message |

### Error Response

```json
{
  "success": false,
  "message": "product not found",
  "error": {
    "code": "not_found",
    "message": "product not found",
    "details": null
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

## Ponytail Notes

Intentional ceilings — things we chose not to build yet.

| Area | Ceiling | Upgrade Path |
|------|---------|-------------|
| **i18n** | Indonesian only (landing, SEO, contact) | react-helmet-async for dynamic meta, i18n for EN |
| **SSR/SSG** | Vite static SPA | Next.js/Remix for product page SEO |
| **CMS** | Hardcoded landing content | CMS integration for dynamic sections |
| **OG Image** | Static SVG | Dynamic OG per page (@vercel/og) |
| **Pagination** | Landing: 8 items, no pagination | Full pagination if landing becomes secondary catalog |
| **Animations** | CSS keyframes only | Framer Motion / GSAP for complex sequences |
| **Carousel** | Static hero, no slider | Carousel on desgin request |
| **Image Opt** | Single srcset | `<picture>` + multiple densities |
| **State Mgmt** | React context, no Redux | Zustand/Redux if cross-component state grows |
| **API Pagination** | cursor-based on orders, offset on products | Unify to cursor-based |
| **Rate Limiting** | Not implemented | Middleware + Redis rate limiter |
| **Email** | Not implemented | Mailgun/SendGrid for order confirmations |
| **File Upload** | Local filesystem | S3/GCS for production |

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
| `DB_PASSWORD` | _(required)_ | PostgreSQL password |
| `DB_NAME` | `go_ecommerce_api` | Database name |
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `JWT_SECRET` | _(required)_ | JWT signing secret |
| `JWT_EXPIRES_IN` | `24h` | Token expiry duration |

Full list of connection pool and timeout variables in `.env.example`.

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

## Contributors

- [@kodokbakar](https://github.com/kodokbakar) — Backend Developer, DevOps, Frontend Developer, Everything On This project
- [@bintangrobbany](https://github.com/bintangrobbany) — Frontend Developer

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

[@kodokbakar](https://github.com/kodokbakar) · [@bintangrobbany](https://github.com/bintangrobbany)

</div>
