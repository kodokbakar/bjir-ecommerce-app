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
	getAllFunc      func(ctx context.Context) ([]models.Product, error)
	getByIDFunc     func(ctx context.Context, id string) (*models.Product, error)
	getBySlugFunc   func(ctx context.Context, slug string) (*models.Product, error)
	updateFunc      func(ctx context.Context, id string, input services.UpdateProductInput) (*models.Product, error)
	deleteFunc      func(ctx context.Context, id string) error
	uploadImageFunc func(ctx context.Context, input services.UploadProductImageInput) (*models.Product, error)
}

func (f *fakeProductService) Create(ctx context.Context, input services.CreateProductInput) (*models.Product, error) {
	return f.createFunc(ctx, input)
}

func (f *fakeProductService) GetAll(ctx context.Context) ([]models.Product, error) {
	return f.getAllFunc(ctx)
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
		getAllFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				*newTestProduct(),
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
