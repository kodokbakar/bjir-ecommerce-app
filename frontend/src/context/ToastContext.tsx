import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";

import {
  ToastContext,
  type ShowToastInput,
  type ShowToastOptions,
  type ToastMessage,
} from "./toast";

const DEFAULT_TOAST_DURATION_MS = 4000;

function createToastID(): string {
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const timersRef = useRef<Map<string, number>>(new Map());

  const dismissToast = useCallback((toastID: string) => {
    const timerID = timersRef.current.get(toastID);

    if (timerID !== undefined) {
      window.clearTimeout(timerID);
      timersRef.current.delete(toastID);
    }

    setToasts((currentToasts) =>
      currentToasts.filter((toast) => toast.id !== toastID),
    );
  }, []);

  const showToast = useCallback(
    (toast: ShowToastInput, options: ShowToastOptions = {}) => {
      const nextToast: ToastMessage = {
        id: createToastID(),
        ...toast,
      };

      setToasts((currentToasts) => [nextToast, ...currentToasts]);

      const duration = options.duration ?? DEFAULT_TOAST_DURATION_MS;

      if (duration > 0) {
        const timerID = window.setTimeout(() => {
          dismissToast(nextToast.id);
        }, duration);

        timersRef.current.set(nextToast.id, timerID);
      }

      return nextToast.id;
    },
    [dismissToast],
  );

  useEffect(() => {
    const timers = timersRef.current;

    return () => {
      timers.forEach((timerID) => {
        window.clearTimeout(timerID);
      });

      timers.clear();
    };
  }, []);

  const value = useMemo(
    () => ({
      toasts,
      showToast,
      dismissToast,
    }),
    [dismissToast, showToast, toasts],
  );

  return (
    <ToastContext.Provider value={value}>{children}</ToastContext.Provider>
  );
}
