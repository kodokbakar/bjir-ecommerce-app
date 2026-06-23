import { useCallback, useEffect, useRef, useState } from "react";

import {
  getCartCount,
  getCartErrorMessage,
  subscribeToCartChanges,
} from "../services/cartService";

interface UseCartCountResult {
  count: number;
  error: string | null;
  refresh: () => Promise<void>;
}

export function useCartCount(): UseCartCountResult {
  const isMountedRef = useRef(false);

  const [count, setCount] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    try {
      const nextCount = await getCartCount();

      if (isMountedRef.current) {
        setCount(nextCount);
        setError(null);
      }
    } catch (refreshError) {
      if (isMountedRef.current) {
        setCount(0);
        setError(
          getCartErrorMessage(refreshError, "Failed to load cart count."),
        );
      }
    }
  }, []);

  useEffect(() => {
    isMountedRef.current = true;

    const initialRefreshID = window.setTimeout(() => {
      void refresh();
    }, 0);

    const unsubscribe = subscribeToCartChanges(() => {
      void refresh();
    });

    return () => {
      window.clearTimeout(initialRefreshID);
      isMountedRef.current = false;
      unsubscribe();
    };
  }, [refresh]);

  return {
    count,
    error,
    refresh,
  };
}
