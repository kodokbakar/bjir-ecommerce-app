package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type fakeCategoryRepository struct {
	createFunc           func(ctx context.Context, category *models.Category) error
	findAllFunc          func(ctx context.Context) ([]models.Category, error)
	findByIDFunc         func(ctx context.Context, id string) (*models.Category, error)
	findBySlugFunc       func(ctx context.Context, slug string) (*models.Category, error)
	existsByNameFunc     func(ctx context.Context, name string, excludeID string) (bool, error)
	existsBySlugFunc     func(ctx context.Context, slug string, excludeID string) (bool, error)
	hasProductsFunc      func(ctx context.Context, categoryID string) (bool, error)
	hasChildrenFunc      func(ctx context.Context, categoryID string) (bool, error)
	updateFunc           func(ctx context.Context, category *models.Category) error
	findAllPaginatedFunc func(ctx context.Context, filter repository.CategoryListFilter) ([]models.Category, int, error)
	deleteFunc           func(ctx context.Context, id string) error
}

func newFakeCategoryRepository() *fakeCategoryRepository {
	now := time.Now()

	return &fakeCategoryRepository{
		createFunc: func(ctx context.Context, category *models.Category) error {
			category.ID = "category-id"
			category.CreatedAt = now
			category.UpdatedAt = now
			return nil
		},
		findAllFunc: func(ctx context.Context) ([]models.Category, error) {
			return []models.Category{}, nil
		},
		findByIDFunc: func(ctx context.Context, id string) (*models.Category, error) {
			return &models.Category{
				ID:        id,
				Name:      "Electronics",
				Slug:      "electronics",
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
		findBySlugFunc: func(ctx context.Context, slug string) (*models.Category, error) {
			return &models.Category{
				ID:        "category-id",
				Name:      "Electronics",
				Slug:      slug,
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
		existsByNameFunc: func(ctx context.Context, name string, excludeID string) (bool, error) {
			return false, nil
		},
		existsBySlugFunc: func(ctx context.Context, slug string, excludeID string) (bool, error) {
			return false, nil
		},
		hasProductsFunc: func(ctx context.Context, categoryID string) (bool, error) {
			return false, nil
		},
		hasChildrenFunc: func(ctx context.Context, categoryID string) (bool, error) {
			return false, nil
		},
		updateFunc: func(ctx context.Context, category *models.Category) error {
			category.UpdatedAt = now
			return nil
		},
		findAllPaginatedFunc: func(ctx context.Context, filter repository.CategoryListFilter) ([]models.Category, int, error) {
			return []models.Category{}, 0, nil
		},
		deleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}
}

func (f *fakeCategoryRepository) Create(ctx context.Context, category *models.Category) error {
	return f.createFunc(ctx, category)
}

func (f *fakeCategoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	return f.findAllFunc(ctx)
}

func (f *fakeCategoryRepository) FindAllPaginated(ctx context.Context, filter repository.CategoryListFilter) ([]models.Category, int, error) {
	return f.findAllPaginatedFunc(ctx, filter)
}

func (f *fakeCategoryRepository) FindByID(ctx context.Context, id string) (*models.Category, error) {
	return f.findByIDFunc(ctx, id)
}

func (f *fakeCategoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	return f.findBySlugFunc(ctx, slug)
}

func (f *fakeCategoryRepository) ExistsByName(ctx context.Context, name string, excludeID string) (bool, error) {
	return f.existsByNameFunc(ctx, name, excludeID)
}

func (f *fakeCategoryRepository) ExistsBySlug(ctx context.Context, slug string, excludeID string) (bool, error) {
	return f.existsBySlugFunc(ctx, slug, excludeID)
}

func (f *fakeCategoryRepository) HasProducts(ctx context.Context, categoryID string) (bool, error) {
	return f.hasProductsFunc(ctx, categoryID)
}

func (f *fakeCategoryRepository) HasChildren(ctx context.Context, categoryID string) (bool, error) {
	return f.hasChildrenFunc(ctx, categoryID)
}

func (f *fakeCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	return f.updateFunc(ctx, category)
}

func (f *fakeCategoryRepository) Delete(ctx context.Context, id string) error {
	return f.deleteFunc(ctx, id)
}

func TestCategoryService_Create_Success(t *testing.T) {
	repo := newFakeCategoryRepository()
	service := NewCategoryService(repo)

	category, err := service.Create(context.Background(), CreateCategoryInput{
		Name:        "  Electronic   Devices  ",
		Description: "Electronic products",
		ImageURL:    "https://example.com/electronics.jpg",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if category.ID == "" {
		t.Fatal("expected category ID to be set")
	}

	if category.Name != "Electronic Devices" {
		t.Fatalf("expected normalized name Electronic Devices, got %s", category.Name)
	}

	if category.Slug != "electronic-devices" {
		t.Fatalf("expected slug electronic-devices, got %s", category.Slug)
	}
}

func TestCategoryService_Create_DuplicateName(t *testing.T) {
	repo := newFakeCategoryRepository()
	repo.existsByNameFunc = func(ctx context.Context, name string, excludeID string) (bool, error) {
		return true, nil
	}

	service := NewCategoryService(repo)

	_, err := service.Create(context.Background(), CreateCategoryInput{
		Name: "Electronics",
	})

	if !errors.Is(err, models.ErrCategoryAlreadyExists) {
		t.Fatalf("expected ErrCategoryAlreadyExists, got %v", err)
	}
}

func TestCategoryService_Create_DuplicateSlug(t *testing.T) {
	repo := newFakeCategoryRepository()
	repo.existsBySlugFunc = func(ctx context.Context, slug string, excludeID string) (bool, error) {
		return true, nil
	}

	service := NewCategoryService(repo)

	_, err := service.Create(context.Background(), CreateCategoryInput{
		Name: "Electronics",
	})

	if !errors.Is(err, models.ErrCategoryAlreadyExists) {
		t.Fatalf("expected ErrCategoryAlreadyExists, got %v", err)
	}
}

func TestCategoryService_Create_InvalidInput(t *testing.T) {
	service := NewCategoryService(newFakeCategoryRepository())

	tests := []struct {
		name  string
		input CreateCategoryInput
	}{
		{
			name: "empty name",
			input: CreateCategoryInput{
				Name: "",
			},
		},
		{
			name: "too short name",
			input: CreateCategoryInput{
				Name: "ab",
			},
		},
		{
			name: "too long name",
			input: CreateCategoryInput{
				Name: strings.Repeat("a", 101),
			},
		},
		{
			name: "invalid image url",
			input: CreateCategoryInput{
				Name:     "Electronics",
				ImageURL: "not-a-url",
			},
		},
		{
			name: "unsupported image url scheme",
			input: CreateCategoryInput{
				Name:     "Electronics",
				ImageURL: "ftp://example.com/image.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Create(context.Background(), tt.input)
			if !errors.Is(err, models.ErrInvalidCategoryInput) {
				t.Fatalf("expected ErrInvalidCategoryInput, got %v", err)
			}
		})
	}
}

func TestCategoryService_Create_WithParent(t *testing.T) {
	parentID := "parent-id"
	parentChecked := false

	repo := newFakeCategoryRepository()
	repo.findByIDFunc = func(ctx context.Context, id string) (*models.Category, error) {
		if id == parentID {
			parentChecked = true
		}

		return &models.Category{
			ID:   id,
			Name: "Parent",
			Slug: "parent",
		}, nil
	}

	service := NewCategoryService(repo)

	category, err := service.Create(context.Background(), CreateCategoryInput{
		ParentID: &parentID,
		Name:     "Phones",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !parentChecked {
		t.Fatal("expected parent category to be checked")
	}

	if category.ParentID == nil || *category.ParentID != parentID {
		t.Fatal("expected parent_id to be set")
	}
}

func TestCategoryService_Update_Success(t *testing.T) {
	service := NewCategoryService(newFakeCategoryRepository())

	category, err := service.Update(context.Background(), "category-id", UpdateCategoryInput{
		Name:        "Updated Electronics",
		Description: "Updated description",
		ImageURL:    "https://example.com/updated.jpg",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if category.ID != "category-id" {
		t.Fatalf("expected category-id, got %s", category.ID)
	}

	if category.Slug != "updated-electronics" {
		t.Fatalf("expected slug updated-electronics, got %s", category.Slug)
	}
}

func TestCategoryService_Update_NotFound(t *testing.T) {
	repo := newFakeCategoryRepository()
	repo.findByIDFunc = func(ctx context.Context, id string) (*models.Category, error) {
		return nil, models.ErrCategoryNotFound
	}

	service := NewCategoryService(repo)

	_, err := service.Update(context.Background(), "missing-id", UpdateCategoryInput{
		Name: "Electronics",
	})

	if !errors.Is(err, models.ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestCategoryService_Update_SelfParent(t *testing.T) {
	service := NewCategoryService(newFakeCategoryRepository())

	parentID := "category-id"

	_, err := service.Update(context.Background(), "category-id", UpdateCategoryInput{
		ParentID: &parentID,
		Name:     "Electronics",
	})

	if !errors.Is(err, models.ErrInvalidCategoryInput) {
		t.Fatalf("expected ErrInvalidCategoryInput, got %v", err)
	}
}

func TestCategoryService_Delete_Success(t *testing.T) {
	service := NewCategoryService(newFakeCategoryRepository())

	err := service.Delete(context.Background(), "category-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCategoryService_Delete_WithProducts(t *testing.T) {
	repo := newFakeCategoryRepository()
	repo.hasProductsFunc = func(ctx context.Context, categoryID string) (bool, error) {
		return true, nil
	}

	service := NewCategoryService(repo)

	err := service.Delete(context.Background(), "category-id")

	if !errors.Is(err, models.ErrCategoryHasProducts) {
		t.Fatalf("expected ErrCategoryHasProducts, got %v", err)
	}
}

func TestCategoryService_Delete_WithChildren(t *testing.T) {
	repo := newFakeCategoryRepository()
	repo.hasChildrenFunc = func(ctx context.Context, categoryID string) (bool, error) {
		return true, nil
	}

	service := NewCategoryService(repo)

	err := service.Delete(context.Background(), "category-id")

	if !errors.Is(err, models.ErrCategoryHasChildren) {
		t.Fatalf("expected ErrCategoryHasChildren, got %v", err)
	}
}

func TestCategoryService_BuildCategoryTree(t *testing.T) {
	parentID := "parent-id"

	categories := []models.Category{
		{
			ID:   parentID,
			Name: "Electronics",
			Slug: "electronics",
		},
		{
			ID:       "child-id",
			ParentID: &parentID,
			Name:     "Phones",
			Slug:     "phones",
		},
	}

	tree := buildCategoryTree(categories)

	if len(tree) != 1 {
		t.Fatalf("expected 1 root category, got %d", len(tree))
	}

	if len(tree[0].Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(tree[0].Children))
	}

	if tree[0].Children[0].Name != "Phones" {
		t.Fatalf("expected child name Phones, got %s", tree[0].Children[0].Name)
	}
}

func TestCategoryService_BuildCategoryTree_OrphanBecomesRoot(t *testing.T) {
	missingParentID := "missing-parent-id"

	categories := []models.Category{
		{
			ID:       "orphan-id",
			ParentID: &missingParentID,
			Name:     "Orphan",
			Slug:     "orphan",
		},
	}

	tree := buildCategoryTree(categories)

	if len(tree) != 1 {
		t.Fatalf("expected orphan category to become root, got %d roots", len(tree))
	}

	if tree[0].ID != "orphan-id" {
		t.Fatalf("expected orphan-id, got %s", tree[0].ID)
	}
}

func TestCategoryService_BuildCategoryTree_CycleDoesNotPanic(t *testing.T) {
	categoryAID := "category-a"
	categoryBID := "category-b"

	categories := []models.Category{
		{
			ID:       categoryAID,
			ParentID: &categoryBID,
			Name:     "Category A",
			Slug:     "category-a",
		},
		{
			ID:       categoryBID,
			ParentID: &categoryAID,
			Name:     "Category B",
			Slug:     "category-b",
		},
	}

	tree := buildCategoryTree(categories)

	if len(tree) != 0 {
		t.Fatalf("expected no root categories for pure cycle, got %d", len(tree))
	}
}

func TestCategoryService_GetAll_WithPagination(t *testing.T) {
	repo := newFakeCategoryRepository()

	repo.findAllPaginatedFunc = func(ctx context.Context, filter repository.CategoryListFilter) ([]models.Category, int, error) {
		if filter.Limit != 10 {
			t.Fatalf("expected limit 10, got %d", filter.Limit)
		}

		if filter.Offset != 10 {
			t.Fatalf("expected offset 10, got %d", filter.Offset)
		}

		return []models.Category{}, 25, nil
	}

	service := NewCategoryService(repo)

	result, err := service.GetAll(context.Background(), CategoryListInput{
		Page:  2,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != 2 {
		t.Fatalf("expected page 2, got %d", result.Page)
	}

	if result.Limit != 10 {
		t.Fatalf("expected limit 10, got %d", result.Limit)
	}

	if result.Total != 25 {
		t.Fatalf("expected total 25, got %d", result.Total)
	}

	if result.TotalPages != 3 {
		t.Fatalf("expected total_pages 3, got %d", result.TotalPages)
	}
}

func TestCategoryService_GetAll_InvalidPaginationDefaults(t *testing.T) {
	repo := newFakeCategoryRepository()

	repo.findAllPaginatedFunc = func(ctx context.Context, filter repository.CategoryListFilter) ([]models.Category, int, error) {
		if filter.Limit != DefaultCategoryLimit {
			t.Fatalf("expected default limit %d, got %d", DefaultCategoryLimit, filter.Limit)
		}

		if filter.Offset != 0 {
			t.Fatalf("expected offset 0, got %d", filter.Offset)
		}

		return []models.Category{}, 0, nil
	}

	service := NewCategoryService(repo)

	result, err := service.GetAll(context.Background(), CategoryListInput{
		Page:  -1,
		Limit: -1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != DefaultCategoryPage {
		t.Fatalf("expected default page %d, got %d", DefaultCategoryPage, result.Page)
	}

	if result.Limit != DefaultCategoryLimit {
		t.Fatalf("expected default limit %d, got %d", DefaultCategoryLimit, result.Limit)
	}
}

func TestCategoryService_GetAll_LimitCappedAtMax(t *testing.T) {
	repo := newFakeCategoryRepository()

	repo.findAllPaginatedFunc = func(ctx context.Context, filter repository.CategoryListFilter) ([]models.Category, int, error) {
		if filter.Limit != MaxCategoryLimit {
			t.Fatalf("expected max limit %d, got %d", MaxCategoryLimit, filter.Limit)
		}

		return []models.Category{}, 0, nil
	}

	service := NewCategoryService(repo)

	result, err := service.GetAll(context.Background(), CategoryListInput{
		Page:  1,
		Limit: MaxCategoryLimit + 1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Limit != MaxCategoryLimit {
		t.Fatalf("expected max limit %d, got %d", MaxCategoryLimit, result.Limit)
	}
}
