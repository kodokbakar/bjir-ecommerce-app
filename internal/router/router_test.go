package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/auth"
	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/handlers"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type fakeCategoryService struct {
	createFunc    func(ctx context.Context, input services.CreateCategoryInput) (*models.Category, error)
	getAllFunc    func(ctx context.Context, input services.CategoryListInput) (*services.CategoryListResult, error)
	getByIDFunc   func(ctx context.Context, id string) (*models.Category, error)
	getBySlugFunc func(ctx context.Context, slug string) (*models.Category, error)
	updateFunc    func(ctx context.Context, id string, input services.UpdateCategoryInput) (*models.Category, error)
	deleteFunc    func(ctx context.Context, id string) error
}

type fakeRouterProductService struct{}

func (f *fakeRouterProductService) Create(ctx context.Context, input services.CreateProductInput) (*models.Product, error) {
	return &models.Product{
		ID:         "product-id",
		CategoryID: input.CategoryID,
		Name:       input.Name,
		Slug:       "product-slug",
		Price:      input.Price,
		Stock:      input.Stock,
		IsActive:   true,
	}, nil
}

func (f *fakeRouterProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	return &models.Product{
		ID:         id,
		CategoryID: "category-id",
		Name:       "Product",
		Slug:       "product",
		Price:      10000,
		Stock:      10,
		IsActive:   true,
	}, nil
}

func (f *fakeRouterProductService) Update(ctx context.Context, id string, input services.UpdateProductInput) (*models.Product, error) {
	return &models.Product{
		ID:         id,
		CategoryID: input.CategoryID,
		Name:       input.Name,
		Slug:       "product-slug",
		Price:      input.Price,
		Stock:      input.Stock,
		IsActive:   true,
	}, nil
}

func (f *fakeRouterProductService) UploadImage(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error) {
	return &models.Product{
		ID:         input.ProductID,
		CategoryID: "category-id",
		Name:       "Product",
		Slug:       "product",
		Price:      10000,
		Stock:      10,
		ImageURL:   "/uploads/products/test.png",
		IsActive:   true,
	}, nil
}

