import { Link } from "react-router-dom";

import ProductImage from "./ProductImage";
import { getCartItemPrice, getCartItemSubtotal } from "../services/cartService";
import type { Cart, CartItem } from "../types/cart";
import { formatRupiah, getProductImage } from "../utils/product";

interface OrderSummaryProps {
  cart: Cart;
}

function getItemName(item: CartItem): string {
  return item.product?.name || "Unavailable product";
}

function getItemSlug(item: CartItem): string {
  return item.product?.slug || "";
}

function getItemImage(item: CartItem): string {
  return item.product ? getProductImage(item.product) : "";
}

function OrderSummary({ cart }: OrderSummaryProps) {
  const itemCount = cart.items.reduce(
    (total, item) => total + item.quantity,
    0,
  );

  return (
    <section className="order-summary" aria-labelledby="order-summary-title">
      <div className="order-summary-header">
        <span className="cart-summary-label">Order Summary</span>
        <h2 id="order-summary-title">Final shelf check.</h2>
        <p>
          {itemCount} item{itemCount === 1 ? "" : "s"} ready to become an order.
        </p>
      </div>

      <div className="order-summary-list">
        {cart.items.map((item) => {
          const productName = getItemName(item);
          const productSlug = getItemSlug(item);
          const productPath = productSlug
            ? `/products/${productSlug}`
            : "/products";
          const imagePath = getItemImage(item);

          return (
            <article className="order-summary-item" key={item.id}>
              <Link
                className="order-summary-image-link"
                to={productPath}
                aria-label={`View ${productName}`}
              >
                <ProductImage
                  key={`${item.id}-${imagePath || ""}`}
                  className="order-summary-image"
                  src={imagePath}
                  alt={productName}
                  width={160}
                  height={140}
                />
              </Link>

              <div className="order-summary-item-main">
                <Link className="order-summary-item-title" to={productPath}>
                  {productName}
                </Link>

                <div className="order-summary-item-meta">
                  <span>{formatRupiah(getCartItemPrice(item))}</span>
                  <span>Qty {item.quantity}</span>
                </div>
              </div>

              <strong className="order-summary-item-subtotal">
                {formatRupiah(getCartItemSubtotal(item))}
              </strong>
            </article>
          );
        })}
      </div>

      <div className="order-summary-total">
        <span>Total</span>
        <strong>{formatRupiah(cart.total_price)}</strong>
      </div>
    </section>
  );
}

export default OrderSummary;
