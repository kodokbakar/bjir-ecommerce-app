import React from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

// ── Colour tokens dari PM ──────────────────────────────────────────────────
const C = {
    primary:       "#A67B7B",
    primaryDark:   "#8f6464",
    primaryLight:  "#b98e8e",
    secondary:     "#E7D7C9",
    accent:        "#36454F",
    textDark:      "#22303a",
    textMuted:     "#7a6e6e",
    textLabel:     "#4a3535",
    border:        "#d4bfb0",
    pillText:      "rgba(231,215,201,0.88)",
    pillBg:        "rgba(255,255,255,0.10)",
    pillBorder:    "rgba(255,255,255,0.18)",
    heroDeco1:     "#b98e8e",
    heroDeco2:     "#8f6464",
    heroDeco3:     "#c9a0a0",
} as const;

// ── Daftar Menu Navigasi ───────────────────────────────────────────────────
const NAV_MENU = [
    { name: "Dashboard", path: "/dashboard", icon: "M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" },
    { name: "Produk", path: "/products", icon: "M20 16.28A2 2 0 0 1 18 18H6a2 2 0 0 1-2-2v-8.56A2 2 0 0 1 6 5.44h12a2 2 0 0 1 2 2.28z" },
    { name: "Pesanan", path: "/orders", icon: "M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" },
    { name: "Pengaturan", path: "/settings", icon: "M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" },
];

// ── Antarmuka Props ────────────────────────────────────────────────────────
interface LayoutProps {
    children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
    const location = useLocation();
    const navigate = useNavigate();
    
    // Asumsi useAuth menyediakan fungsi logout dan data user
    const { logout, user } = useAuth(); 

    const handleLogout = () => {
        logout();
        navigate("/login");
    };

    return (
        <div style={{ display: "flex", minHeight: "100vh", background: "#f8f6f4" }}>
            
            {/* ── SIDEBAR ── */}
            <aside style={{
                width: 260,
                background: "#fff",
                borderRight: `1px solid ${C.border}`,
                display: "flex",
                flexDirection: "column",
                position: "fixed",
                top: 0,
                bottom: 0,
                left: 0,
                zIndex: 10,
            }}>
                {/* Logo Sidebar */}
                <div style={{
                    height: 70,
                    display: "flex",
                    alignItems: "center",
                    padding: "0 24px",
                    borderBottom: `1px solid ${C.border}`,
                    gap: 12
                }}>
                    <div style={{ width: 32, height: 32, background: C.primary, borderRadius: 8, display: "flex", alignItems: "center", justifyContent: "center" }}>
                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10" />
                            <path d="M12 2c0 5.5-4 10-10 10" />
                        </svg>
                    </div>
                    <span style={{ fontSize: 16, fontWeight: 600, color: C.textDark }}>Bjir-Ecommerce</span>
                </div>

                {/* Navigasi */}
                <nav style={{ flex: 1, padding: "20px 16px", display: "flex", flexDirection: "column", gap: 8 }}>
                    {NAV_MENU.map((menu) => {
                        const isActive = location.pathname.startsWith(menu.path);
                        return (
                            <Link key={menu.name} to={menu.path} style={{
                                display: "flex", alignItems: "center", gap: 12, padding: "10px 14px",
                                textDecoration: "none", borderRadius: 8,
                                background: isActive ? C.secondary : "transparent",
                                color: isActive ? C.textDark : C.textMuted,
                                fontWeight: isActive ? 500 : 400,
                                transition: "all 0.2s"
                            }}>
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                                    <path d={menu.icon} />
                                    {menu.name === "Dashboard" && <polyline points="9 22 9 12 15 12 15 22" />}
                                    {menu.name === "Pengaturan" && <circle cx="12" cy="12" r="3" />}
                                </svg>
                                <span style={{ fontSize: 14 }}>{menu.name}</span>
                            </Link>
                        );
                    })}
                </nav>
            </aside>

            {/* ── MAIN WRAPPER (Kanan) ── */}
            <div style={{
                flex: 1,
                marginLeft: 260, // Sesuai dengan lebar sidebar
                display: "flex",
                flexDirection: "column",
                minHeight: "100vh"
            }}>
                
                {/* ── HEADER ── */}
                <header style={{
                    height: 70,
                    background: "#fff",
                    borderBottom: `1px solid ${C.border}`,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "space-between",
                    padding: "0 32px",
                    position: "sticky",
                    top: 0,
                    zIndex: 9
                }}>
                    <div>
                        <h2 style={{ fontSize: 18, fontWeight: 500, color: C.textDark, margin: 0 }}>
                            {/* Menampilkan nama halaman secara dinamis (opsional) */}
                            {NAV_MENU.find(m => location.pathname.startsWith(m.path))?.name || "Halaman"}
                        </h2>
                    </div>
                    
                    <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
                        <span style={{ fontSize: 13, color: C.textMuted }}>
                            Halo, <strong style={{ color: C.textDark }}>{user?.name || "Pengguna"}</strong>
                        </span>
                        
                        {/* Tombol Logout */}
                        <button
                            onClick={handleLogout}
                            style={{
                                background: "none", border: `1px solid ${C.border}`, borderRadius: 8,
                                padding: "6px 12px", cursor: "pointer", display: "flex", alignItems: "center", gap: 6,
                                color: C.primaryDark, transition: "background 0.2s"
                            }}
                        >
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" /><polyline points="16 17 21 12 16 7" /><line x1="21" y1="12" x2="9" y2="12" />
                            </svg>
                            <span style={{ fontSize: 13, fontWeight: 500 }}>Keluar</span>
                        </button>
                    </div>
                </header>

                {/* ── KONTEN HALAMAN ── */}
                <main style={{
                    flex: 1,
                    padding: "32px",
                    display: "flex",
                    flexDirection: "column",
                }}>
                    {/* Children akan me-render komponen spesifik seperti Dashboard, Produk, dll */}
                    {children}
                </main>

                {/* ── FOOTER ── */}
                <footer style={{
                    padding: "20px 32px",
                    borderTop: `1px solid ${C.border}`,
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    color: C.textMuted,
                    fontSize: 12
                }}>
                    <span>&copy; {new Date().getFullYear()} Bjir-Ecommerce. Hak cipta dilindungi.</span>
                    <span>Versi 1.0.0</span>
                </footer>
                
            </div>
        </div>
    );
};

export default Layout;