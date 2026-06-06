package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type CartService interface {
	AddItem(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error)
	GetCart(ctx context.Context, userID string) (*models.Cart, error)
	UpdateItem(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error)
	DeleteItem(ctx context.Context, userID string, itemID string) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) AddItem(ctx context.Context, userID string, productID string, quantity int) (*models.CartItem, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", models.ErrInvalidCartInput)
	}

	productID = strings.TrimSpace(productID)
	if productID == "" {
		return nil, fmt.Errorf("%w: product id is required", models.ErrInvalidCartInput)
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than 0", models.ErrInvalidCartInput)
	}

	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	existingItem, err := s.cartRepo.FindByUserAndProduct(ctx, userID, productID)
	if err != nil && !errors.Is(err, models.ErrCartItemNotFound) {
		return nil, err
	}

	if existingItem != nil {
		newQuantity := existingItem.Quantity + quantity
		if newQuantity > product.Stock {
			return nil, fmt.Errorf("%w: quantity exceeds product stock", models.ErrInvalidCartInput)
		}

		return s.cartRepo.UpdateQuantity(ctx, existingItem.ID, userID, newQuantity)
	}

	if quantity > product.Stock {
		return nil, fmt.Errorf("%w: quantity exceeds product stock", models.ErrInvalidCartInput)
	}

	item := &models.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	if err := s.cartRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *cartService) GetCart(ctx context.Context, userID string) (*models.Cart, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", models.ErrInvalidCartInput)
	}

	items, err := s.cartRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.Subtotal
	}

	return &models.Cart{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}

func (s *cartService) UpdateItem(ctx context.Context, userID string, itemID string, quantity int) (*models.CartItem, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", models.ErrInvalidCartInput)
	}

	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return nil, fmt.Errorf("%w: cart item id is required", models.ErrInvalidCartInput)
	}

	if quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be greater than 0", models.ErrInvalidCartInput)
	}

	existingItem, err := s.cartRepo.FindByID(ctx, itemID, userID)
	if err != nil {
		return nil, err
	}

	if existingItem.Product == nil {
		return nil, fmt.Errorf("%w: product is required", models.ErrInvalidCartInput)
	}

	if quantity > existingItem.Product.Stock {
		return nil, fmt.Errorf("%w: quantity exceeds product stock", models.ErrInvalidCartInput)
	}

	item, err := s.cartRepo.UpdateQuantity(ctx, itemID, userID, quantity)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *cartService) DeleteItem(ctx context.Context, userID string, itemID string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return fmt.Errorf("%w: user id is required", models.ErrInvalidCartInput)
	}

	itemID = strings.TrimSpace(itemID)
	if itemID == "" {
		return fmt.Errorf("%w: cart item id is required", models.ErrInvalidCartInput)
	}

	return s.cartRepo.Delete(ctx, itemID, userID)
}
