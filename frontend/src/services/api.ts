import axios from "axios";

import type {
  Category,
  Product,
  ProductListParams,
  ProductListResponse,
} from "../types/product";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL?.trim() || "/api";

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

type DataEnvelope<T> = {
  data: T;
};

function isDataEnvelope<T>(value: unknown): value is DataEnvelope<T> {
  return typeof value === "object" && value !== null && "data" in value;
}

function unwrapData<T>(value: T | DataEnvelope<T>): T {
  return isDataEnvelope<T>(value) ? value.data : value;
}

function cleanParams<T extends object>(params?: T): Partial<T> {
  return Object.fromEntries(
    Object.entries(params ?? {}).filter(([, value]) => {
      return value !== undefined && value !== null && value !== "";
    }),
  ) as Partial<T>;
}

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");

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
      localStorage.removeItem("token");

      if (window.location.pathname !== "/login") {
        window.location.href = "/login";
      }
    }

    return Promise.reject(error);
  },
);

export const productApi = {
  async list(params?: ProductListParams): Promise<ProductListResponse> {
    const response = await api.get<ProductListResponse>("/v1/products", {
      params: cleanParams(params),
    });

    return response.data;
  },

  async getById(id: string): Promise<Product> {
    const response = await api.get<Product | DataEnvelope<Product>>(
      `/v1/products/${encodeURIComponent(id)}`,
    );

    return unwrapData(response.data);
  },

  async getBySlug(slug: string): Promise<Product> {
    const response = await api.get<Product | DataEnvelope<Product>>(
      `/v1/products/slug/${encodeURIComponent(slug)}`,
    );

    return unwrapData(response.data);
  },
};

export const categoryApi = {
  async list(): Promise<Category[]> {
    const response = await api.get<Category[] | DataEnvelope<Category[]>>(
      "/v1/categories",
    );

    return unwrapData(response.data);
  },

  async getById(id: string): Promise<Category> {
    const response = await api.get<Category | DataEnvelope<Category>>(
      `/v1/categories/${encodeURIComponent(id)}`,
    );

    return unwrapData(response.data);
  },

  async getBySlug(slug: string): Promise<Category> {
    const response = await api.get<Category | DataEnvelope<Category>>(
      `/v1/categories/slug/${encodeURIComponent(slug)}`,
    );

    return unwrapData(response.data);
  },
};

export default api;