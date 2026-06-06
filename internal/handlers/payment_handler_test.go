package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/middleware"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/services"
)

type fakePaymentService struct {
	payOrderFunc func(ctx context.Context, input services.PayOrderInput) (*models.Payment, error)
}

func (f *fakePaymentService) PayOrder(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
	return f.payOrderFunc(ctx, input)
}

func setupPaymentRouter(service services.PaymentService, withUser bool) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewPaymentHandler(service)

	router.POST("/api/v1/payments/pay", func(c *gin.Context) {
		if withUser {
			c.Set(middleware.ContextUserID, "user-id")
		}

		handler.PayOrder(c)
	})

	return router
}

func TestPaymentHandler_PayOrder_Success(t *testing.T) {
	service := &fakePaymentService{
		payOrderFunc: func(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
			if input.UserID != "user-id" {
				t.Fatalf("expected user-id, got %s", input.UserID)
			}

			if input.OrderID != "order-id" {
				t.Fatalf("expected order-id, got %s", input.OrderID)
			}

			if input.Method != "bank_transfer" {
				t.Fatalf("expected bank_transfer, got %s", input.Method)
			}

			return &models.Payment{
				ID:            "payment-id",
				OrderID:       input.OrderID,
				Provider:      models.PaymentProviderMock,
				PaymentMethod: models.PaymentMethodBankTransfer,
				TransactionID: "PAY-TEST",
				Amount:        30000000,
				Status:        models.PaymentStatusPaid,
			}, nil
		},
	}

	router := setupPaymentRouter(service, true)

	body := `{
		"order_id": "order-id",
		"method": "bank_transfer"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/pay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "payment successful") {
		t.Fatalf("expected success message, got body: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "PAY-TEST") {
		t.Fatalf("expected transaction id in body, got: %s", w.Body.String())
	}
}

func TestPaymentHandler_PayOrder_WithoutUserContext(t *testing.T) {
	service := &fakePaymentService{
		payOrderFunc: func(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupPaymentRouter(service, false)

	body := `{
		"order_id": "order-id",
		"method": "bank_transfer"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/pay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestPaymentHandler_PayOrder_InvalidBody(t *testing.T) {
	service := &fakePaymentService{
		payOrderFunc: func(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupPaymentRouter(service, true)

	body := `{
		"method": "bank_transfer"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/pay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestPaymentHandler_PayOrder_OrderNotFound(t *testing.T) {
	service := &fakePaymentService{
		payOrderFunc: func(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
			return nil, models.ErrOrderNotFound
		},
	}

	router := setupPaymentRouter(service, true)

	body := `{
		"order_id": "missing-order-id",
		"method": "ewallet"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/pay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestPaymentHandler_PayOrder_OrderNotPayable(t *testing.T) {
	service := &fakePaymentService{
		payOrderFunc: func(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
			return nil, models.ErrOrderNotPayable
		},
	}

	router := setupPaymentRouter(service, true)

	body := `{
		"order_id": "order-id",
		"method": "credit_card"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/pay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestPaymentHandler_PayOrder_InternalError(t *testing.T) {
	service := &fakePaymentService{
		payOrderFunc: func(ctx context.Context, input services.PayOrderInput) (*models.Payment, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupPaymentRouter(service, true)

	body := `{
		"order_id": "order-id",
		"method": "credit_card"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/pay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}
