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

func SetupRouter(
	jwtManager *jwtAuth.JWTManager,
	authHandler *handlers.AuthHandler,
	categoryHandler *handlers.CategoryHandler,
	productHandler *handlers.ProductHandler,
	cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler,
	paymentHandler *handlers.PaymentHandler,
) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", handlers.HealthCheck)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")

	authRoutes := api.Group("/auth")
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	protected.GET("/me", handlers.Me)

	cartRoutes := protected.Group("/cart")
	cartRoutes.GET("", cartHandler.GetCart)
	cartRoutes.POST("/items", cartHandler.AddCartItem)
	cartRoutes.PUT("/items/:id", cartHandler.UpdateCartItem)
	cartRoutes.DELETE("/items/:id", cartHandler.DeleteCartItem)

	orderRoutes := protected.Group("/orders")
	orderRoutes.GET("", orderHandler.GetMyOrders)
	orderRoutes.POST("/checkout", orderHandler.Checkout)
	orderRoutes.GET("/:id", orderHandler.GetMyOrderDetail)

	paymentRoutes := protected.Group("/payments")
	paymentRoutes.POST("/pay", paymentHandler.PayOrder)

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager))
	admin.Use(middleware.RequireRole("admin"))

	admin.GET("/ping", handlers.AdminPing)
	admin.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)

	categoryRoutes := api.Group("/categories")
	categoryRoutes.GET("", categoryHandler.GetAllCategories)
	categoryRoutes.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
	categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)

	adminCategoryRoutes := categoryRoutes.Group("")
	adminCategoryRoutes.Use(middleware.AuthMiddleware(jwtManager))
	adminCategoryRoutes.Use(middleware.RequireRole("admin"))

	adminCategoryRoutes.POST("", categoryHandler.CreateCategory)
	adminCategoryRoutes.PUT("/:id", categoryHandler.UpdateCategory)
	adminCategoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)

	productRoutes := api.Group("/products")
	productRoutes.GET("", productHandler.GetAllProducts)
	productRoutes.GET("/slug/:slug", productHandler.GetProductBySlug)
	productRoutes.GET("/:id", productHandler.GetProductByID)

	adminProductRoutes := productRoutes.Group("")
	adminProductRoutes.Use(middleware.AuthMiddleware(jwtManager))
	adminProductRoutes.Use(middleware.RequireRole("admin"))

	adminProductRoutes.POST("", productHandler.CreateProduct)
	adminProductRoutes.POST("/:id/image", productHandler.UploadProductImage)
	adminProductRoutes.PUT("/:id", productHandler.UpdateProduct)
	adminProductRoutes.DELETE("/:id", productHandler.DeleteProduct)

	return r
}
