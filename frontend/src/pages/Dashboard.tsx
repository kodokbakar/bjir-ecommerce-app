import { useEffect, useMemo, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import {
  AlertTriangle,
  Boxes,
  Clock3,
  PackageCheck,
  ShoppingBag,
  Sparkles,
  Truck,
} from "lucide-react";

import ProductImage from "../components/ProductImage";
import { useAuth } from "../hooks/useAuth";
import {
  getActiveOrderCount,
  getDashboardProducts,
  getOrderTotal,
  getRecentOrders,
  type DashboardOrder,
} from "../services/dashboardService";
import type { Product } from "../types/product";
import { formatRupiah, getProductImage } from "../utils/product";

interface DashboardState {
  products: Product[];
  productTotal: number;
  orders: DashboardOrder[];
  orderTotal: number;
}

interface CheckoutSuccessState {
  orderID: string;
  orderNumber?: string;
}

interface DashboardLocationState {
  checkoutSuccess?: CheckoutSuccessState;
}

function getCheckoutSuccessState(value: unknown): CheckoutSuccessState | null {
  if (!value || typeof value !== "object" || !("checkoutSuccess" in value)) {
    return null;
  }

  const state = value as DashboardLocationState;
  const checkoutSuccess = state.checkoutSuccess;

  if (!checkoutSuccess || typeof checkoutSuccess.orderID !== "string") {
    return null;
  }

  return checkoutSuccess;
}

function DashboardSkeleton() {
  return (
    <section className="grid gap-6" aria-label="Loading dashboard">
      <div className="h-48 border-4 border-[var(--color-brutal-ink)] bg-white shadow-[6px_6px_0_var(--color-brutal-ink)]">
        <div className="h-full animate-pulse bg-[linear-gradient(90deg,rgba(23,20,18,0.08),rgba(23,20,18,0.2),rgba(23,20,18,0.08))] bg-[length:240%_100%]" />
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {Array.from({ length: 4 }, (_, index) => (
          <div
            className="h-32 border-4 border-[var(--color-brutal-ink)] bg-white shadow-[4px_4px_0_var(--color-brutal-ink)]"
            key={index}
          />
        ))}
      </div>
    </section>
  );
}

function getOrderLabel(order: DashboardOrder): string {
  return `#${order.id.slice(0, 8).toUpperCase()}`;
}

function getStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    pending: "Menunggu",
    paid: "Dibayar",
    shipped: "Dikirim",
    delivered: "Selesai",
    canceled: "Batal",
    cancelled: "Batal",
  };

  return labels[status] ?? status;
}

