package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type fakeOrderService struct {
	checkoutFunc func(ctx context.Context, userID string) (*models.Order, error)
}

func (f *fakeOrderService) Checkout(ctx context.Context, userID string) (*models.Order, error) {
	return f.checkoutFunc(ctx, userID)
}

func setupOrderRouter(service *fakeOrderService, withUserContext bool) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewOrderHandler(service)

	if withUserContext {
		router.Use(func(c *gin.Context) {
			c.Set(middleware.ContextUserID, "user-id")
			c.Set(middleware.ContextUserEmail, "user@example.com")
			c.Set(middleware.ContextUserRole, "customer")
			c.Next()
		})
	}

	router.POST("/api/v1/orders/checkout", handler.Checkout)

	return router
}

func newTestOrder() *models.Order {
	now := time.Now()

	return &models.Order{
		ID:          "order-id",
		UserID:      "user-id",
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
	}
}

func TestOrderHandler_Checkout_Success(t *testing.T) {
	service := &fakeOrderService{
		checkoutFunc: func(ctx context.Context, userID string) (*models.Order, error) {
			if userID != "user-id" {
				t.Fatalf("expected user-id, got %s", userID)
			}

			return newTestOrder(), nil
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "checkout successful") {
		t.Fatalf("expected success message, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "ORD-TEST") {
		t.Fatalf("expected order number in response, got: %s", w.Body.String())
	}
}

func TestOrderHandler_Checkout_WithoutUserContext_ReturnsUnauthorized(t *testing.T) {
	service := &fakeOrderService{
		checkoutFunc: func(ctx context.Context, userID string) (*models.Order, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupOrderRouter(service, false)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_Checkout_EmptyCart(t *testing.T) {
	service := &fakeOrderService{
		checkoutFunc: func(ctx context.Context, userID string) (*models.Order, error) {
			return nil, models.ErrCartEmpty
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_Checkout_InsufficientStock(t *testing.T) {
	service := &fakeOrderService{
		checkoutFunc: func(ctx context.Context, userID string) (*models.Order, error) {
			return nil, models.ErrInsufficientStock
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_Checkout_InternalError(t *testing.T) {
	service := &fakeOrderService{
		checkoutFunc: func(ctx context.Context, userID string) (*models.Order, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}
