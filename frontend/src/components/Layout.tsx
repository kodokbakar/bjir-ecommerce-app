import React, { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";

import { useAuth } from "../hooks/useAuth";
import { listCategories } from "../services/productService";
import { C } from "../styles/tokens";
import type { Category } from "../types/product";

const NAV_MENU = [
  {
    name: "Dashboard",
    path: "/dashboard",
    icon: "M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z",
  },
  {
    name: "Produk",
    path: "/products",
    icon: "M20 16.28A2 2 0 0 1 18 18H6a2 2 0 0 1-2-2v-8.56A2 2 0 0 1 6 5.44h12a2 2 0 0 1 2 2.28z",
  },
  {
    name: "Pesanan",
    path: "/orders",
    icon: "M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z",
  },
  {
    name: "Pengaturan",
    path: "/settings",
    icon: "M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z",
  },
];

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const location = useLocation();
  const navigate = useNavigate();
  const { logout, user } = useAuth();

  const [isCollapsed, setIsCollapsed] = useState(false);
  const [categories, setCategories] = useState<Category[]>([]);

  const sidebarWidth = isCollapsed ? 76 : 260;

  useEffect(() => {
    let isMounted = true;

    async function loadCategories() {
      try {
        const result = await listCategories();

        if (isMounted) {
          setCategories(result.slice(0, 6));
        }
      } catch (error) {
        console.error("Failed to load layout categories:", error);

        if (isMounted) {
          setCategories([]);
        }
      }
    }

    loadCategories();

    return () => {
      isMounted = false;
    };
  }, []);

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <div style={{ display: "flex", minHeight: "100vh", background: "#f8f6f4" }}>
      <aside
        style={{
          width: sidebarWidth,
          background: "#fff",
          borderRight: `1px solid ${C.border}`,
          display: "flex",
          flexDirection: "column",
          position: "fixed",
          top: 0,
          bottom: 0,
          left: 0,
          zIndex: 10,
          transition: "width 0.2s ease-in-out",
          overflow: "hidden",
        }}
      >
        <div
          style={{
            height: 70,
            display: "flex",
            alignItems: "center",
            padding: isCollapsed ? "0" : "0 24px",
            justifyContent: isCollapsed ? "center" : "flex-start",
            borderBottom: `1px solid ${C.border}`,
            gap: 12,
            transition: "padding 0.2s ease-in-out",
          }}
        >
          <div
            style={{
              flexShrink: 0,
              width: 32,
              height: 32,
              background: C.primary,
              borderRadius: 8,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
            }}
          >
            <svg
              width="18"
              height="18"
              viewBox="0 0 24 24"
              fill="none"
              stroke="#fff"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10" />
              <path d="M12 2c0 5.5-4 10-10 10" />
            </svg>
          </div>

          {!isCollapsed && (
            <span
              style={{
                fontSize: 16,
                fontWeight: 600,
                color: C.textDark,
                whiteSpace: "nowrap",
              }}
            >
              Bjir-Ecommerce
            </span>
          )}
        </div>

        <nav
          style={{
            flex: 1,
            padding: isCollapsed ? "20px 8px" : "20px 16px",
            display: "flex",
            flexDirection: "column",
            gap: 8,
            transition: "padding 0.2s",
            overflowY: "auto",
          }}
        >
          {NAV_MENU.map((menu) => {
            const isActive = location.pathname.startsWith(menu.path);

            return (
              <Link
                key={menu.name}
                to={menu.path}
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 12,
                  padding: "10px 14px",
                  justifyContent: isCollapsed ? "center" : "flex-start",
                  textDecoration: "none",
                  borderRadius: 8,
                  background: isActive ? C.secondary : "transparent",
                  color: isActive ? C.textDark : C.textMuted,
                  fontWeight: isActive ? 500 : 400,
                  transition: "all 0.2s",
                }}
                title={isCollapsed ? menu.name : ""}
              >
                <svg
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="1.8"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  style={{ flexShrink: 0 }}
                >
                  <path d={menu.icon} />
                  {menu.name === "Dashboard" && <polyline points="9 22 9 12 15 12 15 22" />}
                  {menu.name === "Pengaturan" && <circle cx="12" cy="12" r="3" />}
                </svg>

                {!isCollapsed && (
                  <span style={{ fontSize: 14, whiteSpace: "nowrap" }}>{menu.name}</span>
                )}
              </Link>
            );
          })}

          {!isCollapsed && (
            <div
              style={{
                marginTop: 14,
                paddingTop: 14,
                borderTop: `1px solid ${C.border}`,
              }}
            >
              <p
                style={{
                  margin: "0 0 8px",
                  color: C.textMuted,
                  fontSize: 11,
                  fontWeight: 700,
                  letterSpacing: "0.08em",
                  textTransform: "uppercase",
                }}
              >
                Categories
              </p>

              <Link
                to="/products"
                style={{
                  display: "block",
                  padding: "8px 12px",
                  borderRadius: 8,
                  color: location.pathname === "/products" ? C.primaryDark : C.textMuted,
                  background: location.pathname === "/products" ? C.secondary : "transparent",
                  fontSize: 13,
                  fontWeight: location.pathname === "/products" ? 700 : 500,
                  textDecoration: "none",
                }}
              >
                All Products
              </Link>

              {categories.map((category) => {
                const categoryPath = `/categories/${category.slug}`;
                const isActive = location.pathname === categoryPath;

                return (
                  <Link
                    key={category.id}
                    to={categoryPath}
                    style={{
                      display: "block",
                      padding: "8px 12px",
                      borderRadius: 8,
                      color: isActive ? C.primaryDark : C.textMuted,
                      background: isActive ? C.secondary : "transparent",
                      fontSize: 13,
                      fontWeight: isActive ? 700 : 500,
                      textDecoration: "none",
                      whiteSpace: "nowrap",
                      overflow: "hidden",
                      textOverflow: "ellipsis",
                    }}
                  >
                    {category.name}
                  </Link>
                );
              })}
            </div>
          )}
        </nav>
      </aside>

      <div
        style={{
          flex: 1,
          marginLeft: sidebarWidth,
          display: "flex",
          flexDirection: "column",
          minHeight: "100vh",
          transition: "margin-left 0.2s ease-in-out",
        }}
      >
        <header
          style={{
            height: 70,
            background: "#fff",
            borderBottom: `1px solid ${C.border}`,
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            padding: "0 32px",
            position: "sticky",
            top: 0,
            zIndex: 9,
          }}
        >
          <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
            <button
              onClick={() => setIsCollapsed(!isCollapsed)}
              style={{
                background: "none",
                border: "none",
                cursor: "pointer",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                padding: 8,
                borderRadius: 8,
                color: C.textDark,
                transition: "background 0.2s",
              }}
              onMouseEnter={(event) => {
                event.currentTarget.style.background = "#f0ece9";
              }}
              onMouseLeave={(event) => {
                event.currentTarget.style.background = "none";
              }}
              aria-label="Toggle Sidebar"
            >
              <svg
                width="22"
                height="22"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <line x1="3" y1="12" x2="21" y2="12" />
                <line x1="3" y1="6" x2="21" y2="6" />
                <line x1="3" y1="18" x2="21" y2="18" />
              </svg>
            </button>

            <h2 style={{ fontSize: 18, fontWeight: 500, color: C.textDark, margin: 0 }}>
              {NAV_MENU.find((item) => location.pathname.startsWith(item.path))?.name ||
                (location.pathname.startsWith("/categories") ? "Categories" : "Halaman")}
            </h2>
          </div>

          <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
            <span style={{ fontSize: 13, color: C.textMuted }}>
              Halo, <strong style={{ color: C.textDark }}>{user?.name || "Pengguna"}</strong>
            </span>

            <button
              onClick={handleLogout}
              style={{
                background: "none",
                border: `1px solid ${C.border}`,
                borderRadius: 8,
                padding: "6px 12px",
                cursor: "pointer",
                display: "flex",
                alignItems: "center",
                gap: 6,
                color: C.primaryDark,
                transition: "background 0.2s",
              }}
            >
              <svg
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                <polyline points="16 17 21 12 16 7" />
                <line x1="21" y1="12" x2="9" y2="12" />
              </svg>
              <span style={{ fontSize: 13, fontWeight: 500 }}>Keluar</span>
            </button>
          </div>
        </header>

        <main
          style={{
            flex: 1,
            padding: "32px",
            display: "flex",
            flexDirection: "column",
          }}
        >
          {children}
        </main>

        <footer
          style={{
            padding: "20px 32px",
            borderTop: `1px solid ${C.border}`,
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            color: C.textMuted,
            fontSize: 12,
          }}
        >
          <span>&copy; {new Date().getFullYear()} Bjir E-commerce. Hak cipta dilindungi.</span>
          <span>Versi 1.0.0</span>
        </footer>
      </div>
    </div>
  );
};

export { default } from "./layout/Layout";
