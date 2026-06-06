package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type fakeOrderRepository struct {
	checkoutFunc          func(ctx context.Context, userID string) (*models.Order, error)
	findAllByUserIDFunc   func(ctx context.Context, userID string, filter repository.OrderListFilter) ([]models.Order, int, error)
	findByIDAndUserIDFunc func(ctx context.Context, orderID string, userID string) (*models.Order, error)
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
		findAllByUserIDFunc: func(ctx context.Context, userID string, filter repository.OrderListFilter) ([]models.Order, int, error) {
			return []models.Order{
				{
					ID:          "order-id",
					UserID:      userID,
					OrderNumber: "ORD-TEST",
					Status:      models.OrderStatusPending,
					TotalAmount: 30000000,
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			}, 1, nil
		},
		findByIDAndUserIDFunc: func(ctx context.Context, orderID string, userID string) (*models.Order, error) {
			return &models.Order{
				ID:          orderID,
				UserID:      userID,
				OrderNumber: "ORD-TEST",
				Status:      models.OrderStatusPending,
				TotalAmount: 30000000,
				Items: []models.OrderItem{
					{
						ID:          "order-item-id",
						OrderID:     orderID,
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

func (f *fakeOrderRepository) FindAllByUserID(ctx context.Context, userID string, filter repository.OrderListFilter) ([]models.Order, int, error) {
	return f.findAllByUserIDFunc(ctx, userID, filter)
}

func (f *fakeOrderRepository) FindByIDAndUserID(ctx context.Context, orderID string, userID string) (*models.Order, error) {
	return f.findByIDAndUserIDFunc(ctx, orderID, userID)
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

func TestOrderService_GetMyOrders_Success(t *testing.T) {
	repo := newFakeOrderRepository()
	repo.findAllByUserIDFunc = func(ctx context.Context, userID string, filter repository.OrderListFilter) ([]models.Order, int, error) {
		if userID != "user-id" {
			t.Fatalf("expected user-id, got %s", userID)
		}

		if filter.Limit != 10 {
			t.Fatalf("expected limit 10, got %d", filter.Limit)
		}

		if filter.Offset != 10 {
			t.Fatalf("expected offset 10, got %d", filter.Offset)
		}

		return []models.Order{
			{
				ID:          "order-id",
				UserID:      userID,
				OrderNumber: "ORD-TEST",
				Status:      models.OrderStatusPending,
				TotalAmount: 30000000,
			},
		}, 25, nil
	}

	service := NewOrderService(repo)

	result, err := service.GetMyOrders(context.Background(), "user-id", OrderListInput{
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

	if len(result.Orders) != 1 {
		t.Fatalf("expected 1 order, got %d", len(result.Orders))
	}
}

func TestOrderService_GetMyOrders_InvalidPaginationDefaults(t *testing.T) {
	repo := newFakeOrderRepository()
	repo.findAllByUserIDFunc = func(ctx context.Context, userID string, filter repository.OrderListFilter) ([]models.Order, int, error) {
		if filter.Limit != DefaultOrderLimit {
			t.Fatalf("expected default limit %d, got %d", DefaultOrderLimit, filter.Limit)
		}

		if filter.Offset != 0 {
			t.Fatalf("expected offset 0, got %d", filter.Offset)
		}

		return []models.Order{}, 0, nil
	}

	service := NewOrderService(repo)

	result, err := service.GetMyOrders(context.Background(), "user-id", OrderListInput{
		Page:  -1,
		Limit: -1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != DefaultOrderPage {
		t.Fatalf("expected default page %d, got %d", DefaultOrderPage, result.Page)
	}

	if result.Limit != DefaultOrderLimit {
		t.Fatalf("expected default limit %d, got %d", DefaultOrderLimit, result.Limit)
	}
}

func TestOrderService_GetMyOrders_EmptyUserID(t *testing.T) {
	service := NewOrderService(newFakeOrderRepository())

	_, err := service.GetMyOrders(context.Background(), "", OrderListInput{})

	if !errors.Is(err, models.ErrInvalidOrderInput) {
		t.Fatalf("expected ErrInvalidOrderInput, got %v", err)
	}
}

func TestOrderService_GetMyOrderDetail_Success(t *testing.T) {
	service := NewOrderService(newFakeOrderRepository())

	order, err := service.GetMyOrderDetail(context.Background(), "user-id", "order-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if order.ID != "order-id" {
		t.Fatalf("expected order-id, got %s", order.ID)
	}

	if order.UserID != "user-id" {
		t.Fatalf("expected user-id, got %s", order.UserID)
	}

	if len(order.Items) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(order.Items))
	}
}

func TestOrderService_GetMyOrderDetail_NotFound(t *testing.T) {
	repo := newFakeOrderRepository()
	repo.findByIDAndUserIDFunc = func(ctx context.Context, orderID string, userID string) (*models.Order, error) {
		return nil, models.ErrOrderNotFound
	}

	service := NewOrderService(repo)

	_, err := service.GetMyOrderDetail(context.Background(), "user-id", "missing-id")

	if !errors.Is(err, models.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestOrderService_GetMyOrderDetail_EmptyOrderID(t *testing.T) {
	service := NewOrderService(newFakeOrderRepository())

	_, err := service.GetMyOrderDetail(context.Background(), "user-id", "")

	if !errors.Is(err, models.ErrInvalidOrderInput) {
		t.Fatalf("expected ErrInvalidOrderInput, got %v", err)
	}
}
