package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type CreateProductInput struct {
	CategoryID  string
	Name        string
	Description string
	Price       float64
	Stock       int
	ImageURL    string
}

type UpdateProductInput struct {
	CategoryID  string
	Name        string
	Description string
	Price       float64
	Stock       int
	ImageURL    string
}

type UploadProductImageInput struct {
	ProductID   string
	FileName    string
	Size        int64
	ContentType string
	File        io.Reader
}

type ReorderProductImageInput struct {
	ID        string `json:"id"`
	SortOrder int    `json:"sort_order"`
}

type ReorderProductImagesInput struct {
	ProductID string
	Images    []ReorderProductImageInput
}

type ProductService interface {
	Create(ctx context.Context, input CreateProductInput) (*models.Product, error)
	GetAll(ctx context.Context, input ProductListInput) (*ProductListResult, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	Update(ctx context.Context, id string, input UpdateProductInput) (*models.Product, error)
	UploadImage(ctx context.Context, input UploadProductImageInput) (*models.Product, error)
	Delete(ctx context.Context, id string) error

	GetImages(ctx context.Context, productID string) ([]models.ProductImage, error)
	UploadGalleryImage(ctx context.Context, input UploadProductImageInput) (*models.ProductImage, error)
	DeleteImage(ctx context.Context, productID string, imageID string) error
	ReorderImages(ctx context.Context, input ReorderProductImagesInput) ([]models.ProductImage, error)
	SetPrimaryImage(ctx context.Context, productID string, imageID string) (*models.ProductImage, error)
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	productCache ProductCache
}

type ProductListInput struct {
	CategoryID   string
	CategorySlug string
	Search       string
	Page         int
	Limit        int
	SortBy       string
	SortOrder    string
}

type ProductListResult struct {
	Products     []models.Product
	Page         int
	Limit        int
	Total        int
	TotalPages   int
	SortBy       string
	SortOrder    string
	CategoryID   string
	CategorySlug string
	Search       string
}

const MaxProductImageSize int64 = 5 << 20 // 5MB file size limit
const MaxProductImagesPerProduct = 10

const (
	DefaultProductPage  = 1
	DefaultProductLimit = 20
	MaxProductLimit     = 100

	DefaultProductSortBy    = "created_at"
	DefaultProductSortOrder = "desc"

	ProductSortByCreatedAt = "created_at"
	ProductSortByPrice     = "price"
	ProductSortByName      = "name"

	ProductSortOrderAsc  = "asc"
	ProductSortOrderDesc = "desc"
)

var (
	productImageUploadDir  = filepath.Join("uploads", "products")
	productImagePublicPath = "/uploads/products"
)

func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
) ProductService {
	return NewProductServiceWithCache(productRepo, categoryRepo, nil)
}

func NewProductServiceWithCache(
	productRepo repository.ProductRepository,
	categoryRepo repository.CategoryRepository,
	productCache ProductCache,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		productCache: productCache,
	}
}

