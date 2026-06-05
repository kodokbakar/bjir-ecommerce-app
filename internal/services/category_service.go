package services

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

const (
	DefaultCategoryPage  = 1
	DefaultCategoryLimit = 20
	MaxCategoryLimit     = 100
)

type CreateCategoryInput struct {
	ParentID    *string
	Name        string
	Description string
	ImageURL    string
}

type UpdateCategoryInput struct {
	ParentID    *string
	Name        string
	Description string
	ImageURL    string
}

type CategoryListInput struct {
	Page  int
	Limit int
}

type CategoryListResult struct {
	Categories []models.Category
	Page       int
	Limit      int
	Total      int
	TotalPages int
}

type CategoryService interface {
	Create(ctx context.Context, input CreateCategoryInput) (*models.Category, error)
	GetAll(ctx context.Context, input CategoryListInput) (*CategoryListResult, error)
	GetByID(ctx context.Context, id string) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
	Update(ctx context.Context, id string, input UpdateCategoryInput) (*models.Category, error)
	Delete(ctx context.Context, id string) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) Create(ctx context.Context, input CreateCategoryInput) (*models.Category, error) {
	parentID, name, description, imageURL, err := normalizeCategoryInput(input.ParentID, input.Name, input.Description, input.ImageURL)
	if err != nil {
		return nil, err
	}

	if parentID != nil {
		if _, err := s.categoryRepo.FindByID(ctx, *parentID); err != nil {
			return nil, err
		}
	}

	slug := slugify(name)
	if slug == "" {
		return nil, fmt.Errorf("%w: category slug is invalid", models.ErrInvalidCategoryInput)
	}

	nameExists, err := s.categoryRepo.ExistsByName(ctx, name, "")
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, models.ErrCategoryAlreadyExists
	}

	slugExists, err := s.categoryRepo.ExistsBySlug(ctx, slug, "")
	if err != nil {
		return nil, err
	}
	if slugExists {
		return nil, models.ErrCategoryAlreadyExists
	}

	category := &models.Category{
		ParentID:    parentID,
		Name:        name,
		Slug:        slug,
		Description: description,
		ImageURL:    imageURL,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetAll(ctx context.Context, input CategoryListInput) (*CategoryListResult, error) {
	page := input.Page
	if page < 1 {
		page = DefaultCategoryPage
	}

	limit := input.Limit
	if limit < 1 {
		limit = DefaultCategoryLimit
	}

	if limit > MaxCategoryLimit {
		limit = MaxCategoryLimit
	}

	offset := (page - 1) * limit

	categories, total, err := s.categoryRepo.FindAllPaginated(ctx, repository.CategoryListFilter{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &CategoryListResult{
		Categories: categories,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *categoryService) GetByID(ctx context.Context, id string) (*models.Category, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: category id is required", models.ErrInvalidCategoryInput)
	}

	return s.categoryRepo.FindByID(ctx, id)
}

func (s *categoryService) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, fmt.Errorf("%w: category slug is required", models.ErrInvalidCategoryInput)
	}

	return s.categoryRepo.FindBySlug(ctx, slug)
}

func (s *categoryService) Update(ctx context.Context, id string, input UpdateCategoryInput) (*models.Category, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("%w: category id is required", models.ErrInvalidCategoryInput)
	}

	if _, err := s.categoryRepo.FindByID(ctx, id); err != nil {
		return nil, err
	}

	parentID, name, description, imageURL, err := normalizeCategoryInput(input.ParentID, input.Name, input.Description, input.ImageURL)
	if err != nil {
		return nil, err
	}

	if parentID != nil {
		if *parentID == id {
			return nil, fmt.Errorf("%w: category cannot be its own parent", models.ErrInvalidCategoryInput)
		}

		if _, err := s.categoryRepo.FindByID(ctx, *parentID); err != nil {
			return nil, err
		}
	}

	slug := slugify(name)
	if slug == "" {
		return nil, fmt.Errorf("%w: category slug is invalid", models.ErrInvalidCategoryInput)
	}

	nameExists, err := s.categoryRepo.ExistsByName(ctx, name, id)
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, models.ErrCategoryAlreadyExists
	}

	slugExists, err := s.categoryRepo.ExistsBySlug(ctx, slug, id)
	if err != nil {
		return nil, err
	}
	if slugExists {
		return nil, models.ErrCategoryAlreadyExists
	}

	category := &models.Category{
		ID:          id,
		ParentID:    parentID,
		Name:        name,
		Slug:        slug,
		Description: description,
		ImageURL:    imageURL,
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("%w: category id is required", models.ErrInvalidCategoryInput)
	}

	if _, err := s.categoryRepo.FindByID(ctx, id); err != nil {
		return err
	}

	hasProducts, err := s.categoryRepo.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return models.ErrCategoryHasProducts
	}

	hasChildren, err := s.categoryRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return models.ErrCategoryHasChildren
	}

	return s.categoryRepo.Delete(ctx, id)
}

