import type { InputHTMLAttributes, ReactNode } from "react";

interface FormFieldProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "id"> {
  id: string;
  label: string;
  error?: string;
  rightSlot?: ReactNode;
}

function FormField({
  id,
  label,
  error,
  rightSlot,
  ...inputProps
}: FormFieldProps) {
  const errorId = `${id}-error`;

  return (
    <div className="grid gap-2">
      <label
        className="text-xs font-black uppercase tracking-[0.12em] text-[var(--color-text-label)]"
        htmlFor={id}
      >
        {label}
      </label>

      <div
        className={[
          "flex min-h-12 items-center gap-2 border-2 bg-white px-3 shadow-[3px_3px_0_var(--color-brutal-ink)] transition",
          error
            ? "border-[var(--color-stock-out)]"
            : "border-[var(--color-brutal-ink)] focus-within:border-[var(--color-brutal-hot)] focus-within:shadow-[5px_5px_0_var(--color-brutal-ink)]",
        ].join(" ")}
      >
        <input
          {...inputProps}
          id={id}
          name={id}
          className="min-w-0 flex-1 border-0 bg-transparent text-sm font-bold text-[var(--color-brutal-ink)] outline-none placeholder:text-[var(--color-text-muted)]"
          aria-invalid={Boolean(error)}
          aria-describedby={error ? errorId : undefined}
        />

        {rightSlot}
      </div>

      {error && (
        <p className="m-0 text-xs font-bold text-[var(--color-stock-out)]" id={errorId}>
          {error}
        </p>
      )}
    </div>
  );
}

export default FormField;
