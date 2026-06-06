package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type OrderService interface {
	Checkout(ctx context.Context, userID string) (*models.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{orderRepo: orderRepo}
}

func (s *orderService) Checkout(ctx context.Context, userID string) (*models.Order, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", models.ErrInvalidOrderInput)
	}

	order, err := s.orderRepo.Checkout(ctx, userID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
