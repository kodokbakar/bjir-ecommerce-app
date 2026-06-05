package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/pashagolub/pgxmock/v5"
)

func TestCategoryRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	now := time.Now()
	category := &models.Category{
		Name:        "Electronics",
		Slug:        "electronics",
		Description: "Electronic products",
		ImageURL:    "https://example.com/electronics.jpg",
	}

	rows := pgxmock.NewRows([]string{
		"id",
		"parent_id",
		"name",
		"slug",
		"description",
		"image_url",
		"created_at",
		"updated_at",
	}).AddRow(
		"category-id",
		"",
		"Electronics",
		"electronics",
		"Electronic products",
		"https://example.com/electronics.jpg",
		now,
		now,
	)

	mock.ExpectQuery("INSERT INTO categories").
		WithArgs(nil, category.Name, category.Slug, category.Description, category.ImageURL).
		WillReturnRows(rows)

	err = repo.Create(context.Background(), category)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if category.ID != "category-id" {
		t.Fatalf("expected category id to be set")
	}

	if category.ParentID != nil {
		t.Fatalf("expected parent_id nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_FindAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id",
		"parent_id",
		"name",
		"slug",
		"description",
		"image_url",
		"created_at",
		"updated_at",
	}).AddRow(
		"category-id-1",
		"",
		"Electronics",
		"electronics",
		"Electronic products",
		"https://example.com/electronics.jpg",
		now,
		now,
	)

	mock.ExpectQuery("FROM categories").
		WillReturnRows(rows)

	categories, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(categories))
	}

	if categories[0].Name != "Electronics" {
		t.Fatalf("expected Electronics, got %s", categories[0].Name)
	}

	if categories[0].ParentID != nil {
		t.Fatalf("expected parent_id nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_FindByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	mock.ExpectQuery("FROM categories").
		WithArgs("missing-id").
		WillReturnError(pgx.ErrNoRows)

	category, err := repo.FindByID(context.Background(), "missing-id")
	if category != nil {
		t.Fatalf("expected nil category")
	}

	if !errors.Is(err, models.ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_FindBySlug_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id",
		"parent_id",
		"name",
		"slug",
		"description",
		"image_url",
		"created_at",
		"updated_at",
	}).AddRow(
		"category-id",
		"",
		"Electronics",
		"electronics",
		"Electronic products",
		"https://example.com/electronics.jpg",
		now,
		now,
	)

	mock.ExpectQuery("FROM categories").
		WithArgs("electronics").
		WillReturnRows(rows)

	category, err := repo.FindBySlug(context.Background(), "electronics")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if category.Slug != "electronics" {
		t.Fatalf("expected slug electronics, got %s", category.Slug)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_ExistsByName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("Electronics", "").
		WillReturnRows(rows)

	exists, err := repo.ExistsByName(context.Background(), "Electronics", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Fatalf("expected exists true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_HasProducts(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("category-id").
		WillReturnRows(rows)

	hasProducts, err := repo.HasProducts(context.Background(), "category-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !hasProducts {
		t.Fatalf("expected hasProducts true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_HasChildren(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("category-id").
		WillReturnRows(rows)

	hasChildren, err := repo.HasChildren(context.Background(), "category-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !hasChildren {
		t.Fatalf("expected hasChildren true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	mock.ExpectExec("DELETE FROM categories").
		WithArgs("category-id").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), "category-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_Delete_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	mock.ExpectExec("DELETE FROM categories").
		WithArgs("missing-id").
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.Delete(context.Background(), "missing-id")
	if !errors.Is(err, models.ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCategoryRepository_FindAllPaginated(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCategoryRepository(mock)

	now := time.Now()

	countRows := pgxmock.NewRows([]string{"count"}).AddRow(25)

	rows := pgxmock.NewRows([]string{
		"id",
		"parent_id",
		"name",
		"slug",
		"description",
		"image_url",
		"created_at",
		"updated_at",
	}).AddRow(
		"category-id",
		"",
		"Electronics",
		"electronics",
		"Electronic products",
		"https://example.com/electronics.jpg",
		now,
		now,
	)

	mock.ExpectQuery("SELECT COUNT").
		WillReturnRows(countRows)

	mock.ExpectQuery("FROM categories").
		WithArgs(10, 10).
		WillReturnRows(rows)

	categories, total, err := repo.FindAllPaginated(context.Background(), CategoryListFilter{
		Limit:  10,
		Offset: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 25 {
		t.Fatalf("expected total 25, got %d", total)
	}

	if len(categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(categories))
	}

	if categories[0].Name != "Electronics" {
		t.Fatalf("expected Electronics, got %s", categories[0].Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
