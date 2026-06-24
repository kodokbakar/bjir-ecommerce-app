import api, { getResponseErrorMessage } from "./api";
import type {
  Category,
  CategoryInput,
  CategoryListMeta,
  CategoryListParams,
  CategoryListResponse,
} from "../types/product";

interface ApiDataResponse<TData> {
  success?: boolean;
  message?: string;
  error?: string;
  details?: string;
  data: TData;
}

interface ApiListResponse<TData> extends ApiDataResponse<TData> {
  meta?: Partial<CategoryListMeta>;
}

const DEFAULT_CATEGORY_META: CategoryListMeta = {
  page: 1,
  limit: 10,
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
  meta: Partial<CategoryListMeta> | undefined,
  fallback: CategoryListMeta,
): CategoryListMeta {
  return {
    page: meta?.page ?? fallback.page,
    limit: meta?.limit ?? fallback.limit,
    total: meta?.total ?? fallback.total,
    total_pages: meta?.total_pages ?? fallback.total_pages,
  };
}

function unwrapCategoryList(
  payload: Category[] | ApiListResponse<Category[]>,
  fallbackMeta: CategoryListMeta,
): CategoryListResponse {
  if (Array.isArray(payload)) {
    return {
      data: payload,
      meta: {
        ...fallbackMeta,
        total: payload.length,
        total_pages: payload.length > 0 ? 1 : 0,
      },
    };
  }

  return {
    data: payload.data,
    meta: normalizeMeta(payload.meta, fallbackMeta),
  };
}

function unwrapCategory(
  payload: Category | ApiDataResponse<Category>,
): Category {
  if (typeof payload === "object" && payload !== null && "data" in payload) {
    return payload.data;
  }

  return payload;
}

export async function listAdminCategories(
  params: CategoryListParams = {},
): Promise<CategoryListResponse> {
  const fallbackMeta: CategoryListMeta = {
    ...DEFAULT_CATEGORY_META,
    page: params.page ?? DEFAULT_CATEGORY_META.page,
    limit: params.limit ?? DEFAULT_CATEGORY_META.limit,
  };

  const response = await api.get<Category[] | ApiListResponse<Category[]>>(
    "/v1/categories",
    {
      params: cleanParams(params),
    },
  );

  return unwrapCategoryList(response.data, fallbackMeta);
}

export async function createCategory(input: CategoryInput): Promise<Category> {
  const response = await api.post<Category | ApiDataResponse<Category>>(
    "/v1/admin/categories",
    input,
  );

  return unwrapCategory(response.data);
}

export async function updateCategory(
  id: string,
  input: CategoryInput,
): Promise<Category> {
  const response = await api.put<Category | ApiDataResponse<Category>>(
    `/v1/admin/categories/${encodeURIComponent(id)}`,
    input,
  );

  return unwrapCategory(response.data);
}

export async function deleteCategory(id: string): Promise<void> {
  await api.delete(`/v1/admin/categories/${encodeURIComponent(id)}`);
}

export function getCategoryErrorMessage(
  error: unknown,
  fallback: string,
): string {
  return getResponseErrorMessage(error, fallback);
}