func (s *productService) Create(ctx context.Context, input CreateProductInput) (*models.Product, error) {
	categoryID, name, description, imageURL, err := normalizeProductInput(
		input.CategoryID,
		input.Name,
		input.Description,
		input.ImageURL,
		input.Price,
		input.Stock,
	)
	if err != nil {
		return nil, err
	}

	if _, err := s.categoryRepo.FindByID(ctx, categoryID); err != nil {
		if err == models.ErrCategoryNotFound {
			return nil, fmt.Errorf("%w: category not found", models.ErrInvalidProductInput)
		}

		return nil, err
	}

	slug := slugify(name)
	if slug == "" {
		return nil, fmt.Errorf("%w: product slug is invalid", models.ErrInvalidProductInput)
	}

	slugExists, err := s.productRepo.ExistsBySlug(ctx, slug, "")
	if err != nil {
		return nil, err
	}

	if slugExists {
		return nil, models.ErrProductAlreadyExists
	}

	product := &models.Product{
		CategoryID:  categoryID,
		Name:        name,
		Slug:        slug,
		Description: description,
		Price:       input.Price,
		Stock:       input.Stock,
		ImageURL:    imageURL,
		IsActive:    true,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	s.invalidateProductListCache(ctx)

	return product, nil
}

func (s *productService) GetAll(ctx context.Context, input ProductListInput) (*ProductListResult, error) {
	categoryID := strings.TrimSpace(input.CategoryID)
	categorySlug := strings.TrimSpace(input.CategorySlug)
	search := strings.TrimSpace(input.Search)

	page := input.Page
	if page < 1 {
		page = DefaultProductPage
	}

	limit := input.Limit
	if limit < 1 {
		limit = DefaultProductLimit
	}

	if limit > MaxProductLimit {
		limit = MaxProductLimit
	}

	sortBy, sortOrder, err := normalizeProductSort(input.SortBy, input.SortOrder)
	if err != nil {
		return nil, err
	}

	normalizedInput := ProductListInput{
		CategoryID:   categoryID,
		CategorySlug: categorySlug,
		Search:       search,
		Page:         page,
		Limit:        limit,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	}

	cacheKey := buildProductListCacheKey(normalizedInput)

	if s.productCache != nil {
		cachedResult, err := s.productCache.GetProductList(ctx, cacheKey)
		if err == nil && cachedResult != nil {
			return cachedResult, nil
		}
	}

	offset := (page - 1) * limit

	products, total, err := s.productRepo.FindAll(ctx, repository.ProductListFilter{
		CategoryID:   categoryID,
		CategorySlug: categorySlug,
		Search:       search,
		Limit:        limit,
		Offset:       offset,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	result := &ProductListResult{
		Products:     products,
		Page:         page,
		Limit:        limit,
		Total:        total,
		TotalPages:   totalPages,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		CategoryID:   categoryID,
		CategorySlug: categorySlug,
		Search:       search,
	}

	if s.productCache != nil {
		_ = s.productCache.SetProductList(ctx, cacheKey, result, ProductListCacheTTL)
	}

	return result, nil
}

func (s *productService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	product, err := s.productRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.attachProductImages(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("%w: product slug is required", models.ErrInvalidProductInput)
	}

	product, err := s.productRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if err := s.attachProductImages(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) Update(ctx context.Context, id string, input UpdateProductInput) (*models.Product, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	if _, err := s.productRepo.FindByID(ctx, id); err != nil {
		return nil, err
	}

	categoryID, name, description, imageURL, err := normalizeProductInput(
		input.CategoryID,
		input.Name,
		input.Description,
		input.ImageURL,
		input.Price,
		input.Stock,
	)
	if err != nil {
		return nil, err
	}

	if _, err := s.categoryRepo.FindByID(ctx, categoryID); err != nil {
		if err == models.ErrCategoryNotFound {
			return nil, fmt.Errorf("%w: category not found", models.ErrInvalidProductInput)
		}

		return nil, err
	}

	slug := slugify(name)
	if slug == "" {
		return nil, fmt.Errorf("%w: product slug is invalid", models.ErrInvalidProductInput)
	}

	slugExists, err := s.productRepo.ExistsBySlug(ctx, slug, id)
	if err != nil {
		return nil, err
	}

	if slugExists {
		return nil, models.ErrProductAlreadyExists
	}

	product := &models.Product{
		ID:          id,
		CategoryID:  categoryID,
		Name:        name,
		Slug:        slug,
		Description: description,
		Price:       input.Price,
		Stock:       input.Stock,
		ImageURL:    imageURL,
		IsActive:    true,
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	s.invalidateProductListCache(ctx)

	return product, nil
}

func (s *productService) UploadImage(ctx context.Context, input UploadProductImageInput) (*models.Product, error) {
	image, err := s.UploadGalleryImage(ctx, input)
	if err != nil {
		return nil, err
	}

	if image.IsPrimary {
		if _, err := s.productRepo.SetPrimaryProductImage(ctx, image.ProductID, image.ID); err != nil {
			s.cleanupCreatedProductImage(ctx, image)
			return nil, err
		}
	} else {
		if err := s.productRepo.SyncProductPrimaryImageURL(ctx, image.ProductID); err != nil {
			s.cleanupCreatedProductImage(ctx, image)
			return nil, err
		}
	}

	product, err := s.GetByID(ctx, image.ProductID)
	if err != nil {
		return nil, err
	}

	s.invalidateProductListCache(ctx)

	return product, nil
}

func (s *productService) GetImages(ctx context.Context, productID string) ([]models.ProductImage, error) {
	productID = strings.TrimSpace(productID)
	if productID == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return nil, err
	}

	return s.productRepo.FindImagesByProductID(ctx, productID)
}

func (s *productService) UploadGalleryImage(ctx context.Context, input UploadProductImageInput) (*models.ProductImage, error) {
	productID := strings.TrimSpace(input.ProductID)
	if productID == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	if input.File == nil {
		return nil, fmt.Errorf("%w: image file is required", models.ErrInvalidProductInput)
	}

	if input.Size <= 0 {
		return nil, fmt.Errorf("%w: image file is required", models.ErrInvalidProductInput)
	}

	if input.Size > MaxProductImageSize {
		return nil, fmt.Errorf("%w: image file must be at most 5MB", models.ErrInvalidProductInput)
	}

	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return nil, err
	}

	imageCount, err := s.productRepo.CountImagesByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if imageCount >= MaxProductImagesPerProduct {
		return nil, fmt.Errorf("%w: product images must be at most %d", models.ErrInvalidProductInput, MaxProductImagesPerProduct)
	}

	data, err := io.ReadAll(io.LimitReader(input.File, MaxProductImageSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	if int64(len(data)) > MaxProductImageSize {
		return nil, fmt.Errorf("%w: image file must be at most 5MB", models.ErrInvalidProductInput)
	}

	extension, err := detectProductImageExtension(data)
	if err != nil {
		return nil, err
	}

	fileName, err := generateProductImageFileName(extension)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(productImageUploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	filePath := filepath.Join(productImageUploadDir, fileName)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to save image file: %w", err)
	}

	imageURL := productImagePublicPath + "/" + fileName

	image := &models.ProductImage{
		ProductID: productID,
		ImageURL:  imageURL,
		SortOrder: imageCount,
		IsPrimary: imageCount == 0,
	}

	if err := s.productRepo.CreateProductImage(ctx, image); err != nil {
		_ = os.Remove(filePath)
		return nil, err
	}

	if image.IsPrimary {
		if _, err := s.productRepo.SetPrimaryProductImage(ctx, productID, image.ID); err != nil {
			_ = os.Remove(filePath)
			return nil, err
		}
	}

	s.invalidateProductListCache(ctx)

	return image, nil
}

func (s *productService) DeleteImage(ctx context.Context, productID string, imageID string) error {
	productID = strings.TrimSpace(productID)
	imageID = strings.TrimSpace(imageID)

	if productID == "" {
		return fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	if imageID == "" {
		return fmt.Errorf("%w: image id is required", models.ErrInvalidProductInput)
	}

	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return err
	}

	images, err := s.productRepo.FindImagesByProductID(ctx, productID)
	if err != nil {
		return err
	}

	var imageURL string
	for _, image := range images {
		if image.ID == imageID {
			imageURL = image.ImageURL
			break
		}
	}

	if imageURL == "" {
		return models.ErrProductImageNotFound
	}

	if err := s.productRepo.DeleteProductImage(ctx, productID, imageID); err != nil {
		return err
	}

	if err := s.productRepo.SyncProductPrimaryImageURL(ctx, productID); err != nil {
		return err
	}

	removeProductImageFile(imageURL)

	s.invalidateProductListCache(ctx)

	return nil
}

func (s *productService) ReorderImages(ctx context.Context, input ReorderProductImagesInput) ([]models.ProductImage, error) {
	productID := strings.TrimSpace(input.ProductID)
	if productID == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	if len(input.Images) == 0 {
		return nil, fmt.Errorf("%w: images is required", models.ErrInvalidProductInput)
	}

	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return nil, err
	}

	seen := make(map[string]bool, len(input.Images))
	sortOrders := make([]repository.ProductImageSortOrder, 0, len(input.Images))

	for _, image := range input.Images {
		imageID := strings.TrimSpace(image.ID)
		if imageID == "" {
			return nil, fmt.Errorf("%w: image id is required", models.ErrInvalidProductInput)
		}

		if seen[imageID] {
			return nil, fmt.Errorf("%w: duplicate image id", models.ErrInvalidProductInput)
		}

		if image.SortOrder < 0 {
			return nil, fmt.Errorf("%w: sort_order must be greater than or equal to 0", models.ErrInvalidProductInput)
		}

		seen[imageID] = true

		sortOrders = append(sortOrders, repository.ProductImageSortOrder{
			ID:        imageID,
			SortOrder: image.SortOrder,
		})
	}

	if err := s.productRepo.BulkUpdateProductImageSortOrders(ctx, productID, sortOrders); err != nil {
		return nil, err
	}

	images, err := s.productRepo.FindImagesByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	s.invalidateProductListCache(ctx)

	return images, nil
}

func (s *productService) SetPrimaryImage(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
	productID = strings.TrimSpace(productID)
	imageID = strings.TrimSpace(imageID)

	if productID == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	if imageID == "" {
		return nil, fmt.Errorf("%w: image id is required", models.ErrInvalidProductInput)
	}

	if _, err := s.productRepo.FindByID(ctx, productID); err != nil {
		return nil, err
	}

	image, err := s.productRepo.SetPrimaryProductImage(ctx, productID, imageID)
	if err != nil {
		return nil, err
	}

	s.invalidateProductListCache(ctx)

	return image, nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	s.invalidateProductListCache(ctx)

	return s.productRepo.Delete(ctx, id)
}

func (s *productService) invalidateProductListCache(ctx context.Context) {
	if s.productCache == nil {
		return
	}

	_ = s.productCache.InvalidateProductList(ctx)
}

func (s *productService) attachProductImages(ctx context.Context, product *models.Product) error {
	if product == nil {
		return nil
	}

	images, err := s.productRepo.FindImagesByProductID(ctx, product.ID)
	if err != nil {
		return err
	}

	product.Images = images

	for _, image := range images {
		if image.IsPrimary {
			product.ImageURL = image.ImageURL
			return nil
		}
	}

	if len(images) > 0 {
		product.ImageURL = images[0].ImageURL
	}

	return nil
}

func (s *productService) cleanupCreatedProductImage(ctx context.Context, image *models.ProductImage) {
	if image == nil {
		return
	}

	if image.ProductID != "" && image.ID != "" {
		_ = s.productRepo.DeleteProductImage(ctx, image.ProductID, image.ID)
		_ = s.productRepo.SyncProductPrimaryImageURL(ctx, image.ProductID)
	}

	removeProductImageFile(image.ImageURL)
}

func removeProductImageFile(imageURL string) {
	filePath, ok := productImageURLToFilePath(imageURL)
	if !ok {
		return
	}

	_ = os.Remove(filePath)
}

func productImageURLToFilePath(imageURL string) (string, bool) {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return "", false
	}

	publicPath := strings.TrimRight(productImagePublicPath, "/")
	prefix := publicPath + "/"

	if !strings.HasPrefix(imageURL, prefix) {
		return "", false
	}

	fileName := strings.TrimPrefix(imageURL, prefix)
	if fileName == "" {
		return "", false
	}

	cleanFileName := filepath.Clean(fileName)
	if cleanFileName != fileName {
		return "", false
	}

	if strings.Contains(cleanFileName, "/") || strings.Contains(cleanFileName, "\\") {
		return "", false
	}

	if strings.HasPrefix(cleanFileName, "..") || filepath.IsAbs(cleanFileName) {
		return "", false
	}

	return filepath.Join(productImageUploadDir, cleanFileName), true
}

func normalizeProductInput(
	categoryID string,
	name string,
	description string,
	imageURL string,
	price float64,
	stock int,
) (string, string, string, string, error) {
	categoryID = strings.TrimSpace(categoryID)
	name = normalizeSpaces(name)
	description = strings.TrimSpace(description)
	imageURL = strings.TrimSpace(imageURL)

	if categoryID == "" {
		return "", "", "", "", fmt.Errorf("%w: category_id is required", models.ErrInvalidProductInput)
	}

	if name == "" {
		return "", "", "", "", fmt.Errorf("%w: name is required", models.ErrInvalidProductInput)
	}

	if len(name) < 3 {
		return "", "", "", "", fmt.Errorf("%w: name must be at least 3 characters", models.ErrInvalidProductInput)
	}

	if len(name) > 150 {
		return "", "", "", "", fmt.Errorf("%w: name must be at most 150 characters", models.ErrInvalidProductInput)
	}

	if len(description) > 2000 {
		return "", "", "", "", fmt.Errorf("%w: description must be at most 2000 characters", models.ErrInvalidProductInput)
	}

	if price <= 0 {
		return "", "", "", "", fmt.Errorf("%w: price must be greater than 0", models.ErrInvalidProductInput)
	}

	if stock < 0 {
		return "", "", "", "", fmt.Errorf("%w: stock must be greater than or equal to 0", models.ErrInvalidProductInput)
	}

	if len(imageURL) > 2048 {
		return "", "", "", "", fmt.Errorf("%w: image_url must be at most 2048 characters", models.ErrInvalidProductInput)
	}

	if imageURL != "" {
		parsedURL, err := url.Parse(imageURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			return "", "", "", "", fmt.Errorf("%w: image_url must be a valid URL", models.ErrInvalidProductInput)
		}

		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return "", "", "", "", fmt.Errorf("%w: image_url must use http or https", models.ErrInvalidProductInput)
		}
	}

	return categoryID, name, description, imageURL, nil
}

func detectProductImageExtension(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("%w: image file is empty", models.ErrInvalidProductInput)
	}

	if isWebP(data) {
		return ".webp", nil
	}

	contentType := http.DetectContentType(data)

	switch contentType {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	default:
		return "", fmt.Errorf("%w: image file must be jpg, png, or webp", models.ErrInvalidProductInput)
	}
}

func isWebP(data []byte) bool {
	return len(data) >= 12 &&
		bytes.Equal(data[0:4], []byte("RIFF")) &&
		bytes.Equal(data[8:12], []byte("WEBP"))
}

func generateProductImageFileName(extension string) (string, error) {
	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate image filename: %w", err)
	}

	return hex.EncodeToString(randomBytes) + extension, nil
}

func normalizeProductSort(sortBy string, sortOrder string) (string, string, error) {
	sortBy = strings.ToLower(strings.TrimSpace(sortBy))
	sortOrder = strings.ToLower(strings.TrimSpace(sortOrder))

	if sortBy == "" {
		sortBy = DefaultProductSortBy
	}

	if sortOrder == "" {
		sortOrder = DefaultProductSortOrder
	}

	switch sortBy {
	case ProductSortByCreatedAt, ProductSortByPrice, ProductSortByName:
	default:
		return "", "", fmt.Errorf("%w: sort_by must be created_at, price, or name", models.ErrInvalidProductInput)
	}

	switch sortOrder {
	case ProductSortOrderAsc, ProductSortOrderDesc:
	default:
		return "", "", fmt.Errorf("%w: sort_order must be asc or desc", models.ErrInvalidProductInput)
	}

	return sortBy, sortOrder, nil
}
