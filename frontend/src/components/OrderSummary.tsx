import OrderSummaryRow, { type OrderSummaryRowData } from "./OrderSummaryRow";
import {
  getCartItemPrice,
  getCartItemStock,
  getCartItemSubtotal,
} from "../services/cartService";
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

function normalizeOrderSummaryItem(
  item: OrderSummaryItem,
): OrderSummaryRowData {
  if (isOrderItem(item)) {
    return {
      id: item.id,
      name: item.product_name,
      quantity: item.quantity,
      unitPrice: item.price,
      subtotal: item.subtotal,
    };
  }

  const imagePath = item.product ? getProductImage(item.product) : "";
  const slug = item.product?.slug || "";

  return {
    id: item.id,
    name: item.product?.name || "Unavailable product",
    quantity: item.quantity,
    unitPrice: getCartItemPrice(item),
    subtotal: getCartItemSubtotal(item),
    stock: getCartItemStock(item),
    productPath: slug ? `/products/${slug}` : "/products",
    imagePath,
    cartItem: item,
  };
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
  const rows = items.map(normalizeOrderSummaryItem);
  const itemCount = rows.reduce((total, row) => total + row.quantity, 0);

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

      {rows.length === 0 ? (
        <div className="order-summary-empty">Belum ada item</div>
      ) : (
        <div role="table" aria-label="Order items">
          <div className="order-summary-table-head" role="row">
            <span role="columnheader">Product</span>
            <span role="columnheader">Qty</span>
            <span role="columnheader">Price</span>
            <span role="columnheader">Subtotal</span>
            {editable && <span role="columnheader">Action</span>}
          </div>

          <div className="order-summary-list" role="rowgroup">
            {rows.map((row) => (
              <OrderSummaryRow
                key={row.id}
                row={row}
                editable={editable}
                isUpdating={updatingItemID === row.id}
                isRemoving={removingItemID === row.id}
                onQuantityChange={onQuantityChange}
                onRemove={onRemove}
              />
            ))}
          </div>
        </div>
      )}

      <div className="order-summary-total">
        <span>Total</span>
        <strong>{formatRupiah(totalPrice)}</strong>
      </div>
    </section>
  );
}

export default OrderSummary;
