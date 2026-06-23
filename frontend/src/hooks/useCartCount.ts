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
          getCartErrorMessage(refreshError, "Failed to refresh cart count."),
        );
      }
    }
  }, []);

  useEffect(() => {
    isMountedRef.current = true;
    let isActive = true;

    async function loadInitialCount() {
      try {
        const nextCount = await getCartCount();

        if (isActive) {
          setCount(nextCount);
          setError(null);
        }
      } catch (loadError) {
        if (isActive) {
          setCount(0);
          setError(
            getCartErrorMessage(loadError, "Failed to load cart count."),
          );
        }
      }
    }

    loadInitialCount();

    const unsubscribe = subscribeToCartChanges(() => {
      void refresh();
    });

    return () => {
      isActive = false;
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
