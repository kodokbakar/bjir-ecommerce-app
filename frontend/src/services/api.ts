import axios from "axios";

import type {
  Category,
  Product,
  ProductInput,
  ProductListMeta,
  ProductListParams,
  ProductListResponse,
} from "../types/product";

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL?.trim() || "/api";

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

type ApiDataResponse<T> = {
  success?: boolean;
  message?: string;
  data: T;
};

type ApiListResponse<TData, TMeta> = ApiDataResponse<TData> & {
  meta: TMeta;
};

function isApiDataResponse<T>(value: unknown): value is ApiDataResponse<T> {
  return typeof value === "object" && value !== null && "data" in value;
}

function unwrapData<T>(value: T | ApiDataResponse<T>): T {
  return isApiDataResponse<T>(value) ? value.data : value;
}

function cleanParams<T extends object>(params?: T): Partial<T> {
  return Object.fromEntries(
    Object.entries(params ?? {}).filter(([, value]) => {
      return value !== undefined && value !== null && value !== "";
    }),
  ) as Partial<T>;
}

export function readStoredToken(): string | null {
  return localStorage.getItem("token") || sessionStorage.getItem("token");
}

export function clearAuthStorage(): void {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  sessionStorage.removeItem("token");
  sessionStorage.removeItem("user");
}

export function getApiOrigin(): string {
  if (!API_BASE_URL || API_BASE_URL.startsWith("/")) {
    return "";
  }

  try {
    return new URL(API_BASE_URL).origin;
  } catch {
    return "";
  }
}

export function getResponseErrorMessage(
  error: unknown,
  fallback: string,
): string {
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

api.interceptors.request.use(
  (config) => {
    const token = readStoredToken();

    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => Promise.reject(error),
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (axios.isAxiosError(error) && error.response?.status === 401) {
      clearAuthStorage();

      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }

    return Promise.reject(error);
  },
);

export const productApi = {
  async list(params?: ProductListParams): Promise<ProductListResponse> {
    const response = await api.get<ApiListResponse<Product[], ProductListMeta>>(
      "/v1/products",
      {
        params: cleanParams(params),
      },
    );

    return {
      data: response.data.data,
      meta: response.data.meta,
    };
  },

  async getById(id: string): Promise<Product> {
    const response = await api.get<Product | ApiDataResponse<Product>>(
      `/v1/products/${encodeURIComponent(id)}`,
    );

    return unwrapData(response.data);
  },

  async getBySlug(slug: string): Promise<Product> {
    const response = await api.get<Product | ApiDataResponse<Product>>(
      `/v1/products/slug/${encodeURIComponent(slug)}`,
    );

    return unwrapData(response.data);
  },

  async create(input: ProductInput): Promise<Product> {
    const response = await api.post<ApiDataResponse<Product>>(
      "/v1/products",
      input,
    );

    return response.data.data;
  },

  async update(id: string, input: ProductInput): Promise<Product> {
    const response = await api.put<ApiDataResponse<Product>>(
      `/v1/products/${encodeURIComponent(id)}`,
      input,
    );

    return response.data.data;
  },

  async uploadImage(id: string, file: File): Promise<void> {
    const formData = new FormData();
    formData.append("file", file);

    await api.post(`/v1/products/${encodeURIComponent(id)}/image`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
  },

  async deleteById(id: string): Promise<void> {
    await api.delete(`/v1/products/${encodeURIComponent(id)}`);
  },
};

export const categoryApi = {
  async list(): Promise<Category[]> {
    const response = await api.get<Category[] | ApiDataResponse<Category[]>>(
      "/v1/categories",
    );

    return unwrapData(response.data);
  },

  async getById(id: string): Promise<Category> {
    const response = await api.get<Category | ApiDataResponse<Category>>(
      `/v1/categories/${encodeURIComponent(id)}`,
    );

    return unwrapData(response.data);
  },

  async getBySlug(slug: string): Promise<Category> {
    const response = await api.get<Category | ApiDataResponse<Category>>(
      `/v1/categories/slug/${encodeURIComponent(slug)}`,
    );

    return unwrapData(response.data);
  },
};

export default api;
