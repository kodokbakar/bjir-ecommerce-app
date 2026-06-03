package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kodokbakar/go-ecommerce-api/docs"
	jwtAuth "github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/handlers"
	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
)

func SetupRouter(jwtManager *jwtAuth.JWTManager, authHandler *handlers.AuthHandler) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", handlers.HealthCheck)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	protected.GET("/me", handlers.Me)

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager))
	admin.Use(middleware.RequireRole("admin"))

	admin.GET("/ping", handlers.AdminPing)

	return r
}
