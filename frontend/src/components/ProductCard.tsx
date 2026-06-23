import { useState } from "react";
import { Link } from "react-router-dom";

import { getImageUrl } from "../services/productService";
import { C } from "../styles/tokens";
import type { Product } from "../types/product";

interface ProductCardProps {
  product: Product;
}

type StockState = {
  label: string;
  className: string;
};

function formatRupiah(value: number): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);
}

function getStockState(stock: number): StockState {
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

function getProductImage(product: Product): string {
  const galleryImage = product.images?.[0]?.image_url;
  return product.image_url || galleryImage || "";
}

function ProductPlaceholder() {
  return (
    <span className="product-card-placeholder" aria-hidden="true">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.4">
        <path d="M4 7h16v12H4z" />
        <path d="M8 7a4 4 0 0 1 8 0" />
        <path d="M8 13h8" />
      </svg>
    </span>
  );
}

function ProductCard({ product }: ProductCardProps) {
  const [imageFailed, setImageFailed] = useState(false);

  const stockState = getStockState(product.stock);
  const imagePath = getProductImage(product);
  const imageUrl = imagePath ? getImageUrl(imagePath) : "";
  const categoryName = product.category?.name || "Uncategorized";

  return (
    <Link
      className="product-card"
      to={`/products/${product.slug}`}
      aria-label={`View product ${product.name}`}
    >
      <div className="product-card-media">
        <span className={`product-card-stock ${stockState.className}`}>
          {stockState.label}
        </span>

        {imageUrl && !imageFailed ? (
          <img
            src={imageUrl}
            alt={product.name}
            loading="lazy"
            onError={() => setImageFailed(true)}
          />
        ) : (
          <ProductPlaceholder />
        )}
      </div>

      <div className="product-card-body">
        <span className="product-card-category">{categoryName}</span>

        <h3 className="product-card-title">{product.name}</h3>

        <div className="product-card-footer">
          <span className="product-card-price" style={{ color: C.primaryDark }}>
            {formatRupiah(product.price)}
          </span>

          <span className="product-card-stock-count">
            {product.stock > 0 ? `${product.stock} left` : "Sold out"}
          </span>
        </div>
      </div>
    </Link>
  );
}

export default ProductCard;