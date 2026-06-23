import type { FormEvent, ReactNode } from "react";

interface AuthFormProps {
  title: string;
  subtitle: string;
  error?: string | null;
  success?: string | null;
  isSubmitting: boolean;
  submitLabel: string;
  submittingLabel: string;
  children: ReactNode;
  footer: ReactNode;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
}

function AuthForm({
  title,
  subtitle,
  error,
  success,
  isSubmitting,
  submitLabel,
  submittingLabel,
  children,
  footer,
  onSubmit,
}: AuthFormProps) {
  return (
    <div className="mx-auto w-full max-w-md">
      <div className="mb-7">
        <span className="mb-3 inline-flex border-2 border-[var(--color-brutal-ink)] bg-[var(--color-brutal-accent)] px-3 py-1 text-xs font-black uppercase tracking-[0.14em] text-[var(--color-brutal-ink)] shadow-[3px_3px_0_var(--color-brutal-ink)]">
          Buyer Access
        </span>
        <h1 className="m-0 text-4xl font-black uppercase leading-[0.9] tracking-[-0.07em] text-[var(--color-brutal-ink)]">
          {title}
        </h1>
        <p className="mt-3 text-sm font-bold leading-6 text-[var(--color-text-muted)]">
          {subtitle}
        </p>
      </div>

      {error && (
        <div
          className="mb-5 border-2 border-[var(--color-stock-out)] bg-[#ffe1d7] px-4 py-3 text-sm font-bold text-[var(--color-stock-out)] shadow-[3px_3px_0_var(--color-brutal-ink)]"
          role="alert"
        >
          {error}
        </div>
      )}

      {success && (
        <div className="mb-5 border-2 border-[var(--color-stock-in)] bg-[#dcfce7] px-4 py-3 text-sm font-bold text-[var(--color-stock-in)] shadow-[3px_3px_0_var(--color-brutal-ink)]">
          {success}
        </div>
      )}

      <form className="grid gap-4" noValidate onSubmit={onSubmit}>
        {children}

        <button
          className="mt-2 flex min-h-12 items-center justify-center gap-2 border-2 border-[var(--color-brutal-ink)] bg-[var(--color-primary)] px-4 text-sm font-black uppercase tracking-[0.12em] text-white shadow-[4px_4px_0_var(--color-brutal-ink)] transition hover:-translate-x-0.5 hover:-translate-y-0.5 hover:shadow-[6px_6px_0_var(--color-brutal-ink)] disabled:cursor-not-allowed disabled:opacity-70"
          type="submit"
          disabled={isSubmitting}
        >
          {isSubmitting && (
            <span
              className="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white"
              aria-hidden="true"
            />
          )}
          {isSubmitting ? submittingLabel : submitLabel}
        </button>
      </form>

      <div className="mt-6 text-center text-sm font-bold text-[var(--color-text-muted)]">
        {footer}
      </div>
    </div>
  );
}

export default AuthForm;
