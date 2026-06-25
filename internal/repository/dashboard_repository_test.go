package repository

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v5"
)

func TestDashboardRepository_GetStats_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewDashboardRepository(mock)

	rows := pgxmock.NewRows([]string{
		"total_orders",
		"total_revenue",
		"pending_orders",
		"completed_today",
		"revenue_today",
		"total_products",
		"total_categories",
	}).AddRow(
		150,
		45000000.0,
		12,
		8,
		2500000.0,
		85,
		10,
	)

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	stats, err := repo.GetStats(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if stats.TotalOrders != 150 {
		t.Fatalf("expected total orders 150, got %d", stats.TotalOrders)
	}

	if stats.TotalRevenue != 45000000.0 {
		t.Fatalf("expected total revenue 45000000, got %f", stats.TotalRevenue)
	}

	if stats.PendingOrders != 12 {
		t.Fatalf("expected pending orders 12, got %d", stats.PendingOrders)
	}

	if stats.CompletedToday != 8 {
		t.Fatalf("expected completed today 8, got %d", stats.CompletedToday)
	}

	if stats.RevenueToday != 2500000.0 {
		t.Fatalf("expected revenue today 2500000, got %f", stats.RevenueToday)
	}

	if stats.TotalProducts != 85 {
		t.Fatalf("expected total products 85, got %d", stats.TotalProducts)
	}

	if stats.TotalCategories != 10 {
		t.Fatalf("expected total categories 10, got %d", stats.TotalCategories)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestDashboardRepository_GetStats_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := NewDashboardRepository(mock)

	mock.ExpectQuery("SELECT").
		WillReturnError(context.Canceled)

	stats, err := repo.GetStats(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}

	if stats != nil {
		t.Fatal("expected nil stats")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
