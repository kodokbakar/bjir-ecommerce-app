# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2026-06-07

### Added

- JWT authentication with configurable expiry
- Role-based access control (customer/admin)
- Category CRUD with slug-based lookup
- Product CRUD with categories, images, stock tracking
- Shopping cart with stock validation and quantity management
- Order management with checkout flow and status tracking
- Mock payment endpoint with status lifecycle
- Product search with category filters, pagination, and sorting
- Redis caching for product listings with cache invalidation
- Request logging middleware with slog (auto log-level)
- Panic recovery middleware with JSON error responses
- Input validation on all endpoints
- Global error handler with standardized error codes
- Swagger/OpenAPI auto-generated documentation
- Docker multi-stage build with Alpine
- Docker Compose for local development (app + PostgreSQL + Redis)
- GitHub Actions CI pipeline (go vet + go test)
- Deployed to Railway with auto-deploy on push to main
