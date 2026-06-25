import { useEffect, useState } from "react";
import { Link, useLocation } from "react-router-dom";
import {
  ArrowLeft,
  Grid2X2,
  Home,
  Menu,
  Package,
  ShoppingBag,
  Tags,
  X,
  type LucideIcon,
} from "lucide-react";

interface AdminLayoutProps {
  children: React.ReactNode;
}

interface AdminNavLink {
  label: string;
  path: string;
  Icon: LucideIcon;
}

const ADMIN_NAV_LINKS: AdminNavLink[] = [
  {
    label: "Dashboard",
    path: "/admin/dashboard",
    Icon: Home,
  },
  {
    label: "Products",
    path: "/admin/products",
    Icon: Package,
  },
  {
    label: "Categories",
    path: "/admin/categories",
    Icon: Tags,
  },
  {
    label: "Orders",
    path: "/admin/orders",
    Icon: ShoppingBag,
  },
];

const ADMIN_SIDEBAR_ID = "admin-sidebar";

function isActivePath(pathname: string, path: string): boolean {
  if (path === "/admin/dashboard") {
    return pathname === "/admin" || pathname === "/admin/dashboard";
  }

  return pathname.startsWith(path);
}

interface AdminBrandProps {
  onClick?: () => void;
}

function AdminBrand({ onClick }: AdminBrandProps) {
  return (
    <Link className="admin-brand" to="/admin/dashboard" onClick={onClick}>
      <span className="admin-brand-mark">
        <Grid2X2 className="h-5 w-5" aria-hidden="true" />
      </span>
      <span>
        <strong>Admin</strong>
        <small>Bjir Control</small>
      </span>
    </Link>
  );
}

function AdminLayout({ children }: AdminLayoutProps) {
  const location = useLocation();
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  function openSidebar() {
    setIsSidebarOpen(true);
  }

  function closeSidebar() {
    setIsSidebarOpen(false);
  }

  useEffect(() => {
    if (!isSidebarOpen) {
      return;
    }

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    return () => {
      document.body.style.overflow = previousOverflow;
    };
  }, [isSidebarOpen]);

  return (
    <div className={`admin-layout ${isSidebarOpen ? "is-sidebar-open" : ""}`}>
      <header className="admin-mobile-header">
        <AdminBrand onClick={closeSidebar} />

        <button
          className="admin-menu-button"
          type="button"
          aria-controls={ADMIN_SIDEBAR_ID}
          aria-expanded={isSidebarOpen}
          onClick={openSidebar}
        >
          <Menu className="h-5 w-5" aria-hidden="true" />
          Menu
        </button>
      </header>

      <button
        className="admin-sidebar-backdrop"
        type="button"
        aria-label="Close admin navigation"
        onClick={closeSidebar}
      />

      <aside
        className="admin-sidebar"
        id={ADMIN_SIDEBAR_ID}
        aria-label="Admin navigation"
      >
        <div className="admin-sidebar-head">
          <AdminBrand />

          <button
            className="admin-sidebar-close"
            type="button"
            aria-label="Close admin navigation"
            onClick={closeSidebar}
          >
            <X className="h-5 w-5" aria-hidden="true" />
          </button>
        </div>

        <nav className="admin-nav">
          {ADMIN_NAV_LINKS.map((item) => {
            const isActive = isActivePath(location.pathname, item.path);

            return (
              <Link
                className={`admin-nav-link ${isActive ? "is-active" : ""}`}
                key={item.path}
                to={item.path}
                aria-current={isActive ? "page" : undefined}
                onClick={closeSidebar}
              >
                <item.Icon className="h-5 w-5" aria-hidden="true" />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>

        <Link className="admin-back-link" to="/dashboard" onClick={closeSidebar}>
          <ArrowLeft className="h-4 w-4" aria-hidden="true" />
          Customer app
        </Link>
      </aside>

      <main className="admin-content">{children}</main>
    </div>
  );
}

export default AdminLayout;
