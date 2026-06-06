package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type fakeOrderRepository struct {
	checkoutFunc func(ctx context.Context, userID string) (*models.Order, error)
}

func newFakeOrderRepository() *fakeOrderRepository {
	now := time.Now()

	return &fakeOrderRepository{
		checkoutFunc: func(ctx context.Context, userID string) (*models.Order, error) {
			return &models.Order{
				ID:          "order-id",
				UserID:      userID,
				OrderNumber: "ORD-TEST",
				Status:      models.OrderStatusPending,
				TotalAmount: 30000000,
				Items: []models.OrderItem{
					{
						ID:          "order-item-id",
						OrderID:     "order-id",
						ProductID:   "product-id",
						ProductName: "iPhone 15",
						Quantity:    2,
						Price:       15000000,
						Subtotal:    30000000,
						CreatedAt:   now,
					},
				},
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
	}
}

func (f *fakeOrderRepository) Checkout(ctx context.Context, userID string) (*models.Order, error) {
	return f.checkoutFunc(ctx, userID)
}

func TestOrderService_Checkout_Success(t *testing.T) {
	service := NewOrderService(newFakeOrderRepository())

	order, err := service.Checkout(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.ID == "" {
		t.Fatal("expected order id")
	}

	if order.UserID != "user-id" {
		t.Fatalf("expected user-id, got %s", order.UserID)
	}

	if order.Status != models.OrderStatusPending {
		t.Fatalf("expected pending, got %s", order.Status)
	}

	if order.TotalAmount != 30000000 {
		t.Fatalf("expected total 30000000, got %f", order.TotalAmount)
	}

	if len(order.Items) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(order.Items))
	}
}

func TestOrderService_Checkout_EmptyUserID(t *testing.T) {
	service := NewOrderService(newFakeOrderRepository())

	_, err := service.Checkout(context.Background(), "")

	if !errors.Is(err, models.ErrInvalidOrderInput) {
		t.Fatalf("expected ErrInvalidOrderInput, got %v", err)
	}
}

func TestOrderService_Checkout_EmptyCart(t *testing.T) {
	repo := newFakeOrderRepository()
	repo.checkoutFunc = func(ctx context.Context, userID string) (*models.Order, error) {
		return nil, models.ErrCartEmpty
	}

	service := NewOrderService(repo)

	_, err := service.Checkout(context.Background(), "user-id")

	if !errors.Is(err, models.ErrCartEmpty) {
		t.Fatalf("expected ErrCartEmpty, got %v", err)
	}
}

func TestOrderService_Checkout_InsufficientStock(t *testing.T) {
	repo := newFakeOrderRepository()
	repo.checkoutFunc = func(ctx context.Context, userID string) (*models.Order, error) {
		return nil, models.ErrInsufficientStock
	}

	service := NewOrderService(repo)

	_, err := service.Checkout(context.Background(), "user-id")

	if !errors.Is(err, models.ErrInsufficientStock) {
		t.Fatalf("expected ErrInsufficientStock, got %v", err)
	}
}
