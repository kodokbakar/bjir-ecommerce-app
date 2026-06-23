import { categoryApi, productApi } from "./api";

import type {
  Category,
  Product,
  ProductListParams,
  ProductListResponse,
} from "../types/product";
import { getImageUrl } from "../utils/image";

export { getImageUrl };

export class ProductServiceValidationError extends Error {
  readonly fieldName: string;

  constructor(fieldName: string) {
    super(`${fieldName} is required`);
    this.name = "ProductServiceValidationError";
    this.fieldName = fieldName;
  }
}

function requireValue(value: string, fieldName: string): string {
  const trimmed = value.trim();

  if (!trimmed) {
    throw new ProductServiceValidationError(fieldName);
  }

  return trimmed;
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
