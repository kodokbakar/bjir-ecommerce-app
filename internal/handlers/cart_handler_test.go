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

type fakeCartService struct {
	addItemFunc    func(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error)
	getCartFunc    func(ctx context.Context, userID string) (*models.Cart, error)
	updateItemFunc func(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error)
	deleteItemFunc func(ctx context.Context, userID string, itemID string) error
}

func (f *fakeCartService) AddItem(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
	return f.addItemFunc(ctx, userID, productID, quantity)
}

func (f *fakeCartService) GetCart(ctx context.Context, userID string) (*models.Cart, error) {
	return f.getCartFunc(ctx, userID)
}

func (f *fakeCartService) UpdateItem(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error) {
	return f.updateItemFunc(ctx, userID, itemID, quantity)
}

func (f *fakeCartService) DeleteItem(ctx context.Context, userID string, itemID string) error {
	return f.deleteItemFunc(ctx, userID, itemID)
}

func setupCartRouter(service *fakeCartService, withUser bool) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	if withUser {
		router.Use(func(c *gin.Context) {
			c.Set(middleware.ContextUserID, "user-id")
			c.Set(middleware.ContextUserRole, "customer")
			c.Next()
		})
	}

	handler := NewCartHandler(service)

	router.POST("/api/v1/cart/items", handler.AddCartItem)
	router.GET("/api/v1/cart", handler.GetCart)
	router.PUT("/api/v1/cart/items/:id", handler.UpdateCartItem)
	router.DELETE("/api/v1/cart/items/:id", handler.DeleteCartItem)

	return router
}

func newTestCartItem() models.CartItem {
	now := time.Now()

	return models.CartItem{
		ID:        "cart-item-id",
		UserID:    "user-id",
		ProductID: "product-id",
		Quantity:  2,
		Subtotal:  30000000,
		CreatedAt: now,
		UpdatedAt: now,
		Product: &models.Product{
			ID:         "product-id",
			CategoryID: "category-id",
			Name:       "iPhone 15",
			Slug:       "iphone-15",
			Price:      15000000,
			Stock:      10,
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}
}

func TestCartHandler_AddCartItem_Success(t *testing.T) {
	service := &fakeCartService{
		addItemFunc: func(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
			if userID != "user-id" {
				t.Fatalf("expected user-id, got %s", userID)
			}

			if productID != "product-id" {
				t.Fatalf("expected product-id, got %s", productID)
			}

			if quantity != 2 {
				t.Fatalf("expected quantity 2, got %d", quantity)
			}

			item := newTestCartItem()
			item.Quantity = 2

			return &item, nil
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", strings.NewReader(`{"product_id":"product-id","quantity":2}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "cart item added successfully") {
		t.Fatalf("expected success message, got: %s", w.Body.String())
	}
}

func TestCartHandler_AddCartItem_WithoutUserContext_ReturnsUnauthorized(t *testing.T) {
	service := &fakeCartService{
		addItemFunc: func(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupCartRouter(service, false)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", strings.NewReader(`{"product_id":"product-id","quantity":2}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_AddCartItem_InvalidBody(t *testing.T) {
	service := &fakeCartService{
		addItemFunc: func(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", strings.NewReader(`{"quantity":2}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_AddCartItem_ProductNotFound(t *testing.T) {
	service := &fakeCartService{
		addItemFunc: func(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
			return nil, models.ErrProductNotFound
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", strings.NewReader(`{"product_id":"missing-product-id","quantity":2}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_AddCartItem_InvalidQuantity(t *testing.T) {
	service := &fakeCartService{
		addItemFunc: func(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
			return nil, models.ErrInvalidCartInput
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cart/items", strings.NewReader(`{"product_id":"product-id","quantity":0}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_GetCart_Success(t *testing.T) {
	service := &fakeCartService{
		getCartFunc: func(ctx context.Context, userID string) (*models.Cart, error) {
			if userID != "user-id" {
				t.Fatalf("expected user-id, got %s", userID)
			}

			return &models.Cart{
				Items:      []models.CartItem{newTestCartItem()},
				TotalPrice: 30000000,
			}, nil
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total_price":30000000`) {
		t.Fatalf("expected total_price in response, got: %s", w.Body.String())
	}
}

func TestCartHandler_GetCart_WithoutUserContext_ReturnsUnauthorized(t *testing.T) {
	service := &fakeCartService{
		getCartFunc: func(ctx context.Context, userID string) (*models.Cart, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupCartRouter(service, false)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_UpdateCartItem_Success(t *testing.T) {
	service := &fakeCartService{
		updateItemFunc: func(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error) {
			if userID != "user-id" {
				t.Fatalf("expected user-id, got %s", userID)
			}

			if itemID != "cart-item-id" {
				t.Fatalf("expected cart-item-id, got %s", itemID)
			}

			if quantity != 3 {
				t.Fatalf("expected quantity 3, got %d", quantity)
			}

			item := newTestCartItem()
			item.Quantity = 3
			item.Subtotal = 45000000

			return &item, nil
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items/cart-item-id", strings.NewReader(`{"quantity":3}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"quantity":3`) {
		t.Fatalf("expected quantity 3 in response, got: %s", w.Body.String())
	}
}

func TestCartHandler_UpdateCartItem_WithoutUserContext_ReturnsUnauthorized(t *testing.T) {
	service := &fakeCartService{
		updateItemFunc: func(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupCartRouter(service, false)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items/cart-item-id", strings.NewReader(`{"quantity":3}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_UpdateCartItem_InvalidBody(t *testing.T) {
	service := &fakeCartService{
		updateItemFunc: func(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error) {
			t.Fatal("expected service not to be called")
			return nil, nil
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items/cart-item-id", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_UpdateCartItem_NotFound(t *testing.T) {
	service := &fakeCartService{
		updateItemFunc: func(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error) {
			return nil, models.ErrCartItemNotFound
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cart/items/missing-id", strings.NewReader(`{"quantity":2}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_DeleteCartItem_Success(t *testing.T) {
	service := &fakeCartService{
		deleteItemFunc: func(ctx context.Context, userID string, itemID string) error {
			if userID != "user-id" {
				t.Fatalf("expected user-id, got %s", userID)
			}

			if itemID != "cart-item-id" {
				t.Fatalf("expected cart-item-id, got %s", itemID)
			}

			return nil
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/items/cart-item-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_DeleteCartItem_WithoutUserContext_ReturnsUnauthorized(t *testing.T) {
	service := &fakeCartService{
		deleteItemFunc: func(ctx context.Context, userID string, itemID string) error {
			t.Fatal("expected service not to be called")
			return nil
		},
	}

	router := setupCartRouter(service, false)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/items/cart-item-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_DeleteCartItem_NotFound(t *testing.T) {
	service := &fakeCartService{
		deleteItemFunc: func(ctx context.Context, userID string, itemID string) error {
			return models.ErrCartItemNotFound
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/cart/items/missing-id", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d. body: %s", w.Code, w.Body.String())
	}
}

func TestCartHandler_InternalError(t *testing.T) {
	service := &fakeCartService{
		getCartFunc: func(ctx context.Context, userID string) (*models.Cart, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupCartRouter(service, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cart", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}
