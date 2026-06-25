package services

import (
	"context"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

type DashboardService interface {
	GetStats(ctx context.Context) (*models.DashboardStats, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	return s.dashboardRepo.GetStats(ctx)
}
