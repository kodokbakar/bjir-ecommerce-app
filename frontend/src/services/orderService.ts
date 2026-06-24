import api, { getResponseErrorMessage } from "./api";
import type {
  Order,
  OrderListMeta,
  OrderListParams,
  OrderListResponse,
} from "../types/order";

interface ApiListResponse<TData> {
  success?: boolean;
  message?: string;
  error?: string;
  details?: string;
  data: TData;
  meta?: Partial<OrderListMeta>;
}

interface ApiDataResponse<TData> {
  success?: boolean;
  message?: string;
  error?: string;
  details?: string;
  data: TData;
}

const DEFAULT_ORDER_META: OrderListMeta = {
  page: 1,
  limit: 8,
  total: 0,
  total_pages: 0,
};

function cleanParams<T extends object>(params?: T): Partial<T> {
  return Object.fromEntries(
    Object.entries(params ?? {}).filter(([, value]) => {
      return value !== undefined && value !== null && value !== "";
    }),
  ) as Partial<T>;
}

function normalizeMeta(
  meta: Partial<OrderListMeta> | undefined,
): OrderListMeta {
  return {
    page: meta?.page ?? DEFAULT_ORDER_META.page,
    limit: meta?.limit ?? DEFAULT_ORDER_META.limit,
    total: meta?.total ?? DEFAULT_ORDER_META.total,
    total_pages: meta?.total_pages ?? DEFAULT_ORDER_META.total_pages,
  };
}

export function getOrderErrorMessage(error: unknown, fallback: string): string {
  return getResponseErrorMessage(error, fallback);
}

export async function listOrders(
  params: OrderListParams = {},
): Promise<OrderListResponse> {
  const response = await api.get<ApiListResponse<Order[]>>("/v1/orders", {
    params: cleanParams(params),
  });

  return {
    data: response.data.data,
    meta: normalizeMeta(response.data.meta),
  };
}

export async function listAdminOrders(
  params: OrderListParams = {},
): Promise<OrderListResponse> {
  const response = await api.get<ApiListResponse<Order[]>>("/v1/admin/orders", {
    params: cleanParams(params),
  });

  return {
    data: response.data.data,
    meta: normalizeMeta(response.data.meta),
  };
}

export async function getOrderById(orderID: string): Promise<Order> {
  const response = await api.get<ApiDataResponse<Order>>(
    `/v1/orders/${encodeURIComponent(orderID)}`,
  );

  return response.data.data;
}
