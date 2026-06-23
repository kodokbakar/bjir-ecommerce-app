import { Eye, EyeOff } from "lucide-react";

interface PasswordToggleProps {
  show: boolean;
  onToggle: () => void;
}

function PasswordToggle({ show, onToggle }: PasswordToggleProps) {
  return (
    <button
      className="grid h-8 w-8 shrink-0 place-items-center rounded-lg text-[var(--color-text-muted)] transition hover:bg-[var(--color-secondary)] hover:text-[var(--color-brutal-ink)]"
      type="button"
      onClick={onToggle}
      aria-label={show ? "Sembunyikan kata sandi" : "Tampilkan kata sandi"}
    >
      {show ? (
        <EyeOff className="h-4 w-4" aria-hidden="true" />
      ) : (
        <Eye className="h-4 w-4" aria-hidden="true" />
      )}
    </button>
  );
}

export default PasswordToggle;
