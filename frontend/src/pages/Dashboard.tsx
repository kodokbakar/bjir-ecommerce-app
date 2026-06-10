import React from "react";
import { useAuth } from "../context/AuthContext";

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
} as const;

// ── Data Dummy untuk Pembeli ───────────────────────────────────────────────
const BUYER_STATS = [
    { label: "Koin Bjir-Ku", value: "12.500 Poin", icon: "M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z", sub: "Gunakan saat checkout" },
    { label: "Pesanan Aktif", value: "2 Paket", icon: "M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z", sub: "Sedang dalam pengiriman" },
    { label: "Voucher Tersedia", value: "5 Voucher", icon: "M2 7a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V7z", sub: "Ada diskon ongkir 100%" },
    { label: "Keranjang Belanja", value: "3 Barang", icon: "M9 22a1 1 0 1 0 0-2 1 1 0 0 0 0 2zm7 0a1 1 0 1 0 0-2 1 1 0 0 0 0 2zM1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6", sub: "Belum dibayar" },
];

const RECOMMENDATIONS = [
    { id: 1, name: "Kemeja Flanel Premium Earth-Tone", price: "Rp 185.000", rating: "4.8", img: "👔" },
    { id: 2, name: "Sepatu Sneakers Klasik Minimalis", price: "Rp 320.000", rating: "4.9", img: "👟" },
    { id: 3, name: "Tas Ransel Kanvas Unisex", price: "Rp 145.000", rating: "4.7", img: "🎒" },
    { id: 4, name: "Topi Corduroy Vintage Brown", price: "Rp 65.000", rating: "4.6", img: "🧢" },
];

