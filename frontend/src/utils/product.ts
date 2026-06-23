import type { Product } from "../types/product";

export type ProductStockState = {
  label: string;
  className: string;
};

export function formatRupiah(value: number): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);
}

export function getStockState(stock: number): ProductStockState {
  if (stock <= 0) {
    return {
      label: "Out of Stock",
      className: "is-out",
    };
  }

  if (stock <= 5) {
    return {
      label: "Low Stock",
      className: "is-low",
    };
  }

  return {
    label: "In Stock",
    className: "is-in",
  };
}

export function getProductImage(product: Product): string {
  const galleryImage = product.images?.[0]?.image_url;
  return product.image_url || galleryImage || "";
}
