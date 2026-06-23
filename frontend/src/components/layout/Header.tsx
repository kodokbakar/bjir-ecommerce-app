import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import {
  ChevronDown,
  LogOut,
  Menu,
  PanelLeftClose,
  PanelLeftOpen,
  ShoppingCart,
  User,
} from "lucide-react";

import CartBadge from "../CartBadge";
import { useAuth } from "../../hooks/useAuth";
import { getBreadcrumbs, getPageTitle } from "./navigation";

interface HeaderProps {
  isCollapsed: boolean;
  cartCount: number;
  onToggleCollapse: () => void;
  onOpenMobile: () => void;
}

function getInitials(name?: string): string {
  const fallback = "P";
  const value = name?.trim();

  if (!value) {
    return fallback;
  }

  return (
    value
      .split(/\s+/)
      .slice(0, 2)
      .map((part) => part[0]?.toUpperCase() ?? "")
      .join("") || fallback
  );
}

function Header({
  isCollapsed,
  cartCount,
  onToggleCollapse,
  onOpenMobile,
}: HeaderProps) {
  const location = useLocation();
  const navigate = useNavigate();
  const { logout, user } = useAuth();

  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  const title = getPageTitle(location.pathname);
  const breadcrumbs = getBreadcrumbs(location.pathname);

  useEffect(() => {
    if (!isUserMenuOpen) {
      return;
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setIsUserMenuOpen(false);
      }
    }

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isUserMenuOpen]);

  function handleLogout() {
    logout();
    setIsUserMenuOpen(false);
    navigate("/login");
  }

  return (
    <header className="sticky top-0 z-20 flex h-[70px] items-center justify-between border-b border-[var(--color-border)] bg-white px-4 sm:px-6 lg:px-8">
      <div className="flex min-w-0 items-center gap-3">
        <button
          className="grid h-10 w-10 place-items-center rounded-xl text-[var(--color-text-dark)] transition hover:bg-[var(--color-secondary)] md:hidden"
          type="button"
          onClick={onOpenMobile}
          aria-label="Open navigation menu"
        >
          <Menu className="h-5 w-5" aria-hidden="true" />
        </button>

        <button
          className="hidden h-10 w-10 place-items-center rounded-xl text-[var(--color-text-dark)] transition hover:bg-[var(--color-secondary)] md:grid"
          type="button"
          onClick={onToggleCollapse}
          aria-label="Toggle sidebar"
          aria-expanded={!isCollapsed}
        >
          {isCollapsed ? (
            <PanelLeftOpen className="h-5 w-5" aria-hidden="true" />
          ) : (
            <PanelLeftClose className="h-5 w-5" aria-hidden="true" />
          )}
        </button>

        <div className="min-w-0">
          <h2 className="m-0 truncate text-lg font-extrabold text-[var(--color-text-dark)]">
            {title}
          </h2>

          <nav
            className="hidden items-center gap-2 text-xs font-bold text-[var(--color-text-muted)] sm:flex"
            aria-label="Breadcrumb"
          >
            {breadcrumbs.map((item, index) => (
              <span
                className="flex items-center gap-2"
                key={`${item.label}-${index}`}
              >
                {index > 0 && <span aria-hidden="true">/</span>}

                {item.path ? (
                  <Link
                    className="text-[var(--color-text-muted)] no-underline hover:text-[var(--color-primary-dark)]"
                    to={item.path}
                  >
                    {item.label}
                  </Link>
                ) : (
                  <span aria-current="page">{item.label}</span>
                )}
              </span>
            ))}
          </nav>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <Link
          className="relative grid h-10 w-10 place-items-center rounded-xl border-2 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-surface)] text-[var(--color-text-dark)] shadow-[3px_3px_0_var(--color-brutal-ink)] transition hover:-translate-x-0.5 hover:-translate-y-0.5 hover:bg-[var(--color-brutal-accent)] hover:shadow-[4px_4px_0_var(--color-brutal-ink)]"
          to="/cart"
          aria-label={
            cartCount > 0
              ? `Open cart, ${cartCount} item${cartCount === 1 ? "" : "s"}`
              : "Open cart"
          }
        >
          <ShoppingCart className="h-5 w-5" aria-hidden="true" />
          <CartBadge count={cartCount} />
        </Link>

        <div className="relative">
          <button
            className="flex items-center gap-3 rounded-2xl border-2 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-surface)] px-2 py-1.5 text-left shadow-[3px_3px_0_var(--color-brutal-ink)] transition hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-[4px_4px_0_var(--color-brutal-ink)]"
            type="button"
            onClick={() => setIsUserMenuOpen((current) => !current)}
            aria-haspopup="menu"
            aria-expanded={isUserMenuOpen}
          >
            <span className="grid h-8 w-8 place-items-center rounded-xl bg-[var(--color-primary)] text-xs font-black text-white">
              {getInitials(user?.name)}
            </span>

            <span className="hidden min-w-0 sm:block">
              <span className="block max-w-36 truncate text-xs font-bold text-[var(--color-text-muted)]">
                Halo,
              </span>
              <span className="block max-w-36 truncate text-sm font-black text-[var(--color-text-dark)]">
                {user?.name || "Pengguna"}
              </span>
            </span>

            <ChevronDown className="hidden h-4 w-4 text-[var(--color-text-muted)] sm:block" />
          </button>

          {isUserMenuOpen && (
            <div
              className="absolute right-0 mt-3 w-56 border-2 border-[var(--color-brutal-ink)] bg-white p-2 shadow-[5px_5px_0_var(--color-brutal-ink)]"
              role="menu"
            >
              <div className="border-b border-[var(--color-border)] px-3 py-2">
                <p className="m-0 truncate text-sm font-black text-[var(--color-text-dark)]">
                  {user?.name || "Pengguna"}
                </p>
                <p className="m-0 truncate text-xs font-bold text-[var(--color-text-muted)]">
                  {user?.email || "buyer@example.test"}
                </p>
              </div>

              <Link
                className="mt-2 flex min-h-10 items-center gap-2 rounded-xl px-3 text-sm font-bold text-[var(--color-text-dark)] no-underline hover:bg-[var(--color-secondary)]"
                to="/profile"
                role="menuitem"
                onClick={() => setIsUserMenuOpen(false)}
              >
                <User className="h-4 w-4" aria-hidden="true" />
                Profile
              </Link>

              <button
                className="flex min-h-10 w-full items-center gap-2 rounded-xl px-3 text-left text-sm font-black text-[var(--color-stock-out)] hover:bg-[#ffe1d7]"
                type="button"
                role="menuitem"
                onClick={handleLogout}
              >
                <LogOut className="h-4 w-4" aria-hidden="true" />
                Logout
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}

export default Header;
