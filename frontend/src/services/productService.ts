import { categoryApi, productApi } from "./api";

import type {
  Category,
  Product,
  ProductListParams,
  ProductListResponse,
} from "../types/product";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL?.trim() || "";

function getApiOrigin(): string {
  if (!API_BASE_URL || API_BASE_URL === "/api") {
    return "";
  }

  return API_BASE_URL.replace(/\/api\/?$/, "").replace(/\/$/, "");
}

function requireValue(value: string, fieldName: string): string {
  const trimmed = value.trim();

  if (!trimmed) {
    throw new Error(`${fieldName} is required`);
  }

  return trimmed;
}

export function getImageUrl(path: string): string {
  const trimmed = path.trim();

  if (!trimmed) {
    return "";
  }

  if (/^https?:\/\//i.test(trimmed)) {
    return trimmed;
  }

  const imagePath = trimmed.replace(/^\/+/, "");
  const uploadPath = imagePath.startsWith("uploads/")
    ? `/${imagePath}`
    : `/uploads/${imagePath}`;

  return `${getApiOrigin()}${uploadPath}`;
}

export async function listProducts(
  params?: ProductListParams,
): Promise<ProductListResponse> {
  return productApi.list(params);
}

export async function getProductById(id: string): Promise<Product> {
  return productApi.getById(requireValue(id, "Product ID"));
}

export async function getProductBySlug(slug: string): Promise<Product> {
  return productApi.getBySlug(requireValue(slug, "Product slug"));
}

export async function listCategories(): Promise<Category[]> {
  return categoryApi.list();
}

export async function getCategoryById(id: string): Promise<Category> {
  return categoryApi.getById(requireValue(id, "Category ID"));
}

export async function getCategoryBySlug(slug: string): Promise<Category> {
  return categoryApi.getBySlug(requireValue(slug, "Category slug"));
}

export const productService = {
  listProducts,
  getProductById,
  getProductBySlug,
  listCategories,
  getCategoryById,
  getCategoryBySlug,
  getImageUrl,
};