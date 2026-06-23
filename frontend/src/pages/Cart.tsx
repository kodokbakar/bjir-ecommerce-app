import { useEffect, useMemo, useRef, useState } from "react";
import { Link } from "react-router-dom";

import OrderSummary from "../components/OrderSummary";
import {
  getCart,
  getCartErrorMessage,
  getCartItemStock,
  normalizeCart,
  removeCartItem,
  updateCartItem,
} from "../services/cartService";
import type { Cart as CartModel, CartItem } from "../types/cart";
import { formatRupiah } from "../utils/product";

const CART_SKELETON_COUNT = 3;
const NOTICE_TIMEOUT_MS = 4000;

type CartNotice = {
  type: "success" | "error";
  message: string;
} | null;

function getItemName(item: CartItem): string {
  return item.product?.name || "Unavailable product";
}

function CartSkeleton() {
  return (
    <div className="cart-items-list" aria-label="Loading cart">
      {Array.from({ length: CART_SKELETON_COUNT }, (_, index) => (
        <div className="cart-skeleton-row" key={index}>
          <div className="cart-skeleton-image" />
          <div className="cart-skeleton-copy">
            <div className="cart-skeleton-line short" />
            <div className="cart-skeleton-line" />
            <div className="cart-skeleton-line tiny" />
          </div>
        </div>
      ))}
    </div>
  );
}

