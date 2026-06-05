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
	getAllFunc    func(ctx context.Context) ([]models.Category, error)
	getByIDFunc   func(ctx context.Context, id string) (*models.Category, error)
	getBySlugFunc func(ctx context.Context, slug string) (*models.Category, error)
	updateFunc    func(ctx context.Context, id string, input services.UpdateCategoryInput) (*models.Category, error)
	deleteFunc    func(ctx context.Context, id string) error
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

func (f *fakeCategoryService) GetAll(ctx context.Context) ([]models.Category, error) {
	if f.getAllFunc != nil {
		return f.getAllFunc(ctx)
	}

	now := time.Now()

	return []models.Category{
		{
			ID:          "category-id",
			Name:        "Electronics",
			Slug:        "electronics",
			Description: "Electronic products",
			ImageURL:    "https://example.com/electronics.jpg",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
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

func setupRouterForCategoryAuthTest() (*gin.Engine, *auth.JWTManager) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager(config.JWTConfig{
		Secret:    "test-secret",
		ExpiresIn: time.Hour,
		Issuer:    "go-ecommerce-api-test",
	})

	authHandler := handlers.NewAuthHandler(nil)
	categoryHandler := handlers.NewCategoryHandler(&fakeCategoryService{})

	return SetupRouter(jwtManager, authHandler, categoryHandler), jwtManager
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
