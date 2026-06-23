import { Link } from "react-router-dom";
import { Minus, Plus, Trash2 } from "lucide-react";

import ProductImage from "./ProductImage";
import type { CartItem } from "../types/cart";
import { formatRupiah } from "../utils/product";

export interface OrderSummaryRowData {
  id: string;
  name: string;
  quantity: number;
  unitPrice: number;
  subtotal: number;
  stock?: number;
  productPath?: string;
  imagePath?: string;
  cartItem?: CartItem;
}

interface OrderSummaryRowProps {
  row: OrderSummaryRowData;
  editable: boolean;
  isUpdating: boolean;
  isRemoving: boolean;
  onQuantityChange?: (item: CartItem, nextQuantity: number) => void;
  onRemove?: (item: CartItem) => void;
}

function OrderSummaryRow({
  row,
  editable,
  isUpdating,
  isRemoving,
  onQuantityChange,
  onRemove,
}: OrderSummaryRowProps) {
  const canEdit = editable && Boolean(row.cartItem);
  const canDecrease = canEdit && row.quantity > 1 && !isUpdating && !isRemoving;
  const canIncrease =
    canEdit &&
    typeof row.stock === "number" &&
    row.quantity < row.stock &&
    !isUpdating &&
    !isRemoving;

  return (
    <article className="order-summary-item" role="row">
      <div className="order-summary-product-cell" role="cell">
        {row.productPath ? (
          <Link
            className="order-summary-image-link"
            to={row.productPath}
            aria-label={`View ${row.name}`}
          >
            <ProductImage
              key={`${row.id}-${row.imagePath || ""}`}
              className="order-summary-image"
              src={row.imagePath || ""}
              alt={row.name}
              width={160}
              height={140}
            />
          </Link>
        ) : (
          <span className="order-summary-image-placeholder" aria-hidden="true">
            #
          </span>
        )}

        <div className="order-summary-item-main">
          {row.productPath ? (
            <Link className="order-summary-item-title" to={row.productPath}>
              {row.name}
            </Link>
          ) : (
            <span className="order-summary-item-title">{row.name}</span>
          )}

          <div className="order-summary-item-meta">
            <span>{formatRupiah(row.unitPrice)}</span>
            {typeof row.stock === "number" && <span>{row.stock} in stock</span>}
          </div>
        </div>
      </div>

      {canEdit && row.cartItem ? (
        <div role="cell">
          <div className="cart-quantity-control order-summary-quantity">
            <button
              type="button"
              onClick={() =>
                onQuantityChange?.(row.cartItem!, row.quantity - 1)
              }
              disabled={!canDecrease}
              aria-label={`Decrease ${row.name} quantity`}
            >
              <Minus className="h-4 w-4" aria-hidden="true" />
            </button>

            <span aria-live="polite">{row.quantity}</span>

            <button
              type="button"
              onClick={() =>
                onQuantityChange?.(row.cartItem!, row.quantity + 1)
              }
              disabled={!canIncrease}
              aria-label={`Increase ${row.name} quantity`}
            >
              <Plus className="h-4 w-4" aria-hidden="true" />
            </button>
          </div>
        </div>
      ) : (
        <span className="order-summary-qty" role="cell">
          Qty {row.quantity}
        </span>
      )}

      <span className="order-summary-price" role="cell">
        {formatRupiah(row.unitPrice)}
      </span>

      <strong className="order-summary-item-subtotal" role="cell">
        {formatRupiah(row.subtotal)}
      </strong>

      {editable && (
        <div role="cell">
          {row.cartItem && (
            <button
              className="cart-remove-button order-summary-remove"
              type="button"
              onClick={() => onRemove?.(row.cartItem!)}
              disabled={isUpdating || isRemoving}
            >
              <Trash2 className="h-4 w-4" aria-hidden="true" />
              {isRemoving ? "Removing..." : "Remove"}
            </button>
          )}
        </div>
      )}
    </article>
  );
}

export default OrderSummaryRow;
