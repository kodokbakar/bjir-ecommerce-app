import {
  AlertTriangle,
  CheckCircle2,
  Info,
  ShieldAlert,
  X,
  type LucideIcon,
} from "lucide-react";

import { useToast, type ToastType } from "../context/toast"

const TOAST_ICONS: Record<ToastType, LucideIcon> = {
  success: CheckCircle2,
  error: AlertTriangle,
  warning: ShieldAlert,
  info: Info,
};

const TOAST_TITLES: Record<ToastType, string> = {
  success: "Success",
  error: "Error",
  warning: "Warning",
  info: "Info",
};

function ToastViewport() {
  const { toasts, dismissToast } = useToast();

  if (toasts.length === 0) {
    return null;
  }

  return (
    <section
      className="toast-viewport"
      aria-label="Notifications"
      aria-live="polite"
    >
      {toasts.map((toast) => {
        const Icon = TOAST_ICONS[toast.type];
        const title = toast.title ?? TOAST_TITLES[toast.type];
        const isAssertive = toast.type === "error" || toast.type === "warning";

        return (
          <article
            className={`toast-card is-${toast.type}`}
            key={toast.id}
            role={isAssertive ? "alert" : "status"}
          >
            <span className="toast-icon" aria-hidden="true">
              <Icon className="h-5 w-5" />
            </span>

            <div className="toast-copy">
              <strong>{title}</strong>
              <p>{toast.message}</p>
            </div>

            <button
              type="button"
              aria-label={`Dismiss ${title.toLowerCase()} notification`}
              onClick={() => dismissToast(toast.id)}
            >
              <X className="h-4 w-4" aria-hidden="true" />
            </button>
          </article>
        );
      })}
    </section>
  );
}

export default ToastViewport;
