import { useEffect, useState } from "react";
import { getResponseErrorMessage } from "../../services/api";
import {
  AlertTriangle,
  CheckCircle2,
  Clock3,
  Package,
  RefreshCw,
  ShoppingBag,
  Tags,
  WalletCards,
  type LucideIcon,
} from "lucide-react";

import {
  getAdminDashboardStats,
  type AdminDashboardStats,
} from "../../services/dashboardService";
import { formatRupiah } from "../../utils/product";

interface AdminStatCardProps {
  label: string;
  value: string;
  description: string;
  Icon: LucideIcon;
  accent: "hot" | "yellow" | "blue" | "green";
}

const EMPTY_STATS: AdminDashboardStats = {
  totalOrders: 0,
  totalRevenue: 0,
  pendingOrders: 0,
  completedToday: 0,
  revenueToday: 0,
  totalProducts: 0,
  totalCategories: 0,
};

function formatNumber(value: number): string {
  return new Intl.NumberFormat("id-ID").format(value);
}

function AdminStatCard({
  label,
  value,
  description,
  Icon,
  accent,
}: AdminStatCardProps) {
  return (
    <article className={`admin-stat-card is-${accent}`}>
      <div className="admin-stat-card-icon">
        <Icon className="h-6 w-6" aria-hidden="true" />
      </div>

      <div className="admin-stat-card-copy">
        <span>{label}</span>
        <strong>{value}</strong>
        <p>{description}</p>
      </div>
    </article>
  );
}

function AdminDashboardSkeleton() {
  return (
    <div className="admin-dashboard-grid" aria-label="Loading dashboard stats">
      {Array.from({ length: 7 }, (_, index) => (
        <div className="admin-dashboard-skeleton-card" key={index}>
          <div className="admin-dashboard-skeleton-icon" />
          <div className="admin-dashboard-skeleton-copy">
            <div className="admin-dashboard-skeleton-line tiny" />
            <div className="admin-dashboard-skeleton-line" />
            <div className="admin-dashboard-skeleton-line short" />
          </div>
        </div>
      ))}
    </div>
  );
}

function AdminDashboard() {
  const [stats, setStats] = useState<AdminDashboardStats>(EMPTY_STATS);
  const [isLoading, setIsLoading] = useState(true);
  const [reloadKey, setReloadKey] = useState(0);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let isActive = true;

    async function loadStats() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await getAdminDashboardStats();

        if (isActive) {
          setStats(result);
        }
      } catch (loadError) {
        if (isActive) {
          setStats(EMPTY_STATS);
          setError(
            getResponseErrorMessage(
              loadError,
              "Admin dashboard could not be loaded.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadStats();

    return () => {
      isActive = false;
    };
  }, [reloadKey]);

  function handleRetry() {
    setReloadKey((current) => current + 1);
  }

  return (
    <section className="admin-page" aria-labelledby="admin-dashboard-title">
      <header className="admin-page-header">
        <span>Admin Dashboard</span>
        <h1 id="admin-dashboard-title">Control room.</h1>
        <p>
          Real-time overview for order volume, delivered revenue, fulfillment
          pressure, and catalog size.
        </p>
      </header>

      {error && (
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
        <AdminDashboardSkeleton />
      ) : (
        <div className="admin-dashboard-grid">
          <AdminStatCard
            label="Total Orders"
            value={formatNumber(stats.totalOrders)}
            description="All orders created by customers."
            Icon={ShoppingBag}
            accent="hot"
          />

          <AdminStatCard
            label="Total Revenue"
            value={formatRupiah(stats.totalRevenue)}
            description="Delivered-order revenue."
            Icon={WalletCards}
            accent="green"
          />

          <AdminStatCard
            label="Pending Orders"
            value={formatNumber(stats.pendingOrders)}
            description="Orders waiting for payment or review."
            Icon={Clock3}
            accent="yellow"
          />

          <AdminStatCard
            label="Completed Today"
            value={formatNumber(stats.completedToday)}
            description="Delivered orders created today."
            Icon={CheckCircle2}
            accent="blue"
          />

          <AdminStatCard
            label="Revenue Today"
            value={formatRupiah(stats.revenueToday)}
            description="Delivered revenue recorded today."
            Icon={WalletCards}
            accent="green"
          />

          <AdminStatCard
            label="Products"
            value={formatNumber(stats.totalProducts)}
            description="Products in the catalog."
            Icon={Package}
            accent="hot"
          />

          <AdminStatCard
            label="Categories"
            value={formatNumber(stats.totalCategories)}
            description="Storefront taxonomy entries."
            Icon={Tags}
            accent="yellow"
          />
        </div>
      )}
    </section>
  );
}

export default AdminDashboard;
