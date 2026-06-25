import { createContext, useContext } from "react";

export type ToastType = "success" | "error" | "warning" | "info";

export interface ToastMessage {
  id: string;
  type: ToastType;
  title?: string;
  message: string;
}

export interface ShowToastInput {
  type: ToastType;
  title?: string;
  message: string;
}

export interface ShowToastOptions {
  duration?: number;
}

export interface ToastContextValue {
  toasts: ToastMessage[];
  showToast: (toast: ShowToastInput, options?: ShowToastOptions) => string;
  dismissToast: (toastID: string) => void;
}

export const ToastContext = createContext<ToastContextValue | undefined>(
  undefined,
);

export function useToast() {
  const context = useContext(ToastContext);

  if (!context) {
    throw new Error("useToast must be used inside ToastProvider");
  }

  return context;
}
