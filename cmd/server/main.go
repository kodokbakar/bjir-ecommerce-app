package main

import (
	"context"
	"log"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/database"
	"github.com/kodokbakar/go-ecommerce-api/internal/handlers"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
	"github.com/kodokbakar/go-ecommerce-api/internal/router"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

// @title Go E-Commerce API
// @version 1.0
// @description REST API for e-commerce platform built with Go, Gin, PostgreSQL, Redis, and JWT authentication.
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.RunMigrations(cfg.Database); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	ctx := context.Background()

	dbPool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}
	defer dbPool.Close()

	log.Println("PostgreSQL connected successfully")

	redisClient, err := database.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
		}
	}()

	log.Println("Redis connected successfully")

	jwtManager := auth.NewJWTManager(cfg.JWT)

	userRepository := repository.NewUserRepository(dbPool)
	authService := services.NewAuthService(userRepository, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	categoryRepository := repository.NewCategoryRepository(dbPool)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	r := router.SetupRouter(jwtManager, authHandler, categoryHandler)

	log.Printf("Server running on port %s", cfg.App.Port)

	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
