import { useEffect, useMemo, useState, type FormEvent } from "react";
import { useSearchParams } from "react-router-dom";
import {
  AlertTriangle,
  ArrowRight,
  Edit3,
  PackageCheck,
  RefreshCw,
  Save,
  Search,
  X,
} from "lucide-react";

import {
  getOrderErrorMessage,
  listAdminOrders,
  updateOrderStatus,
} from "../../services/orderService";
import type {
  Order,
  OrderListMeta,
  OrderListParams,
  OrderStatus,
} from "../../types/order";
import { formatDisplayDate } from "../../utils/date";
import { formatRupiah } from "../../utils/product";
import EmptyState from "../../components/EmptyState";
import { useToast } from "../../context/toast";

const ADMIN_ORDER_LIMIT = 10;

const EMPTY_META: OrderListMeta = {
  page: 1,
  limit: ADMIN_ORDER_LIMIT,
  total: 0,
  total_pages: 0,
};

const ORDER_STATUS_OPTIONS: Array<{
  label: string;
  value: OrderStatus | "";
}> = [
  { label: "All statuses", value: "" },
  { label: "Pending", value: "pending" },
  { label: "Paid", value: "paid" },
  { label: "Shipped", value: "shipped" },
  { label: "Delivered", value: "delivered" },
  { label: "Cancelled", value: "cancelled" },
];

const ORDER_STATUS_CHANGE_OPTIONS: Array<{
  label: string;
  value: OrderStatus;
}> = [
  { label: "Pending", value: "pending" },
  { label: "Paid", value: "paid" },
  { label: "Shipped", value: "shipped" },
  { label: "Delivered", value: "delivered" },
  { label: "Cancelled", value: "cancelled" },
];

