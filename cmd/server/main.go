package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/database"
	"github.com/kodokbakar/go-ecommerce-api/internal/handlers"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
	"github.com/kodokbakar/go-ecommerce-api/internal/router"
	appserver "github.com/kodokbakar/go-ecommerce-api/internal/server"
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

	log.Println("PostgreSQL connected successfully")

	redisClient, err := database.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Printf("failed to connect redis, continuing without cache: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}

	jwtManager := auth.NewJWTManager(cfg.JWT)

	userRepository := repository.NewUserRepository(dbPool)
	authService := services.NewAuthService(userRepository, jwtManager)
	authHandler := handlers.NewAuthHandler(authService)

	categoryRepository := repository.NewCategoryRepository(dbPool)
	categoryService := services.NewCategoryService(categoryRepository)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	productRepository := repository.NewProductRepository(dbPool)

	var productCache services.ProductCache
	if redisClient != nil {
		productCache = services.NewRedisProductCache(redisClient)
	}

	productService := services.NewProductServiceWithCache(productRepository, categoryRepository, productCache)
	productHandler := handlers.NewProductHandler(productService)

	cartRepository := repository.NewCartRepository(dbPool)
	cartService := services.NewCartService(cartRepository, productRepository)
	cartHandler := handlers.NewCartHandler(cartService)

	orderRepository := repository.NewOrderRepository(dbPool)
	orderService := services.NewOrderService(orderRepository)
	orderHandler := handlers.NewOrderHandler(orderService)

	paymentRepository := repository.NewPaymentRepository(dbPool)
	paymentService := services.NewPaymentService(paymentRepository)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	r := router.SetupRouter(jwtManager, authHandler, categoryHandler, productHandler, cartHandler, orderHandler, paymentHandler)

	httpServer := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server running on port %s", cfg.App.Port)

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
			return
		}

		serverErrors <- nil
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}

	case <-shutdownCtx.Done():
		log.Println("shutdown signal received")

		closeFuncs := []appserver.CloseFunc{
			func() error {
				if redisClient == nil {
					return nil
				}

				return redisClient.Close()
			},
			func() error {
				dbPool.Close()
				return nil
			},
		}

		if err := appserver.GracefulShutdown(context.Background(), httpServer, appserver.GracefulShutdownOptions{
			Timeout: cfg.App.ShutdownTimeout,
			OnShutdownStart: func() {
				handlers.SetShuttingDown(true)
			},
			CloseFuncs: closeFuncs,
		}); err != nil {
			log.Printf("server stopped with shutdown error: %v", err)
			return
		}

		log.Println("server stopped gracefully")
	}
}
