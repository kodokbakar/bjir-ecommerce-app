import api, { getResponseErrorMessage } from "./api";
import { listProducts } from "./productService";
import type { Product, ProductListResponse } from "../types/product";

interface ApiListResponse<TData> {
  success?: boolean;
  message?: string;
  data: TData;
  meta?: {
    page?: number;
    limit?: number;
    total?: number;
    total_pages?: number;
  };
}

interface ApiDataResponse<TData> {
  success?: boolean;
  message?: string;
  data: TData;
}

interface AdminDashboardStatsResponse {
  total_orders: number;
  total_revenue: number;
  pending_orders: number;
  completed_today: number;
  revenue_today: number;
  total_products: number;
  total_categories: number;
}

export interface AdminDashboardStats {
  totalOrders: number;
  totalRevenue: number;
  pendingOrders: number;
  completedToday: number;
  revenueToday: number;
  totalProducts: number;
  totalCategories: number;
}

export interface DashboardOrder {
  id: string;
  status: string;
  total?: number;
  total_price?: number;
  grand_total?: number;
  created_at?: string;
  updated_at?: string;
}

export interface DashboardOrderResult {
  data: DashboardOrder[];
  total: number;
}

function mapAdminDashboardStats(
  stats: AdminDashboardStatsResponse,
): AdminDashboardStats {
  return {
    totalOrders: stats.total_orders,
    totalRevenue: stats.total_revenue,
    pendingOrders: stats.pending_orders,
    completedToday: stats.completed_today,
    revenueToday: stats.revenue_today,
    totalProducts: stats.total_products,
    totalCategories: stats.total_categories,
  };
}

export async function getAdminDashboardStats(): Promise<AdminDashboardStats> {
  const response = await api.get<ApiDataResponse<AdminDashboardStatsResponse>>(
    "/v1/admin/dashboard",
  );

  return mapAdminDashboardStats(response.data.data);
}

export function getDashboardErrorMessage(
  error: unknown,
  fallback: string,
): string {
  return getResponseErrorMessage(error, fallback);
}

export async function getDashboardProducts(): Promise<ProductListResponse> {
  return listProducts({
    page: 1,
    limit: 4,
    sort_by: "created_at",
    sort_order: "desc",
  });
}

export async function getRecentOrders(
  limit = 3,
): Promise<DashboardOrderResult> {
  const response = await api.get<ApiListResponse<DashboardOrder[]>>(
    "/v1/orders",
    {
      params: {
        page: 1,
        limit,
      },
    },
  );

  return {
    data: response.data.data,
    total: response.data.meta?.total ?? response.data.data.length,
  };
}

export function getActiveOrderCount(orders: DashboardOrder[]): number {
  return orders.filter((order) => {
    return ["pending", "paid", "shipped"].includes(order.status);
  }).length;
}

export function getOrderTotal(order: DashboardOrder): number | null {
  return order.grand_total ?? order.total_price ?? order.total ?? null;
}

export type { Product };
