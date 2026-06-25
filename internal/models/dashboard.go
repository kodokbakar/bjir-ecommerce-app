package models

type DashboardStats struct {
	TotalOrders     int     `json:"total_orders"`
	TotalRevenue    float64 `json:"total_revenue"`
	PendingOrders   int     `json:"pending_orders"`
	CompletedToday  int     `json:"completed_today"`
	RevenueToday    float64 `json:"revenue_today"`
	TotalProducts   int     `json:"total_products"`
	TotalCategories int     `json:"total_categories"`
}
