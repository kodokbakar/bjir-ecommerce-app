import api from "./api";
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

export async function getDashboardProducts(): Promise<ProductListResponse> {
  return listProducts({
    page: 1,
    limit: 4,
    sort_by: "created_at",
    sort_order: "desc",
  });
}

export async function getRecentOrders(limit = 3): Promise<DashboardOrderResult> {
  const response = await api.get<ApiListResponse<DashboardOrder[]>>("/v1/orders", {
    params: {
      page: 1,
      limit,
    },
  });

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
