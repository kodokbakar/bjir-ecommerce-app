import { Link, useLocation } from "react-router-dom";
import {
  ArrowLeft,
  Grid2X2,
  Home,
  Package,
  ShoppingBag,
  Tags,
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
    path: "/admin",
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

function isActivePath(pathname: string, path: string): boolean {
  if (path === "/admin") {
    return pathname === "/admin";
  }

  return pathname.startsWith(path);
}

function AdminLayout({ children }: AdminLayoutProps) {
  const location = useLocation();

  return (
    <div className="admin-layout">
      <aside className="admin-sidebar" aria-label="Admin navigation">
        <Link className="admin-brand" to="/admin">
          <span className="admin-brand-mark">
            <Grid2X2 className="h-5 w-5" aria-hidden="true" />
          </span>
          <span>
            <strong>Admin</strong>
            <small>Bjir Control</small>
          </span>
        </Link>

        <nav className="admin-nav">
          {ADMIN_NAV_LINKS.map((item) => {
            const isActive = isActivePath(location.pathname, item.path);

            return (
              <Link
                className={`admin-nav-link ${isActive ? "is-active" : ""}`}
                key={item.path}
                to={item.path}
                aria-current={isActive ? "page" : undefined}
              >
                <item.Icon className="h-5 w-5" aria-hidden="true" />
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>

        <Link className="admin-back-link" to="/dashboard">
          <ArrowLeft className="h-4 w-4" aria-hidden="true" />
          Customer app
        </Link>
      </aside>

      <main className="admin-content">{children}</main>
    </div>
  );
}

export default AdminLayout;
