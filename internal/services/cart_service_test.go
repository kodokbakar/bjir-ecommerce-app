package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type fakeCartRepository struct {
	createFunc               func(ctx context.Context, item *models.CartItem) error
	findAllByUserIDFunc      func(ctx context.Context, userID string) ([]models.CartItem, error)
	findByIDFunc             func(ctx context.Context, id string, userID string) (*models.CartItem, error)
	findByUserAndProductFunc func(ctx context.Context, userID string, productID string) (*models.CartItem, error)
	updateQuantityFunc       func(ctx context.Context, id string, userID string, quantity int) (*models.CartItem, error)
	deleteFunc               func(ctx context.Context, id string, userID string) error
}

func newFakeCartRepository() *fakeCartRepository {
	now := time.Now()

	item := models.CartItem{
		ID:        "cart-item-id",
		UserID:    "user-id",
		ProductID: "product-id",
		Quantity:  2,
		Subtotal:  30000000,
		CreatedAt: now,
		UpdatedAt: now,
		Product: &models.Product{
			ID:       "product-id",
			Name:     "iPhone 15",
			Slug:     "iphone-15",
			Price:    15000000,
			Stock:    10,
			IsActive: true,
		},
	}

	return &fakeCartRepository{
		createFunc: func(ctx context.Context, item *models.CartItem) error {
			item.ID = "cart-item-id"
			item.CreatedAt = now
			item.UpdatedAt = now
			item.Product = &models.Product{
				ID:       item.ProductID,
				Name:     "iPhone 15",
				Slug:     "iphone-15",
				Price:    15000000,
				Stock:    10,
				IsActive: true,
			}
			item.Subtotal = item.Product.Price * float64(item.Quantity)
			return nil
		},

		findAllByUserIDFunc: func(ctx context.Context, userID string) ([]models.CartItem, error) {
			return []models.CartItem{item}, nil
		},

		findByIDFunc: func(ctx context.Context, id string, userID string) (*models.CartItem, error) {
			return &item, nil
		},

		findByUserAndProductFunc: func(ctx context.Context, userID string, productID string) (*models.CartItem, error) {
			return nil, models.ErrCartItemNotFound
		},

		updateQuantityFunc: func(ctx context.Context, id string, userID string, quantity int) (*models.CartItem, error) {
			item.Quantity = quantity
			item.Subtotal = item.Product.Price * float64(quantity)
			return &item, nil
		},

		deleteFunc: func(ctx context.Context, id string, userID string) error {
			return nil
		},
	}
}

func (f *fakeCartRepository) Create(ctx context.Context, item *models.CartItem) error {
	return f.createFunc(ctx, item)
}

func (f *fakeCartRepository) FindByID(ctx context.Context, id string, userID string) (*models.CartItem, error) {
	return f.findByIDFunc(ctx, id, userID)
}

func (f *fakeCartRepository) FindByUserAndProduct(ctx context.Context, userID string, productID string) (*models.CartItem, error) {
	return f.findByUserAndProductFunc(ctx, userID, productID)
}

func (f *fakeCartRepository) FindAllByUserID(ctx context.Context, userID string) ([]models.CartItem, error) {
	return f.findAllByUserIDFunc(ctx, userID)
}

func (f *fakeCartRepository) UpdateQuantity(ctx context.Context, id string, userID string, quantity int) (*models.CartItem, error) {
	return f.updateQuantityFunc(ctx, id, userID, quantity)
}

func (f *fakeCartRepository) Delete(ctx context.Context, id string, userID string) error {
	return f.deleteFunc(ctx, id, userID)
}

func TestCartService_AddItem_NewItemSuccess(t *testing.T) {
	repo := newFakeCartRepository()

	service := NewCartService(repo, newFakeProductRepository())

	item, err := service.AddItem(context.Background(), "user-id", "product-id", 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID == "" {
		t.Fatal("expected cart item id to be set")
	}

	if item.Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", item.Quantity)
	}

	if item.Subtotal != 30000000 {
		t.Fatalf("expected subtotal 30000000, got %f", item.Subtotal)
	}
}

