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

func newCartItemRows(now time.Time) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"user_id",
		"product_id",
		"quantity",
		"created_at",
		"updated_at",
		"product_id",
		"category_id",
		"name",
		"slug",
		"description",
		"price",
		"stock",
		"image_url",
		"is_active",
		"product_created_at",
		"product_updated_at",
	}).AddRow(
		"cart-item-id",
		"user-id",
		"product-id",
		2,
		now,
		now,
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
	)
}

func TestCartRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	now := time.Now()

	item := &models.CartItem{
		UserID:    "user-id",
		ProductID: "product-id",
		Quantity:  2,
	}

	mock.ExpectQuery("INSERT INTO carts").
		WithArgs(item.UserID, item.ProductID, item.Quantity).
		WillReturnRows(newCartItemRows(now))

	err = repo.Create(context.Background(), item)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID != "cart-item-id" {
		t.Fatalf("expected cart-item-id, got %s", item.ID)
	}

	if item.Subtotal != 30000000 {
		t.Fatalf("expected subtotal 30000000, got %f", item.Subtotal)
	}

	if item.Product == nil {
		t.Fatal("expected product to be preloaded")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_AddOrIncrement(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	now := time.Now()

	mock.ExpectQuery("INSERT INTO carts").
		WithArgs("user-id", "product-id", 2).
		WillReturnRows(newCartItemRows(now))

	item, err := repo.AddOrIncrement(context.Background(), "user-id", "product-id", 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID != "cart-item-id" {
		t.Fatalf("expected cart-item-id, got %s", item.ID)
	}

	if item.Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", item.Quantity)
	}

	if item.Product == nil {
		t.Fatal("expected product to be preloaded")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_AddOrIncrement_InsufficientStock(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	mock.ExpectQuery("INSERT INTO carts").
		WithArgs("user-id", "product-id", 99).
		WillReturnError(pgx.ErrNoRows)

	item, err := repo.AddOrIncrement(context.Background(), "user-id", "product-id", 99)
	if item != nil {
		t.Fatal("expected nil item")
	}

	if !errors.Is(err, models.ErrInvalidCartInput) {
		t.Fatalf("expected ErrInvalidCartInput, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_FindByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	now := time.Now()

	mock.ExpectQuery("FROM carts").
		WithArgs("cart-item-id", "user-id").
		WillReturnRows(newCartItemRows(now))

	item, err := repo.FindByID(context.Background(), "cart-item-id", "user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID != "cart-item-id" {
		t.Fatalf("expected cart-item-id, got %s", item.ID)
	}

	if item.UserID != "user-id" {
		t.Fatalf("expected user-id, got %s", item.UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_FindByUserAndProduct(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	now := time.Now()

	mock.ExpectQuery("FROM carts").
		WithArgs("user-id", "product-id").
		WillReturnRows(newCartItemRows(now))

	item, err := repo.FindByUserAndProduct(context.Background(), "user-id", "product-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ProductID != "product-id" {
		t.Fatalf("expected product-id, got %s", item.ProductID)
	}

	if item.UserID != "user-id" {
		t.Fatalf("expected user-id, got %s", item.UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_FindAllByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	now := time.Now()

	mock.ExpectQuery("FROM carts").
		WithArgs("user-id").
		WillReturnRows(newCartItemRows(now))

	items, err := repo.FindAllByUserID(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if items[0].Subtotal != 30000000 {
		t.Fatalf("expected subtotal 30000000, got %f", items[0].Subtotal)
	}

	if items[0].Product == nil {
		t.Fatal("expected product to be preloaded")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_UpdateQuantity(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	now := time.Now()

	mock.ExpectQuery("UPDATE carts").
		WithArgs("cart-item-id", "user-id", 2).
		WillReturnRows(newCartItemRows(now))

	item, err := repo.UpdateQuantity(context.Background(), "cart-item-id", "user-id", 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID != "cart-item-id" {
		t.Fatalf("expected cart-item-id, got %s", item.ID)
	}

	if item.Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", item.Quantity)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_UpdateQuantity_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	mock.ExpectQuery("UPDATE carts").
		WithArgs("missing-id", "user-id", 2).
		WillReturnError(pgx.ErrNoRows)

	mock.ExpectQuery("FROM carts").
		WithArgs("missing-id", "user-id").
		WillReturnError(pgx.ErrNoRows)

	item, err := repo.UpdateQuantity(context.Background(), "missing-id", "user-id", 2)
	if item != nil {
		t.Fatal("expected nil item")
	}

	if !errors.Is(err, models.ErrCartItemNotFound) {
		t.Fatalf("expected ErrCartItemNotFound, got %v", err)
	}

	if err == nil {
		t.Fatal("expected error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	mock.ExpectExec("DELETE FROM carts").
		WithArgs("cart-item-id", "user-id").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(context.Background(), "cart-item-id", "user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCartRepository_Delete_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewCartRepository(mock)

	mock.ExpectExec("DELETE FROM carts").
		WithArgs("missing-id", "user-id").
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.Delete(context.Background(), "missing-id", "user-id")
	if !errors.Is(err, models.ErrCartItemNotFound) {
		t.Fatalf("expected ErrCartItemNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
