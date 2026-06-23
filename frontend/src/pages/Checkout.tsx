import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";

import {
  checkoutCart,
  getCart,
  getCartErrorMessage,
} from "../services/cartService";
import type { Cart } from "../types/cart";
import type { Order } from "../types/order";
import { formatRupiah } from "../utils/product";

type CheckoutState = "review" | "success";

function Checkout() {
  const [cart, setCart] = useState<Cart | null>(null);
  const [order, setOrder] = useState<Order | null>(null);
  const [state, setState] = useState<CheckoutState>("review");
  const [isLoading, setIsLoading] = useState(true);
  const [isCheckingOut, setIsCheckingOut] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const itemCount = useMemo(() => {
    return cart?.items.reduce((total, item) => total + item.quantity, 0) ?? 0;
  }, [cart?.items]);

  useEffect(() => {
    let isMounted = true;

    async function loadCart() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await getCart();

        if (isMounted) {
          setCart(result);
        }
      } catch (loadError) {
        if (isMounted) {
          setCart(null);
          setError(getCartErrorMessage(loadError, "Failed to load checkout summary."));
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadCart();

    return () => {
      isMounted = false;
    };
  }, []);

  async function handleCheckout() {
    setIsCheckingOut(true);
    setError(null);

    try {
      const result = await checkoutCart();
      setOrder(result);
      setState("success");
    } catch (checkoutError) {
      setError(
        getCartErrorMessage(
          checkoutError,
          "Checkout failed. Please review stock and try again.",
        ),
      );
    } finally {
      setIsCheckingOut(false);
    }
  }

  if (isLoading) {
    return (
      <section className="cart-page" aria-label="Loading checkout">
        <div className="cart-skeleton-row">
          <div className="cart-skeleton-copy">
            <div className="cart-skeleton-line short" />
            <div className="cart-skeleton-line" />
            <div className="cart-skeleton-line tiny" />
          </div>
        </div>
      </section>
    );
  }

  if (state === "success" && order) {
    return (
      <section className="cart-page" aria-labelledby="checkout-success-title">
        <div className="cart-state checkout-success">
          <div>
            <span className="cart-summary-label">Order Locked</span>
            <h1 id="checkout-success-title">Checkout successful.</h1>
            <p>
              Order <strong>{order.order_number}</strong> is now waiting with status{" "}
              <strong>{order.status}</strong>.
            </p>
            <p className="checkout-total">{formatRupiah(order.total_amount)}</p>
            <Link className="cart-primary-button" to="/orders">
              View orders
            </Link>
          </div>
        </div>
      </section>
    );
  }

  const hasItems = Boolean(cart?.items.length);

  return (
    <section className="cart-page" aria-labelledby="checkout-title">
      <header className="cart-hero checkout-hero">
        <span className="products-eyebrow">Checkout</span>
        <h1 className="cart-title" id="checkout-title">
          Final counter.
        </h1>
        <p className="cart-copy">
          This checkout uses your current cart and creates a pending order.
        </p>
      </header>

      {error && (
        <div className="cart-notice is-error" role="alert">
          {error}
        </div>
      )}

      {!hasItems ? (
        <div className="cart-state">
          <div>
            <h2>No items to checkout.</h2>
            <p>Add products to your cart before creating an order.</p>
            <Link className="cart-primary-button" to="/products">
              Belanja Sekarang
            </Link>
          </div>
        </div>
      ) : (
        <div className="cart-summary checkout-panel">
          <span className="cart-summary-label">Checkout Summary</span>

          <div className="cart-summary-row">
            <span>Items</span>
            <strong>{itemCount}</strong>
          </div>

          <div className="cart-summary-row">
            <span>Total</span>
            <strong>{formatRupiah(cart?.total_price ?? 0)}</strong>
          </div>

          <button
            className="cart-checkout-button"
            type="button"
            onClick={handleCheckout}
            disabled={isCheckingOut}
          >
            {isCheckingOut ? "Creating order..." : "Place order"}
          </button>

          <Link className="cart-secondary-link" to="/cart">
            Back to cart
          </Link>
        </div>
      )}
    </section>
  );
}

export default Checkout;