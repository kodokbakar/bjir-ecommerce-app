package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type fakeProductRepository struct {
	createFunc                           func(ctx context.Context, product *models.Product) error
	findAllFunc                          func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error)
	findByIDFunc                         func(ctx context.Context, id string) (*models.Product, error)
	findBySlugFunc                       func(ctx context.Context, slug string) (*models.Product, error)
	existsBySlugFunc                     func(ctx context.Context, slug string, excludeID string) (bool, error)
	updateFunc                           func(ctx context.Context, product *models.Product) error
	deleteFunc                           func(ctx context.Context, id string) error
	updateImageURLFunc                   func(ctx context.Context, id string, imageURL string) (*models.Product, error)
	findImagesByProductIDFunc            func(ctx context.Context, productID string) ([]models.ProductImage, error)
	countImagesByProductIDFunc           func(ctx context.Context, productID string) (int, error)
	createProductImageFunc               func(ctx context.Context, image *models.ProductImage) error
	deleteProductImageFunc               func(ctx context.Context, productID string, imageID string) error
	updateProductImageSortOrderFunc      func(ctx context.Context, productID string, imageID string, sortOrder int) error
	bulkUpdateProductImageSortOrdersFunc func(ctx context.Context, productID string, images []repository.ProductImageSortOrder) error
	setPrimaryProductImageFunc           func(ctx context.Context, productID string, imageID string) (*models.ProductImage, error)
	syncProductPrimaryImageURLFunc       func(ctx context.Context, productID string) error
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
		findImagesByProductIDFunc: func(ctx context.Context, productID string) ([]models.ProductImage, error) {
			return []models.ProductImage{}, nil
		},
		countImagesByProductIDFunc: func(ctx context.Context, productID string) (int, error) {
			return 0, nil
		},
		createProductImageFunc: func(ctx context.Context, image *models.ProductImage) error {
			image.ID = "image-id"
			image.CreatedAt = now
			image.UpdatedAt = now
			return nil
		},
		deleteProductImageFunc: func(ctx context.Context, productID string, imageID string) error {
			return nil
		},
		updateProductImageSortOrderFunc: func(ctx context.Context, productID string, imageID string, sortOrder int) error {
			return nil
		},
		bulkUpdateProductImageSortOrdersFunc: func(ctx context.Context, productID string, images []repository.ProductImageSortOrder) error {
			return nil
		},
		setPrimaryProductImageFunc: func(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
			return &models.ProductImage{
				ID:        imageID,
				ProductID: productID,
				ImageURL:  "/uploads/products/test.png",
				SortOrder: 0,
				IsPrimary: true,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
		syncProductPrimaryImageURLFunc: func(ctx context.Context, productID string) error {
			return nil
		},
	}
}

func (f *fakeProductRepository) Create(ctx context.Context, product *models.Product) error {
	return f.createFunc(ctx, product)
}

