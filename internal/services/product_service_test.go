package services

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type fakeProductRepository struct {
	createFunc         func(ctx context.Context, product *models.Product) error
	findAllFunc        func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error)
	findByIDFunc       func(ctx context.Context, id string) (*models.Product, error)
	findBySlugFunc     func(ctx context.Context, slug string) (*models.Product, error)
	existsBySlugFunc   func(ctx context.Context, slug string, excludeID string) (bool, error)
	updateFunc         func(ctx context.Context, product *models.Product) error
	deleteFunc         func(ctx context.Context, id string) error
	updateImageURLFunc func(ctx context.Context, id string, imageURL string) (*models.Product, error)
}

func newFakeProductRepository() *fakeProductRepository {
	now := time.Now()

	return &fakeProductRepository{
		createFunc: func(ctx context.Context, product *models.Product) error {
			product.ID = "product-id"
			product.IsActive = true
			product.CreatedAt = now
			product.UpdatedAt = now
			return nil
		},
		findAllFunc: func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
			return []models.Product{
				{
					ID:         "product-id",
					CategoryID: "category-id",
					Name:       "iPhone 15",
					Slug:       "iphone-15",
					Price:      15000000,
					Stock:      10,
					IsActive:   true,
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			}, 1, nil
		},
		findByIDFunc: func(ctx context.Context, id string) (*models.Product, error) {
			return &models.Product{
				ID:         id,
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Slug:       "iphone-15",
				Price:      15000000,
				Stock:      10,
				IsActive:   true,
				CreatedAt:  now,
				UpdatedAt:  now,
			}, nil
		},
		findBySlugFunc: func(ctx context.Context, slug string) (*models.Product, error) {
			return &models.Product{
				ID:         "product-id",
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Slug:       slug,
				Price:      15000000,
				Stock:      10,
				IsActive:   true,
				CreatedAt:  now,
				UpdatedAt:  now,
			}, nil
		},
		existsBySlugFunc: func(ctx context.Context, slug string, excludeID string) (bool, error) {
			return false, nil
		},
		updateFunc: func(ctx context.Context, product *models.Product) error {
			product.UpdatedAt = now
			return nil
		},
		updateImageURLFunc: func(ctx context.Context, id string, imageURL string) (*models.Product, error) {
			return &models.Product{
				ID:         id,
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Slug:       "iphone-15",
				Price:      15000000,
				Stock:      10,
				ImageURL:   imageURL,
				IsActive:   true,
				CreatedAt:  now,
				UpdatedAt:  now,
			}, nil
		},
		deleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
}

func (f *fakeProductRepository) Create(ctx context.Context, product *models.Product) error {
	return f.createFunc(ctx, product)
}

func (f *fakeProductRepository) FindAll(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
	return f.findAllFunc(ctx, filter)
}

func (f *fakeProductRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	return f.findByIDFunc(ctx, id)
}

func (f *fakeProductRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	return f.findBySlugFunc(ctx, slug)
}

func (f *fakeProductRepository) ExistsBySlug(ctx context.Context, slug string, excludeID string) (bool, error) {
	return f.existsBySlugFunc(ctx, slug, excludeID)
}

func (f *fakeProductRepository) Update(ctx context.Context, product *models.Product) error {
	return f.updateFunc(ctx, product)
}

func (f *fakeProductRepository) UpdateImageURL(ctx context.Context, id string, imageURL string) (*models.Product, error) {
	return f.updateImageURLFunc(ctx, id, imageURL)
}

func (f *fakeProductRepository) Delete(ctx context.Context, id string) error {
	return f.deleteFunc(ctx, id)
}

func TestProductService_Create_Success(t *testing.T) {
	productRepo := newFakeProductRepository()
	categoryRepo := newFakeCategoryRepository()

	service := NewProductService(productRepo, categoryRepo)

	product, err := service.Create(context.Background(), CreateProductInput{
		CategoryID:  "category-id",
		Name:        "  iPhone   15  ",
		Description: "Apple smartphone",
		Price:       15000000,
		Stock:       10,
		ImageURL:    "https://example.com/iphone.jpg",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID == "" {
		t.Fatal("expected product id to be set")
	}

	if product.Name != "iPhone 15" {
		t.Fatalf("expected normalized name iPhone 15, got %s", product.Name)
	}

	if product.Slug != "iphone-15" {
		t.Fatalf("expected slug iphone-15, got %s", product.Slug)
	}
}

func TestProductService_GetAll_Success(t *testing.T) {
	productRepo := newFakeProductRepository()
	categoryRepo := newFakeCategoryRepository()

	service := NewProductService(productRepo, categoryRepo)

	result, err := service.GetAll(context.Background(), ProductListInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != DefaultProductPage {
		t.Fatalf("expected page %d, got %d", DefaultProductPage, result.Page)
	}

	if result.Limit != DefaultProductLimit {
		t.Fatalf("expected limit %d, got %d", DefaultProductLimit, result.Limit)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}

	if result.TotalPages != 1 {
		t.Fatalf("expected total_pages 1, got %d", result.TotalPages)
	}

	if len(result.Products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(result.Products))
	}

	if result.SortBy != DefaultProductSortBy {
		t.Fatalf("expected sort_by %s, got %s", DefaultProductSortBy, result.SortBy)
	}

	if result.SortOrder != DefaultProductSortOrder {
		t.Fatalf("expected sort_order %s, got %s", DefaultProductSortOrder, result.SortOrder)
	}
}

func TestProductService_Create_CategoryNotFound(t *testing.T) {
	productRepo := newFakeProductRepository()
	categoryRepo := newFakeCategoryRepository()
	categoryRepo.findByIDFunc = func(ctx context.Context, id string) (*models.Category, error) {
		return nil, models.ErrCategoryNotFound
	}

	service := NewProductService(productRepo, categoryRepo)

	_, err := service.Create(context.Background(), CreateProductInput{
		CategoryID: "missing-category",
		Name:       "iPhone 15",
		Price:      15000000,
		Stock:      10,
	})

	if !errors.Is(err, models.ErrInvalidProductInput) {
		t.Fatalf("expected ErrInvalidProductInput, got %v", err)
	}
}

func TestProductService_Create_DuplicateSlug(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.existsBySlugFunc = func(ctx context.Context, slug string, excludeID string) (bool, error) {
		return true, nil
	}

	categoryRepo := newFakeCategoryRepository()

	service := NewProductService(productRepo, categoryRepo)

	_, err := service.Create(context.Background(), CreateProductInput{
		CategoryID: "category-id",
		Name:       "iPhone 15",
		Price:      15000000,
		Stock:      10,
	})

	if !errors.Is(err, models.ErrProductAlreadyExists) {
		t.Fatalf("expected ErrProductAlreadyExists, got %v", err)
	}
}

func TestProductService_Create_InvalidInput(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	tests := []struct {
		name  string
		input CreateProductInput
	}{
		{
			name: "empty category id",
			input: CreateProductInput{
				CategoryID: "",
				Name:       "iPhone 15",
				Price:      15000000,
				Stock:      10,
			},
		},
		{
			name: "empty name",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       "",
				Price:      15000000,
				Stock:      10,
			},
		},
		{
			name: "too short name",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       "ab",
				Price:      15000000,
				Stock:      10,
			},
		},
		{
			name: "too long name",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       strings.Repeat("a", 151),
				Price:      15000000,
				Stock:      10,
			},
		},
		{
			name: "zero price",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Price:      0,
				Stock:      10,
			},
		},
		{
			name: "negative price",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Price:      -1,
				Stock:      10,
			},
		},
		{
			name: "negative stock",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Price:      15000000,
				Stock:      -1,
			},
		},
		{
			name: "invalid image url",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Price:      15000000,
				Stock:      10,
				ImageURL:   "not-a-url",
			},
		},
		{
			name: "unsupported image url scheme",
			input: CreateProductInput{
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Price:      15000000,
				Stock:      10,
				ImageURL:   "ftp://example.com/image.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Create(context.Background(), tt.input)
			if !errors.Is(err, models.ErrInvalidProductInput) {
				t.Fatalf("expected ErrInvalidProductInput, got %v", err)
			}
		})
	}
}