func TestCartService_AddItem_ExistingItemAddsQuantity(t *testing.T) {
	repo := newFakeCartRepository()

	repo.findByUserAndProductFunc = func(ctx context.Context, userID string, productID string) (*models.CartItem, error) {
		item := models.CartItem{
			ID:        "cart-item-id",
			UserID:    userID,
			ProductID: productID,
			Quantity:  2,
			Product: &models.Product{
				ID:       productID,
				Name:     "iPhone 15",
				Slug:     "iphone-15",
				Price:    15000000,
				Stock:    10,
				IsActive: true,
			},
		}

		return &item, nil
	}

	var updatedQuantity int
	repo.updateQuantityFunc = func(ctx context.Context, id string, userID string, quantity int) (*models.CartItem, error) {
		updatedQuantity = quantity

		item := models.CartItem{
			ID:        id,
			UserID:    userID,
			ProductID: "product-id",
			Quantity:  quantity,
			Product: &models.Product{
				ID:       "product-id",
				Name:     "iPhone 15",
				Slug:     "iphone-15",
				Price:    15000000,
				Stock:    10,
				IsActive: true,
			},
		}
		item.Subtotal = item.Product.Price * float64(quantity)

		return &item, nil
	}

	service := NewCartService(repo, newFakeProductRepository())

	item, err := service.AddItem(context.Background(), "user-id", "product-id", 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedQuantity != 5 {
		t.Fatalf("expected updated quantity 5, got %d", updatedQuantity)
	}

	if item.Quantity != 5 {
		t.Fatalf("expected quantity 5, got %d", item.Quantity)
	}
}

func TestCartService_AddItem_InvalidQuantity(t *testing.T) {
	service := NewCartService(newFakeCartRepository(), newFakeProductRepository())

	_, err := service.AddItem(context.Background(), "user-id", "product-id", 0)

	if !errors.Is(err, models.ErrInvalidCartInput) {
		t.Fatalf("expected ErrInvalidCartInput, got %v", err)
	}
}

func TestCartService_AddItem_ProductNotFound(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.findByIDFunc = func(ctx context.Context, id string) (*models.Product, error) {
		return nil, models.ErrProductNotFound
	}

	service := NewCartService(newFakeCartRepository(), productRepo)

	_, err := service.AddItem(context.Background(), "user-id", "missing-product-id", 1)

	if !errors.Is(err, models.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestCartService_AddItem_InsufficientStock(t *testing.T) {
	productRepo := newFakeProductRepository()
	productRepo.findByIDFunc = func(ctx context.Context, id string) (*models.Product, error) {
		return &models.Product{
			ID:       id,
			Name:     "iPhone 15",
			Slug:     "iphone-15",
			Price:    15000000,
			Stock:    1,
			IsActive: true,
		}, nil
	}

	service := NewCartService(newFakeCartRepository(), productRepo)

	_, err := service.AddItem(context.Background(), "user-id", "product-id", 2)

	if !errors.Is(err, models.ErrInvalidCartInput) {
		t.Fatalf("expected ErrInvalidCartInput, got %v", err)
	}
}

func TestCartService_GetCart_Success(t *testing.T) {
	service := NewCartService(newFakeCartRepository(), newFakeProductRepository())

	cart, err := service.GetCart(context.Background(), "user-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(cart.Items))
	}

	if cart.TotalPrice != 30000000 {
		t.Fatalf("expected total price 30000000, got %f", cart.TotalPrice)
	}
}

func TestCartService_GetCart_InvalidUserID(t *testing.T) {
	service := NewCartService(newFakeCartRepository(), newFakeProductRepository())

	_, err := service.GetCart(context.Background(), "")

	if !errors.Is(err, models.ErrInvalidCartInput) {
		t.Fatalf("expected ErrInvalidCartInput, got %v", err)
	}
}

func TestCartService_UpdateItem_Success(t *testing.T) {
	service := NewCartService(newFakeCartRepository(), newFakeProductRepository())

	item, err := service.UpdateItem(context.Background(), "user-id", "cart-item-id", 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.Quantity != 3 {
		t.Fatalf("expected quantity 3, got %d", item.Quantity)
	}

	if item.Subtotal != 45000000 {
		t.Fatalf("expected subtotal 45000000, got %f", item.Subtotal)
	}
}

func TestCartService_UpdateItem_InvalidQuantity(t *testing.T) {
	service := NewCartService(newFakeCartRepository(), newFakeProductRepository())

	_, err := service.UpdateItem(context.Background(), "user-id", "cart-item-id", 0)

	if !errors.Is(err, models.ErrInvalidCartInput) {
		t.Fatalf("expected ErrInvalidCartInput, got %v", err)
	}
}

func TestCartService_DeleteItem_Success(t *testing.T) {
	service := NewCartService(newFakeCartRepository(), newFakeProductRepository())

	err := service.DeleteItem(context.Background(), "user-id", "cart-item-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCartService_DeleteItem_NotFound(t *testing.T) {
	repo := newFakeCartRepository()
	repo.deleteFunc = func(ctx context.Context, id string, userID string) error {
		return models.ErrCartItemNotFound
	}

	service := NewCartService(repo, newFakeProductRepository())

	err := service.DeleteItem(context.Background(), "user-id", "missing-id")

	if !errors.Is(err, models.ErrCartItemNotFound) {
		t.Fatalf("expected ErrCartItemNotFound, got %v", err)
	}
}
