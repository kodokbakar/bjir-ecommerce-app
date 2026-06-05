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

type ProductService interface {
	Create(ctx context.Context, input CreateProductInput) (*models.Product, error)
	GetAll(ctx context.Context, input ProductListInput) (*ProductListResult, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	Update(ctx context.Context, id string, input UpdateProductInput) (*models.Product, error)
	UploadImage(ctx context.Context, input UploadProductImageInput) (*models.Product, error)
	Delete(ctx context.Context, id string) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
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
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
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
	if limit == 0 {
		limit = DefaultProductLimit
	}

	if limit < 1 {
		return nil, fmt.Errorf("%w: limit must be greater than 0", models.ErrInvalidProductInput)
	}

	if limit > MaxProductLimit {
		return nil, fmt.Errorf("%w: limit must be at most 100", models.ErrInvalidProductInput)
	}

	sortBy, sortOrder, err := normalizeProductSort(input.SortBy, input.SortOrder)
	if err != nil {
		return nil, err
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

	return &ProductListResult{
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
	}, nil
}

func (s *productService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	return s.productRepo.FindByID(ctx, id)
}

func (s *productService) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("%w: product slug is required", models.ErrInvalidProductInput)
	}

	return s.productRepo.FindBySlug(ctx, slug)
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

	return product, nil
}

func (s *productService) UploadImage(ctx context.Context, input UploadProductImageInput) (*models.Product, error) {
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

	product, err := s.productRepo.UpdateImageURL(ctx, productID, imageURL)
	if err != nil {
		_ = os.Remove(filePath)
		return nil, err
	}

	return product, nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: product id is required", models.ErrInvalidProductInput)
	}

	return s.productRepo.Delete(ctx, id)
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
