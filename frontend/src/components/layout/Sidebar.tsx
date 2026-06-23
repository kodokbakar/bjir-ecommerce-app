import { useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { Package } from "lucide-react";

import NavItem from "./NavItem";
import { Grid2X2, MAIN_NAV_ITEMS } from "./navigation";
import { useLayoutCategories } from "./useLayoutCategories";

interface SidebarProps {
  isCollapsed: boolean;
  isMobileOpen: boolean;
  onCloseMobile: () => void;
}

interface SidebarContentProps {
  isCollapsed: boolean;
  onNavigate?: () => void;
}

function Brand({ isCollapsed }: { isCollapsed: boolean }) {
  return (
    <Link
      className={[
        "flex h-[70px] items-center gap-3 border-b border-[var(--color-border)] px-4 text-[var(--color-text-dark)] no-underline",
        isCollapsed ? "justify-center" : "justify-start",
      ].join(" ")}
      to="/dashboard"
    >
      <span className="grid h-9 w-9 shrink-0 place-items-center rounded-xl bg-[var(--color-primary)] text-white shadow-[3px_3px_0_var(--color-brutal-ink)]">
        <Grid2X2 className="h-5 w-5" aria-hidden="true" />
      </span>

      {!isCollapsed && (
        <span className="truncate text-base font-black tracking-tight">
          Bjir-Ecommerce
        </span>
      )}
    </Link>
  );
}

function SidebarContent({ isCollapsed, onNavigate }: SidebarContentProps) {
  const location = useLocation();
  const { categories, error, isLoading } = useLayoutCategories();

  return (
    <>
      <Brand isCollapsed={isCollapsed} />

      <nav
        className={[
          "flex flex-1 flex-col gap-2 overflow-y-auto px-3 py-5",
          isCollapsed ? "items-stretch" : "",
        ].join(" ")}
        aria-label="Main navigation"
      >
        {MAIN_NAV_ITEMS.map((item) => {
          const isActive = location.pathname.startsWith(item.path);

          return (
            <NavItem
              key={item.path}
              to={item.path}
              label={item.label}
              Icon={item.Icon}
              isActive={isActive}
              isCollapsed={isCollapsed}
              onClick={onNavigate}
            />
          );
        })}

        {!isCollapsed && (
          <div className="mt-4 border-t border-[var(--color-border)] pt-4">
            <p className="mb-2 px-3 text-[11px] font-black uppercase tracking-[0.14em] text-[var(--color-text-muted)]">
              Categories
            </p>

            <Link
              className={[
                "mb-1 flex min-h-10 items-center gap-2 rounded-xl px-3 text-sm font-bold no-underline transition",
                location.pathname === "/products"
                  ? "bg-[var(--color-secondary)] text-[var(--color-primary-dark)]"
                  : "text-[var(--color-text-muted)] hover:bg-[var(--color-brutal-accent)] hover:text-[var(--color-brutal-ink)]",
              ].join(" ")}
              to="/products"
              onClick={onNavigate}
            >
              <Package className="h-4 w-4" aria-hidden="true" />
              <span>All Products</span>
            </Link>

            {isLoading && (
              <p className="px-3 py-2 text-xs font-bold text-[var(--color-text-muted)]">
                Loading categories...
              </p>
            )}

            {error && (
              <p className="px-3 py-2 text-xs font-bold text-[var(--color-stock-out)]">
                {error}
              </p>
            )}

            {!isLoading &&
              !error &&
              categories.map((category) => {
                const categoryPath = `/categories/${category.slug}`;
                const isActive = location.pathname === categoryPath;

                return (
                  <Link
                    key={category.id}
                    className={[
                      "block min-h-9 truncate rounded-xl px-3 py-2 text-sm font-bold no-underline transition",
                      isActive
                        ? "bg-[var(--color-secondary)] text-[var(--color-primary-dark)]"
                        : "text-[var(--color-text-muted)] hover:bg-[var(--color-brutal-accent)] hover:text-[var(--color-brutal-ink)]",
                    ].join(" ")}
                    to={categoryPath}
                    onClick={onNavigate}
                  >
                    {category.name}
                  </Link>
                );
              })}
          </div>
        )}
      </nav>
    </>
  );
}

function Sidebar({ isCollapsed, isMobileOpen, onCloseMobile }: SidebarProps) {
  useEffect(() => {
    if (!isMobileOpen) {
      return;
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        onCloseMobile();
      }
    }

    window.addEventListener("keydown", handleKeyDown);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isMobileOpen, onCloseMobile]);

  return (
    <>
      <aside
        className={[
          "fixed inset-y-0 left-0 z-30 hidden flex-col overflow-hidden border-r border-[var(--color-border)] bg-white transition-[width] duration-300 md:flex",
          "md:w-[76px]",
          isCollapsed ? "lg:w-[76px]" : "lg:w-[260px]",
        ].join(" ")}
      >
        <SidebarContent isCollapsed={isCollapsed} />
      </aside>

      {isMobileOpen && (
        <div className="fixed inset-0 z-50 md:hidden" role="dialog" aria-modal="true">
          <button
            className="absolute inset-0 bg-black/45"
            type="button"
            aria-label="Close navigation drawer"
            onClick={onCloseMobile}
          />

          <aside className="relative flex h-full w-[82vw] max-w-[320px] animate-[products-page-enter_220ms_ease-out] flex-col overflow-hidden border-r-4 border-[var(--color-brutal-ink)] bg-white shadow-[8px_0_0_var(--color-brutal-ink)]">
            <SidebarContent isCollapsed={false} onNavigate={onCloseMobile} />
          </aside>
        </div>
      )}
    </>
  );
}

export default Sidebar;
