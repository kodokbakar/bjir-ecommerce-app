import { useEffect, useMemo, useState, type ChangeEvent } from "react";
import { useSearchParams } from "react-router-dom";

import Pagination from "../components/Pagination";
import ProductCard from "../components/ProductCard";
import { listCategories, listProducts } from "../services/productService";
import type {
  Category,
  Product,
  ProductListMeta,
  ProductListParams,
} from "../types/product";

const PRODUCT_LIMIT = 12;

type SortOption = "relevance" | "price-asc" | "price-desc" | "newest";

const SORT_OPTIONS: Array<{ value: SortOption; label: string }> = [
  { value: "relevance", label: "Relevance" },
  { value: "price-asc", label: "Price: Low to High" },
  { value: "price-desc", label: "Price: High to Low" },
  { value: "newest", label: "Newest" },
];

function getPositiveNumber(value: string | null, fallback: number): number {
  const parsed = Number(value);

  if (!Number.isInteger(parsed) || parsed < 1) {
    return fallback;
  }

  return parsed;
}

function getSortOption(searchParams: URLSearchParams): SortOption {
  const sortBy = searchParams.get("sort_by");
  const sortOrder = searchParams.get("sort_order");

  if (sortBy === "price" && sortOrder === "asc") {
    return "price-asc";
  }

  if (sortBy === "price" && sortOrder === "desc") {
    return "price-desc";
  }

  if (sortBy === "created_at" && sortOrder === "desc") {
    return "newest";
  }

  return "relevance";
}

function getSortParams(sort: SortOption): Pick<ProductListParams, "sort_by" | "sort_order"> {
  if (sort === "price-asc") {
    return {
      sort_by: "price",
      sort_order: "asc",
    };
  }

  if (sort === "price-desc") {
    return {
      sort_by: "price",
      sort_order: "desc",
    };
  }

  if (sort === "newest") {
    return {
      sort_by: "created_at",
      sort_order: "desc",
    };
  }

  return {};
}

function buildMetaFallback(page: number): ProductListMeta {
  return {
    page,
    limit: PRODUCT_LIMIT,
    total: 0,
    total_pages: 0,
    sort_by: "",
    sort_order: "",
    category_id: "",
    category: "",
    search: "",
  };
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  return "Failed to load products. Please try again.";
}

function ProductsSkeleton() {
  return (
    <div className="products-grid" aria-label="Loading products">
      {Array.from({ length: PRODUCT_LIMIT }, (_, index) => (
        <div className="products-skeleton-card" key={index}>
          <div className="products-skeleton-media" />
          <div className="products-skeleton-body">
            <div className="products-skeleton-line tiny" />
            <div className="products-skeleton-line" />
            <div className="products-skeleton-line short" />
          </div>
        </div>
      ))}
    </div>
  );
}

