package handlers

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type fakeProductService struct {
	createFunc      func(ctx context.Context, input services.CreateProductInput) (*models.Product, error)
	getAllFunc      func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error)
	getByIDFunc     func(ctx context.Context, id string) (*models.Product, error)
	getBySlugFunc   func(ctx context.Context, slug string) (*models.Product, error)
	updateFunc      func(ctx context.Context, id string, input services.UpdateProductInput) (*models.Product, error)
	deleteFunc      func(ctx context.Context, id string) error
	uploadImageFunc func(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error)
}

func (f *fakeProductService) Create(ctx context.Context, input services.CreateProductInput) (*models.Product, error) {
	return f.createFunc(ctx, input)
}

func (f *fakeProductService) GetAll(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
	return f.getAllFunc(ctx, input)
}

func (f *fakeProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	return f.getByIDFunc(ctx, id)
}

func (f *fakeProductService) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	return f.getBySlugFunc(ctx, slug)
}

func (f *fakeProductService) Update(ctx context.Context, id string, input services.UpdateProductInput) (*models.Product, error) {
	return f.updateFunc(ctx, id, input)
}

func (f *fakeProductService) UploadImage(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error) {
	return f.uploadImageFunc(ctx, input)
}

func (f *fakeProductService) Delete(ctx context.Context, id string) error {
	return f.deleteFunc(ctx, id)
}

func setupProductRouter(service services.ProductService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewProductHandler(service)

	router.POST("/api/v1/products", handler.CreateProduct)
	router.POST("/api/v1/products/:id/image", handler.UploadProductImage)
	router.GET("/api/v1/products", handler.GetAllProducts)
	router.GET("/api/v1/products/slug/:slug", handler.GetProductBySlug)
	router.GET("/api/v1/products/:id", handler.GetProductByID)
	router.PUT("/api/v1/products/:id", handler.UpdateProduct)
	router.DELETE("/api/v1/products/:id", handler.DeleteProduct)

	return router
}

