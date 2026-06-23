import { useEffect, useMemo, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { AlertTriangle, Clock3, PackageCheck, ShoppingBag } from "lucide-react";

import { getOrderErrorMessage, listOrders } from "../services/orderService";
import type { Order, OrderListMeta, OrderStatus } from "../types/order";
import { formatRupiah } from "../utils/product";

const ORDER_LIMIT = 8;

const EMPTY_META: OrderListMeta = {
  page: 1,
  limit: ORDER_LIMIT,
  total: 0,
  total_pages: 0,
};

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

function getPositivePage(value: string | null): number {
  const parsed = Number(value);

  if (!Number.isInteger(parsed) || parsed < 1) {
    return 1;
  }

  return parsed;
}

function formatOrderDate(value?: string): string {
  if (!value) {
    return "Tanggal belum tersedia";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "Tanggal belum tersedia";
  }

  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(date);
}

function getStatusLabel(status: OrderStatus): string {
  return STATUS_LABELS[status] ?? status;
}

function getStatusClassName(status: OrderStatus): string {
  return STATUS_BADGE_CLASS_NAMES[status] ?? "is-pending";
}

function OrdersSkeleton() {
  return (
    <div className="orders-list" aria-label="Loading orders">
      {Array.from({ length: 4 }, (_, index) => (
        <div className="orders-skeleton-card" key={index}>
          <div className="orders-skeleton-line short" />
          <div className="orders-skeleton-line" />
          <div className="orders-skeleton-line tiny" />
        </div>
      ))}
    </div>
  );
}

function MyOrders() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = getPositivePage(searchParams.get("page"));

  const [orders, setOrders] = useState<Order[]>([]);
  const [meta, setMeta] = useState<OrderListMeta>(EMPTY_META);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const hasOrders = orders.length > 0;

  const pageSummary = useMemo(() => {
    if (meta.total === 0) {
      return "0 order";
    }

    const start = (meta.page - 1) * meta.limit + 1;
    const end = Math.min(meta.page * meta.limit, meta.total);

    return `${start}-${end} dari ${meta.total} order`;
  }, [meta]);

  useEffect(() => {
    let isActive = true;

    async function loadOrders() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await listOrders({
          page,
          limit: ORDER_LIMIT,
        });

        if (isActive) {
          setOrders(result.data);
          setMeta(result.meta);
        }
      } catch (loadError) {
        if (isActive) {
          setOrders([]);
          setMeta({
            ...EMPTY_META,
            page,
          });
          setError(
            getOrderErrorMessage(
              loadError,
              "Data pesanan belum bisa dimuat. Coba lagi sebentar.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadOrders();

    return () => {
      isActive = false;
    };
  }, [page]);

  function handlePageChange(nextPage: number) {
    const safePage = Math.max(1, nextPage);

    setSearchParams({
      page: String(safePage),
    });
  }

  return (
    <section className="orders-page" aria-labelledby="orders-title">
      <header className="orders-hero">
        <span className="products-eyebrow">My Orders</span>
        <h1 className="orders-title" id="orders-title">
          Pesanan kamu.
        </h1>
        <p className="orders-copy">
          Track order number, payment state, shipping movement, and total spend
          without adding filters before they are needed.
        </p>
      </header>

      {error && (
        <div className="orders-notice" role="alert">
          <AlertTriangle className="h-5 w-5" aria-hidden="true" />
          <span>{error}</span>
        </div>
      )}

      {isLoading ? (
        <OrdersSkeleton />
      ) : !hasOrders ? (
        <div className="orders-empty">
          <div>
            <Clock3 className="mx-auto mb-3 h-10 w-10" aria-hidden="true" />
            <h2>Belum ada pesanan</h2>
            <p>
              Mulai belanja dari katalog dan order pertama kamu akan muncul di
              sini.
            </p>
            <Link className="cart-primary-button" to="/products">
              Belanja Sekarang
            </Link>
          </div>
        </div>
      ) : (
        <>
          <div className="orders-status-line">
            <span>{pageSummary}</span>
            <span>
              Page {meta.page} / {Math.max(meta.total_pages, 1)}
            </span>
          </div>

          <div className="orders-table" role="table" aria-label="My orders">
            <div className="orders-table-head" role="row">
              <span role="columnheader">Order</span>
              <span role="columnheader">Date</span>
              <span role="columnheader">Total</span>
              <span role="columnheader">Status</span>
            </div>

            <div className="orders-list" role="rowgroup">
              {orders.map((order) => (
                <Link
                  className="orders-row"
                  key={order.id}
                  to={`/orders/${order.id}`}
                  role="row"
                  aria-label={`Open order ${order.order_number}`}
                >
                  <span className="orders-order-cell" role="cell">
                    <span className="orders-icon-box">
                      <PackageCheck className="h-5 w-5" aria-hidden="true" />
                    </span>

                    <span>
                      <strong>{order.order_number}</strong>
                      <small>{order.items?.length ?? 0} item</small>
                    </span>
                  </span>

                  <span className="orders-date-cell" role="cell">
                    {formatOrderDate(order.created_at)}
                  </span>

                  <strong className="orders-total-cell" role="cell">
                    {formatRupiah(order.total_amount)}
                  </strong>

                  <span role="cell">
                    <span
                      className={`orders-status-badge ${getStatusClassName(
                        order.status,
                      )}`}
                    >
                      {getStatusLabel(order.status)}
                    </span>
                  </span>
                </Link>
              ))}
            </div>
          </div>

          <nav className="orders-pagination" aria-label="Order pagination">
            <button
              className="pagination-button"
              type="button"
              disabled={meta.page <= 1}
              onClick={() => handlePageChange(meta.page - 1)}
            >
              Previous
            </button>

            <span>
              {meta.page} / {Math.max(meta.total_pages, 1)}
            </span>

            <button
              className="pagination-button"
              type="button"
              disabled={meta.page >= meta.total_pages}
              onClick={() => handlePageChange(meta.page + 1)}
            >
              Next
            </button>
          </nav>
        </>
      )}

      <aside className="orders-footnote">
        <ShoppingBag className="h-5 w-5" aria-hidden="true" />
        <span>
          No filter yet. Add status filter later when buyers ask for it.
        </span>
      </aside>
    </section>
  );
}

export default MyOrders;
