package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type fakeDashboardService struct {
	getStatsFunc func(ctx context.Context) (*models.DashboardStats, error)
}

func (f *fakeDashboardService) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	return f.getStatsFunc(ctx)
}

func setupDashboardRouter(service *fakeDashboardService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewDashboardHandler(service)

	router.GET("/api/v1/admin/dashboard", handler.GetStats)

	return router
}

func TestDashboardHandler_GetStats_Success(t *testing.T) {
	service := &fakeDashboardService{
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

	router := setupDashboardRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "admin dashboard retrieved successfully") {
		t.Fatalf("expected success message, got: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), `"total_orders":150`) {
		t.Fatalf("expected total orders in response, got: %s", w.Body.String())
	}
}

func TestDashboardHandler_GetStats_Error(t *testing.T) {
	service := &fakeDashboardService{
		getStatsFunc: func(ctx context.Context) (*models.DashboardStats, error) {
			return nil, errors.New("database error")
		},
	}

	router := setupDashboardRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/dashboard", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d. body: %s", w.Code, w.Body.String())
	}
}
