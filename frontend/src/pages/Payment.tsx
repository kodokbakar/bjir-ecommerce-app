import { useEffect, useState, type FormEvent } from "react";
import { Link, useSearchParams } from "react-router-dom";
import {
  AlertTriangle,
  ArrowLeft,
  Banknote,
  CheckCircle2,
  CreditCard,
  Loader2,
  ReceiptText,
  Smartphone,
  type LucideIcon,
} from "lucide-react";

import { getOrderById, getOrderErrorMessage } from "../services/orderService";
import {
  getPaymentErrorMessage,
  payOrder,
} from "../services/paymentService";
import type { Order } from "../types/order";
import type { PaymentMethod, PaymentResult } from "../types/payment";
import { formatRupiah } from "../utils/product";

interface PaymentMethodOption {
  value: PaymentMethod;
  label: string;
  description: string;
  Icon: LucideIcon;
}

const PAYMENT_METHODS: PaymentMethodOption[] = [
  {
    value: "bank_transfer",
    label: "Bank Transfer",
    description: "Mock transfer confirmation from the backend.",
    Icon: Banknote,
  },
  {
    value: "credit_card",
    label: "Credit Card",
    description: "No card fields here. The backend mocks the payment.",
    Icon: CreditCard,
  },
  {
    value: "ewallet",
    label: "E-Wallet",
    description: "Mock wallet payment for this order.",
    Icon: Smartphone,
  },
];

