package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
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

func (f *fakeCategoryService) Create(ctx context.Context, input services.CreateCategoryInput) (*models.Category, error) {
	return f.createFunc(ctx, input)
}

func (f *fakeCategoryService) GetAll(ctx context.Context, input services.CategoryListInput) (*services.CategoryListResult, error) {
	if f.getAllFunc != nil {
		return f.getAllFunc(ctx, input)
	}

	now := time.Now()

	return &services.CategoryListResult{
		Categories: []models.Category{
			{
				ID:          "category-id",
				Name:        "Electronics",
				Slug:        "electronics",
				Description: "Electronic products",
				ImageURL:    "https://example.com/electronics.jpg",
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

func (f *fakeCategoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	return f.getByIDFunc(ctx, id)
}

func (f *fakeCategoryService) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	if f.getBySlugFunc == nil {
		return nil, models.ErrCategoryNotFound
	}

	return f.getBySlugFunc(ctx, slug)
}

func (f *fakeCategoryService) Update(ctx context.Context, id string, input services.UpdateCategoryInput) (*models.Category, error) {
	return f.updateFunc(ctx, id, input)
}

func (f *fakeCategoryService) Delete(ctx context.Context, id string) error {
	return f.deleteFunc(ctx, id)
}

func setupCategoryRouter(service services.CategoryService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewCategoryHandler(service)

	router.POST("/api/v1/categories", handler.CreateCategory)
	router.GET("/api/v1/categories", handler.GetAllCategories)

	router.GET("/api/v1/categories/slug/:slug", handler.GetCategoryBySlug)

	router.GET("/api/v1/categories/:id", handler.GetCategoryByID)
	router.PUT("/api/v1/categories/:id", handler.UpdateCategory)
	router.DELETE("/api/v1/categories/:id", handler.DeleteCategory)

	return router
}

func TestCategoryHandler_CreateCategory_Success(t *testing.T) {
	now := time.Now()

	service := &fakeCategoryService{
		createFunc: func(ctx context.Context, input services.CreateCategoryInput) (*models.Category, error) {
			if input.Name != "Electronics" {
				t.Fatalf("expected Electronics, got %s", input.Name)
			}

			return &models.Category{
				ID:          "category-id",
				Name:        "Electronics",
				Slug:        "electronics",
				Description: "Electronic products",
				ImageURL:    "https://example.com/electronics.jpg",
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}

	router := setupCategoryRouter(service)

	body := `{
		"name": "Electronics",
		"description": "Electronic products",
		"image_url": "https://example.com/electronics.jpg"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "category created successfully") {
		t.Fatalf("expected success message, got body: %s", w.Body.String())
	}
}

func TestCategoryHandler_CreateCategory_Duplicate(t *testing.T) {
	service := &fakeCategoryService{
		createFunc: func(ctx context.Context, input services.CreateCategoryInput) (*models.Category, error) {
			return nil, models.ErrCategoryAlreadyExists
		},
	}

	router := setupCategoryRouter(service)

	body := `{
		"name": "Electronics",
		"description": "Electronic products"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_GetAllCategories_Success(t *testing.T) {
	now := time.Now()

	service := &fakeCategoryService{
		getAllFunc: func(ctx context.Context, input services.CategoryListInput) (*services.CategoryListResult, error) {
			if input.Page != 1 {
				t.Fatalf("expected page 1, got %d", input.Page)
			}

			if input.Limit != 20 {
				t.Fatalf("expected limit 20, got %d", input.Limit)
			}

			return &services.CategoryListResult{
				Categories: []models.Category{
					{
						ID:          "category-id",
						Name:        "Electronics",
						Slug:        "electronics",
						Description: "Electronic products",
						ImageURL:    "https://example.com/electronics.jpg",
						CreatedAt:   now,
						UpdatedAt:   now,
					},
				},
				Page:       1,
				Limit:      20,
				Total:      1,
				TotalPages: 1,
			}, nil
		},
	}

	router := setupCategoryRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "Electronics") {
		t.Fatalf("expected response body to contain category name, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"page":1`) {
		t.Fatalf("expected response to contain page meta, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"limit":20`) {
		t.Fatalf("expected response to contain limit meta, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total":1`) {
		t.Fatalf("expected response to contain total meta, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total_pages":1`) {
		t.Fatalf("expected response to contain total_pages meta, got: %s", w.Body.String())
	}
}

func TestCategoryHandler_GetCategoryByID_NotFound(t *testing.T) {
	service := &fakeCategoryService{
		getByIDFunc: func(ctx context.Context, id string) (*models.Category, error) {
			if id != "missing-id" {
				t.Fatalf("expected missing-id, got %s", id)
			}

			return nil, models.ErrCategoryNotFound
		},
	}

	router := setupCategoryRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/missing-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_GetCategoryBySlug_Success(t *testing.T) {
	now := time.Now()

	service := &fakeCategoryService{
		getBySlugFunc: func(ctx context.Context, slug string) (*models.Category, error) {
			if slug != "electronics" {
				t.Fatalf("expected electronics, got %s", slug)
			}

			return &models.Category{
				ID:          "category-id",
				Name:        "Electronics",
				Slug:        "electronics",
				Description: "Electronic products",
				ImageURL:    "https://example.com/electronics.jpg",
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}

	router := setupCategoryRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories/slug/electronics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "Electronics") {
		t.Fatalf("expected response body to contain category name, got: %s", w.Body.String())
	}
}

func TestCategoryHandler_UpdateCategory_Success(t *testing.T) {
	now := time.Now()

	service := &fakeCategoryService{
		updateFunc: func(ctx context.Context, id string, input services.UpdateCategoryInput) (*models.Category, error) {
			if id != "category-id" {
				t.Fatalf("expected category-id, got %s", id)
			}

			return &models.Category{
				ID:          id,
				Name:        "Updated Electronics",
				Slug:        "updated-electronics",
				Description: "Updated description",
				ImageURL:    "https://example.com/updated.jpg",
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}

	router := setupCategoryRouter(service)

	body := `{
		"name": "Updated Electronics",
		"description": "Updated description",
		"image_url": "https://example.com/updated.jpg"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/v1/categories/category-id", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_DeleteCategory_HasProducts(t *testing.T) {
	service := &fakeCategoryService{
		deleteFunc: func(ctx context.Context, id string) error {
			return models.ErrCategoryHasProducts
		},
	}

	router := setupCategoryRouter(service)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/categories/category-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_InternalError(t *testing.T) {
	service := &fakeCategoryService{
		getAllFunc: func(ctx context.Context, input services.CategoryListInput) (*services.CategoryListResult, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupCategoryRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCategoryHandler_GetAllCategories_WithPagination(t *testing.T) {
	service := &fakeCategoryService{
		getAllFunc: func(ctx context.Context, input services.CategoryListInput) (*services.CategoryListResult, error) {
			if input.Page != 2 {
				t.Fatalf("expected page 2, got %d", input.Page)
			}

			if input.Limit != 10 {
				t.Fatalf("expected limit 10, got %d", input.Limit)
			}

			return &services.CategoryListResult{
				Categories: []models.Category{},
				Page:       2,
				Limit:      10,
				Total:      25,
				TotalPages: 3,
			}, nil
		},
	}

	router := setupCategoryRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories?page=2&limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"page":2`) {
		t.Fatalf("expected response to contain page 2, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"limit":10`) {
		t.Fatalf("expected response to contain limit 10, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total":25`) {
		t.Fatalf("expected response to contain total 25, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total_pages":3`) {
		t.Fatalf("expected response to contain total_pages 3, got: %s", w.Body.String())
	}
}

func TestCategoryHandler_GetAllCategories_InvalidPaginationDefaults(t *testing.T) {
	service := &fakeCategoryService{
		getAllFunc: func(ctx context.Context, input services.CategoryListInput) (*services.CategoryListResult, error) {
			if input.Page != 1 {
				t.Fatalf("expected page default 1, got %d", input.Page)
			}

			if input.Limit != 20 {
				t.Fatalf("expected limit default 20, got %d", input.Limit)
			}

			return &services.CategoryListResult{
				Categories: []models.Category{},
				Page:       1,
				Limit:      20,
				Total:      0,
				TotalPages: 0,
			}, nil
		},
	}

	router := setupCategoryRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories?page=abc&limit=-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}