func TestProductService_Update_Success(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	product, err := service.Update(context.Background(), "product-id", UpdateProductInput{
		CategoryID:  "category-id",
		Name:        "iPhone 15 Pro",
		Description: "Apple smartphone pro",
		Price:       18000000,
		Stock:       5,
		ImageURL:    "https://example.com/iphone-pro.jpg",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != "product-id" {
		t.Fatalf("expected product-id, got %s", product.ID)
	}

	if product.Slug != "iphone-15-pro" {
		t.Fatalf("expected iphone-15-pro, got %s", product.Slug)
	}
}

func TestProductService_Update_NotFound(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.findByIDFunc = func(ctx context.Context, id string) (*models.Product, error) {
		return nil, models.ErrProductNotFound
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	_, err := service.Update(context.Background(), "missing-id", UpdateProductInput{
		CategoryID: "category-id",
		Name:       "iPhone 15",
		Price:      15000000,
		Stock:      10,
	})

	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_Delete_Success(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	err := service.Delete(context.Background(), "product-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func validPNGBytes() []byte {
	return []byte{
		0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 'I', 'H', 'D', 'R',
	}
}

func TestProductService_UploadImage_Success(t *testing.T) {
	oldUploadDir := productImageUploadDir
	productImageUploadDir = t.TempDir()
	defer func() {
		productImageUploadDir = oldUploadDir
	}()

	var savedImageURL string

	productRepo := newFakeProductRepository()
	productRepo.updateImageURLFunc = func(ctx context.Context, id string, imageURL string) (*models.Product, error) {
		savedImageURL = imageURL

		return &models.Product{
			ID:         id,
			CategoryID: "category-id",
			Name:       "iPhone 15",
			Slug:       "iphone-15",
			Price:      15000000,
			Stock:      10,
			ImageURL:   imageURL,
			IsActive:   true,
		}, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	product, err := service.UploadImage(context.Background(), UploadProductImageInput{
		ProductID:   "product-id",
		FileName:    "test.png",
		Size:        int64(len(validPNGBytes())),
		ContentType: "image/png",
		File:        bytes.NewReader(validPNGBytes()),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ImageURL == "" {
		t.Fatal("expected image_url to be set")
	}

	if !strings.HasPrefix(product.ImageURL, "/uploads/products/") {
		t.Fatalf("expected image_url to start with /uploads/products/, got %s", product.ImageURL)
	}

	fileName := strings.TrimPrefix(savedImageURL, "/uploads/products/")
	filePath := filepath.Join(productImageUploadDir, fileName)

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected uploaded file to exist, got %v", err)
	}
}

func TestProductService_UploadImage_InvalidFileType(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	_, err := service.UploadImage(context.Background(), UploadProductImageInput{
		ProductID:   "product-id",
		FileName:    "test.txt",
		Size:        int64(len([]byte("hello world"))),
		ContentType: "text/plain",
		File:        bytes.NewReader([]byte("hello world")),
	})

	if !errors.Is(err, models.ErrInvalidProductInput) {
		t.Fatalf("expected ErrInvalidProductInput, got %v", err)
	}
}

func TestProductService_UploadImage_FileTooLarge(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	_, err := service.UploadImage(context.Background(), UploadProductImageInput{
		ProductID:   "product-id",
		FileName:    "large.png",
		Size:        MaxProductImageSize + 1,
		ContentType: "image/png",
		File:        bytes.NewReader(validPNGBytes()),
	})

	if !errors.Is(err, models.ErrInvalidProductInput) {
		t.Fatalf("expected ErrInvalidProductInput, got %v", err)
	}
}

func TestProductService_UploadImage_ProductNotFound(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.findByIDFunc = func(ctx context.Context, id string) (*models.Product, error) {
		return nil, models.ErrProductNotFound
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	_, err := service.UploadImage(context.Background(), UploadProductImageInput{
		ProductID:   "missing-id",
		FileName:    "test.png",
		Size:        int64(len(validPNGBytes())),
		ContentType: "image/png",
		File:        bytes.NewReader(validPNGBytes()),
	})

	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_GetAll_WithCategoryFilter(t *testing.T) {
	productRepo := newFakeProductRepository()
	categoryRepo := newFakeCategoryRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.CategoryID != "category-id" {
			t.Fatalf("expected category-id, got %s", filter.CategoryID)
		}

		if filter.Limit != DefaultProductLimit {
			t.Fatalf("expected limit %d, got %d", DefaultProductLimit, filter.Limit)
		}

		if filter.SortBy != DefaultProductSortBy {
			t.Fatalf("expected sort_by %s, got %s", DefaultProductSortBy, filter.SortBy)
		}

		if filter.SortOrder != DefaultProductSortOrder {
			t.Fatalf("expected sort_order %s, got %s", DefaultProductSortOrder, filter.SortOrder)
		}

		if filter.Offset != 0 {
			t.Fatalf("expected offset 0, got %d", filter.Offset)
		}

		return []models.Product{
			{
				ID:         "product-id",
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Slug:       "iphone-15",
				Price:      15000000,
				Stock:      10,
				IsActive:   true,
			},
		}, 1, nil
	}

	service := NewProductService(productRepo, categoryRepo)

	result, err := service.GetAll(context.Background(), ProductListInput{
		CategoryID: "category-id",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}
}

func TestProductService_GetAll_WithSearch(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.Search != "phone" {
			t.Fatalf("expected search phone, got %s", filter.Search)
		}

		if filter.CategoryID != "" {
			t.Fatalf("expected empty category_id, got %s", filter.CategoryID)
		}

		return []models.Product{
			{
				ID:         "product-id",
				CategoryID: "category-id",
				Name:       "iPhone 15",
				Slug:       "iphone-15",
				Price:      15000000,
				Stock:      10,
				IsActive:   true,
			},
		}, 1, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		Search: " phone ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Search != "phone" {
		t.Fatalf("expected result search phone, got %s", result.Search)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}
}

func TestProductService_GetAll_WithSearchAndCategoryID(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.Search != "phone" {
			t.Fatalf("expected search phone, got %s", filter.Search)
		}

		if filter.CategoryID != "category-id" {
			t.Fatalf("expected category-id, got %s", filter.CategoryID)
		}

		return []models.Product{}, 0, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		Search:     "phone",
		CategoryID: "category-id",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Search != "phone" {
		t.Fatalf("expected result search phone, got %s", result.Search)
	}

	if result.CategoryID != "category-id" {
		t.Fatalf("expected result category-id, got %s", result.CategoryID)
	}
}

func TestProductService_GetAll_InvalidCategoryIDReturnsEmptyResult(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.CategoryID != "invalid-category-id" {
			t.Fatalf("expected invalid-category-id, got %s", filter.CategoryID)
		}

		return []models.Product{}, 0, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		CategoryID: "invalid-category-id",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Total != 0 {
		t.Fatalf("expected total 0, got %d", result.Total)
	}

	if len(result.Products) != 0 {
		t.Fatalf("expected empty products, got %d", len(result.Products))
	}
}

func TestProductService_GetAll_WithCategorySlug(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.CategorySlug != "phones" {
			t.Fatalf("expected category slug phones, got %s", filter.CategorySlug)
		}

		return []models.Product{}, 0, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		CategorySlug: " phones ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.CategorySlug != "phones" {
		t.Fatalf("expected result category slug phones, got %s", result.CategorySlug)
	}
}

func TestProductService_GetAll_WithPagination(t *testing.T) {
	productRepo := newFakeProductRepository()
	categoryRepo := newFakeCategoryRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.Limit != 10 {
			t.Fatalf("expected limit 10, got %d", filter.Limit)
		}

		if filter.Offset != 20 {
			t.Fatalf("expected offset 20, got %d", filter.Offset)
		}

		return []models.Product{}, 45, nil
	}

	service := NewProductService(productRepo, categoryRepo)

	result, err := service.GetAll(context.Background(), ProductListInput{
		Page:  3,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != 3 {
		t.Fatalf("expected page 3, got %d", result.Page)
	}

	if result.Limit != 10 {
		t.Fatalf("expected limit 10, got %d", result.Limit)
	}

	if result.Total != 45 {
		t.Fatalf("expected total 45, got %d", result.Total)
	}

	if result.TotalPages != 5 {
		t.Fatalf("expected total_pages 5, got %d", result.TotalPages)
	}
}

func TestProductService_GetAll_InvalidPagination(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	tests := []struct {
		name  string
		input ProductListInput
	}{
		{
			name: "negative limit",
			input: ProductListInput{
				Page:  1,
				Limit: -1,
			},
		},
		{
			name: "limit too large",
			input: ProductListInput{
				Page:  1,
				Limit: MaxProductLimit + 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetAll(context.Background(), tt.input)
			if !errors.Is(err, models.ErrInvalidProductInput) {
				t.Fatalf("expected ErrInvalidProductInput, got %v", err)
			}
		})
	}
}

func TestProductService_GetAll_InvalidPageDefaultsToOne(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.Offset != 0 {
			t.Fatalf("expected offset 0 for invalid page default, got %d", filter.Offset)
		}

		return []models.Product{}, 0, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		Page:  -1,
		Limit: 20,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != 1 {
		t.Fatalf("expected page default to 1, got %d", result.Page)
	}
}

func TestProductService_GetAll_WithSort(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.SortBy != ProductSortByPrice {
			t.Fatalf("expected sort_by price, got %s", filter.SortBy)
		}

		if filter.SortOrder != ProductSortOrderAsc {
			t.Fatalf("expected sort_order asc, got %s", filter.SortOrder)
		}

		return []models.Product{}, 0, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		Page:      1,
		Limit:     20,
		SortBy:    "price",
		SortOrder: "asc",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.SortBy != "price" {
		t.Fatalf("expected result sort_by price, got %s", result.SortBy)
	}

	if result.SortOrder != "asc" {
		t.Fatalf("expected result sort_order asc, got %s", result.SortOrder)
	}
}

func TestProductService_GetAll_InvalidSort(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	tests := []struct {
		name  string
		input ProductListInput
	}{
		{
			name: "invalid sort_by",
			input: ProductListInput{
				Page:      1,
				Limit:     20,
				SortBy:    "stock",
				SortOrder: "asc",
			},
		},
		{
			name: "invalid sort_order",
			input: ProductListInput{
				Page:      1,
				Limit:     20,
				SortBy:    "price",
				SortOrder: "sideways",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetAll(context.Background(), tt.input)
			if !errors.Is(err, models.ErrInvalidProductInput) {
				t.Fatalf("expected ErrInvalidProductInput, got %v", err)
			}
		})
	}
}