func newTestProduct() *models.Product {
	now := time.Now()

	return &models.Product{
		ID:         "product-id",
		CategoryID: "category-id",
		Category: &models.Category{
			ID:   "category-id",
			Name: "Phones",
			Slug: "phones",
		},
		Name:        "iPhone 15",
		Slug:        "iphone-15",
		Description: "Apple smartphone",
		Price:       15000000.0,
		Stock:       10,
		ImageURL:    "https://example.com/iphone.jpg",
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	service := &fakeProductService{
		createFunc: func(ctx context.Context, input services.CreateProductInput) (*models.Product, error) {
			if input.CategoryID != "category-id" {
				t.Fatalf("expected category-id, got %s", input.CategoryID)
			}

			if input.Name != "iPhone 15" {
				t.Fatalf("expected iPhone 15, got %s", input.Name)
			}

			return newTestProduct(), nil
		},
	}

	router := setupProductRouter(service)

	body := `{
		"category_id": "category-id",
		"name": "iPhone 15",
		"description": "Apple smartphone",
		"price": 15000000,
		"stock": 10,
		"image_url": "https://example.com/iphone.jpg"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "product created successfully") {
		t.Fatalf("expected success message, got body: %s", w.Body.String())
	}
}

func TestProductHandler_CreateProduct_InvalidInput(t *testing.T) {
	service := &fakeProductService{
		createFunc: func(ctx context.Context, input services.CreateProductInput) (*models.Product, error) {
			return nil, models.ErrInvalidProductInput
		},
	}

	router := setupProductRouter(service)

	body := `{
		"category_id": "category-id",
		"name": "iPhone 15",
		"price": -1,
		"stock": 10
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_GetProductByID_Success(t *testing.T) {
	service := &fakeProductService{
		getByIDFunc: func(ctx context.Context, id string) (*models.Product, error) {
			if id != "product-id" {
				t.Fatalf("expected product-id, got %s", id)
			}

			return newTestProduct(), nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/product-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "iPhone 15") {
		t.Fatalf("expected response to contain product name, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetProductByID_NotFound(t *testing.T) {
	service := &fakeProductService{
		getByIDFunc: func(ctx context.Context, id string) (*models.Product, error) {
			return nil, models.ErrProductNotFound
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/missing-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_UpdateProduct_Success(t *testing.T) {
	service := &fakeProductService{
		updateFunc: func(ctx context.Context, id string, input services.UpdateProductInput) (*models.Product, error) {
			if id != "product-id" {
				t.Fatalf("expected product-id, got %s", id)
			}

			product := newTestProduct()
			product.Name = "iPhone 15 Pro"
			product.Slug = "iphone-15-pro"
			product.Price = 18000000.0

			return product, nil
		},
	}

	router := setupProductRouter(service)

	body := `{
		"category_id": "category-id",
		"name": "iPhone 15 Pro",
		"description": "Apple smartphone pro",
		"price": 18000000,
		"stock": 5,
		"image_url": "https://example.com/iphone-pro.jpg"
	}`

	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/product-id", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "iPhone 15 Pro") {
		t.Fatalf("expected updated product name, got: %s", w.Body.String())
	}
}

func TestProductHandler_DeleteProduct_Success(t *testing.T) {
	service := &fakeProductService{
		deleteFunc: func(ctx context.Context, id string) error {
			if id != "product-id" {
				t.Fatalf("expected product-id, got %s", id)
			}

			return nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/product-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d. body: %s", w.Code, w.Body.String())
	}

	if w.Body.String() != "" {
		t.Fatalf("expected empty body, got: %s", w.Body.String())
	}
}

func TestProductHandler_InternalError(t *testing.T) {
	service := &fakeProductService{
		getByIDFunc: func(ctx context.Context, id string) (*models.Product, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/product-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_Success(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.CategoryID != "" {
				t.Fatalf("expected empty category_id, got %s", input.CategoryID)
			}

			if input.Page != services.DefaultProductPage {
				t.Fatalf("expected default page %d, got %d", services.DefaultProductPage, input.Page)
			}

			if input.Limit != services.DefaultProductLimit {
				t.Fatalf("expected default limit %d, got %d", services.DefaultProductLimit, input.Limit)
			}

			if input.SortBy != "" {
				t.Fatalf("expected empty sort_by from query, got %s", input.SortBy)
			}

			if input.SortOrder != "" {
				t.Fatalf("expected empty sort_order from query, got %s", input.SortOrder)
			}

			return &services.ProductListResult{
				Products: []models.Product{
					*newTestProduct(),
				},
				Page:         1,
				Limit:        20,
				Total:        1,
				TotalPages:   1,
				SortBy:       services.DefaultProductSortBy,
				SortOrder:    services.DefaultProductSortOrder,
				Search:       input.Search,
				CategoryID:   input.CategoryID,
				CategorySlug: input.CategorySlug,
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "iPhone 15") {
		t.Fatalf("expected response to contain product name, got: %s", w.Body.String())
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

	if !strings.Contains(w.Body.String(), `"sort_by":"created_at"`) {
		t.Fatalf("expected response to contain sort_by created_at, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"sort_order":"desc"`) {
		t.Fatalf("expected response to contain sort_order desc, got: %s", w.Body.String())
	}

}

func TestProductHandler_GetProductBySlug_Success(t *testing.T) {
	service := &fakeProductService{
		getBySlugFunc: func(ctx context.Context, slug string) (*models.Product, error) {
			if slug != "iphone-15" {
				t.Fatalf("expected iphone-15, got %s", slug)
			}

			return newTestProduct(), nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/slug/iphone-15", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "iPhone 15") {
		t.Fatalf("expected response to contain product name, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetProductBySlug_NotFound(t *testing.T) {
	service := &fakeProductService{
		getBySlugFunc: func(ctx context.Context, slug string) (*models.Product, error) {
			return nil, models.ErrProductNotFound
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/slug/missing-product", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func newMultipartImageRequest(t *testing.T, url string, fieldName string, fileName string, content []byte) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	if _, err := part.Write(content); err != nil {
		t.Fatalf("failed to write form file: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func TestProductHandler_UploadProductImage_Success(t *testing.T) {
	pngFile := []byte{
		0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
	}

	service := &fakeProductService{
		uploadImageFunc: func(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error) {
			if input.ProductID != "product-id" {
				t.Fatalf("expected product-id, got %s", input.ProductID)
			}

			if input.File == nil {
				t.Fatal("expected uploaded file")
			}

			product := newTestProduct()
			product.ImageURL = "/uploads/products/test.png"

			return product, nil
		},
	}

	router := setupProductRouter(service)

	req := newMultipartImageRequest(
		t,
		"/api/v1/products/product-id/image",
		"image",
		"test.png",
		pngFile,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "/uploads/products/test.png") {
		t.Fatalf("expected response to contain uploaded image URL, got: %s", w.Body.String())
	}
}

func TestProductHandler_UploadProductImage_MissingFile(t *testing.T) {
	service := &fakeProductService{
		uploadImageFunc: func(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error) {
			t.Fatal("service should not be called when image file is missing")
			return nil, nil
		},
	}

	router := setupProductRouter(service)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/product-id/image", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_UploadProductImage_InvalidFileType(t *testing.T) {
	service := &fakeProductService{
		uploadImageFunc: func(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error) {
			return nil, models.ErrInvalidProductInput
		},
	}

	router := setupProductRouter(service)

	req := newMultipartImageRequest(
		t,
		"/api/v1/products/product-id/image",
		"image",
		"test.txt",
		[]byte("hello world"),
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_UploadProductImage_ProductNotFound(t *testing.T) {
	service := &fakeProductService{
		uploadImageFunc: func(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error) {
			return nil, models.ErrProductNotFound
		},
	}

	router := setupProductRouter(service)

	pngFile := []byte{
		0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
	}

	req := newMultipartImageRequest(
		t,
		"/api/v1/products/missing-id/image",
		"image",
		"test.png",
		pngFile,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_WithCategoryFilter(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.CategoryID != "category-id" {
				t.Fatalf("expected category-id, got %s", input.CategoryID)
			}

			if input.Page != 2 {
				t.Fatalf("expected page 2, got %d", input.Page)
			}

			if input.Limit != 10 {
				t.Fatalf("expected limit 10, got %d", input.Limit)
			}

			return &services.ProductListResult{
				Products: []models.Product{
					*newTestProduct(),
				},
				Page:         2,
				Limit:        10,
				Total:        21,
				TotalPages:   3,
				Search:       input.Search,
				CategoryID:   input.CategoryID,
				CategorySlug: input.CategorySlug,
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?category_id=category-id&page=2&limit=10", nil)
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

	if !strings.Contains(w.Body.String(), `"total":21`) {
		t.Fatalf("expected response to contain total 21, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total_pages":3`) {
		t.Fatalf("expected response to contain total_pages 3, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_InvalidPageDefaultsToOne(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.Page != 1 {
				t.Fatalf("expected page to default to 1, got %d", input.Page)
			}

			return &services.ProductListResult{
				Products:   []models.Product{},
				Page:       1,
				Limit:      20,
				Total:      0,
				TotalPages: 0,
				SortBy:     services.DefaultProductSortBy,
				SortOrder:  services.DefaultProductSortOrder,
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?page=abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"page":1`) {
		t.Fatalf("expected response page to default to 1, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_InvalidLimitDefaults(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.Limit != 20 {
				t.Fatalf("expected limit default 20, got %d", input.Limit)
			}

			return &services.ProductListResult{
				Products:   []models.Product{},
				Page:       1,
				Limit:      20,
				Total:      0,
				TotalPages: 0,
				SortBy:     services.DefaultProductSortBy,
				SortOrder:  services.DefaultProductSortOrder,
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?limit=abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_LimitTooLargeCappedAtMax(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.Limit != 100 {
				t.Fatalf("expected limit capped at 100, got %d", input.Limit)
			}

			return &services.ProductListResult{
				Products:   []models.Product{},
				Page:       1,
				Limit:      100,
				Total:      0,
				TotalPages: 0,
				SortBy:     services.DefaultProductSortBy,
				SortOrder:  services.DefaultProductSortOrder,
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?limit=101", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_WithSort(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.SortBy != "price" {
				t.Fatalf("expected sort_by price, got %s", input.SortBy)
			}

			if input.SortOrder != "asc" {
				t.Fatalf("expected sort_order asc, got %s", input.SortOrder)
			}

			return &services.ProductListResult{
				Products: []models.Product{
					*newTestProduct(),
				},
				Page:       1,
				Limit:      20,
				Total:      1,
				TotalPages: 1,
				SortBy:     "price",
				SortOrder:  "asc",
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?sort_by=price&sort_order=asc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"sort_by":"price"`) {
		t.Fatalf("expected response to contain sort_by price, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"sort_order":"asc"`) {
		t.Fatalf("expected response to contain sort_order asc, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_InvalidSort(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			return nil, models.ErrInvalidProductInput
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?sort_by=stock", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_WithSearch(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.Search != "phone" {
				t.Fatalf("expected search phone, got %s", input.Search)
			}

			return &services.ProductListResult{
				Products: []models.Product{
					*newTestProduct(),
				},
				Page:       1,
				Limit:      20,
				Total:      1,
				TotalPages: 1,
				SortBy:     services.DefaultProductSortBy,
				SortOrder:  services.DefaultProductSortOrder,
				Search:     "phone",
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?search=phone", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"search":"phone"`) {
		t.Fatalf("expected response to contain search phone, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_WithSearchAndCategoryID(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.Search != "phone" {
				t.Fatalf("expected search phone, got %s", input.Search)
			}

			if input.CategoryID != "category-id" {
				t.Fatalf("expected category-id, got %s", input.CategoryID)
			}

			return &services.ProductListResult{
				Products:   []models.Product{},
				Page:       1,
				Limit:      20,
				Total:      0,
				TotalPages: 0,
				SortBy:     services.DefaultProductSortBy,
				SortOrder:  services.DefaultProductSortOrder,
				Search:     "phone",
				CategoryID: "category-id",
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?search=phone&category_id=category-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"search":"phone"`) {
		t.Fatalf("expected response to contain search phone, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"category_id":"category-id"`) {
		t.Fatalf("expected response to contain category_id, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_WithCategorySlug(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.CategorySlug != "phones" {
				t.Fatalf("expected category slug phones, got %s", input.CategorySlug)
			}

			return &services.ProductListResult{
				Products:     []models.Product{},
				Page:         1,
				Limit:        20,
				Total:        0,
				TotalPages:   0,
				SortBy:       services.DefaultProductSortBy,
				SortOrder:    services.DefaultProductSortOrder,
				CategorySlug: "phones",
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?category=phones", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"category":"phones"`) {
		t.Fatalf("expected response to contain category phones, got: %s", w.Body.String())
	}
}

func TestProductHandler_GetAllProducts_InvalidCategoryIDReturnsEmptyResult(t *testing.T) {
	service := &fakeProductService{
		getAllFunc: func(ctx context.Context, input services.ProductListInput) (*services.ProductListResult, error) {
			if input.CategoryID != "invalid-category-id" {
				t.Fatalf("expected invalid-category-id, got %s", input.CategoryID)
			}

			return &services.ProductListResult{
				Products:   []models.Product{},
				Page:       1,
				Limit:      20,
				Total:      0,
				TotalPages: 0,
				SortBy:     services.DefaultProductSortBy,
				SortOrder:  services.DefaultProductSortOrder,
				CategoryID: "invalid-category-id",
			}, nil
		},
	}

	router := setupProductRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products?category_id=invalid-category-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total":0`) {
		t.Fatalf("expected response total 0, got: %s", w.Body.String())
	}
}
