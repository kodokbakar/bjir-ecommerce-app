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
	GetMyOrders(ctx context.Context, userID string, input OrderListInput) (*OrderListResult, error)
	GetMyOrderDetail(ctx context.Context, userID string, orderID string) (*models.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status string) (*models.Order, error)
}

const (
	DefaultOrderPage  = 1
	DefaultOrderLimit = 20
	MaxOrderLimit     = 100
)

type OrderListInput struct {
	Page  int
	Limit int
}

type OrderListResult struct {
	Orders     []models.Order `json:"orders"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
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

func (s *orderService) GetMyOrders(ctx context.Context, userID string, input OrderListInput) (*OrderListResult, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", models.ErrInvalidOrderInput)
	}

	page := input.Page
	if page < 1 {
		page = DefaultOrderPage
	}

	limit := input.Limit
	if limit < 1 {
		limit = DefaultOrderLimit
	}

	if limit > MaxOrderLimit {
		limit = MaxOrderLimit
	}

	offset := (page - 1) * limit

	orders, total, err := s.orderRepo.FindAllByUserID(ctx, userID, repository.OrderListFilter{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &OrderListResult{
		Orders:     orders,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *orderService) GetMyOrderDetail(ctx context.Context, userID string, orderID string) (*models.Order, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, fmt.Errorf("%w: user id is required", models.ErrInvalidOrderInput)
	}

	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order id is required", models.ErrInvalidOrderInput)
	}

	order, err := s.orderRepo.FindByIDAndUserID(ctx, orderID, userID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) UpdateStatus(ctx context.Context, orderID string, status string) (*models.Order, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return nil, fmt.Errorf("%w: order id is required", models.ErrInvalidOrderInput)
	}

	nextStatus, err := normalizeOrderStatus(status)
	if err != nil {
		return nil, err
	}

	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.Status == nextStatus {
		return order, nil
	}

	if !isValidOrderStatusTransition(order.Status, nextStatus) {
		return nil, fmt.Errorf(
			"%w: cannot change status from %s to %s",
			models.ErrInvalidOrderStatusTransition,
			order.Status,
			nextStatus,
		)
	}

	updatedOrder, err := s.orderRepo.UpdateStatus(ctx, orderID, order.Status, nextStatus)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func normalizeOrderStatus(status string) (string, error) {
	status = strings.ToLower(strings.TrimSpace(status))

	switch status {
	case models.OrderStatusPending:
		return models.OrderStatusPending, nil
	case models.OrderStatusPaid:
		return models.OrderStatusPaid, nil
	case models.OrderStatusShipped:
		return models.OrderStatusShipped, nil
	case models.OrderStatusDelivered:
		return models.OrderStatusDelivered, nil
	case models.OrderStatusCancelled:
		return models.OrderStatusCancelled, nil
	case "canceled":
		return models.OrderStatusCancelled, nil
	default:
		return "", fmt.Errorf("%w: unsupported status %q", models.ErrInvalidOrderStatus, status)
	}
}

func isValidOrderStatusTransition(currentStatus string, nextStatus string) bool {
	if currentStatus == nextStatus {
		return true
	}

	switch currentStatus {
	case models.OrderStatusPending:
		return nextStatus == models.OrderStatusPaid ||
			nextStatus == models.OrderStatusCancelled
	case models.OrderStatusPaid:
		return nextStatus == models.OrderStatusShipped
	case models.OrderStatusShipped:
		return nextStatus == models.OrderStatusDelivered
	default:
		return false
	}
}
