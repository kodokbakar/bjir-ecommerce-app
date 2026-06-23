import {
  useCallback,
  useEffect,
  useRef,
  useState,
  type ChangeEvent,
  type FormEvent,
} from "react";
import { Link, useNavigate } from "react-router-dom";

import OrderSummary from "../components/OrderSummary";
import {
  checkoutCart,
  getCart,
  getCartErrorMessage,
  normalizeCart,
} from "../services/cartService";
import type { Cart } from "../types/cart";
import type { CheckoutInput } from "../types/order";

const EMPTY_CHECKOUT_FORM: Required<CheckoutInput> = {
  shipping_address: "",
  notes: "",
};

function Checkout() {
  const navigate = useNavigate();
  const isMountedRef = useRef(false);

  const [cart, setCart] = useState<Cart | null>(null);
  const [form, setForm] = useState(EMPTY_CHECKOUT_FORM);
  const [isLoading, setIsLoading] = useState(true);
  const [isCheckingOut, setIsCheckingOut] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const hasItems = Boolean(cart?.items.length);

  const loadCheckoutCart = useCallback(async () => {
    setIsLoading(true);
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
          getCartErrorMessage(loadError, "Failed to load checkout summary."),
        );
      }
    } finally {
      if (isMountedRef.current) {
        setIsLoading(false);
      }
    }
  }, []);

  useEffect(() => {
    isMountedRef.current = true;

    const loadTimerID = window.setTimeout(() => {
      void loadCheckoutCart();
    }, 0);

    return () => {
      window.clearTimeout(loadTimerID);
      isMountedRef.current = false;
    };
  }, [loadCheckoutCart]);

  function handleFieldChange(
    event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) {
    const { name, value } = event.target;

    setForm((currentForm) => ({
      ...currentForm,
      [name]: value,
    }));
  }

  async function handleCheckout(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!cart || cart.items.length === 0) {
      setError("Keranjang kosong. Tambahkan produk sebelum checkout.");
      return;
    }

    setIsCheckingOut(true);
    setError(null);

    try {
      const order = await checkoutCart({
        shipping_address: form.shipping_address,
        notes: form.notes,
      });

      navigate("/dashboard", {
        replace: true,
        state: {
          checkoutSuccess: {
            orderID: order.id,
            orderNumber: order.order_number,
          },
        },
      });
    } catch (checkoutError) {
      setError(
        getCartErrorMessage(
          checkoutError,
          "Checkout failed. Keranjang kosong, stok tidak mencukupi, atau produk tidak ditemukan.",
        ),
      );
    } finally {
      if (isMountedRef.current) {
        setIsCheckingOut(false);
      }
    }
  }

  if (isLoading) {
    return (
      <section className="cart-page" aria-label="Loading checkout">
        <div className="checkout-skeleton-shell">
          <div className="cart-skeleton-row">
            <div className="cart-skeleton-copy">
              <div className="cart-skeleton-line short" />
              <div className="cart-skeleton-line" />
              <div className="cart-skeleton-line tiny" />
            </div>
          </div>

          <div className="cart-skeleton-row">
            <div className="cart-skeleton-copy">
              <div className="cart-skeleton-line short" />
              <div className="cart-skeleton-line" />
              <div className="cart-skeleton-line tiny" />
            </div>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="cart-page" aria-labelledby="checkout-title">
      <header className="cart-hero checkout-hero">
        <span className="products-eyebrow">Checkout</span>
        <h1 className="cart-title" id="checkout-title">
          Final counter.
        </h1>
        <p className="cart-copy">
          Review the order, add delivery notes, then convert the cart into a
          pending order.
        </p>
      </header>

      {error && (
        <div className="cart-notice is-error" role="alert">
          {error}
        </div>
      )}

      {!hasItems || !cart ? (
        <div className="cart-state">
          <div>
            <h2>Keranjang kosong.</h2>
            <p>Add products to your cart before creating an order.</p>
            <Link className="cart-primary-button" to="/products">
              Belanja Sekarang
            </Link>
          </div>
        </div>
      ) : (
        <div className="checkout-shell">
          <form className="checkout-form" onSubmit={handleCheckout}>
            <div className="checkout-form-header">
              <span className="cart-summary-label">Shipping Details</span>
              <h2>Where should this go?</h2>
              <p>Add delivery details and notes before placing your order.</p>
            </div>

            <label className="checkout-field">
              <span>Shipping address</span>
              <input
                name="shipping_address"
                value={form.shipping_address}
                onChange={handleFieldChange}
                placeholder="Street, city, postal code"
                autoComplete="shipping street-address"
                disabled={isCheckingOut}
              />
            </label>

            <label className="checkout-field">
              <span>Notes</span>
              <textarea
                name="notes"
                value={form.notes}
                onChange={handleFieldChange}
                placeholder="Optional delivery note, size preference, or reminder"
                rows={5}
                disabled={isCheckingOut}
              />
            </label>

            <button
              className="cart-checkout-button"
              type="submit"
              disabled={isCheckingOut}
            >
              {isCheckingOut ? "Creating order..." : "Place order"}
            </button>

            <Link className="cart-secondary-link" to="/cart">
              Back to cart
            </Link>
          </form>

          <OrderSummary
            items={cart?.items ?? []}
            totalPrice={cart?.total_price ?? 0}
          />
        </div>
      )}
    </section>
  );
}

export default Checkout;
