package services

import (
	"context"
	"errors"
	"testing"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type fakeDashboardRepository struct {
	getStatsFunc func(ctx context.Context) (*models.DashboardStats, error)
}

func (f *fakeDashboardRepository) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	return f.getStatsFunc(ctx)
}

func TestDashboardService_GetStats_Success(t *testing.T) {
	repo := &fakeDashboardRepository{
		getStatsFunc: func(ctx context.Context) (*models.DashboardStats, error) {
			return &models.DashboardStats{
				TotalOrders:     150,
				TotalRevenue:    45000000,
				PendingOrders:   12,
				CompletedToday:  8,
				RevenueToday:    2500000,
				TotalProducts:   85,
				TotalCategories: 10,
			}, nil
		},
	}

	service := NewDashboardService(repo)

	stats, err := service.GetStats(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if stats.TotalOrders != 150 {
		t.Fatalf("expected total orders 150, got %d", stats.TotalOrders)
	}
}

func TestDashboardService_GetStats_Error(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakeDashboardRepository{
		getStatsFunc: func(ctx context.Context) (*models.DashboardStats, error) {
			return nil, expectedErr
		},
	}

	service := NewDashboardService(repo)

	stats, err := service.GetStats(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}

	if stats != nil {
		t.Fatal("expected nil stats")
	}
}
