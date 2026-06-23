import axios from "axios";

import api from "./api";
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
  if (axios.isAxiosError(error)) {
    const responseData = error.response?.data as
      | {
          message?: string;
          error?: string;
          details?: string;
        }
      | undefined;

    return (
      responseData?.details ||
      responseData?.message ||
      responseData?.error ||
      fallback
    );
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return fallback;
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
