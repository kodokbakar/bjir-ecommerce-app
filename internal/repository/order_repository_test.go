package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

func newCheckoutCartRows(now time.Time, stock int) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"cart_item_id",
		"user_id",
		"product_id",
		"quantity",
		"product_id",
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
	}).AddRow(
		"cart-item-id",
		"user-id",
		"product-id",
		2,
		"product-id",
		"category-id",
		"iPhone 15",
		"iphone-15",
		"Apple smartphone",
		15000000.0,
		stock,
		"https://example.com/iphone.jpg",
		true,
		now,
		now,
	)
}

func newOrderRows(now time.Time) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"user_id",
		"order_number",
		"status",
		"total_amount",
		"shipping_address",
		"notes",
		"created_at",
		"updated_at",
	}).AddRow(
		"order-id",
		"user-id",
		"ORD-TEST",
		models.OrderStatusPending,
		30000000.0,
		"",
		"",
		now,
		now,
	)
}

func newOrderItemRows(now time.Time) *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"order_id",
		"product_id",
		"product_name",
		"quantity",
		"price",
		"subtotal",
		"created_at",
	}).AddRow(
		"order-item-id",
		"order-id",
		"product-id",
		"iPhone 15",
		2,
		15000000.0,
		30000000.0,
		now,
	)
}

func TestOrderRepository_Checkout_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewOrderRepository(mock)

	now := time.Now()

	mock.ExpectBegin()

	mock.ExpectQuery("FROM carts").
		WithArgs("user-id").
		WillReturnRows(newCheckoutCartRows(now, 10))

	mock.ExpectExec("UPDATE products").
		WithArgs(2, "product-id").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectQuery("INSERT INTO orders").
		WithArgs(
			"user-id",
			pgxmock.AnyArg(),
			models.OrderStatusPending,
			30000000.0,
			"",
			"",
		).
		WillReturnRows(newOrderRows(now))

	mock.ExpectQuery("INSERT INTO order_items").
		WithArgs(
			"order-id",
			"product-id",
			"iPhone 15",
			2,
			15000000.0,
			30000000.0,
		).
		WillReturnRows(newOrderItemRows(now))

	mock.ExpectExec("DELETE FROM carts").
		WithArgs("user-id").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	mock.ExpectCommit()

	order, err := repo.Checkout(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.ID != "order-id" {
		t.Fatalf("expected order-id, got %s", order.ID)
	}

	if order.TotalAmount != 30000000.0 {
		t.Fatalf("expected total amount 30000000, got %f", order.TotalAmount)
	}

	if len(order.Items) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(order.Items))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestOrderRepository_Checkout_EmptyCart(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewOrderRepository(mock)

	mock.ExpectBegin()

	rows := pgxmock.NewRows([]string{
		"cart_item_id",
		"user_id",
		"product_id",
		"quantity",
		"product_id",
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
	})

	mock.ExpectQuery("FROM carts").
		WithArgs("user-id").
		WillReturnRows(rows)

	mock.ExpectRollback()

	order, err := repo.Checkout(context.Background(), "user-id")
	if order != nil {
		t.Fatal("expected nil order")
	}

	if !errors.Is(err, models.ErrCartEmpty) {
		t.Fatalf("expected ErrCartEmpty, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestOrderRepository_Checkout_InsufficientStock_RollbackBeforeCreateOrder(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewOrderRepository(mock)

	now := time.Now()

	mock.ExpectBegin()

	mock.ExpectQuery("FROM carts").
		WithArgs("user-id").
		WillReturnRows(newCheckoutCartRows(now, 1))

	mock.ExpectExec("UPDATE products").
		WithArgs(2, "product-id").
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	mock.ExpectRollback()

	order, err := repo.Checkout(context.Background(), "user-id")
	if order != nil {
		t.Fatal("expected nil order")
	}

	if !errors.Is(err, models.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
