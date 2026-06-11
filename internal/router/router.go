package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/kodokbakar/go-ecommerce-api/docs"
	jwtAuth "github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/handlers"
	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
	"github.com/kodokbakar/go-ecommerce-api/internal/response"
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

	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())

	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "route not found", nil)
	})

	r.GET("/health", handlers.HealthCheck)
	r.HEAD("/health", handlers.HealthCheck)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")

	bodyLimitConfig := middleware.LoadBodyLimitConfigFromEnv()

	authRoutes := api.Group("/auth")
	authRoutes.Use(middleware.BodySizeLimit(bodyLimitConfig.Auth))
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))

	protected.GET("/me", handlers.Me)

	if cartHandler != nil {
		cartRoutes := protected.Group("/cart")
		cartRoutes.GET("", cartHandler.GetCart)
		cartRoutes.POST("/items", middleware.BodySizeLimit(bodyLimitConfig.API), cartHandler.AddCartItem)
		cartRoutes.PUT("/items/:id", middleware.BodySizeLimit(bodyLimitConfig.API), cartHandler.UpdateCartItem)
		cartRoutes.DELETE("/items/:id", cartHandler.DeleteCartItem)
	}

	if orderHandler != nil {
		orderRoutes := protected.Group("/orders")
		orderRoutes.POST("/checkout", middleware.BodySizeLimit(bodyLimitConfig.API), orderHandler.Checkout)
		orderRoutes.GET("", orderHandler.GetMyOrders)
		orderRoutes.GET("/:id", orderHandler.GetMyOrderDetail)
	}

	if paymentHandler != nil {
		paymentRoutes := protected.Group("/payments")
		paymentRoutes.POST("/pay", middleware.BodySizeLimit(bodyLimitConfig.API), paymentHandler.PayOrder)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager))
	admin.Use(middleware.RequireRole("admin"))

	admin.GET("/ping", handlers.AdminPing)

	if orderHandler != nil {
		admin.PATCH("/orders/:id/status", middleware.BodySizeLimit(bodyLimitConfig.API), orderHandler.UpdateOrderStatus)
	}

	categoryRoutes := api.Group("/categories")
	categoryRoutes.GET("", categoryHandler.GetAllCategories)
	categoryRoutes.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
	categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)

	adminCategoryRoutes := categoryRoutes.Group("")
	adminCategoryRoutes.Use(middleware.AuthMiddleware(jwtManager))
	adminCategoryRoutes.Use(middleware.RequireRole("admin"))

	adminCategoryRoutes.POST("", middleware.BodySizeLimit(bodyLimitConfig.API), categoryHandler.CreateCategory)
	adminCategoryRoutes.PUT("/:id", middleware.BodySizeLimit(bodyLimitConfig.API), categoryHandler.UpdateCategory)
	adminCategoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)

	productRoutes := api.Group("/products")
	productRoutes.GET("", productHandler.GetAllProducts)
	productRoutes.GET("/slug/:slug", productHandler.GetProductBySlug)
	productRoutes.GET("/:id/images", productHandler.GetProductImages)
	productRoutes.GET("/:id", productHandler.GetProductByID)

	adminProductRoutes := productRoutes.Group("")
	adminProductRoutes.Use(middleware.AuthMiddleware(jwtManager))
	adminProductRoutes.Use(middleware.RequireRole("admin"))

	adminProductRoutes.POST("", middleware.BodySizeLimit(bodyLimitConfig.API), productHandler.CreateProduct)
	adminProductRoutes.POST("/:id/image", middleware.BodySizeLimit(bodyLimitConfig.Upload), productHandler.UploadProductImage)
	adminProductRoutes.POST("/:id/images", middleware.BodySizeLimit(bodyLimitConfig.Upload), productHandler.UploadProductGalleryImage)
	adminProductRoutes.DELETE("/:id/images/:image_id", productHandler.DeleteProductImage)
	adminProductRoutes.PATCH("/:id/images/reorder", middleware.BodySizeLimit(bodyLimitConfig.API), productHandler.ReorderProductImages)
	adminProductRoutes.PATCH("/:id/images/:image_id/primary", productHandler.SetPrimaryProductImage)
	adminProductRoutes.PUT("/:id", middleware.BodySizeLimit(bodyLimitConfig.API), productHandler.UpdateProduct)
	adminProductRoutes.DELETE("/:id", productHandler.DeleteProduct)

	return r
}
