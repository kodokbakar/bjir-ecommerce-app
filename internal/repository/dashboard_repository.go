package repository

import (
	"context"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type DashboardRepository interface {
	GetStats(ctx context.Context) (*models.DashboardStats, error)
}

type dashboardRepository struct {
	db PgxQuerier
}

func NewDashboardRepository(db PgxQuerier) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM orders)::int AS total_orders,
			COALESCE((
				SELECT SUM(total_amount)
				FROM orders
				WHERE status = 'delivered'
			), 0)::float8 AS total_revenue,
			(SELECT COUNT(*) FROM orders WHERE status = 'pending')::int AS pending_orders,
			(
				SELECT COUNT(*)
				FROM orders
				WHERE status = 'delivered'
				AND created_at >= CURRENT_DATE
				AND created_at < CURRENT_DATE + INTERVAL '1 day'
			)::int AS completed_today,
			COALESCE((
				SELECT SUM(total_amount)
				FROM orders
				WHERE status = 'delivered'
				AND created_at >= CURRENT_DATE
				AND created_at < CURRENT_DATE + INTERVAL '1 day'
			), 0)::float8 AS revenue_today,
			(SELECT COUNT(*) FROM products)::int AS total_products,
			(SELECT COUNT(*) FROM categories)::int AS total_categories
	`

	stats := &models.DashboardStats{}

	if err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalOrders,
		&stats.TotalRevenue,
		&stats.PendingOrders,
		&stats.CompletedToday,
		&stats.RevenueToday,
		&stats.TotalProducts,
		&stats.TotalCategories,
	); err != nil {
		return nil, err
	}

	return stats, nil
}
