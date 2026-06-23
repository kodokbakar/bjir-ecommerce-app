import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import {
  AlertTriangle,
  ArrowLeft,
  CreditCard,
  MapPin,
  PackageCheck,
} from "lucide-react";

import OrderSummary from "../components/OrderSummary";
import { getOrderById, getOrderErrorMessage } from "../services/orderService";
import type { Order, OrderStatus } from "../types/order";
import { formatRupiah } from "../utils/product";
import { formatDisplayDate } from "../utils/date";

const STATUS_LABELS: Record<OrderStatus, string> = {
  pending: "Menunggu Pembayaran",
  paid: "Dibayar",
  shipped: "Dikirim",
  delivered: "Selesai",
  cancelled: "Dibatalkan",
};

const STATUS_BADGE_CLASS_NAMES: Record<OrderStatus, string> = {
  pending: "is-pending",
  paid: "is-paid",
  shipped: "is-shipped",
  delivered: "is-delivered",
  cancelled: "is-cancelled",
};

function getStatusLabel(status: OrderStatus): string {
  return STATUS_LABELS[status] ?? status;
}

function getStatusClassName(status: OrderStatus): string {
  return STATUS_BADGE_CLASS_NAMES[status] ?? "is-pending";
}

function OrderDetailSkeleton() {
  return (
    <section className="order-detail-page" aria-label="Loading order detail">
      <div className="order-detail-skeleton hero" />
      <div className="order-detail-skeleton line" />
      <div className="order-detail-skeleton panel" />
    </section>
  );
}

function OrderDetail() {
  const { id } = useParams<{ id: string }>();

  const [order, setOrder] = useState<Order | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const paymentPath = useMemo(() => {
    if (!order) {
      return "/payment";
    }

    return `/payment?order_id=${encodeURIComponent(order.id)}`;
  }, [order]);

  const hasShippingInfo = Boolean(order?.shipping_address || order?.notes);

  useEffect(() => {
    let isActive = true;

    async function loadOrder() {
      if (!id) {
        setError("Order ID tidak ditemukan.");
        setIsLoading(false);
        return;
      }

      setIsLoading(true);
      setError(null);

      try {
        const result = await getOrderById(id);

        if (isActive) {
          setOrder(result);
        }
      } catch (loadError) {
        if (isActive) {
          setOrder(null);
          setError(
            getOrderErrorMessage(
              loadError,
              "Detail pesanan belum bisa dimuat. Coba lagi sebentar.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadOrder();

    return () => {
      isActive = false;
    };
  }, [id]);

  if (isLoading) {
    return <OrderDetailSkeleton />;
  }

  if (error || !order) {
    return (
      <section
        className="order-detail-page"
        aria-labelledby="order-detail-title"
      >
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
            <h2 id="order-detail-title">Order tidak ditemukan.</h2>
            <p>{error || "Pesanan tidak tersedia."}</p>
            <Link className="cart-primary-button" to="/orders">
              Lihat Pesanan
            </Link>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="order-detail-page" aria-labelledby="order-detail-title">
      <Link className="order-detail-back" to="/orders">
        <ArrowLeft className="h-4 w-4" aria-hidden="true" />
        Back to orders
      </Link>

      <header className="order-detail-hero">
        <div className="order-detail-hero-copy">
          <span className="products-eyebrow">Order Detail</span>
          <h1 className="order-detail-title" id="order-detail-title">
            {order.order_number}
          </h1>
          <p>
            Created on <strong>{formatDisplayDate(order.created_at)}</strong>
          </p>
        </div>

        <span
          className={`orders-status-badge ${getStatusClassName(order.status)}`}
        >
          {getStatusLabel(order.status)}
        </span>
      </header>

      <div className="order-detail-shell">
        <div className="order-detail-main">
          <OrderSummary
            items={order.items ?? []}
            totalPrice={order.total_amount}
          />

          {hasShippingInfo && (
            <section
              className="order-detail-panel"
              aria-labelledby="shipping-title"
            >
              <div className="order-detail-panel-heading">
                <MapPin className="h-5 w-5" aria-hidden="true" />
                <h2 id="shipping-title">Shipping details</h2>
              </div>

              {order.shipping_address && (
                <div className="order-detail-info-block">
                  <strong>Shipping address</strong>
                  <p>{order.shipping_address}</p>
                </div>
              )}

              {order.notes && (
                <div className="order-detail-info-block">
                  <strong>Notes</strong>
                  <p>{order.notes}</p>
                </div>
              )}
            </section>
          )}
        </div>

        <aside className="order-detail-side">
          <section className="order-detail-panel">
            <div className="order-detail-panel-heading">
              <PackageCheck className="h-5 w-5" aria-hidden="true" />
              <h2>Receipt</h2>
            </div>

            <div className="order-detail-meta-row">
              <span>Order number</span>
              <strong>{order.order_number}</strong>
            </div>

            <div className="order-detail-meta-row">
              <span>Status</span>
              <strong>{getStatusLabel(order.status)}</strong>
            </div>

            <div className="order-detail-meta-row">
              <span>Total</span>
              <strong>{formatRupiah(order.total_amount)}</strong>
            </div>

            <div className="order-detail-meta-row">
              <span>Created</span>
              <strong>{formatDisplayDate(order.created_at)}</strong>
            </div>
          </section>

          <section className="order-detail-panel">
            <div className="order-detail-panel-heading">
              <CreditCard className="h-5 w-5" aria-hidden="true" />
              <h2>Payment</h2>
            </div>

            {order.status === "pending" ? (
              <>
                <p className="order-detail-payment-copy">
                  This order is waiting for payment. Continue when you are
                  ready.
                </p>
                <Link className="cart-checkout-button" to={paymentPath}>
                  Bayar Sekarang
                </Link>
              </>
            ) : (
              <p className="order-detail-payment-copy">
                Payment action is not available for this order status.
              </p>
            )}
          </section>
        </aside>
      </div>
    </section>
  );
}

export default OrderDetail;