func normalizeCategoryInput(parentID *string, name, description, imageURL string) (*string, string, string, string, error) {
	normalizedParentID := normalizeOptionalString(parentID)
	name = normalizeSpaces(name)
	description = strings.TrimSpace(description)
	imageURL = strings.TrimSpace(imageURL)

	if name == "" {
		return nil, "", "", "", fmt.Errorf("%w: name is required", models.ErrInvalidCategoryInput)
	}

	if len(name) < 3 {
		return nil, "", "", "", fmt.Errorf("%w: name must be at least 3 characters", models.ErrInvalidCategoryInput)
	}

	if len(name) > 100 {
		return nil, "", "", "", fmt.Errorf("%w: name must be at most 100 characters", models.ErrInvalidCategoryInput)
	}

	if len(description) > 1000 {
		return nil, "", "", "", fmt.Errorf("%w: description must be at most 1000 characters", models.ErrInvalidCategoryInput)
	}

	if len(imageURL) > 2048 {
		return nil, "", "", "", fmt.Errorf("%w: image_url must be at most 2048 characters", models.ErrInvalidCategoryInput)
	}

	if imageURL != "" {
		parsedURL, err := url.Parse(imageURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			return nil, "", "", "", fmt.Errorf("%w: image_url must be a valid URL", models.ErrInvalidCategoryInput)
		}

		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			return nil, "", "", "", fmt.Errorf("%w: image_url must use http or https", models.ErrInvalidCategoryInput)
		}
	}

	return normalizedParentID, name, description, imageURL, nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizeSpaces(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))

	nonAlphaNum := regexp.MustCompile(`[^a-z0-9]+`)
	value = nonAlphaNum.ReplaceAllString(value, "-")

	value = strings.Trim(value, "-")

	multipleDashes := regexp.MustCompile(`-+`)
	value = multipleDashes.ReplaceAllString(value, "-")

	return value
}

func buildCategoryTree(categories []models.Category) []models.Category {
	nodes := make(map[string]models.Category, len(categories))
	childrenByParent := make(map[string][]string)
	rootIDs := make([]string, 0)

	for _, category := range categories {
		category.Children = nil
		nodes[category.ID] = category
	}

	for _, category := range categories {
		if category.ParentID == nil {
			rootIDs = append(rootIDs, category.ID)
			continue
		}

		if _, ok := nodes[*category.ParentID]; !ok {
			rootIDs = append(rootIDs, category.ID)
			continue
		}

		childrenByParent[*category.ParentID] = append(childrenByParent[*category.ParentID], category.ID)
	}

	var build func(id string, visited map[string]bool) models.Category
	build = func(id string, visited map[string]bool) models.Category {
		category := nodes[id]
		category.Children = nil

		if visited[id] {
			return category
		}

		visited[id] = true

		for _, childID := range childrenByParent[id] {
			category.Children = append(category.Children, build(childID, visited))
		}

		return category
	}

	result := make([]models.Category, 0, len(rootIDs))
	for _, rootID := range rootIDs {
		result = append(result, build(rootID, map[string]bool{}))
	}

	return result
}