func (f *fakeRouterProductService) GetImages(ctx context.Context, productID string) ([]models.ProductImage, error) {
	return []models.ProductImage{
		{
			ID:        "image-id",
			ProductID: productID,
			ImageURL:  "/uploads/products/test.png",
			SortOrder: 0,
			IsPrimary: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (f *fakeRouterProductService) UploadGalleryImage(ctx context.Context, input services.UploadProductImageInput) (*models.ProductImage, error) {
	return &models.ProductImage{
		ID:        "image-id",
		ProductID: input.ProductID,
		ImageURL:  "/uploads/products/test.png",
		SortOrder: 0,
		IsPrimary: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (f *fakeRouterProductService) DeleteImage(ctx context.Context, productID string, imageID string) error {
	return nil
}

func (f *fakeRouterProductService) ReorderImages(ctx context.Context, input services.ReorderProductImagesInput) ([]models.ProductImage, error) {
	return []models.ProductImage{
		{
			ID:        "image-2",
			ProductID: input.ProductID,
			ImageURL:  "/uploads/products/2.png",
			SortOrder: 0,
			IsPrimary: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "image-1",
			ProductID: input.ProductID,
			ImageURL:  "/uploads/products/1.png",
			SortOrder: 1,
			IsPrimary: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

func (f *fakeRouterProductService) SetPrimaryImage(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
	return &models.ProductImage{
		ID:        imageID,
		ProductID: productID,
		ImageURL:  "/uploads/products/primary.png",
		SortOrder: 0,
		IsPrimary: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (f *fakeRouterProductService) Delete(ctx context.Context, id string) error {
	return nil
}

func (f *fakeCategoryService) Create(ctx context.Context, input services.CreateCategoryInput) (*models.Category, error) {
	if f.createFunc != nil {
		return f.createFunc(ctx, input)
	}

	now := time.Now()

	return &models.Category{
		ID:          "category-id",
		ParentID:    input.ParentID,
		Name:        input.Name,
		Slug:        "electronics",
		Description: input.Description,
		ImageURL:    input.ImageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (f *fakeCategoryService) GetAll(ctx context.Context, input services.CategoryListInput) (*services.CategoryListResult, error) {
	if f.getAllFunc != nil {
		return f.getAllFunc(ctx, input)
	}

	return &services.CategoryListResult{
		Categories: []models.Category{
			{
				ID:   "category-id",
				Name: "Electronics",
				Slug: "electronics",
			},
		},
		Page:       1,
		Limit:      20,
		Total:      1,
		TotalPages: 1,
	}, nil
}

func (f *fakeCategoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	if f.getByIDFunc != nil {
		return f.getByIDFunc(ctx, id)
	}

	now := time.Now()

	return &models.Category{
		ID:          id,
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "Electronic products",
		ImageURL:    "https://example.com/electronics.jpg",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (f *fakeCategoryService) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	if f.getBySlugFunc != nil {
		return f.getBySlugFunc(ctx, slug)
	}

	now := time.Now()

	return &models.Category{
		ID:          "category-id",
		Name:        "Electronics",
		Slug:        slug,
		Description: "Electronic products",
		ImageURL:    "https://example.com/electronics.jpg",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (f *fakeCategoryService) Update(ctx context.Context, id string, input services.UpdateCategoryInput) (*models.Category, error) {
	if f.updateFunc != nil {
		return f.updateFunc(ctx, id, input)
	}

	now := time.Now()

	return &models.Category{
		ID:          id,
		ParentID:    input.ParentID,
		Name:        input.Name,
		Slug:        "electronics",
		Description: input.Description,
		ImageURL:    input.ImageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (f *fakeCategoryService) Delete(ctx context.Context, id string) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(ctx, id)
	}

	return nil
}

type fakeRouterCartService struct{}

func (f *fakeRouterCartService) AddItem(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
	return &models.CartItem{
		ID:        "cart-item-id",
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Product: &models.Product{
			ID:       productID,
			Name:     "Product",
			Slug:     "product",
			Price:    10000,
			Stock:    10,
			IsActive: true,
		},
		Subtotal: 10000 * float64(quantity),
	}, nil
}

func (f *fakeRouterCartService) GetCart(ctx context.Context, userID string) (*models.Cart, error) {
	return &models.Cart{
		Items: []models.CartItem{
			{
				ID:        "cart-item-id",
				UserID:    userID,
				ProductID: "product-id",
				Quantity:  1,
				Product: &models.Product{
					ID:       "product-id",
					Name:     "Product",
					Slug:     "product",
					Price:    10000,
					Stock:    10,
					IsActive: true,
				},
				Subtotal: 10000,
			},
		},
		TotalPrice: 10000,
	}, nil
}

func (f *fakeRouterCartService) UpdateItem(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error) {
	return &models.CartItem{
		ID:        itemID,
		UserID:    userID,
		ProductID: "product-id",
		Quantity:  quantity,
		Product: &models.Product{
			ID:       "product-id",
			Name:     "Product",
			Slug:     "product",
			Price:    10000,
			Stock:    10,
			IsActive: true,
		},
		Subtotal: 10000 * float64(quantity),
	}, nil
}

func (f *fakeRouterCartService) DeleteItem(ctx context.Context, userID string, itemID string) error {
	return nil
}

func setupRouterForCategoryAuthTest() (*gin.Engine, *auth.JWTManager) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager(config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	authHandler := handlers.NewAuthHandler(nil)
	categoryHandler := handlers.NewCategoryHandler(&fakeCategoryService{})
	productHandler := handlers.NewProductHandler(&fakeRouterProductService{})
	cartHandler := handlers.NewCartHandler(&fakeRouterCartService{})
	orderHandler := handlers.NewOrderHandler(&fakeRouterOrderService{})
	paymentHandler := handlers.NewPaymentHandler(&fakeRouterPaymentService{})
	dashboardHandler := handlers.NewDashboardHandler(&fakeRouterDashboardService{})

	return SetupRouter(
		jwtManager,
		authHandler,
		categoryHandler,
		productHandler,
		cartHandler,
		orderHandler,
		paymentHandler,
		dashboardHandler,
	), jwtManager
}

type fakeRouterOrderService struct{}

func (f *fakeRouterOrderService) Checkout(ctx context.Context, userID string) (*models.Order, error) {
	now := time.Now()

	return &models.Order{
		ID:          "order-id",
		UserID:      userID,
		OrderNumber: "ORD-TEST",
		Status:      models.OrderStatusPending,
		TotalAmount: 10000,
		Items: []models.OrderItem{
			{
				ID:          "order-item-id",
				OrderID:     "order-id",
				ProductID:   "product-id",
				ProductName: "Product",
				Quantity:    1,
				Price:       10000,
				Subtotal:    10000,
				CreatedAt:   now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (f *fakeRouterOrderService) GetMyOrders(ctx context.Context, userID string, input services.OrderListInput) (*services.OrderListResult, error) {
	now := time.Now()

	return &services.OrderListResult{
		Orders: []models.Order{
			{
				ID:          "order-id",
				UserID:      userID,
				OrderNumber: "ORD-TEST",
				Status:      models.OrderStatusPending,
				TotalAmount: 10000,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		Page:       1,
		Limit:      20,
		Total:      1,
		TotalPages: 1,
	}, nil
}

func (f *fakeRouterOrderService) GetMyOrderDetail(ctx context.Context, userID string, orderID string) (*models.Order, error) {
	now := time.Now()

	return &models.Order{
		ID:          orderID,
		UserID:      userID,
		OrderNumber: "ORD-TEST",
		Status:      models.OrderStatusPending,
		TotalAmount: 10000,
		Items: []models.OrderItem{
			{
				ID:          "order-item-id",
				OrderID:     orderID,
				ProductID:   "product-id",
				ProductName: "Product",
				Quantity:    1,
				Price:       10000,
				Subtotal:    10000,
				CreatedAt:   now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (f *fakeRouterOrderService) GetAllOrders(ctx context.Context, input services.OrderListInput) (*services.OrderListResult, error) {
	return &services.OrderListResult{
		Orders:     []models.Order{},
		Page:       1,
		Limit:      20,
		Total:      0,
		TotalPages: 0,
	}, nil
}

func TestRBAC_PublicProductRoutes_WithoutToken_ReturnsOK(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "GET products",
			path: "/api/v1/products",
		},
		{
			name: "GET product by id",
			path: "/api/v1/products/product-id",
		},
		{
			name: "GET product by slug",
			path: "/api/v1/products/slug/product",
		},
		{
			name: "GET categories",
			path: "/api/v1/categories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestRBAC_ProtectedMe_WithoutToken_ReturnsUnauthorized(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRBAC_ProtectedMe_WithCustomerToken_ReturnsOK(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("customer-id", "customer@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate customer token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRBAC_AdminPing_WithCustomerToken_ReturnsForbidden(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("customer-id", "customer@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate customer token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestRBAC_AdminPing_WithAdminToken_ReturnsOK(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryAdminRoutes_WithoutToken_ReturnsUnauthorized(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "POST categories",
			method: http.MethodPost,
			path:   "/api/v1/categories",
			body:   `{"name":"Electronics"}`,
		},
		{
			name:   "PUT categories",
			method: http.MethodPut,
			path:   "/api/v1/categories/category-id",
			body:   `{"name":"Updated Electronics"}`,
		},
		{
			name:   "DELETE categories",
			method: http.MethodDelete,
			path:   "/api/v1/categories/category-id",
			body:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestCategoryAdminRoutes_WithCustomerToken_ReturnsForbidden(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("user-id", "customer@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate customer token: %v", err)
	}

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "POST categories",
			method: http.MethodPost,
			path:   "/api/v1/categories",
			body:   `{"name":"Electronics"}`,
		},
		{
			name:   "PUT categories",
			method: http.MethodPut,
			path:   "/api/v1/categories/category-id",
			body:   `{"name":"Updated Electronics"}`,
		},
		{
			name:   "DELETE categories",
			method: http.MethodDelete,
			path:   "/api/v1/categories/category-id",
			body:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusForbidden {
				t.Fatalf("expected status 403, got %d. body: %s", w.Code, w.Body.String())
			}
		})
	}
}

func TestCategoryAdminRoutes_WithAdminToken_AllowsAccess(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		{
			name:           "POST categories",
			method:         http.MethodPost,
			path:           "/api/v1/categories",
			body:           `{"name":"Electronics","description":"Electronic products","image_url":"https://example.com/electronics.jpg"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "PUT categories",
			method:         http.MethodPut,
			path:           "/api/v1/categories/category-id",
			body:           `{"name":"Updated Electronics","description":"Updated description","image_url":"https://example.com/updated.jpg"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DELETE categories",
			method:         http.MethodDelete,
			path:           "/api/v1/categories/category-id",
			body:           "",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d. body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestCategoryPublicRoutes_WithoutToken_ReturnsOK(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	tests := []struct {
		name string
		path string
	}{
		{
			name: "GET categories",
			path: "/api/v1/categories",
		},
		{
			name: "GET category by id",
			path: "/api/v1/categories/category-id",
		},
		{
			name: "GET category by slug",
			path: "/api/v1/categories/slug/electronics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
			}
		})
	}
}

func (f *fakeRouterProductService) GetAll(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
	return &services.ProductListResult{
		Products: []models.Product{
			{
				ID:         "product-id",
				CategoryID: "category-id",
				Name:       "Product",
				Slug:       "product",
				Price:      10000,
				Stock:      10,
				IsActive:   true,
			},
		},
		Page:         1,
		Limit:        20,
		Total:        1,
		TotalPages:   1,
		SortBy:       services.DefaultProductSortBy,
		SortOrder:    services.DefaultProductSortOrder,
		CategoryID:   input.CategoryID,
		CategorySlug: input.CategorySlug,
		Search:       input.Search,
	}, nil
}

func (f *fakeRouterProductService) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	return &models.Product{
		ID:         "product-id",
		CategoryID: "category-id",
		Name:       "Product",
		Slug:       slug,
		Price:      10000,
		Stock:      10,
		IsActive:   true,
	}, nil
}

func (f *fakeRouterOrderService) UpdateStatus(ctx context.Context, orderID string, status string) (*models.Order, error) {
	return &models.Order{
		ID:          orderID,
		UserID:      "user-id",
		OrderNumber: "ORD-TEST",
		Status:      status,
		TotalAmount: 30000000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func TestOrderAdminRoutes_WithoutToken_ReturnsUnauthorized(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/admin/orders/order-id/status",
		strings.NewReader(`{"status":"paid"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderAdminRoutes_WithCustomerToken_ReturnsForbidden(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("customer-id", "customer@example.com", "customer")
	if err != nil {
		t.Fatalf("failed to generate customer token: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/admin/orders/order-id/status",
		strings.NewReader(`{"status":"paid"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderAdminRoutes_WithAdminToken_ReturnsOK(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/admin/orders/order-id/status",
		strings.NewReader(`{"status":"paid"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

type fakeRouterPaymentService struct{}

func (f *fakeRouterPaymentService) PayOrder(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
	return &models.Payment{
		ID:            "payment-id",
		OrderID:       input.OrderID,
		Provider:      models.PaymentProviderMock,
		PaymentMethod: input.Method,
		TransactionID: "PAY-TEST",
		Amount:        30000000,
		Status:        models.PaymentStatusPaid,
	}, nil
}

func TestRouter_UnknownRoute_ReturnsJSONNotFound(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"success":false`) {
		t.Fatalf("expected JSON error response, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"code":"not_found"`) {
		t.Fatalf("expected not_found error code, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"message":"route not found"`) {
		t.Fatalf("expected route not found message, got: %s", w.Body.String())
	}
}

func TestRouter_HealthHead_ReturnsOK(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductImageRoutes_PublicGetImages_ReturnsOK(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/product-id/images", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "product images retrieved successfully") {
		t.Fatalf("expected product images response, got body: %s", w.Body.String())
	}
}

func TestProductImageRoutes_AdminUploadWithoutToken_ReturnsUnauthorized(t *testing.T) {
	r, _ := setupRouterForCategoryAuthTest()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/product-id/images", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductImageRoutes_AdminReorderWithAdminToken_ReturnsOK(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	body := `{
		"images": [
			{"id": "image-1", "sort_order": 1},
			{"id": "image-2", "sort_order": 0}
		]
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/products/product-id/images/reorder",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductImageRoutes_AdminSetPrimaryWithAdminToken_ReturnsOK(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/products/product-id/images/image-id/primary",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductImageRoutes_AdminDeleteWithAdminToken_ReturnsNoContent(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/products/product-id/images/image-id",
		nil,
	)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestDashboardAdminRoute_WithAdminToken_ReturnsOK(t *testing.T) {
	r, jwtManager := setupRouterForCategoryAuthTest()

	token, err := jwtManager.GenerateToken("admin-id", "admin@example.com", "admin")
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "total_orders") {
		t.Fatalf("expected dashboard stats response, got: %s", w.Body.String())
	}
}

type fakeRouterDashboardService struct{}

func (f *fakeRouterDashboardService) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	return &models.DashboardStats{
		TotalOrders:     10,
		TotalRevenue:    1500000,
		PendingOrders:   2,
		CompletedToday:  1,
		RevenueToday:    250000,
		TotalProducts:   8,
		TotalCategories: 4,
	}, nil
}