function Dashboard() {
  const { user } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const [checkoutSuccess] = useState(() =>
    getCheckoutSuccessState(location.state),
  );

  const [dashboard, setDashboard] = useState<DashboardState>({
    products: [],
    productTotal: 0,
    orders: [],
    orderTotal: 0,
  });
  const [isLoading, setIsLoading] = useState(true);
  const [errorNote, setErrorNote] = useState<string | null>(null);
  const [ordersError, setOrdersError] = useState<string | null>(null);

  useEffect(() => {
    if (!checkoutSuccess) {
      return;
    }

    navigate(
      {
        pathname: location.pathname,
        search: location.search,
      },
      {
        replace: true,
        state: null,
      },
    );
  }, [checkoutSuccess, location.pathname, location.search, navigate]);

  useEffect(() => {
    let isMounted = true;

    async function loadDashboard() {
      setIsLoading(true);
      setErrorNote(null);
      setOrdersError(null);

      const [productsResult, ordersResult] = await Promise.allSettled([
        getDashboardProducts(),
        getRecentOrders(3),
      ]);

      if (!isMounted) {
        return;
      }

      const nextDashboard: DashboardState = {
        products: [],
        productTotal: 0,
        orders: [],
        orderTotal: 0,
      };

      if (productsResult.status === "fulfilled") {
        nextDashboard.products = productsResult.value.data;
        nextDashboard.productTotal = productsResult.value.meta.total;
      } else {
        console.error(
          "Failed to load dashboard products:",
          productsResult.reason,
        );
        setErrorNote("Rekomendasi produk belum bisa dimuat.");
      }

      if (ordersResult.status === "fulfilled") {
        nextDashboard.orders = ordersResult.value.data;
        nextDashboard.orderTotal = ordersResult.value.total;
      } else {
        console.error("Failed to load dashboard orders:", ordersResult.reason);
        setOrdersError(
          "Data pesanan belum bisa dimuat. Coba refresh halaman atau buka menu Orders.",
        );
      }

      setDashboard(nextDashboard);
      setIsLoading(false);
    }

    loadDashboard();

    return () => {
      isMounted = false;
    };
  }, []);

  const stats = useMemo(() => {
    const activeOrders = getActiveOrderCount(dashboard.orders);
    const lowStockCount = dashboard.products.filter((product) => {
      return product.stock > 0 && product.stock <= 5;
    }).length;

    return [
      {
        label: "Produk Katalog",
        value: `${dashboard.productTotal} item`,
        sub: "Produk aktif untuk dibeli",
        to: "/products",
        Icon: Boxes,
      },
      {
        label: "Pesanan Aktif",
        value: `${activeOrders} order`,
        sub: `${dashboard.orderTotal} total order kamu`,
        to: "/orders",
        Icon: Truck,
      },
      {
        label: "Stok Menipis",
        value: `${lowStockCount} item`,
        sub: "Dari rekomendasi terbaru",
        to: "/products?sort_by=created_at&sort_order=desc",
        Icon: AlertTriangle,
      },
      {
        label: "Rekomendasi",
        value: `${dashboard.products.length} produk`,
        sub: "Produk terbaru di katalog",
        to: "/products",
        Icon: Sparkles,
      },
    ];
  }, [dashboard]);

  if (isLoading) {
    return <DashboardSkeleton />;
  }

  const latestOrder = dashboard.orders[0];

  return (
    <section className="grid gap-7 animate-[products-page-enter_420ms_ease-out_both]">
      <header className="relative overflow-hidden border-4 border-[var(--color-brutal-ink)] bg-[linear-gradient(135deg,var(--color-primary),var(--color-primary-dark))] p-6 text-[var(--color-brutal-paper)] shadow-[6px_6px_0_var(--color-brutal-ink)] sm:p-8">
        <div className="relative z-10 max-w-2xl">
          <span className="mb-4 inline-flex border-2 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-accent)] px-3 py-1 text-xs font-black uppercase tracking-[0.14em] text-[var(--color-brutal-ink)] shadow-[3px_3px_0_var(--color-brutal-ink)]">
            Buyer Dashboard
          </span>

          <h1 className="m-0 text-4xl font-black uppercase leading-[0.9] tracking-[-0.07em] sm:text-5xl">
            Selamat datang, {user?.name || "Pelanggan"}!
          </h1>

          <p className="mt-4 max-w-xl text-sm font-bold leading-7 text-[rgba(255,248,232,0.88)]">
            Mau belanja apa hari ini? Cek rekomendasi produk terbaru dan
            lanjutkan order yang sedang berjalan.
          </p>
        </div>

        <ShoppingBag
          className="absolute -bottom-5 right-6 h-28 w-28 text-black/15"
          aria-hidden="true"
        />
      </header>

      {checkoutSuccess && (
        <div
          className="flex flex-col gap-4 border-4 border-[var(--color-brutal-ink)] bg-[#dcfce7] p-4 text-[var(--color-brutal-ink)] shadow-[5px_5px_0_var(--color-brutal-ink)] sm:flex-row sm:items-center sm:justify-between"
          role="status"
        >
          <div className="flex min-w-0 items-start gap-3">
            <span className="grid h-12 w-12 shrink-0 place-items-center border-2 border-[var(--color-brutal-ink)] bg-[var(--color-stock-in)] text-white shadow-[3px_3px_0_var(--color-brutal-ink)]">
              <PackageCheck className="h-6 w-6" aria-hidden="true" />
            </span>

            <div className="min-w-0">
              <strong className="block text-lg font-black uppercase tracking-[-0.04em]">
                Checkout successful.
              </strong>
              <p className="m-0 mt-1 text-sm font-bold text-[var(--color-text-muted)]">
                Order{" "}
                <span className="font-black text-[var(--color-brutal-ink)]">
                  {checkoutSuccess.orderNumber ||
                    `#${checkoutSuccess.orderID.slice(0, 8).toUpperCase()}`}
                </span>{" "}
                has been created. Check your order list for the latest status.
              </p>
            </div>
          </div>

          <Link
            className="inline-flex min-h-10 shrink-0 items-center justify-center border-2 border-[var(--color-brutal-ink)] bg-[var(--color-primary)] px-4 text-xs font-black uppercase tracking-[0.1em] text-white no-underline shadow-[3px_3px_0_var(--color-brutal-ink)]"
            to="/orders"
          >
            View orders
          </Link>
        </div>
      )}

      {errorNote && (
        <div className="border-2 border-[var(--color-brutal-ink)] bg-[#ffe1d7] px-4 py-3 text-sm font-bold text-[var(--color-brutal-ink)] shadow-[3px_3px_0_var(--color-brutal-ink)]">
          {errorNote}
        </div>
      )}

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {stats.map((stat) => (
          <Link
            className="group flex min-h-32 flex-col justify-between border-4 border-[var(--color-brutal-ink)] bg-white p-5 text-[var(--color-brutal-ink)] no-underline shadow-[4px_4px_0_var(--color-brutal-ink)] transition hover:-translate-x-1 hover:-translate-y-1 hover:bg-[var(--color-brutal-accent)] hover:shadow-[7px_7px_0_var(--color-brutal-ink)]"
            key={stat.label}
            to={stat.to}
          >
            <div className="flex items-center justify-between gap-3">
              <span className="text-xs font-black uppercase tracking-[0.12em] text-[var(--color-text-muted)] group-hover:text-[var(--color-brutal-ink)]">
                {stat.label}
              </span>
              <stat.Icon className="h-5 w-5" aria-hidden="true" />
            </div>

            <div>
              <strong className="block text-2xl font-black tracking-[-0.05em]">
                {stat.value}
              </strong>
              <span className="mt-1 block text-xs font-bold text-[var(--color-text-muted)] group-hover:text-[var(--color-brutal-ink)]">
                {stat.sub}
              </span>
            </div>
          </Link>
        ))}
      </div>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.4fr)_minmax(280px,0.6fr)]">
        <article className="border-4 border-[var(--color-brutal-ink)] bg-white p-5 shadow-[5px_5px_0_var(--color-brutal-ink)]">
          <div className="mb-4 flex items-center justify-between gap-3">
            <h2 className="m-0 text-xl font-black uppercase tracking-[-0.04em] text-[var(--color-brutal-ink)]">
              Status Pesanan
            </h2>
            <Link
              className="text-xs font-black uppercase tracking-[0.1em] text-[var(--color-primary-dark)] no-underline hover:underline"
              to="/orders"
            >
              Lihat semua
            </Link>
          </div>

          {ordersError ? (
            <div className="grid min-h-32 place-items-center border-2 border-[var(--color-stock-out)] bg-[#ffe1d7] p-5 text-center shadow-[3px_3px_0_var(--color-brutal-ink)]">
              <div>
                <AlertTriangle className="mx-auto mb-2 h-8 w-8 text-[var(--color-stock-out)]" />
                <p className="m-0 text-sm font-black text-[var(--color-brutal-ink)]">
                  {ordersError}
                </p>
              </div>
            </div>
          ) : latestOrder ? (
            <div className="flex flex-col gap-4 border-2 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-paper)] p-4 sm:flex-row sm:items-center">
              <span className="grid h-14 w-14 shrink-0 place-items-center border-2 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-blue)] shadow-[3px_3px_0_var(--color-brutal-ink)]">
                <PackageCheck className="h-7 w-7" aria-hidden="true" />
              </span>

              <div className="min-w-0 flex-1">
                <div className="flex flex-wrap items-center gap-2">
                  <strong className="text-sm font-black text-[var(--color-brutal-ink)]">
                    {getOrderLabel(latestOrder)}
                  </strong>
                  <span className="border-2 border-[var(--color-brutal-ink)] bg-[var(--color-secondary)] px-2 py-0.5 text-[10px] font-black uppercase tracking-[0.1em] text-[var(--color-brutal-ink)]">
                    {getStatusLabel(latestOrder.status)}
                  </span>
                </div>

                <p className="mt-2 text-sm font-bold text-[var(--color-text-muted)]">
                  {getOrderTotal(latestOrder) !== null
                    ? `Total order ${formatRupiah(getOrderTotal(latestOrder) ?? 0)}`
                    : "Order terbaru sedang diproses."}
                </p>
              </div>

              <Link
                className="inline-flex min-h-10 items-center justify-center border-2 border-[var(--color-brutal-ink)] bg-[var(--color-primary)] px-4 text-xs font-black uppercase tracking-[0.1em] text-white no-underline shadow-[3px_3px_0_var(--color-brutal-ink)]"
                to="/orders"
              >
                Detail
              </Link>
            </div>
          ) : (
            <div className="grid min-h-32 place-items-center border-2 border-dashed border-[var(--color-brutal-ink)] bg-[var(--color-brutal-paper)] p-5 text-center">
              <div>
                <Clock3 className="mx-auto mb-2 h-8 w-8 text-[var(--color-text-muted)]" />
                <p className="m-0 text-sm font-bold text-[var(--color-text-muted)]">
                  Belum ada order aktif. Mulai dari katalog produk dulu.
                </p>
              </div>
            </div>
          )}
        </article>

        <aside className="border-4 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-accent)] p-5 shadow-[5px_5px_0_var(--color-brutal-ink)]">
          <h2 className="m-0 text-xl font-black uppercase tracking-[-0.05em] text-[var(--color-brutal-ink)]">
            Deal Hunter Mode
          </h2>
          <p className="mt-3 text-sm font-bold leading-6 text-[var(--color-brutal-ink)]">
            Lompat ke katalog dan urutkan produk untuk cari barang terbaik
            sebelum stok habis.
          </p>
          <Link
            className="mt-5 inline-flex min-h-11 items-center border-2 border-[var(--color-brutal-ink)] bg-white px-4 text-xs font-black uppercase tracking-[0.1em] text-[var(--color-brutal-ink)] no-underline shadow-[3px_3px_0_var(--color-brutal-ink)]"
            to="/products?sort_by=created_at&sort_order=desc"
          >
            Buka katalog
          </Link>
        </aside>
      </div>

      <section className="grid gap-4">
        <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h2 className="m-0 text-2xl font-black uppercase tracking-[-0.05em] text-[var(--color-brutal-ink)]">
              Rekomendasi Produk
            </h2>
            <p className="m-0 mt-1 text-sm font-bold text-[var(--color-text-muted)]">
              Produk terbaru dari API katalog.
            </p>
          </div>

          <Link
            className="text-xs font-black uppercase tracking-[0.1em] text-[var(--color-primary-dark)] no-underline hover:underline"
            to="/products"
          >
            Lihat semua →
          </Link>
        </div>

        {dashboard.products.length === 0 ? (
          <div className="grid min-h-40 place-items-center border-4 border-[var(--color-brutal-ink)] bg-white p-5 text-center shadow-[5px_5px_0_var(--color-brutal-ink)]">
            <p className="m-0 text-sm font-bold text-[var(--color-text-muted)]">
              Belum ada produk untuk direkomendasikan.
            </p>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            {dashboard.products.map((product) => (
              <Link
                className="overflow-hidden border-4 border-[var(--color-brutal-ink)] bg-white text-[var(--color-brutal-ink)] no-underline shadow-[4px_4px_0_var(--color-brutal-ink)] transition hover:-translate-x-1 hover:-translate-y-1 hover:shadow-[7px_7px_0_var(--color-brutal-ink)]"
                key={product.id}
                to={`/products/${product.slug}`}
              >
                <ProductImage
                  className="h-40 w-full border-b-4 border-[var(--color-brutal-ink)] [&_.product-image-element]:object-cover"
                  src={getProductImage(product)}
                  alt={product.name}
                  width={320}
                  height={180}
                />

                <div className="grid gap-2 p-4">
                  <h3 className="m-0 truncate text-sm font-black text-[var(--color-brutal-ink)]">
                    {product.name}
                  </h3>
                  <div className="flex items-center justify-between gap-3">
                    <span className="text-sm font-black text-[var(--color-primary-dark)]">
                      {formatRupiah(product.price)}
                    </span>
                    <span className="text-xs font-bold text-[var(--color-text-muted)]">
                      {product.stock > 0 ? `${product.stock} left` : "Sold out"}
                    </span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </section>
    </section>
  );
}

export default Dashboard;