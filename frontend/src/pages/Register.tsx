import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import api from "../services/api";

// ── Colour tokens (sama dengan Login) ─────────────────────────────────────
const C = {
    primary:       "#A67B7B",
    primaryDark:   "#8f6464",
    primaryLight: "#b98e8e",
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

// ── Sub-komponen: LeftPanel ────────────────────────────────────────────────
const LeftPanel: React.FC = () => (
    <div style={{
        flex: 1,
        background: C.primary,
        display: "flex",
        flexDirection: "column",
        justifyContent: "space-between",
        padding: "40px 36px",
        position: "relative",
        overflow: "hidden",
    }}>
        <span style={{ position:"absolute", top:-50, right:-50, width:180, height:180, borderRadius:"50%", background:C.heroDeco1, opacity:0.4 }} />
        <span style={{ position:"absolute", bottom:-40, left:-40, width:140, height:140, borderRadius:"50%", background:C.heroDeco2, opacity:0.35 }} />
        <span style={{ position:"absolute", bottom:100, right:24, width:80, height:80, borderRadius:"50%", background:C.heroDeco3, opacity:0.3 }} />

        {/* Logo */}
        <div style={{ display:"flex", alignItems:"center", gap:10, position:"relative", zIndex:1 }}>
            <div style={{ width:38, height:38, background:C.secondary, borderRadius:8, display:"flex", alignItems:"center", justifyContent:"center" }}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke={C.primary} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10" />
                    <path d="M12 2c0 5.5-4 10-10 10" />
                </svg>
            </div>
            <span style={{ fontSize:16, fontWeight:500, color:C.secondary, letterSpacing:"0.3px" }}>
                Bjir E-commerce
            </span>
        </div>

        {/* Hero */}
        <div style={{ flex:1, display:"flex", flexDirection:"column", justifyContent:"center", padding:"28px 0", position:"relative", zIndex:1 }}>
            <div style={{ display:"flex", gap:10, marginBottom:20 }}>
                {[
                    { stroke: C.secondary, d: "M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2" },
                    { stroke: C.secondary, d: "M12 3C8.13 3 5 6.13 5 10s3.13 7 7 7 7-3.13 7-7-3.13-7-7-7z" },
                    { stroke: C.accent,    d: "M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" },
                ].map((icon, i) => (
                    <div key={i} style={{
                        width:52, height:52,
                        background: i === 2 ? "rgba(54,69,79,0.18)" : "rgba(255,255,255,0.12)",
                        borderRadius:12,
                        border: `0.5px solid ${i === 2 ? "rgba(54,69,79,0.4)" : "rgba(255,255,255,0.18)"}`,
                        display:"flex", alignItems:"center", justifyContent:"center",
                    }}>
                        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke={icon.stroke} strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                            <path d={icon.d} />
                        </svg>
                    </div>
                ))}
            </div>
            <h1 style={{ fontSize:22, fontWeight:500, color:C.secondary, lineHeight:1.35, margin:"0 0 10px" }}>
                Bergabung<br />Bersama Kami
            </h1>
            <p style={{ fontSize:13, color:"rgba(231,215,201,0.75)", lineHeight:1.65, margin:0 }}>
                Buat akun gratis dan mulai nikmati ribuan produk berkualitas.
                <br />Dapatkan penawaran eksklusif khusus member baru!
            </p>
        </div>

        {/* Trust pills */}
        <div style={{ display:"flex", gap:8, flexWrap:"wrap", position:"relative", zIndex:1 }}>
            {["Gratis daftar", "Promo member", "Belanja aman"].map((label) => (
                <span key={label} style={{
                    fontSize:11, color:C.pillText, background:C.pillBg,
                    border:`0.5px solid ${C.pillBorder}`, borderRadius:99, padding:"4px 12px",
                }}>
                    {label}
                </span>
            ))}
        </div>
    </div>
);

// ── Sub-komponen: InputField (reusable) ────────────────────────────────────
interface InputFieldProps {
    id: string;
    label: string;
    type: string;
    value: string;
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
    placeholder: string;
    autoComplete?: string;
    required?: boolean;
    rightSlot?: React.ReactNode;
}

const InputField: React.FC<InputFieldProps> = ({
    id, label, type, value, onChange, placeholder, autoComplete, required, rightSlot,
}) => {
    const [focused, setFocused] = useState(false);
    return (
        <div style={{ marginBottom:14 }}>
            <label htmlFor={id} style={{ display:"block", fontSize:12, fontWeight:500, color:C.textLabel, marginBottom:5 }}>
                {label}
            </label>
            <div style={{
                display:"flex", alignItems:"center", height:40, background:"#fff",
                border:`1px solid ${focused ? C.primary : C.border}`,
                borderRadius:8, padding:"0 12px", gap:8,
                boxShadow: focused ? `0 0 0 3px rgba(166,123,123,0.18)` : "none",
                transition:"border-color 0.15s, box-shadow 0.15s",
            }}>
                <input
                    id={id} name={id} type={type} value={value} onChange={onChange}
                    onFocus={() => setFocused(true)} onBlur={() => setFocused(false)}
                    placeholder={placeholder} autoComplete={autoComplete} required={required}
                    style={{ flex:1, border:"none", outline:"none", background:"transparent", fontSize:13, color:C.textDark }}
                />
                {rightSlot}
            </div>
        </div>
    );
};

// ── Eye toggle button (reusable) ───────────────────────────────────────────
const EyeToggle: React.FC<{ show: boolean; onToggle: () => void }> = ({ show, onToggle }) => (
    <button
        type="button" onClick={onToggle}
        style={{ background:"none", border:"none", cursor:"pointer", padding:0, color:"#a08888", display:"flex", alignItems:"center" }}
        aria-label={show ? "Sembunyikan kata sandi" : "Tampilkan kata sandi"}
    >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
            {show ? (
                <>
                    <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94" />
                    <path d="M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19" />
                    <line x1="1" y1="1" x2="23" y2="23" />
                </>
            ) : (
                <>
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                    <circle cx="12" cy="12" r="3" />
                </>
            )}
        </svg>
    </button>
);

// ── Sub-komponen: RightPanel ───────────────────────────────────────────────
interface RightPanelProps {
    name: string;
    setName: (v: string) => void;
    email: string;
    setEmail: (v: string) => void;
    password: string;
    setPassword: (v: string) => void;
    confirmPassword: string;
    setConfirmPassword: (v: string) => void;
    error: string | null;
    success: string | null;
    isLoading: boolean;
    onSubmit: (e: React.FormEvent) => void;
}

const RightPanel: React.FC<RightPanelProps> = ({
    name, setName, email, setEmail,
    password, setPassword, confirmPassword, setConfirmPassword,
    error, success, isLoading, onSubmit,
}) => {
    const [showPassword, setShowPassword]               = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    return (
        <div style={{
            flex:1, background:C.secondary, display:"flex",
            flexDirection:"column", justifyContent:"center", padding:"36px 40px",
        }}>
            {/* Header */}
            <div style={{ marginBottom:20 }}>
                <h2 style={{ fontSize:20, fontWeight:500, color:C.textDark, margin:"0 0 4px" }}>
                    Buat akun baru
                </h2>
                <p style={{ fontSize:13, color:C.textMuted, margin:0 }}>
                    Isi data di bawah untuk mulai berbelanja
                </p>
            </div>

            {/* Error */}
            {error && (
                <div style={{
                    background:"#fef2f2", border:"1px solid #fecaca", borderRadius:8,
                    padding:"10px 14px", fontSize:13, color:"#b91c1c", marginBottom:14,
                }}>
                    {error}
                </div>
            )}

            {/* Success */}
            {success && (
                <div style={{
                    background:"#f0fdf4", border:"1px solid #bbf7d0", borderRadius:8,
                    padding:"10px 14px", fontSize:13, color:"#15803d", marginBottom:14,
                }}>
                    {success}
                </div>
            )}

            {/* Form */}
            <form onSubmit={onSubmit}>
                <InputField
                    id="name" label="Nama lengkap" type="text" value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="Nama lengkap Anda" autoComplete="name" required
                />

                <InputField
                    id="email" label="Alamat email" type="email" value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="email@contoh.com" autoComplete="email" required
                />

                <InputField
                    id="password" label="Kata sandi"
                    type={showPassword ? "text" : "password"} value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••" autoComplete="new-password" required
                    rightSlot={<EyeToggle show={showPassword} onToggle={() => setShowPassword(p => !p)} />}
                />

                <InputField
                    id="confirm-password" label="Konfirmasi kata sandi"
                    type={showConfirmPassword ? "text" : "password"} value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="••••••••" autoComplete="new-password" required
                    rightSlot={<EyeToggle show={showConfirmPassword} onToggle={() => setShowConfirmPassword(p => !p)} />}
                />

                {/* Tombol daftar */}
                <button
                    type="submit" disabled={isLoading}
                    style={{
                        width:"100%", height:42, marginTop:4,
                        background: isLoading ? C.primaryLight : C.primary,
                        border:"none", borderRadius:8,
                        display:"flex", alignItems:"center", justifyContent:"center", gap:8,
                        cursor: isLoading ? "not-allowed" : "pointer",
                        transition:"background 0.15s",
                    }}
                >
                    <span style={{ fontSize:14, fontWeight:500, color:C.secondary }}>
                        {isLoading ? "Sedang mendaftar..." : "Daftar sekarang"}
                    </span>
                    {!isLoading && (
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke={C.pillText} strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="5" y1="12" x2="19" y2="12" /><polyline points="12 5 19 12 12 19" />
                        </svg>
                    )}
                </button>
            </form>

            {/* Link masuk */}
            <p style={{ textAlign:"center", marginTop:18, fontSize:12, color:C.textMuted }}>
                Sudah punya akun?{" "}
                <Link to="/login" style={{ color:C.accent, fontWeight:500, textDecoration:"none" }}>
                    Masuk di sini
                </Link>
            </p>
        </div>
    );
};

// ── Komponen utama: Register ───────────────────────────────────────────────
const Register: React.FC = () => {
    const [name, setName]                     = useState("");
    const [email, setEmail]                   = useState("");
    const [password, setPassword]             = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError]                   = useState<string | null>(null);
    const [success, setSuccess]               = useState<string | null>(null);
    const [isLoading, setIsLoading]           = useState(false);

    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);

        if (password !== confirmPassword) {
            setError("Kata sandi dan konfirmasi kata sandi tidak cocok.");
            return;
        }

        setIsLoading(true);

        try {
            await api.post("/v1/auth/register", { name, email, password });

            setSuccess("Pendaftaran berhasil! Mengalihkan ke halaman masuk...");
            setName(""); setEmail(""); setPassword(""); setConfirmPassword("");

            setTimeout(() => navigate("/login"), 2500);
        } catch (err: any) {
            console.error(err);
            if (err.response?.data?.message) {
                setError(err.response.data.message);
            } else {
                setError("Pendaftaran gagal. Silakan periksa kembali data Anda atau coba lagi nanti.");
            }
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div style={{
            display:"flex", minHeight:"100vh", alignItems:"center", justifyContent:"center",
            background:"#ddd0c8", padding:"24px 16px",
        }}>
            {/* Tag style untuk mendaftarkan animasi keyframes secara runtime */}
            <style>{`
                @keyframes registerFadeIn {
                    from {
                        opacity: 0;
                        transform: translateY(12px);
                    }
                    to {
                        opacity: 1;
                        transform: translateY(0);
                    }
                }
            `}</style>

            <div style={{
                display:"flex", width:"100%", maxWidth:900, minHeight:580,
                borderRadius:16, overflow:"hidden", border:`0.5px solid ${C.border}`,
                animation: "registerFadeIn 1s ease-out forwards",
            }}>
                <LeftPanel />
                <RightPanel
                    name={name} setName={setName}
                    email={email} setEmail={setEmail}
                    password={password} setPassword={setPassword}
                    confirmPassword={confirmPassword} setConfirmPassword={setConfirmPassword}
                    error={error} success={success}
                    isLoading={isLoading} onSubmit={handleSubmit}
                />
            </div>
        </div>
    );
};

export default Register;