func TestProductService_Create_InvalidatesProductListCache(t *testing.T) {
	productRepo := newFakeProductRepository()
	cache := newFakeProductCache()

	var invalidated bool
	cache.invalidateFunc = func(ctx context.Context) error {
		invalidated = true
		return nil
	}

	service := NewProductServiceWithCache(productRepo, newFakeCategoryRepository(), cache)

	_, err := service.Create(context.Background(), CreateProductInput{
		CategoryID: "category-id",
		Name:       "iPhone 15",
		Price:      15000000,
		Stock:      10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !invalidated {
		t.Fatal("expected product list cache to be invalidated")
	}
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

func TestProductService_Update_InvalidatesProductListCache(t *testing.T) {
	productRepo := newFakeProductRepository()
	cache := newFakeProductCache()

	var invalidated bool
	cache.invalidateFunc = func(ctx context.Context) error {
		invalidated = true
		return nil
	}

	service := NewProductServiceWithCache(productRepo, newFakeCategoryRepository(), cache)

	_, err := service.Update(context.Background(), "product-id", UpdateProductInput{
		CategoryID: "category-id",
		Name:       "iPhone 15 Pro",
		Price:      18000000,
		Stock:      5,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !invalidated {
		t.Fatal("expected product list cache to be invalidated")
	}
}

func (f *fakeProductRepository) UpdateImageURL(ctx context.Context, id string, imageURL string) (*models.Product, error) {
	return f.updateImageURLFunc(ctx, id, imageURL)
}

func (f *fakeProductRepository) Delete(ctx context.Context, id string) error {
	return f.deleteFunc(ctx, id)
}

func (f *fakeProductRepository) FindImagesByProductID(ctx context.Context, productID string) ([]models.ProductImage, error) {
	return f.findImagesByProductIDFunc(ctx, productID)
}

func (f *fakeProductRepository) CountImagesByProductID(ctx context.Context, productID string) (int, error) {
	return f.countImagesByProductIDFunc(ctx, productID)
}

func (f *fakeProductRepository) CreateProductImage(ctx context.Context, image *models.ProductImage) error {
	return f.createProductImageFunc(ctx, image)
}

func (f *fakeProductRepository) DeleteProductImage(ctx context.Context, productID string, imageID string) error {
	return f.deleteProductImageFunc(ctx, productID, imageID)
}

func (f *fakeProductRepository) UpdateProductImageSortOrder(ctx context.Context, productID string, imageID string, sortOrder int) error {
	return f.updateProductImageSortOrderFunc(ctx, productID, imageID, sortOrder)
}

func (f *fakeProductRepository) BulkUpdateProductImageSortOrders(ctx context.Context, productID string, images []repository.ProductImageSortOrder) error {
	return f.bulkUpdateProductImageSortOrdersFunc(ctx, productID, images)
}

func (f *fakeProductRepository) SetPrimaryProductImage(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
	return f.setPrimaryProductImageFunc(ctx, productID, imageID)
}

func (f *fakeProductRepository) SyncProductPrimaryImageURL(ctx context.Context, productID string) error {
	return f.syncProductPrimaryImageURLFunc(ctx, productID)
}

func TestProductService_Delete_InvalidatesProductListCache(t *testing.T) {
	productRepo := newFakeProductRepository()
	cache := newFakeProductCache()

	var invalidated bool
	cache.invalidateFunc = func(ctx context.Context) error {
		invalidated = true
		return nil
	}

	service := NewProductServiceWithCache(productRepo, newFakeCategoryRepository(), cache)

	err := service.Delete(context.Background(), "product-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !invalidated {
		t.Fatal("expected product list cache to be invalidated")
	}
}

type fakeProductCache struct {
	getFunc        func(ctx context.Context, key string) (*ProductListResult, error)
	setFunc        func(ctx context.Context, key string, result *ProductListResult, ttl time.Duration) error
	invalidateFunc func(ctx context.Context) error
}

func newFakeProductCache() *fakeProductCache {
	return &fakeProductCache{
		getFunc: func(ctx context.Context, key string) (*ProductListResult, error) {
			return nil, nil
		},
		setFunc: func(ctx context.Context, key string, result *ProductListResult, ttl time.Duration) error {
			return nil
		},
		invalidateFunc: func(ctx context.Context) error {
			return nil
		},
	}
}

func (f *fakeProductCache) GetProductList(ctx context.Context, key string) (*ProductListResult, error) {
	return f.getFunc(ctx, key)
}

func (f *fakeProductCache) SetProductList(ctx context.Context, key string, result *ProductListResult, ttl time.Duration) error {
	return f.setFunc(ctx, key, result, ttl)
}

func (f *fakeProductCache) InvalidateProductList(ctx context.Context) error {
	return f.invalidateFunc(ctx)
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

func TestProductService_GetAll_CacheHit(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		t.Fatal("expected repository not to be called on cache hit")
		return nil, 0, nil
	}

	cache := newFakeProductCache()
	cache.getFunc = func(ctx context.Context, key string) (*ProductListResult, error) {
		if !strings.Contains(key, "products:list") {
			t.Fatalf("expected products list cache key, got %s", key)
		}

		return &ProductListResult{
			Products: []models.Product{
				{
					ID:         "cached-product-id",
					CategoryID: "category-id",
					Name:       "Cached Product",
					Slug:       "cached-product",
					Price:      10000,
					Stock:      1,
					IsActive:   true,
				},
			},
			Page:       1,
			Limit:      20,
			Total:      1,
			TotalPages: 1,
			SortBy:     DefaultProductSortBy,
			SortOrder:  DefaultProductSortOrder,
		}, nil
	}

	service := NewProductServiceWithCache(productRepo, newFakeCategoryRepository(), cache)

	result, err := service.GetAll(context.Background(), ProductListInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Products) != 1 {
		t.Fatalf("expected 1 cached product, got %d", len(result.Products))
	}

	if result.Products[0].ID != "cached-product-id" {
		t.Fatalf("expected cached-product-id, got %s", result.Products[0].ID)
	}
}

func TestProductService_GetAll_CacheMissSetsCache(t *testing.T) {
	productRepo := newFakeProductRepository()

	cache := newFakeProductCache()

	var cacheSetCalled bool
	var cacheKey string

	cache.getFunc = func(ctx context.Context, key string) (*ProductListResult, error) {
		return nil, nil
	}

	cache.setFunc = func(ctx context.Context, key string, result *ProductListResult, ttl time.Duration) error {
		cacheSetCalled = true
		cacheKey = key

		if ttl != ProductListCacheTTL {
			t.Fatalf("expected ttl %v, got %v", ProductListCacheTTL, ttl)
		}

		if result.Total != 1 {
			t.Fatalf("expected total 1, got %d", result.Total)
		}

		return nil
	}

	service := NewProductServiceWithCache(productRepo, newFakeCategoryRepository(), cache)

	result, err := service.GetAll(context.Background(), ProductListInput{
		CategoryID: "category-id",
		Search:     "phone",
		Page:       2,
		Limit:      10,
		SortBy:     "price",
		SortOrder:  "asc",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}

	if !cacheSetCalled {
		t.Fatal("expected cache set to be called")
	}

	if !strings.Contains(cacheKey, "category_id=category-id") {
		t.Fatalf("expected cache key to contain category_id, got %s", cacheKey)
	}

	if !strings.Contains(cacheKey, "search=phone") {
		t.Fatalf("expected cache key to contain search, got %s", cacheKey)
	}

	if !strings.Contains(cacheKey, "sort_by=price") {
		t.Fatalf("expected cache key to contain sort_by, got %s", cacheKey)
	}
}

func TestProductService_GetAll_CacheErrorFallsBackToRepository(t *testing.T) {
	productRepo := newFakeProductRepository()

	var repositoryCalled bool
	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		repositoryCalled = true
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

	cache := newFakeProductCache()
	cache.getFunc = func(ctx context.Context, key string) (*ProductListResult, error) {
		return nil, fmt.Errorf("redis unavailable")
	}

	service := NewProductServiceWithCache(productRepo, newFakeCategoryRepository(), cache)

	result, err := service.GetAll(context.Background(), ProductListInput{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repositoryCalled {
		t.Fatal("expected repository to be called when cache returns error")
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
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
	var createdImage models.ProductImage

	productRepo := newFakeProductRepository()
	productRepo.createProductImageFunc = func(ctx context.Context, image *models.ProductImage) error {
		image.ID = "image-id"
		savedImageURL = image.ImageURL
		createdImage = *image
		return nil
	}
	productRepo.setPrimaryProductImageFunc = func(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
		createdImage.IsPrimary = true
		return &createdImage, nil
	}
	productRepo.findImagesByProductIDFunc = func(ctx context.Context, productID string) ([]models.ProductImage, error) {
		createdImage.IsPrimary = true
		return []models.ProductImage{createdImage}, nil
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

	if len(product.Images) != 1 {
		t.Fatalf("expected 1 image, got %d", len(product.Images))
	}

	fileName := strings.TrimPrefix(savedImageURL, "/uploads/products/")
	filePath := filepath.Join(productImageUploadDir, fileName)

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected uploaded file to exist, got %v", err)
	}
}

func TestProductService_GetByID_IncludesImagesAndPrimaryImageURL(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.findImagesByProductIDFunc = func(ctx context.Context, productID string) ([]models.ProductImage, error) {
		return []models.ProductImage{
			{
				ID:        "image-1",
				ProductID: productID,
				ImageURL:  "/uploads/products/secondary.png",
				SortOrder: 1,
				IsPrimary: false,
			},
			{
				ID:        "image-2",
				ProductID: productID,
				ImageURL:  "/uploads/products/primary.png",
				SortOrder: 2,
				IsPrimary: true,
			},
		}, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	product, err := service.GetByID(context.Background(), "product-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ImageURL != "/uploads/products/primary.png" {
		t.Fatalf("expected primary image_url, got %s", product.ImageURL)
	}

	if len(product.Images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(product.Images))
	}
}

func TestProductService_GetImages_ProductNotFound(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.findByIDFunc = func(ctx context.Context, id string) (*models.Product, error) {
		return nil, models.ErrProductNotFound
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	_, err := service.GetImages(context.Background(), "missing-id")
	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestProductService_UploadGalleryImage_Success(t *testing.T) {
	oldUploadDir := productImageUploadDir
	productImageUploadDir = t.TempDir()
	defer func() {
		productImageUploadDir = oldUploadDir
	}()

	var savedImageURL string
	var createdImage models.ProductImage

	productRepo := newFakeProductRepository()
	productRepo.countImagesByProductIDFunc = func(ctx context.Context, productID string) (int, error) {
		return 0, nil
	}
	productRepo.createProductImageFunc = func(ctx context.Context, image *models.ProductImage) error {
		image.ID = "image-id"
		savedImageURL = image.ImageURL
		createdImage = *image
		return nil
	}
	productRepo.setPrimaryProductImageFunc = func(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
		createdImage.IsPrimary = true
		return &createdImage, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	image, err := service.UploadGalleryImage(context.Background(), UploadProductImageInput{
		ProductID:   "product-id",
		FileName:    "test.png",
		Size:        int64(len(validPNGBytes())),
		ContentType: "image/png",
		File:        bytes.NewReader(validPNGBytes()),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if image.ID != "image-id" {
		t.Fatalf("expected image-id, got %s", image.ID)
	}

	if !image.IsPrimary {
		t.Fatal("expected first image to be primary")
	}

	if image.SortOrder != 0 {
		t.Fatalf("expected sort_order 0, got %d", image.SortOrder)
	}

	fileName := strings.TrimPrefix(savedImageURL, "/uploads/products/")
	filePath := filepath.Join(productImageUploadDir, fileName)

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected uploaded file to exist, got %v", err)
	}
}

func TestProductService_UploadGalleryImage_MaxImagesExceeded(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.countImagesByProductIDFunc = func(ctx context.Context, productID string) (int, error) {
		return MaxProductImagesPerProduct, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	_, err := service.UploadGalleryImage(context.Background(), UploadProductImageInput{
		ProductID:   "product-id",
		FileName:    "test.png",
		Size:        int64(len(validPNGBytes())),
		ContentType: "image/png",
		File:        bytes.NewReader(validPNGBytes()),
	})

	if !errors.Is(err, models.ErrInvalidProductInput) {
		t.Fatalf("expected ErrInvalidProductInput, got %v", err)
	}
}

func TestProductService_DeleteImage_Success(t *testing.T) {
	oldUploadDir := productImageUploadDir
	productImageUploadDir = t.TempDir()
	defer func() {
		productImageUploadDir = oldUploadDir
	}()

	filePath := filepath.Join(productImageUploadDir, "test.png")
	if err := os.WriteFile(filePath, []byte("fake image"), 0644); err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	productRepo := newFakeProductRepository()

	var deleted bool
	var synced bool

	productRepo.findImagesByProductIDFunc = func(ctx context.Context, productID string) ([]models.ProductImage, error) {
		if productID != "product-id" {
			t.Fatalf("expected product-id, got %s", productID)
		}

		return []models.ProductImage{
			{
				ID:        "image-id",
				ProductID: productID,
				ImageURL:  "/uploads/products/test.png",
				SortOrder: 0,
				IsPrimary: true,
			},
		}, nil
	}

	productRepo.deleteProductImageFunc = func(ctx context.Context, productID string, imageID string) error {
		if productID != "product-id" {
			t.Fatalf("expected product-id, got %s", productID)
		}

		if imageID != "image-id" {
			t.Fatalf("expected image-id, got %s", imageID)
		}

		deleted = true
		return nil
	}

	productRepo.syncProductPrimaryImageURLFunc = func(ctx context.Context, productID string) error {
		if productID != "product-id" {
			t.Fatalf("expected product-id, got %s", productID)
		}

		synced = true
		return nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	err := service.DeleteImage(context.Background(), "product-id", "image-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !deleted {
		t.Fatal("expected image to be deleted")
	}

	if !synced {
		t.Fatal("expected product primary image_url to be synced")
	}

	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected physical image file to be removed, got err: %v", err)
	}
}

func TestProductService_DeleteImage_DoesNotRemoveExternalImageURL(t *testing.T) {
	oldUploadDir := productImageUploadDir
	productImageUploadDir = t.TempDir()
	defer func() {
		productImageUploadDir = oldUploadDir
	}()

	filePath := filepath.Join(productImageUploadDir, "test.png")
	if err := os.WriteFile(filePath, []byte("fake image"), 0644); err != nil {
		t.Fatalf("failed to create test image: %v", err)
	}

	productRepo := newFakeProductRepository()

	productRepo.findImagesByProductIDFunc = func(ctx context.Context, productID string) ([]models.ProductImage, error) {
		return []models.ProductImage{
			{
				ID:        "image-id",
				ProductID: productID,
				ImageURL:  "https://example.com/uploads/products/test.png",
				SortOrder: 0,
				IsPrimary: true,
			},
		}, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	err := service.DeleteImage(context.Background(), "product-id", "image-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected local file to remain for external image URL, got %v", err)
	}
}

func TestProductService_DeleteImage_ImageNotFound(t *testing.T) {
	productRepo := newFakeProductRepository()

	var deleteCalled bool

	productRepo.findImagesByProductIDFunc = func(ctx context.Context, productID string) ([]models.ProductImage, error) {
		return []models.ProductImage{
			{
				ID:        "other-image-id",
				ProductID: productID,
				ImageURL:  "/uploads/products/other.png",
			},
		}, nil
	}

	productRepo.deleteProductImageFunc = func(ctx context.Context, productID string, imageID string) error {
		deleteCalled = true
		return nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	err := service.DeleteImage(context.Background(), "product-id", "missing-image-id")
	if !errors.Is(err, models.ErrProductImageNotFound) {
		t.Fatalf("expected ErrProductImageNotFound, got %v", err)
	}

	if deleteCalled {
		t.Fatal("expected delete repository not to be called when image is missing")
	}
}

func TestProductService_ReorderImages_Success(t *testing.T) {
	productRepo := newFakeProductRepository()

	var bulkCalled bool
	var receivedSortOrders []repository.ProductImageSortOrder

	productRepo.bulkUpdateProductImageSortOrdersFunc = func(ctx context.Context, productID string, images []repository.ProductImageSortOrder) error {
		if productID != "product-id" {
			t.Fatalf("expected product-id, got %s", productID)
		}

		bulkCalled = true
		receivedSortOrders = images
		return nil
	}

	productRepo.updateProductImageSortOrderFunc = func(ctx context.Context, productID string, imageID string, sortOrder int) error {
		t.Fatal("expected bulk reorder to be used, but UpdateProductImageSortOrder was called")
		return nil
	}

	productRepo.findImagesByProductIDFunc = func(ctx context.Context, productID string) ([]models.ProductImage, error) {
		return []models.ProductImage{
			{
				ID:        "image-2",
				ProductID: productID,
				ImageURL:  "/uploads/products/2.png",
				SortOrder: 0,
			},
			{
				ID:        "image-1",
				ProductID: productID,
				ImageURL:  "/uploads/products/1.png",
				SortOrder: 1,
			},
		}, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	images, err := service.ReorderImages(context.Background(), ReorderProductImagesInput{
		ProductID: "product-id",
		Images: []ReorderProductImageInput{
			{ID: "image-1", SortOrder: 1},
			{ID: "image-2", SortOrder: 0},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !bulkCalled {
		t.Fatal("expected bulk reorder to be called")
	}

	if len(receivedSortOrders) != 2 {
		t.Fatalf("expected 2 sort orders, got %d", len(receivedSortOrders))
	}

	if receivedSortOrders[0].ID != "image-1" || receivedSortOrders[0].SortOrder != 1 {
		t.Fatalf("unexpected first sort order: %#v", receivedSortOrders[0])
	}

	if receivedSortOrders[1].ID != "image-2" || receivedSortOrders[1].SortOrder != 0 {
		t.Fatalf("unexpected second sort order: %#v", receivedSortOrders[1])
	}

	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(images))
	}
}

func TestProductService_ReorderImages_DuplicateID(t *testing.T) {
	service := NewProductService(newFakeProductRepository(), newFakeCategoryRepository())

	_, err := service.ReorderImages(context.Background(), ReorderProductImagesInput{
		ProductID: "product-id",
		Images: []ReorderProductImageInput{
			{ID: "image-1", SortOrder: 0},
			{ID: "image-1", SortOrder: 1},
		},
	})

	if !errors.Is(err, models.ErrInvalidProductInput) {
		t.Fatalf("expected ErrInvalidProductInput, got %v", err)
	}
}

func TestProductService_SetPrimaryImage_Success(t *testing.T) {
	productRepo := newFakeProductRepository()

	var called bool

	productRepo.setPrimaryProductImageFunc = func(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
		if productID != "product-id" {
			t.Fatalf("expected product-id, got %s", productID)
		}

		if imageID != "image-id" {
			t.Fatalf("expected image-id, got %s", imageID)
		}

		called = true

		return &models.ProductImage{
			ID:        imageID,
			ProductID: productID,
			ImageURL:  "/uploads/products/primary.png",
			SortOrder: 1,
			IsPrimary: true,
		}, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	image, err := service.SetPrimaryImage(context.Background(), "product-id", "image-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !called {
		t.Fatal("expected SetPrimaryProductImage to be called")
	}

	if !image.IsPrimary {
		t.Fatal("expected image to be primary")
	}
}

func TestProductService_UploadImage_InvalidatesProductListCache(t *testing.T) {
	oldUploadDir := productImageUploadDir
	productImageUploadDir = t.TempDir()
	defer func() {
		productImageUploadDir = oldUploadDir
	}()

	productRepo := newFakeProductRepository()
	cache := newFakeProductCache()

	var invalidated bool
	cache.invalidateFunc = func(ctx context.Context) error {
		invalidated = true
		return nil
	}

	service := NewProductServiceWithCache(productRepo, newFakeCategoryRepository(), cache)

	_, err := service.UploadImage(context.Background(), UploadProductImageInput{
		ProductID:   "product-id",
		FileName:    "test.png",
		Size:        int64(len(validPNGBytes())),
		ContentType: "image/png",
		File:        bytes.NewReader(validPNGBytes()),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !invalidated {
		t.Fatal("expected product list cache to be invalidated")
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

func TestProductService_UploadImage_SyncFailureCleansUpCreatedImage(t *testing.T) {
	oldUploadDir := productImageUploadDir
	productImageUploadDir = t.TempDir()
	defer func() {
		productImageUploadDir = oldUploadDir
	}()

	productRepo := newFakeProductRepository()

	var savedImageURL string
	var deleted bool

	productRepo.countImagesByProductIDFunc = func(ctx context.Context, productID string) (int, error) {
		return 1, nil
	}

	productRepo.createProductImageFunc = func(ctx context.Context, image *models.ProductImage) error {
		image.ID = "image-id"
		savedImageURL = image.ImageURL
		return nil
	}

	productRepo.syncProductPrimaryImageURLFunc = func(ctx context.Context, productID string) error {
		return errors.New("sync failed")
	}

	productRepo.deleteProductImageFunc = func(ctx context.Context, productID string, imageID string) error {
		if productID != "product-id" {
			t.Fatalf("expected product-id, got %s", productID)
		}

		if imageID != "image-id" {
			t.Fatalf("expected image-id, got %s", imageID)
		}

		deleted = true
		return nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	_, err := service.UploadImage(context.Background(), UploadProductImageInput{
		ProductID:   "product-id",
		FileName:    "test.png",
		Size:        int64(len(validPNGBytes())),
		ContentType: "image/png",
		File:        bytes.NewReader(validPNGBytes()),
	})
	if err == nil {
		t.Fatal("expected sync error")
	}

	if !deleted {
		t.Fatal("expected created image record to be deleted during cleanup")
	}

	fileName := strings.TrimPrefix(savedImageURL, "/uploads/products/")
	filePath := filepath.Join(productImageUploadDir, fileName)

	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected created image file to be removed, got err: %v", err)
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

func TestProductService_GetAll_LimitTooLargeCappedAtMax(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.Limit != MaxProductLimit {
			t.Fatalf("expected max limit %d, got %d", MaxProductLimit, filter.Limit)
		}

		if filter.Offset != 0 {
			t.Fatalf("expected offset 0, got %d", filter.Offset)
		}

		return []models.Product{}, 0, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		Page:  1,
		Limit: MaxProductLimit + 1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Limit != MaxProductLimit {
		t.Fatalf("expected max limit %d, got %d", MaxProductLimit, result.Limit)
	}
}

func TestProductService_GetAll_InvalidLimitDefaults(t *testing.T) {
	productRepo := newFakeProductRepository()

	productRepo.findAllFunc = func(ctx context.Context, filter repository.ProductListFilter) ([]models.Product, int, error) {
		if filter.Limit != DefaultProductLimit {
			t.Fatalf("expected default limit %d, got %d", DefaultProductLimit, filter.Limit)
		}

		if filter.Offset != 0 {
			t.Fatalf("expected offset 0, got %d", filter.Offset)
		}

		return []models.Product{}, 0, nil
	}

	service := NewProductService(productRepo, newFakeCategoryRepository())

	result, err := service.GetAll(context.Background(), ProductListInput{
		Page:  1,
		Limit: -1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Limit != DefaultProductLimit {
		t.Fatalf("expected default limit %d, got %d", DefaultProductLimit, result.Limit)
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