function formatPaymentDate(value?: string): string {
  if (!value) {
    return "Belum tersedia";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "Belum tersedia";
  }

  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function getMethodLabel(method: PaymentMethod): string {
  return PAYMENT_METHODS.find((item) => item.value === method)?.label ?? method;
}

function PaymentSkeleton() {
  return (
    <section className="payment-page" aria-label="Loading payment">
      <div className="payment-skeleton hero" />
      <div className="payment-skeleton panel" />
    </section>
  );
}

function Payment() {
  const [searchParams] = useSearchParams();
  const orderID = searchParams.get("order_id")?.trim() ?? "";

  const [order, setOrder] = useState<Order | null>(null);
  const [method, setMethod] = useState<PaymentMethod>("bank_transfer");
  const [paymentResult, setPaymentResult] = useState<PaymentResult | null>(null);
  const [isLoadingOrder, setIsLoadingOrder] = useState(Boolean(orderID));
  const [isPaying, setIsPaying] = useState(false);
  const [orderError, setOrderError] = useState<string | null>(null);
  const [paymentError, setPaymentError] = useState<string | null>(null);

  useEffect(() => {
    let isActive = true;

    async function loadOrder() {
      if (!orderID) {
        setIsLoadingOrder(false);
        setOrderError("Order ID tidak ditemukan di URL pembayaran.");
        return;
      }

      setIsLoadingOrder(true);
      setOrderError(null);

      try {
        const result = await getOrderById(orderID);

        if (isActive) {
          setOrder(result);
        }
      } catch (loadError) {
        if (isActive) {
          setOrder(null);
          setOrderError(
            getOrderErrorMessage(
              loadError,
              "Order pembayaran belum bisa dimuat. Coba lagi sebentar.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoadingOrder(false);
        }
      }
    }

    loadOrder();

    return () => {
      isActive = false;
    };
  }, [orderID]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!orderID) {
      setPaymentError("Order ID tidak ditemukan di URL pembayaran.");
      return;
    }

    setIsPaying(true);
    setPaymentError(null);
    setPaymentResult(null);

    try {
      const result = await payOrder({
        order_id: orderID,
        method,
      });

      setPaymentResult(result);
    } catch (payError) {
      setPaymentError(
        getPaymentErrorMessage(
          payError,
          "Pembayaran gagal. Order mungkin sudah dibayar atau tidak bisa dibayar.",
        ),
      );
    } finally {
      setIsPaying(false);
    }
  }

  if (isLoadingOrder) {
    return <PaymentSkeleton />;
  }

  if (orderError || !order) {
    return (
      <section className="payment-page" aria-labelledby="payment-title">
        <Link className="order-detail-back" to="/orders">
          <ArrowLeft className="h-4 w-4" aria-hidden="true" />
          Back to orders
        </Link>

        <div className="orders-empty" role="alert">
          <div>
            <AlertTriangle
              className="mx-auto mb-3 h-10 w-10"
              aria-hidden="true"
            />
            <h2 id="payment-title">Payment tidak tersedia.</h2>
            <p>{orderError || "Order pembayaran tidak ditemukan."}</p>
            <Link className="cart-primary-button" to="/orders">
              Lihat Pesanan
            </Link>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="payment-page" aria-labelledby="payment-title">
      <Link className="order-detail-back" to={`/orders/${order.id}`}>
        <ArrowLeft className="h-4 w-4" aria-hidden="true" />
        Back to order
      </Link>

      <header className="payment-hero">
        <div>
          <span className="products-eyebrow">Mock Payment</span>
          <h1 className="payment-title" id="payment-title">
            Bayar order.
          </h1>
          <p className="payment-copy">
            Choose one mocked payment method. The backend handles the fake
            provider flow.
          </p>
        </div>

        <span className="payment-hero-icon" aria-hidden="true">
          <ReceiptText className="h-10 w-10" />
        </span>
      </header>

      <div className="payment-shell">
        <form className="payment-card" onSubmit={handleSubmit}>
          <div className="payment-card-heading">
            <CreditCard className="h-5 w-5" aria-hidden="true" />
            <h2>Payment method</h2>
          </div>

          {paymentError && (
            <div className="payment-notice is-error" role="alert">
              <AlertTriangle className="h-5 w-5" aria-hidden="true" />
              <span>{paymentError}</span>
            </div>
          )}

          {paymentResult && (
            <div className="payment-result" role="status">
              <div className="payment-result-heading">
                <CheckCircle2 className="h-6 w-6" aria-hidden="true" />
                <strong>Payment success.</strong>
              </div>

              <div className="payment-result-grid">
                <span>Status</span>
                <strong>{paymentResult.status}</strong>

                <span>Transaction ID</span>
                <strong>{paymentResult.transaction_id}</strong>

                <span>Paid at</span>
                <strong>{formatPaymentDate(paymentResult.paid_at)}</strong>
              </div>

              <Link className="cart-secondary-link" to={`/orders/${order.id}`}>
                Back to order detail
              </Link>
            </div>
          )}

          <div className="payment-method-list">
            {PAYMENT_METHODS.map((item) => (
              <label
                className={[
                  "payment-method-card",
                  method === item.value ? "is-selected" : "",
                ].join(" ")}
                key={item.value}
              >
                <input
                  type="radio"
                  name="payment_method"
                  value={item.value}
                  checked={method === item.value}
                  disabled={isPaying || Boolean(paymentResult)}
                  onChange={() => setMethod(item.value)}
                />

                <span className="payment-method-icon">
                  <item.Icon className="h-5 w-5" aria-hidden="true" />
                </span>

                <span className="payment-method-copy">
                  <strong>{item.label}</strong>
                  <small>{item.description}</small>
                </span>
              </label>
            ))}
          </div>

          <button
            className="cart-checkout-button"
            type="submit"
            disabled={isPaying || Boolean(paymentResult)}
          >
            {isPaying ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
                Processing...
              </>
            ) : (
              "Bayar Sekarang"
            )}
          </button>
        </form>

        <aside className="payment-card accent">
          <div className="payment-card-heading">
            <ReceiptText className="h-5 w-5" aria-hidden="true" />
            <h2>Amount to pay</h2>
          </div>

          <div className="payment-receipt-row">
            <span>Order number</span>
            <strong>{order.order_number}</strong>
          </div>

          <div className="payment-receipt-row">
            <span>Status</span>
            <strong>{order.status}</strong>
          </div>

          <div className="payment-receipt-total">
            <span>Total amount</span>
            <strong>{formatRupiah(order.total_amount)}</strong>
          </div>

          <div className="payment-receipt-row">
            <span>Selected method</span>
            <strong>{getMethodLabel(method)}</strong>
          </div>
        </aside>
      </div>
    </section>
  );
}

export default Payment;