const STATUS_LABELS: Record<OrderStatus, string> = {
  pending: "Pending",
  paid: "Paid",
  shipped: "Shipped",
  delivered: "Delivered",
  cancelled: "Cancelled",
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

function getStatusParam(value: string | null): OrderStatus | "" {
  const status = value?.trim() ?? "";

  if (
    status === "pending" ||
    status === "paid" ||
    status === "shipped" ||
    status === "delivered" ||
    status === "cancelled"
  ) {
    return status;
  }

  return "";
}

function getStatusLabel(status: OrderStatus): string {
  return STATUS_LABELS[status] ?? status;
}

function getStatusClassName(status: OrderStatus): string {
  return STATUS_BADGE_CLASS_NAMES[status] ?? "is-pending";
}

function getCustomerLabel(order: Order): string {
  return order.user_name || order.user_email || order.user_id;
}

function getAllowedStatusTransitions(status: OrderStatus): OrderStatus[] {
  switch (status) {
    case "pending":
      return ["paid", "cancelled"];
    case "paid":
      return ["shipped", "cancelled"];
    case "shipped":
      return ["delivered"];
    default:
      return [];
  }
}

function getDefaultTargetStatus(order: Order): OrderStatus {
  return getAllowedStatusTransitions(order.status)[0] ?? order.status;
}

function mergeUpdatedOrder(currentOrder: Order, updatedOrder: Order): Order {
  return {
    ...currentOrder,
    ...updatedOrder,
    user_name: updatedOrder.user_name ?? currentOrder.user_name,
    user_email: updatedOrder.user_email ?? currentOrder.user_email,
  };
}

function AdminOrdersSkeleton() {
  return (
    <div className="admin-orders-list" aria-label="Loading admin orders">
      {Array.from({ length: 5 }, (_, index) => (
        <div className="admin-orders-skeleton-row" key={index}>
          <div className="admin-orders-skeleton-line short" />
          <div className="admin-orders-skeleton-line" />
          <div className="admin-orders-skeleton-line tiny" />
        </div>
      ))}
    </div>
  );
}

function AdminOrders() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = getPositivePage(searchParams.get("page"));
  const search = searchParams.get("search")?.trim() ?? "";
  const status = getStatusParam(searchParams.get("status"));

  const [searchInput, setSearchInput] = useState(search);
  const [orders, setOrders] = useState<Order[]>([]);
  const [meta, setMeta] = useState<OrderListMeta>(EMPTY_META);
  const [reloadKey, setReloadKey] = useState(0);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const [targetStatus, setTargetStatus] = useState<OrderStatus>("paid");
  const [isLoading, setIsLoading] = useState(true);
  const [isUpdatingStatus, setIsUpdatingStatus] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { showToast } = useToast();

  const hasOrders = orders.length > 0;
  const allowedTargetStatuses = selectedOrder
    ? getAllowedStatusTransitions(selectedOrder.status)
    : [];
  const canSubmitStatusChange =
    selectedOrder !== null &&
    allowedTargetStatuses.includes(targetStatus) &&
    targetStatus !== selectedOrder.status &&
    !isUpdatingStatus;

  const query = useMemo<OrderListParams>(
    () => ({
      page,
      limit: ADMIN_ORDER_LIMIT,
      status,
      search: search || undefined,
    }),
    [page, search, status],
  );

  const pageSummary = useMemo(() => {
    if (meta.total === 0) {
      return "0 order";
    }

    const start = (meta.page - 1) * meta.limit + 1;
    const end = Math.min(meta.page * meta.limit, meta.total);

    return `${start}-${end} of ${meta.total} orders`;
  }, [meta]);

  useEffect(() => {
    let isActive = true;

    async function loadOrders() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await listAdminOrders(query);

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
              "Admin order list could not be loaded.",
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
  }, [page, query, reloadKey]);

  function updateParams(
    nextPage: number,
    nextSearch: string,
    nextStatus: string,
  ) {
    const params = new URLSearchParams();

    if (nextPage > 1) {
      params.set("page", String(nextPage));
    }

    if (nextSearch.trim()) {
      params.set("search", nextSearch.trim());
    }

    if (nextStatus.trim()) {
      params.set("status", nextStatus.trim());
    }

    setSearchParams(params);
  }

  function handleSearchSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    updateParams(1, searchInput, status);
  }

  function handleClearSearch() {
    setSearchInput("");
    updateParams(1, "", status);
  }

  function handleStatusChange(nextStatus: string) {
    updateParams(1, search, nextStatus);
  }

  function handlePageChange(nextPage: number) {
    updateParams(Math.max(1, nextPage), search, status);
  }

  function handleRetry() {
    setReloadKey((current) => current + 1);
  }

  function openStatusDialog(order: Order) {
    setSelectedOrder(order);
    setTargetStatus(getDefaultTargetStatus(order));
    setError(null);
  }

  function closeStatusDialog() {
    if (isUpdatingStatus) {
      return;
    }

    setSelectedOrder(null);
    setTargetStatus("paid");
  }

  async function handleConfirmStatusChange() {
    if (!selectedOrder || !canSubmitStatusChange) {
      return;
    }

    const currentOrder = selectedOrder;
    const nextStatus = targetStatus;

    setIsUpdatingStatus(true);
    setError(null);

    try {
      const updatedOrder = await updateOrderStatus(currentOrder.id, nextStatus);

      setOrders((currentOrders) =>
        currentOrders.map((order) =>
          order.id === currentOrder.id
            ? mergeUpdatedOrder(order, updatedOrder)
            : order,
        ),
      );

      showToast({
        type: "success",
        message: `${currentOrder.order_number} changed from ${getStatusLabel(
          currentOrder.status,
        )} to ${getStatusLabel(nextStatus)}.`,
      });
      setSelectedOrder(null);
      setTargetStatus("paid");
    } catch (updateError) {
      showToast(
        {
          type: "error",
          message: getOrderErrorMessage(
            updateError,
            "Order status could not be updated. Check the transition and try again.",
          ),
        },
        { duration: 6000 },
      );
    } finally {
      setIsUpdatingStatus(false);
    }
  }

  return (
    <section className="admin-page" aria-labelledby="admin-orders-title">
      <header className="admin-page-header">
        <span>Admin Orders</span>
        <h1 id="admin-orders-title">Orders.</h1>
        <p>
          Review every customer order, filter fulfillment status, and search by
          order number.
        </p>
      </header>

      <div className="admin-products-toolbar">
        <form className="admin-products-search" onSubmit={handleSearchSubmit}>
          <label htmlFor="admin-order-search">Search orders</label>

          <div className="admin-products-search-row">
            <Search className="h-5 w-5" aria-hidden="true" />
            <input
              key={search}
              id="admin-order-search"
              type="search"
              placeholder="Search order number..."
              defaultValue={search}
              onChange={(event) => setSearchInput(event.target.value)}
            />
            {search && (
              <button type="button" onClick={handleClearSearch}>
                Clear
              </button>
            )}
            <button type="submit">Search</button>
          </div>
        </form>

        <div className="admin-orders-filter">
          <label htmlFor="admin-order-status">Status</label>
          <select
            id="admin-order-status"
            value={status}
            onChange={(event) => handleStatusChange(event.target.value)}
          >
            {ORDER_STATUS_OPTIONS.map((option) => (
              <option key={option.label} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      {error && hasOrders && (
        <div className="admin-products-notice is-error" role="alert">
          <AlertTriangle className="h-5 w-5" aria-hidden="true" />
          <span>{error}</span>
          <button type="button" onClick={handleRetry}>
            <RefreshCw className="h-4 w-4" aria-hidden="true" />
            Retry
          </button>
        </div>
      )}

      {isLoading ? (
        <AdminOrdersSkeleton />
      ) : error && !hasOrders ? (
        <EmptyState
          tone="error"
          eyebrow="Order Error"
          title="Order list jammed."
          description={error}
          action={
            <button
              className="admin-products-create-button"
              type="button"
              onClick={handleRetry}
            >
              <RefreshCw className="h-5 w-5" aria-hidden="true" />
              Retry
            </button>
          }
        />
      ) : !hasOrders ? (
        <EmptyState
          eyebrow="Order Ledger"
          title="No orders found."
          description="Try another order number, clear the search, or switch the status filter."
        />
      ) : (
        <>
          <div className="admin-products-status-line">
            <span>{pageSummary}</span>
            <span>
              Page {meta.page} / {Math.max(meta.total_pages, 1)}
            </span>
          </div>

          <div className="admin-orders-table" aria-label="Admin order list">
            <div className="admin-orders-table-head" aria-hidden="true">
              <span>Order</span>
              <span>Customer</span>
              <span>Total</span>
              <span>Status</span>
              <span>Date</span>
              <span>Action</span>
            </div>

            <div className="admin-orders-list">
              {orders.map((order) => (
                <article className="admin-orders-row" key={order.id}>
                  <div className="admin-orders-order-cell">
                    <span className="orders-icon-box">
                      <PackageCheck className="h-5 w-5" aria-hidden="true" />
                    </span>

                    <span>
                      <strong>{order.order_number}</strong>
                      <small>{order.id}</small>
                    </span>
                  </div>

                  <div className="admin-orders-customer-cell">
                    <strong>{getCustomerLabel(order)}</strong>
                    {order.user_email && <small>{order.user_email}</small>}
                  </div>

                  <strong className="admin-products-price">
                    {formatRupiah(order.total_amount)}
                  </strong>

                  <span>
                    <span
                      className={`orders-status-badge ${getStatusClassName(
                        order.status,
                      )}`}
                    >
                      {getStatusLabel(order.status)}
                    </span>
                  </span>

                  <span className="admin-products-muted">
                    {formatDisplayDate(order.created_at, "short-date")}
                  </span>

                  <div className="admin-products-actions">
                    <button
                      type="button"
                      onClick={() => openStatusDialog(order)}
                    >
                      <Edit3 className="h-4 w-4" aria-hidden="true" />
                      Status
                    </button>
                  </div>
                </article>
              ))}
            </div>
          </div>

          <nav
            className="admin-products-pagination"
            aria-label="Admin order pagination"
          >
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

      {selectedOrder && (
        <div className="admin-order-status-modal" role="presentation">
          <button
            className="admin-order-status-backdrop"
            type="button"
            aria-label="Close order status dialog"
            onClick={closeStatusDialog}
          />

          <section
            className="admin-order-status-dialog"
            role="dialog"
            aria-modal="true"
            aria-labelledby="admin-order-status-title"
          >
            <header className="admin-order-status-header">
              <div>
                <span>Change Order Status</span>
                <h2 id="admin-order-status-title">Confirm status.</h2>
              </div>

              <button
                type="button"
                onClick={closeStatusDialog}
                disabled={isUpdatingStatus}
              >
                <X className="h-5 w-5" aria-hidden="true" />
                <span className="sr-only">Close</span>
              </button>
            </header>

            <div className="admin-order-status-body">
              <div className="admin-order-status-route">
                <span
                  className={`orders-status-badge ${getStatusClassName(
                    selectedOrder.status,
                  )}`}
                >
                  {getStatusLabel(selectedOrder.status)}
                </span>

                <ArrowRight className="h-5 w-5" aria-hidden="true" />

                <span
                  className={`orders-status-badge ${getStatusClassName(
                    targetStatus,
                  )}`}
                >
                  {getStatusLabel(targetStatus)}
                </span>
              </div>

              <p>
                Change order <strong>{selectedOrder.order_number}</strong> from{" "}
                <strong>{getStatusLabel(selectedOrder.status)}</strong> to{" "}
                <strong>{getStatusLabel(targetStatus)}</strong>?
              </p>

              {allowedTargetStatuses.length > 0 ? (
                <div className="admin-product-field">
                  <label htmlFor="order-next-status">Next status</label>
                  <select
                    id="order-next-status"
                    value={targetStatus}
                    onChange={(event) =>
                      setTargetStatus(event.target.value as OrderStatus)
                    }
                  >
                    {ORDER_STATUS_CHANGE_OPTIONS.filter((option) =>
                      allowedTargetStatuses.includes(option.value),
                    ).map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </div>
              ) : (
                <div className="admin-order-status-terminal" role="status">
                  This order is already in a terminal state. No further status
                  transition is available.
                </div>
              )}
            </div>

            <footer className="admin-order-status-actions">
              <button
                type="button"
                onClick={closeStatusDialog}
                disabled={isUpdatingStatus}
              >
                Cancel
              </button>
              <button
                type="button"
                disabled={!canSubmitStatusChange}
                onClick={() => void handleConfirmStatusChange()}
              >
                <Save className="h-4 w-4" aria-hidden="true" />
                {isUpdatingStatus ? "Updating..." : "Confirm"}
              </button>
            </footer>
          </section>
        </div>
      )}
    </section>
  );
}

export default AdminOrders;
