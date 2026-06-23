import { Link } from "react-router-dom";
import { Minus, Plus, Trash2 } from "lucide-react";

import ProductImage from "./ProductImage";
import { getCartItemPrice, getCartItemSubtotal } from "../services/cartService";
import type { CartItem } from "../types/cart";
import type { OrderItem } from "../types/order";
import { formatRupiah, getProductImage } from "../utils/product";

export type OrderSummaryItem = CartItem | OrderItem;

interface OrderSummaryProps {
  items: OrderSummaryItem[];
  totalPrice: number;
  editable?: boolean;
  updatingItemID?: string | null;
  removingItemID?: string | null;
  onQuantityChange?: (item: CartItem, nextQuantity: number) => void;
  onRemove?: (item: CartItem) => void;
}

function isOrderItem(item: OrderSummaryItem): item is OrderItem {
  return "product_name" in item;
}

function isCartItem(item: OrderSummaryItem): item is CartItem {
  return !isOrderItem(item);
}

function getItemName(item: OrderSummaryItem): string {
  return isOrderItem(item)
    ? item.product_name
    : item.product?.name || "Unavailable product";
}

function getItemPrice(item: OrderSummaryItem): number {
  return isOrderItem(item) ? item.price : getCartItemPrice(item);
}

function getItemSubtotal(item: OrderSummaryItem): number {
  return isOrderItem(item) ? item.subtotal : getCartItemSubtotal(item);
}

function getCartItemStock(item: CartItem): number {
  return item.product?.stock ?? item.quantity;
}

function OrderSummary({
  items,
  totalPrice,
  editable = false,
  updatingItemID = null,
  removingItemID = null,
  onQuantityChange,
  onRemove,
}: OrderSummaryProps) {
  const itemCount = items.reduce((total, item) => total + item.quantity, 0);

  return (
    <section
      className={`order-summary ${editable ? "is-editable" : "is-readonly"}`}
      aria-labelledby="order-summary-title"
    >
      <div className="order-summary-header">
        <span className="cart-summary-label">Order Summary</span>
        <h2 id="order-summary-title">Final shelf check.</h2>
        <p>
          {itemCount} item{itemCount === 1 ? "" : "s"} in this summary.
        </p>
      </div>

      {items.length === 0 ? (
        <div className="order-summary-empty">Belum ada item</div>
      ) : (
        <>
          <div className="order-summary-table-head" aria-hidden="true">
            <span>Product</span>
            <span>Qty</span>
            <span>Price</span>
            <span>Subtotal</span>
            {editable && <span>Action</span>}
          </div>

          <div className="order-summary-list">
            {items.map((item) => {
              const cartItem = isCartItem(item) ? item : null;
              const name = getItemName(item);
              const slug = cartItem?.product?.slug || "";
              const productPath = slug ? `/products/${slug}` : "/products";
              const imagePath = cartItem?.product
                ? getProductImage(cartItem.product)
                : "";
              const isUpdating = updatingItemID === item.id;
              const isRemoving = removingItemID === item.id;
              const stock = cartItem
                ? getCartItemStock(cartItem)
                : item.quantity;
              const canEdit = editable && Boolean(cartItem);
              const canDecrease =
                canEdit && item.quantity > 1 && !isUpdating && !isRemoving;
              const canIncrease =
                canEdit && item.quantity < stock && !isUpdating && !isRemoving;

              return (
                <article className="order-summary-item" key={item.id}>
                  <div className="order-summary-product-cell">
                    {cartItem ? (
                      <Link
                        className="order-summary-image-link"
                        to={productPath}
                        aria-label={`View ${name}`}
                      >
                        <ProductImage
                          key={`${item.id}-${imagePath || ""}`}
                          className="order-summary-image"
                          src={imagePath}
                          alt={name}
                          width={160}
                          height={140}
                        />
                      </Link>
                    ) : (
                      <span
                        className="order-summary-image-placeholder"
                        aria-hidden="true"
                      >
                        #
                      </span>
                    )}

                    <div className="order-summary-item-main">
                      {cartItem ? (
                        <Link
                          className="order-summary-item-title"
                          to={productPath}
                        >
                          {name}
                        </Link>
                      ) : (
                        <span className="order-summary-item-title">{name}</span>
                      )}

                      <div className="order-summary-item-meta">
                        <span>{formatRupiah(getItemPrice(item))}</span>
                        {cartItem && <span>{stock} in stock</span>}
                      </div>
                    </div>
                  </div>

                  {canEdit && cartItem ? (
                    <div className="cart-quantity-control order-summary-quantity">
                      <button
                        type="button"
                        onClick={() =>
                          onQuantityChange?.(cartItem, item.quantity - 1)
                        }
                        disabled={!canDecrease}
                        aria-label={`Decrease ${name} quantity`}
                      >
                        <Minus className="h-4 w-4" aria-hidden="true" />
                      </button>

                      <span aria-live="polite">{item.quantity}</span>

                      <button
                        type="button"
                        onClick={() =>
                          onQuantityChange?.(cartItem, item.quantity + 1)
                        }
                        disabled={!canIncrease}
                        aria-label={`Increase ${name} quantity`}
                      >
                        <Plus className="h-4 w-4" aria-hidden="true" />
                      </button>
                    </div>
                  ) : (
                    <span className="order-summary-qty">
                      Qty {item.quantity}
                    </span>
                  )}

                  <span className="order-summary-price">
                    {formatRupiah(getItemPrice(item))}
                  </span>

                  <strong className="order-summary-item-subtotal">
                    {formatRupiah(getItemSubtotal(item))}
                  </strong>

                  {editable && cartItem && (
                    <button
                      className="cart-remove-button order-summary-remove"
                      type="button"
                      onClick={() => onRemove?.(cartItem)}
                      disabled={isUpdating || isRemoving}
                    >
                      <Trash2 className="h-4 w-4" aria-hidden="true" />
                      {isRemoving ? "Removing..." : "Remove"}
                    </button>
                  )}
                </article>
              );
            })}
          </div>
        </>
      )}

      <div className="order-summary-total">
        <span>Total</span>
        <strong>{formatRupiah(totalPrice)}</strong>
      </div>
    </section>
  );
}

export default OrderSummary;
