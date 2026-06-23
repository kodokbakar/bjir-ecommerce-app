import { useEffect, useState } from "react";

import { listCategories } from "../../services/productService";
import type { Category } from "../../types/product";

const CATEGORY_LIMIT = 6;

let cachedCategories: Category[] | null = null;
let categoryRequest: Promise<Category[]> | null = null;

function getCachedCategories(): Category[] {
  return cachedCategories?.slice(0, CATEGORY_LIMIT) ?? [];
}

export function useLayoutCategories() {
  const [categories, setCategories] = useState<Category[]>(getCachedCategories);
  const [isLoading, setIsLoading] = useState(cachedCategories === null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (cachedCategories) {
      return;
    }

    let isMounted = true;

    async function loadCategories() {
      try {
        categoryRequest ??= listCategories();

        const result = await categoryRequest;
        cachedCategories = result;

        if (isMounted) {
          setCategories(result.slice(0, CATEGORY_LIMIT));
          setError(null);
        }
      } catch (loadError) {
        console.error("Failed to load layout categories:", loadError);

        categoryRequest = null;

        if (isMounted) {
          setCategories([]);
          setError("Categories unavailable");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    }

    loadCategories();

    return () => {
      isMounted = false;
    };
  }, []);

  return {
    categories,
    error,
    isLoading,
  };
}