function Products() {
  const [searchParams, setSearchParams] = useSearchParams();

  const page = getPositiveNumber(searchParams.get("page"), 1);
  const category = searchParams.get("category") ?? "";
  const search = searchParams.get("search") ?? "";
  const sort = getSortOption(searchParams);

  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [meta, setMeta] = useState<ProductListMeta>(() => buildMetaFallback(page));
  const [isLoadingProducts, setIsLoadingProducts] = useState(true);
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);
  const [categoryError, setCategoryError] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [reloadKey, setReloadKey] = useState(0);
  const [categoryReloadKey, setCategoryReloadKey] = useState(0);

  const hasActiveFilters = Boolean(category || search || sort !== "relevance");

  const productQuery = useMemo<ProductListParams>(() => {
    return {
      page,
      limit: PRODUCT_LIMIT,
      category: category || undefined,
      search: search || undefined,
      ...getSortParams(sort),
    };
  }, [category, page, search, sort]);

  useEffect(() => {
    let isMounted = true;

    async function loadCategories() {
      setIsLoadingCategories(true);
      setCategoryError(null);

      try {
        const result = await listCategories();

        if (isMounted) {
          setCategories(result);
        }
      } catch (loadError) {
        console.error("Failed to load categories:", loadError);

        if (isMounted) {
          setCategories([]);
          setCategoryError(getErrorMessage(loadError));
        }
      } finally {
        if (isMounted) {
          setIsLoadingCategories(false);
        }
      }
    }

    loadCategories();

    return () => {
      isMounted = false;
    };
  }, [categoryReloadKey]);

  useEffect(() => {
    let isMounted = true;

    async function loadProducts() {
      setIsLoadingProducts(true);
      setError(null);

      try {
        const result = await listProducts(productQuery);

        if (isMounted) {
          setProducts(result.data);
          setMeta(result.meta);
        }
      } catch (loadError) {
        if (isMounted) {
          setProducts([]);
          setMeta(buildMetaFallback(page));
          setError(getErrorMessage(loadError));
        }
      } finally {
        if (isMounted) {
          setIsLoadingProducts(false);
        }
      }
    }

    loadProducts();

    return () => {
      isMounted = false;
    };
  }, [page, productQuery, reloadKey]);

  function updateQuery(nextValues: Record<string, string | number | null>) {
    const nextParams = new URLSearchParams(searchParams);

    Object.entries(nextValues).forEach(([key, value]) => {
      if (value === null || value === "") {
        nextParams.delete(key);
        return;
      }

      nextParams.set(key, String(value));
    });

    setSearchParams(nextParams);
  }

  function handleSearchChange(event: ChangeEvent<HTMLInputElement>) {
    updateQuery({
      search: event.target.value,
      page: 1,
    });
  }

  function handleCategoryChange(event: ChangeEvent<HTMLSelectElement>) {
    updateQuery({
      category: event.target.value,
      page: 1,
    });
  }

  function handleSortChange(event: ChangeEvent<HTMLSelectElement>) {
    const nextSort = event.target.value as SortOption;
    const nextParams = getSortParams(nextSort);

    updateQuery({
      sort_by: nextParams.sort_by ?? null,
      sort_order: nextParams.sort_order ?? null,
      page: 1,
    });
  }

  function handleClearFilters() {
    setSearchParams(new URLSearchParams({ page: "1" }));
  }

  function handlePageChange(nextPage: number) {
    updateQuery({
      page: nextPage,
    });
  }

  function handleRetry() {
    setReloadKey((current) => current + 1);
  }

  function handleCategoryRetry() {
    setCategoryReloadKey((current) => current + 1);
  }

  return (
    <section className="products-page" aria-labelledby="products-title">
      <header className="products-hero">
        <span className="products-eyebrow">Product Catalog</span>
        <h1 className="products-title" id="products-title">
          Browse the loud shelf.
        </h1>
        <p className="products-copy">
          A compact catalog for buyers who want the product, the price, and the stock
          status without wandering through decoration.
        </p>
      </header>

      <div className="products-toolbar" aria-label="Product controls">
        <label className="products-field">
          <span className="products-label">Search</span>
          <input
            className="products-input"
            type="search"
            value={search}
            onChange={handleSearchChange}
            placeholder="Search products..."
            aria-label="Search products"
          />
        </label>

        <label className="products-field">
          <span className="products-label">Category</span>
          <select
            className="products-select"
            value={category}
            onChange={handleCategoryChange}
            disabled={isLoadingCategories}
            aria-label="Filter by category"
          >
            <option value="">
              {isLoadingCategories ? "Loading categories..." : "All Categories"}
            </option>
            {categories.map((item) => (
              <option key={item.id} value={item.slug}>
                {item.name}
              </option>
            ))}
          </select>

          {categoryError && (
            <div className="products-category-error" role="status">
              <span>Category list failed to load.</span>
              <button type="button" onClick={handleCategoryRetry}>
                Retry
              </button>
            </div>
          )}
        </label>

        <label className="products-field">
          <span className="products-label">Sort</span>
          <select
            className="products-select"
            value={sort}
            onChange={handleSortChange}
            aria-label="Sort products"
          >
            {SORT_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
      </div>

      <div className="products-status-line">
        <span>
          {isLoadingProducts
            ? "Loading product rows..."
            : `${meta.total} product${meta.total === 1 ? "" : "s"} found`}
          {category ? ` in ${category}` : ""}
          {search ? ` for "${search}"` : ""}
        </span>

        {hasActiveFilters && (
          <button className="products-clear-button" type="button" onClick={handleClearFilters}>
            Clear filters
          </button>
        )}
      </div>

      {isLoadingProducts ? (
        <ProductsSkeleton />
      ) : error ? (
        <div className="products-panel-state" role="alert">
          <div>
            <h2>Catalog jammed.</h2>
            <p>{error}</p>
            <button className="products-retry-button" type="button" onClick={handleRetry}>
              Retry
            </button>
          </div>
        </div>
      ) : products.length === 0 ? (
        <div className="products-panel-state">
          <div>
            <h2>No products found.</h2>
            <p>
              {hasActiveFilters
                ? "Your current search, category, or sorting context returned no products. Reset the filters or try a wider query."
                : "The catalog is empty. Add products from the admin side before exposing this shelf to buyers."}
            </p>
            {hasActiveFilters && (
              <button className="products-retry-button" type="button" onClick={handleClearFilters}>
                Reset catalog
              </button>
            )}
          </div>
        </div>
      ) : (
        <>
          <div className="products-grid">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>

          <Pagination
            page={meta.page}
            limit={meta.limit}
            total={meta.total}
            totalPages={meta.total_pages}
            onPageChange={handlePageChange}
          />
        </>
      )}
    </section>
  );
}

export default Products;
