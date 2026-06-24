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
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type fakeOrderService struct {
	checkoutFunc         func(ctx context.Context, userID string) (*models.Order, error)
	getMyOrdersFunc      func(ctx context.Context, userID string, input services.OrderListInput) (*services.OrderListResult, error)
	getMyOrderDetailFunc func(ctx context.Context, userID string, orderID string) (*models.Order, error)
	updateStatusFunc     func(ctx context.Context, orderID string, status string) (*models.Order, error)
	getAllOrdersFunc     func(ctx context.Context, input services.OrderListInput) (*services.OrderListResult, error)
}

func (f *fakeOrderService) Checkout(ctx context.Context, userID string) (*models.Order, error) {
	return f.checkoutFunc(ctx, userID)
}

func (f *fakeOrderService) GetMyOrders(ctx context.Context, userID string, input services.OrderListInput) (*services.OrderListResult, error) {
	return f.getMyOrdersFunc(ctx, userID, input)
}

func (f *fakeOrderService) GetMyOrderDetail(ctx context.Context, userID string, orderID string) (*models.Order, error) {
	return f.getMyOrderDetailFunc(ctx, userID, orderID)
}

func (f *fakeOrderService) UpdateStatus(ctx context.Context, orderID string, status string) (*models.Order, error) {
	return f.updateStatusFunc(ctx, orderID, status)
}

func (f *fakeOrderService) GetAllOrders(ctx context.Context, input services.OrderListInput) (*services.OrderListResult, error) {
	return f.getAllOrdersFunc(ctx, input)
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

	router.GET("/api/v1/orders", handler.GetMyOrders)
	router.POST("/api/v1/orders/checkout", handler.Checkout)
	router.GET("/api/v1/orders/:id", handler.GetMyOrderDetail)
	router.PATCH("/api/v1/admin/orders/:id/status", handler.UpdateOrderStatus)
	router.GET("/api/v1/admin/orders", handler.GetAllOrders)

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

func TestOrderHandler_Checkout_ProductNotFound(t *testing.T) {
	service := &fakeOrderService{
		checkoutFunc: func(ctx context.Context, userID string) (*models.Order, error) {
			return nil, models.ErrProductNotFound
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders/checkout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
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

func TestOrderHandler_GetMyOrders_Success(t *testing.T) {
	service := &fakeOrderService{
		getMyOrdersFunc: func(ctx context.Context, userID string, input services.OrderListInput) (*services.OrderListResult, error) {
			if userID != "user-id" {
				t.Fatalf("expected user-id, got %s", userID)
			}

			if input.Page != 2 {
				t.Fatalf("expected page 2, got %d", input.Page)
			}

			if input.Limit != 10 {
				t.Fatalf("expected limit 10, got %d", input.Limit)
			}

			return &services.OrderListResult{
				Orders: []models.Order{
					{
						ID:          "order-id",
						UserID:      userID,
						OrderNumber: "ORD-TEST",
						Status:      models.OrderStatusPending,
						TotalAmount: 30000000,
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				},
				Page:       2,
				Limit:      10,
				Total:      25,
				TotalPages: 3,
			}, nil
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders?page=2&limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "ORD-TEST") {
		t.Fatalf("expected order number in response, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"page":2`) {
		t.Fatalf("expected page meta, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total_pages":3`) {
		t.Fatalf("expected total_pages meta, got: %s", w.Body.String())
	}
}

func TestOrderHandler_GetMyOrders_WithoutUserContext_ReturnsUnauthorized(t *testing.T) {
	service := &fakeOrderService{
		getMyOrdersFunc: func(ctx context.Context, userID string, input services.OrderListInput) (*services.OrderListResult, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupOrderRouter(service, false)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_GetMyOrders_InternalError(t *testing.T) {
	service := &fakeOrderService{
		getMyOrdersFunc: func(ctx context.Context, userID string, input services.OrderListInput) (*services.OrderListResult, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_GetMyOrderDetail_Success(t *testing.T) {
	service := &fakeOrderService{
		getMyOrderDetailFunc: func(ctx context.Context, userID string, orderID string) (*models.Order, error) {
			if userID != "user-id" {
				t.Fatalf("expected user-id, got %s", userID)
			}

			if orderID != "order-id" {
				t.Fatalf("expected order-id, got %s", orderID)
			}

			return newTestOrder(), nil
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/order-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "ORD-TEST") {
		t.Fatalf("expected order number in response, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "iPhone 15") {
		t.Fatalf("expected order item in response, got: %s", w.Body.String())
	}
}

func TestOrderHandler_GetMyOrderDetail_NotFound(t *testing.T) {
	service := &fakeOrderService{
		getMyOrderDetailFunc: func(ctx context.Context, userID string, orderID string) (*models.Order, error) {
			return nil, models.ErrOrderNotFound
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/missing-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_GetMyOrderDetail_InternalError(t *testing.T) {
	service := &fakeOrderService{
		getMyOrderDetailFunc: func(ctx context.Context, userID string, orderID string) (*models.Order, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/order-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_GetMyOrderDetail_WithoutUserContext_ReturnsUnauthorized(t *testing.T) {
	service := &fakeOrderService{
		getMyOrderDetailFunc: func(ctx context.Context, userID string, orderID string) (*models.Order, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupOrderRouter(service, false)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/order-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_UpdateOrderStatus_Success(t *testing.T) {
	service := &fakeOrderService{
		updateStatusFunc: func(ctx context.Context, orderID string, status string) (*models.Order, error) {
			if orderID != "order-id" {
				t.Fatalf("expected order-id, got %s", orderID)
			}

			if status != "paid" {
				t.Fatalf("expected paid, got %s", status)
			}

			order := newTestOrder()
			order.Status = models.OrderStatusPaid
			return order, nil
		},
	}

	router := setupOrderRouter(service, true)

	body := `{"status":"paid"}`

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/orders/order-id/status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "order status updated successfully") {
		t.Fatalf("expected success message, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), models.OrderStatusPaid) {
		t.Fatalf("expected paid status, got: %s", w.Body.String())
	}
}

func TestOrderHandler_UpdateOrderStatus_InvalidBody(t *testing.T) {
	service := &fakeOrderService{
		updateStatusFunc: func(ctx context.Context, orderID string, status string) (*models.Order, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/orders/order-id/status", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_UpdateOrderStatus_InvalidTransition(t *testing.T) {
	service := &fakeOrderService{
		updateStatusFunc: func(ctx context.Context, orderID string, status string) (*models.Order, error) {
			return nil, models.ErrInvalidOrderStatusTransition
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/orders/order-id/status", strings.NewReader(`{"status":"shipped"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_UpdateOrderStatus_InvalidStatus(t *testing.T) {
	service := &fakeOrderService{
		updateStatusFunc: func(ctx context.Context, orderID string, status string) (*models.Order, error) {
			return nil, models.ErrInvalidOrderStatus
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/orders/order-id/status", strings.NewReader(`{"status":"unknown"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestOrderHandler_UpdateOrderStatus_NotFound(t *testing.T) {
	service := &fakeOrderService{
		updateStatusFunc: func(ctx context.Context, orderID string, status string) (*models.Order, error) {
			return nil, models.ErrOrderNotFound
		},
	}

	router := setupOrderRouter(service, true)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/orders/missing-id/status", strings.NewReader(`{"status":"paid"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}