function Cart() {
  const isMountedRef = useRef(false);

  const [cart, setCart] = useState<CartModel | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<CartNotice>(null);
  const [updatingItemID, setUpdatingItemID] = useState<string | null>(null);
  const [removingItemID, setRemovingItemID] = useState<string | null>(null);

  useEffect(() => {
    isMountedRef.current = true;
    let isActive = true;

    async function loadInitialCart() {
      try {
        const result = await getCart();

        if (isActive) {
          setCart(normalizeCart(result));
        }
      } catch (loadError) {
        if (isActive) {
          setCart(null);
          setError(
            getCartErrorMessage(
              loadError,
              "Failed to load your cart. Please try again.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadInitialCart();

    return () => {
      isActive = false;
      isMountedRef.current = false;
    };
  }, []);

  useEffect(() => {
    if (!notice) {
      return;
    }

    const timeoutID = window.setTimeout(() => {
      setNotice(null);
    }, NOTICE_TIMEOUT_MS);

    return () => {
      window.clearTimeout(timeoutID);
    };
  }, [notice]);

  const itemCount = useMemo(() => {
    return cart?.items.reduce((total, item) => total + item.quantity, 0) ?? 0;
  }, [cart?.items]);

  const totalPrice = cart?.total_price ?? 0;
  const hasItems = Boolean(cart?.items.length);

  async function loadCart(options: { showLoading?: boolean } = {}) {
    const shouldShowLoading = options.showLoading ?? true;

    if (shouldShowLoading) {
      setIsLoading(true);
    }

    setError(null);

    try {
      const result = await getCart();

      if (isMountedRef.current) {
        setCart(normalizeCart(result));
      }
    } catch (loadError) {
      if (isMountedRef.current) {
        setCart(null);
        setError(
          getCartErrorMessage(
            loadError,
            "Failed to load your cart. Please try again.",
          ),
        );
      }
    } finally {
      if (isMountedRef.current) {
        setIsLoading(false);
      }
    }
  }

  async function handleQuantityChange(item: CartItem, nextQuantity: number) {
    const stock = getCartItemStock(item);

    if (nextQuantity < 1) {
      return;
    }

    if (nextQuantity > stock) {
      setNotice({
        type: "error",
        message: `Only ${stock} unit${stock === 1 ? "" : "s"} available for ${getItemName(item)}.`,
      });
      return;
    }

    setUpdatingItemID(item.id);
    setNotice(null);

    try {
      const updatedItem = await updateCartItem(item.id, nextQuantity);

      setCart((currentCart) => {
        if (!currentCart) {
          return currentCart;
        }

        return normalizeCart({
          items: currentCart.items.map((currentItem) =>
            currentItem.id === updatedItem.id ? updatedItem : currentItem,
          ),
          total_price: 0,
        });
      });

      setNotice({
        type: "success",
        message: `${getItemName(updatedItem)} quantity updated.`,
      });
    } catch (updateError) {
      setNotice({
        type: "error",
        message: getCartErrorMessage(
          updateError,
          "Failed to update quantity. The item may be gone or stock may be insufficient.",
        ),
      });

      loadCart({ showLoading: false });
    } finally {
      setUpdatingItemID(null);
    }
  }

  async function handleRemoveItem(item: CartItem) {
    const isConfirmed = window.confirm(
      `Remove ${getItemName(item)} from your cart?`,
    );

    if (!isConfirmed) {
      return;
    }

    setRemovingItemID(item.id);
    setNotice(null);

    try {
      await removeCartItem(item.id);

      setCart((currentCart) => {
        if (!currentCart) {
          return currentCart;
        }

        return normalizeCart({
          items: currentCart.items.filter(
            (currentItem) => currentItem.id !== item.id,
          ),
          total_price: 0,
        });
      });

      setNotice({
        type: "success",
        message: `${getItemName(item)} removed from cart.`,
      });
    } catch (removeError) {
      setNotice({
        type: "error",
        message: getCartErrorMessage(
          removeError,
          "Failed to remove item. It may already be gone.",
        ),
      });

      loadCart({ showLoading: false });
    } finally {
      setRemovingItemID(null);
    }
  }

  return (
    <section className="cart-page" aria-labelledby="cart-title">
      <header className="cart-hero">
        <span className="products-eyebrow">Cart Counter</span>
        <h1 className="cart-title" id="cart-title">
          Your noisy basket.
        </h1>
        <p className="cart-copy">
          Review the shelf, lock the quantity, and send it to checkout without
          losing the brutal catalog rhythm.
        </p>
      </header>

      {notice && (
        <div
          className={`cart-notice is-${notice.type}`}
          role={notice.type === "error" ? "alert" : "status"}
        >
          {notice.message}
        </div>
      )}

      {isLoading ? (
        <CartSkeleton />
      ) : error ? (
        <div className="cart-state" role="alert">
          <div>
            <h2>Cart jammed.</h2>
            <p>{error}</p>
            <button
              className="cart-primary-button"
              type="button"
              onClick={() => loadCart()}
            >
              Reload cart
            </button>
          </div>
        </div>
      ) : !hasItems ? (
        <div className="cart-state">
          <div>
            <h2>Cart is empty.</h2>
            <p>
              Your basket has no products yet. Go back to the shelf and grab
              something worth checking out.
            </p>
            <Link className="cart-primary-button" to="/products">
              Belanja Sekarang
            </Link>
          </div>
        </div>
      ) : (
        <div className="cart-shell">
          <OrderSummary
            items={cart?.items ?? []}
            totalPrice={totalPrice}
            editable
            updatingItemID={updatingItemID}
            removingItemID={removingItemID}
            onQuantityChange={handleQuantityChange}
            onRemove={handleRemoveItem}
          />

          <aside className="cart-summary" aria-label="Cart summary">
            <span className="cart-summary-label">Receipt Check</span>

            <div className="cart-summary-row">
              <span>Items</span>
              <strong>{itemCount}</strong>
            </div>

            <div className="cart-summary-row">
              <span>Total</span>
              <strong>{formatRupiah(totalPrice)}</strong>
            </div>

            <Link className="cart-checkout-button" to="/checkout">
              Continue to Checkout
            </Link>

            <Link className="cart-secondary-link" to="/products">
              Keep shopping
            </Link>
          </aside>
        </div>
      )}
    </section>
  );
}

export default Cart;
