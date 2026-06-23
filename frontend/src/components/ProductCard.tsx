import { Link } from "react-router-dom";

import ProductImage from "./ProductImage";
import type { Product } from "../types/product";
import { formatRupiah, getProductImage, getStockState } from "../utils/product";

interface ProductCardProps {
  product: Product;
}

function ProductCard({ product }: ProductCardProps) {
  const stockState = getStockState(product.stock);
  const imagePath = getProductImage(product);
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

        <ProductImage
          key={imagePath || product.id}
          className="product-card-image"
          src={imagePath}
          alt={product.name}
          width={400}
          height={300}
        />
      </div>

      <div className="product-card-body">
        <span className="product-card-category">{categoryName}</span>

        <h3 className="product-card-title">{product.name}</h3>

        <div className="product-card-footer">
          <span className="product-card-price">{formatRupiah(product.price)}</span>

          <span className="product-card-stock-count">
            {product.stock > 0 ? `${product.stock} left` : "Sold out"}
          </span>
        </div>
      </div>
    </Link>
  );
}

export default ProductCard;