const Dashboard: React.FC = () => {
    const { user } = useAuth();

    return (
        <div style={{ animation: "buyerDashboardFadeIn 0.5s ease-out forwards" }}>
            <style>{`
                @keyframes buyerDashboardFadeIn {
                    from { opacity: 0; transform: translateY(10px); }
                    to { opacity: 1; transform: translateY(0); }
                }
            `}</style>

            {/* Banner Selamat Datang Pembeli */}
            <div style={{
                background: `linear-gradient(135deg, ${C.primary} 0%, ${C.primaryDark} 100%)`,
                borderRadius: 16,
                padding: "32px 24px",
                color: C.secondary,
                marginBottom: 28,
                position: "relative",
                overflow: "hidden"
            }}>
                <div style={{ position: "relative", zIndex: 2 }}>
                    <h1 style={{ fontSize: 24, fontWeight: 600, margin: "0 0 8px" }}>
                        Selamat Datang, {user?.name || "Pelanggan Setia"}! 
                    </h1>
                    <p style={{ fontSize: 14, color: "rgba(231,215,201,0.85)", margin: 0, maxWidth: 500 }}>
                        Mau belanja apa hari ini? Dapatkan promo gratis ongkir ekstra khusus untuk transaksi pertamamu bulan ini.
                    </p>
                </div>
                <div style={{ fontSize: 80, position: "absolute", right: 30, bottom: -10, opacity: 0.15, userSelect: "none" }}>🛍️</div>
            </div>

            {/* Grid Kartu Ringkasan Akun Pembeli */}
            <div style={{
                display: "grid",
                gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))",
                gap: 20,
                marginBottom: 32
            }}>
                {BUYER_STATS.map((stat, i) => (
                    <div key={i} style={{
                        background: "#fff",
                        padding: 20,
                        borderRadius: 12,
                        border: `1px solid ${C.border}`,
                        display: "flex",
                        flexDirection: "column",
                        gap: 12,
                        boxShadow: "0 1px 3px rgba(0,0,0,0.01)"
                    }}>
                        <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
                            <span style={{ fontSize: 13, color: C.textMuted, fontWeight: 500 }}>{stat.label}</span>
                            <span style={{ color: C.primaryDark }}>
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                    <path d={stat.icon} />
                                </svg>
                            </span>
                        </div>
                        <div>
                            <h3 style={{ fontSize: 20, fontWeight: 600, color: C.textDark, margin: "0 0 4px" }}>
                                {stat.value}
                            </h3>
                            <p style={{ fontSize: 11, color: C.textMuted, margin: 0 }}>
                                {stat.sub}
                            </p>
                        </div>
                    </div>
                ))}
            </div>

            {/* Sesi Utama: Lacak Pesanan Terakhir & Rekomendasi */}
            <div style={{ display: "flex", gap: 24, flexWrap: "wrap" }}>
                
                {/* Bagian Kiri: Status Pengiriman Terakhir */}
                <div style={{ flex: "2 1 400px", background: "#fff", border: `1px solid ${C.border}`, borderRadius: 12, padding: 20 }}>
                    <h3 style={{ margin: "0 0 16px", fontSize: 16, fontWeight: 600, color: C.textDark }}>
                        Status Pengiriman Terakhir
                    </h3>
                    <div style={{ display: "flex", alignItems: "center", background: "#faf8f6", borderRadius: 8, padding: 16, border: `1px solid ${C.border}` }}>
                        <div style={{ fontSize: 28, marginRight: 16 }}>🚚</div>
                        <div style={{ flex: 1 }}>
                            <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
                                <span style={{ fontSize: 13, fontWeight: 600, color: C.textDark }}>#BJIR-882910</span>
                                <span style={{ fontSize: 11, background: C.secondary, color: C.textLabel, padding: "2px 8px", borderRadius: 4, fontWeight: 500 }}>DIKIRIM</span>
                            </div>
                            <p style={{ fontSize: 12, color: C.textMuted, margin: 0 }}>
                                Kurir sedang menuju ke lokasi Anda (Estimasi tiba hari ini).
                            </p>
                        </div>
                    </div>
                </div>

                {/* Bagian Kanan: Voucher Cepat */}
                <div style={{ flex: "1 1 250px", background: "#fff", border: `1px solid ${C.border}`, borderRadius: 12, padding: 20 }}>
                    <h3 style={{ margin: "0 0 16px", fontSize: 16, fontWeight: 600, color: C.textDark }}>
                        Klaim Voucher Spesial
                    </h3>
                    <div style={{
                        border: `1px dashed ${C.primary}`, background: "rgba(166, 123, 123, 0.05)",
                        padding: 12, borderRadius: 8, display: "flex", justifyContent: "space-between", alignItems: "center"
                    }}>
                        <div>
                            <span style={{ fontSize: 12, fontWeight: 600, color: C.primaryDark, display: "block" }}>DISKON 20RB</span>
                            <span style={{ fontSize: 10, color: C.textMuted }}>Min. Belanja Rp 150K</span>
                        </div>
                        <button style={{
                            background: C.primary, color: "#fff", border: "none", borderRadius: 6,
                            padding: "4px 10px", fontSize: 11, fontWeight: 500, cursor: "pointer"
                        }}>
                            Klaim
                        </button>
                    </div>
                </div>

            </div>

            {/* Sesi Rekomendasi Produk */}
            <div style={{ marginTop: 32 }}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 16 }}>
                    <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600, color: C.textDark }}>
                        Rekomendasi Produk Untuk Kamu
                    </h3>
                    <span style={{ fontSize: 12, color: C.primaryDark, cursor: "pointer", fontWeight: 500 }}>Lihat Semua &rarr;</span>
                </div>
                
                <div style={{
                    display: "grid",
                    gridTemplateColumns: "repeat(auto-fit, minmax(180px, 1fr))",
                    gap: 16
                }}>
                    {RECOMMENDATIONS.map((prod) => (
                        <div key={prod.id} style={{
                            background: "#fff", border: `1px solid ${C.border}`, borderRadius: 12,
                            overflow: "hidden", cursor: "pointer", transition: "transform 0.2s"
                        }}>
                            <div style={{ height: 130, background: "#faf8f6", display: "flex", alignItems: "center", justifyContent: "center", fontSize: 48 }}>
                                {prod.img}
                            </div>
                            <div style={{ padding: 12 }}>
                                <h4 style={{ fontSize: 13, fontWeight: 500, color: C.textDark, margin: "0 0 6px", whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
                                    {prod.name}
                                </h4>
                                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                                    <span style={{ fontSize: 13, fontWeight: 600, color: C.primaryDark }}>{prod.price}</span>
                                    <span style={{ fontSize: 11, color: "#eab308", display: "flex", alignItems: "center", gap: 2 }}>
                                        ⭐ {prod.rating}
                                    </span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Dashboard;