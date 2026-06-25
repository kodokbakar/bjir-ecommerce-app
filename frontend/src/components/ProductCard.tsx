import { useState } from "react";
import { Link } from "react-router-dom";

import ProductImage from "./ProductImage";
import { addCartItem, getCartErrorMessage } from "../services/cartService";
import type { Product } from "../types/product";
import { formatRupiah, getProductImage, getStockState } from "../utils/product";

interface ProductCardProps {
  product: Product;
}

type CartActionState = "idle" | "loading" | "success" | "error";

function ProductCard({ product }: ProductCardProps) {
  const stockState = getStockState(product.stock);
  const imagePath = getProductImage(product);
  const categoryName = product.category?.name || "Uncategorized";
  const [cartState, setCartState] = useState<CartActionState>("idle");
  const [cartMessage, setCartMessage] = useState("");

  const isSoldOut = product.stock <= 0;
  const isAdding = cartState === "loading";
  const feedbackID = `product-card-cart-feedback-${product.id}`;

  async function handleAddToCart() {
    if (isSoldOut || isAdding) {
      return;
    }

    setCartState("loading");
    setCartMessage("");

    try {
      await addCartItem(product.id, 1);
      setCartState("success");
      setCartMessage("Added to cart.");
    } catch (error) {
      setCartState("error");
      setCartMessage(
        getCartErrorMessage(error, "Failed to add this product to cart."),
      );
    }
  }

  return (
    <article className="product-card">
      <Link
        className="product-card-link"
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
            width={640}
            height={480}
            sizes="(max-width: 720px) 100vw, (max-width: 1180px) 50vw, 25vw"
          />
        </div>

        <div className="product-card-body">
          <span className="product-card-category">{categoryName}</span>

          <h3 className="product-card-title">{product.name}</h3>

          <div className="product-card-footer">
            <span className="product-card-price">
              {formatRupiah(product.price)}
            </span>

            <span className="product-card-stock-count">
              {product.stock > 0 ? `${product.stock} left` : "Sold out"}
            </span>
          </div>
        </div>
      </Link>

      <div className="product-card-actions">
        <button
          className={[
            "product-card-cart-button",
            cartState === "success" ? "is-success" : "",
          ]
            .filter(Boolean)
            .join(" ")}
          type="button"
          onClick={handleAddToCart}
          disabled={isSoldOut || isAdding}
          aria-describedby={cartMessage ? feedbackID : undefined}
        >
          {isSoldOut ? "Sold Out" : isAdding ? "Adding..." : "Add to Cart"}
        </button>

        {cartMessage && (
          <p
            className={`product-card-cart-feedback is-${cartState}`}
            id={feedbackID}
            role={cartState === "error" ? "alert" : "status"}
          >
            {cartMessage}
          </p>
        )}
      </div>
    </article>
  );
}

export default ProductCard;
