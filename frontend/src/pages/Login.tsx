import { useMemo, useState, type FormEvent } from "react";
import { Link, Navigate, useLocation, useNavigate } from "react-router-dom";

import AuthForm from "../components/auth/AuthForm";
import AuthLayout from "../components/auth/AuthLayout";
import FormField from "../components/auth/FormField";
import PasswordToggle from "../components/auth/PasswordToggle";
import { useAuth } from "../hooks/useAuth";
import { getApiErrorMessage, loginUser } from "../services/authService";

interface LoginLocationState {
  authNotice?: string;
}

type LoginFieldErrors = Partial<Record<"email" | "password", string>>;

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function isLoginLocationState(value: unknown): value is LoginLocationState {
  return typeof value === "object" && value !== null && "authNotice" in value;
}

function validateLoginForm(email: string, password: string): LoginFieldErrors {
  const errors: LoginFieldErrors = {};

  if (!email.trim()) {
    errors.email = "Email wajib diisi.";
  } else if (!EMAIL_PATTERN.test(email.trim())) {
    errors.email = "Format email tidak valid.";
  }

  if (!password) {
    errors.password = "Kata sandi wajib diisi.";
  } else if (password.length < 8) {
    errors.password = "Kata sandi minimal 8 karakter.";
  }

  return errors;
}

function AuthPageLoading() {
  return (
    <div className="grid min-h-screen place-items-center bg-[var(--color-brutal-paper)]">
      <div className="grid place-items-center gap-3 border-4 border-[var(--color-brutal-ink)] bg-white px-8 py-7 shadow-[6px_6px_0_var(--color-brutal-ink)]">
        <span className="h-7 w-7 animate-spin rounded-full border-4 border-[var(--color-brutal-ink)] border-t-[var(--color-brutal-hot)]" />
        <p className="m-0 text-sm font-black uppercase tracking-[0.12em] text-[var(--color-text-muted)]">
          Memeriksa sesi...
        </p>
      </div>
    </div>
  );
}

function Login() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, isAuthenticated, isLoading: isAuthLoading } = useAuth();

  const authNotice = useMemo(() => {
    return isLoginLocationState(location.state) ? location.state.authNotice ?? null : null;
  }, [location.state]);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rememberMe, setRememberMe] = useState(true);
  const [showPassword, setShowPassword] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<LoginFieldErrors>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const nextErrors = validateLoginForm(email, password);
    setFieldErrors(nextErrors);
    setFormError(null);

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    setIsSubmitting(true);

    try {
      const result = await loginUser({
        email: email.trim(),
        password,
      });

      login(result.accessToken, result.user, rememberMe);
      navigate("/dashboard", { replace: true });
    } catch (error) {
      setFormError(getApiErrorMessage(error, "Login gagal. Silakan coba lagi."));
      setIsSubmitting(false);
    }
  }

  if (isAuthLoading) {
    return <AuthPageLoading />;
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <AuthLayout variant="login">
      <AuthForm
        title="Selamat datang"
        subtitle="Masuk untuk lanjut belanja, cek produk, dan lihat rekomendasi terbaru."
        error={formError}
        success={authNotice}
        isSubmitting={isSubmitting}
        submitLabel="Masuk"
        submittingLabel="Sedang masuk..."
        onSubmit={handleSubmit}
        footer={
          <>
            Belum punya akun?{" "}
            <Link
              className="font-black text-[var(--color-primary-dark)] no-underline hover:underline"
              to="/register"
            >
              Daftar sekarang
            </Link>
          </>
        }
      >
        <FormField
          id="email"
          label="Alamat email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="email@contoh.com"
          autoComplete="email"
          error={fieldErrors.email}
        />

        <FormField
          id="password"
          label="Kata sandi"
          type={showPassword ? "text" : "password"}
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          placeholder="••••••••"
          autoComplete="current-password"
          error={fieldErrors.password}
          rightSlot={
            <PasswordToggle
              show={showPassword}
              onToggle={() => setShowPassword((current) => !current)}
            />
          }
        />

        <div className="flex flex-col gap-3 text-sm font-bold text-[var(--color-text-muted)] sm:flex-row sm:items-center sm:justify-between">
          <label className="inline-flex cursor-pointer items-center gap-2">
            <input
              className="h-4 w-4 accent-[var(--color-primary)]"
              type="checkbox"
              checked={rememberMe}
              onChange={(event) => setRememberMe(event.target.checked)}
            />
            Ingat saya
          </label>

          <a
            className="text-[var(--color-primary-dark)] no-underline hover:underline"
            href="#forgot-password"
          >
            Lupa kata sandi?
          </a>
        </div>
      </AuthForm>
    </AuthLayout>
  );
}

export default Login;
