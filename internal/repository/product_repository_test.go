package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

func newProductRows(now time.Time) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"category_id",
		"name",
		"slug",
		"description",
		"price",
		"stock",
		"image_url",
		"is_active",
		"created_at",
		"updated_at",
		"category_id",
		"category_parent_id",
		"category_name",
		"category_slug",
		"category_description",
		"category_image_url",
		"category_created_at",
		"category_updated_at",
	}).AddRow(
		"product-id",
		"category-id",
		"iPhone 15",
		"iphone-15",
		"Apple smartphone",
		15000000.0,
		10,
		"https://example.com/iphone.jpg",
		true,
		now,
		now,
		"category-id",
		"",
		"Phones",
		"phones",
		"Phone products",
		"https://example.com/phones.jpg",
		now,
		now,
	)
}

func TestProductRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	now := time.Now()

	product := &models.Product{
		CategoryID:  "category-id",
		Name:        "iPhone 15",
		Slug:        "iphone-15",
		Description: "Apple smartphone",
		Price:       15000000.0,
		Stock:       10,
		ImageURL:    "https://example.com/iphone.jpg",
	}

	mock.ExpectQuery("INSERT INTO products").
		WithArgs(
			product.CategoryID,
			product.Name,
			product.Slug,
			product.Description,
			product.Price,
			product.Stock,
			product.ImageURL,
		).
		WillReturnRows(newProductRows(now))

	err = repo.Create(context.Background(), product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != "product-id" {
		t.Fatalf("expected product-id, got %s", product.ID)
	}

	if product.Category == nil {
		t.Fatal("expected preloaded category")
	}

	if product.Category.Name != "Phones" {
		t.Fatalf("expected category Phones, got %s", product.Category.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_FindByID_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	now := time.Now()

	mock.ExpectQuery("FROM products").
		WithArgs("product-id").
		WillReturnRows(newProductRows(now))

	product, err := repo.FindByID(context.Background(), "product-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != "product-id" {
		t.Fatalf("expected product-id, got %s", product.ID)
	}

	if product.Category == nil {
		t.Fatal("expected preloaded category")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_FindByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	mock.ExpectQuery("FROM products").
		WithArgs("missing-id").
		WillReturnError(pgx.ErrNoRows)

	product, err := repo.FindByID(context.Background(), "missing-id")
	if product != nil {
		t.Fatal("expected nil product")
	}

	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_ExistsBySlug(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs("iphone-15", "").
		WillReturnRows(rows)

	exists, err := repo.ExistsBySlug(context.Background(), "iphone-15", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !exists {
		t.Fatal("expected exists true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	now := time.Now()

	product := &models.Product{
		ID:          "product-id",
		CategoryID:  "category-id",
		Name:        "iPhone 15 Pro",
		Slug:        "iphone-15-pro",
		Description: "Apple smartphone pro",
		Price:       18000000.0,
		Stock:       5,
		ImageURL:    "https://example.com/iphone-pro.jpg",
	}

	mock.ExpectQuery("UPDATE products").
		WithArgs(
			product.ID,
			product.CategoryID,
			product.Name,
			product.Slug,
			product.Description,
			product.Price,
			product.Stock,
			product.ImageURL,
		).
		WillReturnRows(newProductRows(now))

	err = repo.Update(context.Background(), product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != "product-id" {
		t.Fatalf("expected product-id, got %s", product.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	mock.ExpectExec("UPDATE products").
		WithArgs("product-id").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Delete(context.Background(), "product-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	mock.ExpectExec("UPDATE products").
		WithArgs("missing-id").
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.Delete(context.Background(), "missing-id")
	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_FindAll_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	now := time.Now()

	mock.ExpectQuery("FROM products").
		WillReturnRows(newProductRows(now))

	products, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	if products[0].ID != "product-id" {
		t.Fatalf("expected product-id, got %s", products[0].ID)
	}

	if products[0].Category == nil {
		t.Fatal("expected preloaded category")
	}

	if products[0].Category.Name != "Phones" {
		t.Fatalf("expected category Phones, got %s", products[0].Category.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_FindBySlug_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	now := time.Now()

	mock.ExpectQuery("FROM products").
		WithArgs("iphone-15").
		WillReturnRows(newProductRows(now))

	product, err := repo.FindBySlug(context.Background(), "iphone-15")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}

	if product.Slug != "iphone-15" {
		t.Fatalf("expected slug iphone-15, got %s", product.Slug)
	}

	if product.Category == nil {
		t.Fatal("expected preloaded category")
	}

	if product.Category.Name != "Phones" {
		t.Fatalf("expected category Phones, got %s", product.Category.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_FindBySlug_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	mock.ExpectQuery("FROM products").
		WithArgs("missing-product").
		WillReturnError(pgx.ErrNoRows)

	product, err := repo.FindBySlug(context.Background(), "missing-product")
	if product != nil {
		t.Fatal("expected nil product")
	}

	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_UpdateImageURL_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	now := time.Now()

	mock.ExpectQuery("UPDATE products").
		WithArgs("product-id", "/uploads/products/test.png").
		WillReturnRows(newProductRows(now))

	product, err := repo.UpdateImageURL(context.Background(), "product-id", "/uploads/products/test.png")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product == nil {
		t.Fatal("expected product, got nil")
	}

	if product.Category == nil {
		t.Fatal("expected preloaded category")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestProductRepository_UpdateImageURL_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewProductRepository(mock)

	mock.ExpectQuery("UPDATE products").
		WithArgs("missing-id", "/uploads/products/test.png").
		WillReturnError(pgx.ErrNoRows)

	product, err := repo.UpdateImageURL(context.Background(), "missing-id", "/uploads/products/test.png")
	if product != nil {
		t.Fatal("expected nil product")
	}

	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
