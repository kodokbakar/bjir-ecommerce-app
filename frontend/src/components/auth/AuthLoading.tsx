interface AuthLoadingProps {
  message?: string;
}

function AuthLoading({ message = "Memeriksa sesi..." }: AuthLoadingProps) {
  return (
    <div className="grid min-h-screen w-screen place-items-center bg-[var(--color-brutal-paper)]">
      <div className="grid place-items-center gap-3 border-4 border-[var(--color-brutal-ink)] bg-white px-8 py-7 shadow-[6px_6px_0_var(--color-brutal-ink)]">
        <span
          className="h-8 w-8 animate-spin rounded-full border-4 border-[var(--color-brutal-ink)] border-t-[var(--color-brutal-hot)]"
          aria-hidden="true"
        />
        <p className="m-0 text-sm font-black uppercase tracking-[0.12em] text-[var(--color-text-muted)]">
          {message}
        </p>
      </div>
    </div>
  );
}

export default AuthLoading;
