import { useMemo, useState, type FormEvent } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";

import AuthForm from "../components/auth/AuthForm";
import AuthLayout from "../components/auth/AuthLayout";
import FormField from "../components/auth/FormField";
import PasswordToggle from "../components/auth/PasswordToggle";
import { useAuth } from "../hooks/useAuth";
import { getApiErrorMessage, registerUser } from "../services/authService";
import { getDashboardPath } from "../utils/authRouting";

type RegisterFieldErrors = Partial<
  Record<"name" | "email" | "password" | "confirmPassword", string>
>;

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function validateRegisterForm(
  name: string,
  email: string,
  password: string,
  confirmPassword: string,
): RegisterFieldErrors {
  const errors: RegisterFieldErrors = {};

  if (name.trim().length < 2) {
    errors.name = "Nama minimal 2 karakter.";
  }

  if (!email.trim()) {
    errors.email = "Email wajib diisi.";
  } else if (!EMAIL_PATTERN.test(email.trim())) {
    errors.email = "Format email tidak valid.";
  }

  if (password.length < 8) {
    errors.password = "Kata sandi minimal 8 karakter.";
  } else if (password.length > 72) {
    errors.password = "Kata sandi maksimal 72 karakter.";
  }

  if (!confirmPassword) {
    errors.confirmPassword = "Konfirmasi kata sandi wajib diisi.";
  } else if (password !== confirmPassword) {
    errors.confirmPassword = "Konfirmasi kata sandi tidak cocok.";
  }

  return errors;
}

function getPasswordScore(password: string): number {
  let score = 0;

  if (password.length >= 8) {
    score += 1;
  }

  if (/[A-Z]/.test(password)) {
    score += 1;
  }

  if (/[0-9]/.test(password)) {
    score += 1;
  }

  if (/[^A-Za-z0-9]/.test(password)) {
    score += 1;
  }

  return score;
}

function getPasswordStrengthLabel(score: number): string {
  if (score <= 1) {
    return "Lemah";
  }

  if (score === 2) {
    return "Cukup";
  }

  if (score === 3) {
    return "Bagus";
  }

  return "Kuat";
}

function Register() {
  const navigate = useNavigate();
  const { login, user, isAuthenticated, isLoading: isAuthLoading } = useAuth();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<RegisterFieldErrors>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const passwordScore = useMemo(() => getPasswordScore(password), [password]);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const nextErrors = validateRegisterForm(
      name,
      email,
      password,
      confirmPassword,
    );
    setFieldErrors(nextErrors);
    setFormError(null);

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    setIsSubmitting(true);

    try {
      const result = await registerUser({
        name: name.trim(),
        email: email.trim(),
        password,
      });

      login(result.accessToken, result.user, true);
      navigate(getDashboardPath(result.user), { replace: true });
    } catch (error) {
      setFormError(
        getApiErrorMessage(error, "Gagal mendaftar. Silakan coba lagi."),
      );
      setIsSubmitting(false);
    }
  }

  if (!isAuthLoading && isAuthenticated) {
    return <Navigate to={getDashboardPath(user)} replace />;
  }

  return (
    <AuthLayout variant="register">
      <AuthForm
        title="Buat akun baru"
        subtitle="Daftar sekali, lalu mulai jelajahi katalog produk dan checkout lebih cepat."
        error={formError}
        isSubmitting={isSubmitting}
        submitLabel="Daftar sekarang"
        submittingLabel="Mendaftarkan..."
        onSubmit={handleSubmit}
        footer={
          <>
            Sudah punya akun?{" "}
            <Link
              className="font-black text-[var(--color-primary-dark)] no-underline hover:underline"
              to="/login"
            >
              Masuk di sini
            </Link>
          </>
        }
      >
        <FormField
          id="name"
          label="Nama"
          type="text"
          value={name}
          onChange={(event) => setName(event.target.value)}
          placeholder="Nama Anda"
          autoComplete="name"
          error={fieldErrors.name}
        />

        <FormField
          id="email"
          label="Email"
          type="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          placeholder="email@contoh.com"
          autoComplete="email"
          error={fieldErrors.email}
        />

        <div className="grid gap-2">
          <FormField
            id="password"
            label="Password"
            type={showPassword ? "text" : "password"}
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            placeholder="••••••••"
            autoComplete="new-password"
            error={fieldErrors.password}
            rightSlot={
              <PasswordToggle
                show={showPassword}
                onToggle={() => setShowPassword((current) => !current)}
              />
            }
          />

          {password.length > 0 && (
            <div className="grid gap-2">
              <div className="grid grid-cols-4 gap-2" aria-hidden="true">
                {Array.from({ length: 4 }, (_, index) => (
                  <span
                    className={[
                      "h-2 border border-[var(--color-brutal-ink)]",
                      index < passwordScore
                        ? "bg-[var(--color-primary)]"
                        : "bg-[var(--color-secondary)]",
                    ].join(" ")}
                    key={index}
                  />
                ))}
              </div>
              <p className="m-0 text-xs font-bold text-[var(--color-text-muted)]">
                Kekuatan password: {getPasswordStrengthLabel(passwordScore)}
              </p>
            </div>
          )}
        </div>

        <FormField
          id="confirmPassword"
          label="Konfirmasi Password"
          type={showConfirmPassword ? "text" : "password"}
          value={confirmPassword}
          onChange={(event) => setConfirmPassword(event.target.value)}
          placeholder="••••••••"
          autoComplete="new-password"
          error={fieldErrors.confirmPassword}
          rightSlot={
            <PasswordToggle
              show={showConfirmPassword}
              onToggle={() => setShowConfirmPassword((current) => !current)}
            />
          }
        />
      </AuthForm>
    </AuthLayout>
  );
}

export default Register;
