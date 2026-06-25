import { useEffect, useState } from "react";
import { Link, useParams, useSearchParams } from "react-router-dom";

import Pagination from "../components/Pagination";
import ProductCard from "../components/ProductCard";
import { getCategoryBySlug, listProducts } from "../services/productService";
import type { Category, Product, ProductListMeta } from "../types/product";
import EmptyState from "../components/EmptyState";

const PRODUCT_LIMIT = 12;

type CategoryPageState = "loading" | "ready" | "error" | "not-found";

function getPositiveNumber(value: string | null, fallback: number): number {
  const parsed = Number(value);

  if (!Number.isInteger(parsed) || parsed < 1) {
    return fallback;
  }

  return parsed;
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

function isNotFoundError(error: unknown): boolean {
  if (!error || typeof error !== "object" || !("response" in error)) {
    return false;
  }

  const apiError = error as { response?: { status?: number } };

  return apiError.response?.status === 404;
}

function getErrorMessage(error: unknown): string {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  return "Failed to load category products. Please try again.";
}

function CategoryProductsSkeleton() {
  return (
    <div className="products-grid" aria-label="Loading category products">
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

function CategoryProducts() {
  const { slug } = useParams<{ slug: string }>();
  const [searchParams, setSearchParams] = useSearchParams();

  const page = getPositiveNumber(searchParams.get("page"), 1);

  const [category, setCategory] = useState<Category | null>(null);
  const [products, setProducts] = useState<Product[]>([]);
  const [meta, setMeta] = useState<ProductListMeta>(() =>
    buildMetaFallback(page),
  );
  const [state, setState] = useState<CategoryPageState>("loading");
  const [error, setError] = useState<string | null>(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    let isMounted = true;

    async function loadCategoryProducts() {
      if (!slug) {
        setState("not-found");
        return;
      }

      setState("loading");
      setError(null);

      try {
        const [categoryResult, productResult] = await Promise.all([
          getCategoryBySlug(slug),
          listProducts({
            page,
            limit: PRODUCT_LIMIT,
            category: slug,
          }),
        ]);

        if (isMounted) {
          setCategory(categoryResult);
          setProducts(productResult.data);
          setMeta(productResult.meta);
          setState("ready");
        }
      } catch (loadError) {
        if (!isMounted) {
          return;
        }

        setCategory(null);
        setProducts([]);
        setMeta(buildMetaFallback(page));

        if (isNotFoundError(loadError)) {
          setState("not-found");
          return;
        }

        setError(getErrorMessage(loadError));
        setState("error");
      }
    }

    loadCategoryProducts();

    return () => {
      isMounted = false;
    };
  }, [page, reloadKey, slug]);

  function handlePageChange(nextPage: number) {
    const nextParams = new URLSearchParams(searchParams);
    nextParams.set("page", String(nextPage));
    setSearchParams(nextParams);
  }

  function handleRetry() {
    setReloadKey((current) => current + 1);
  }

  if (state === "not-found") {
    return (
      <section className="category-products-page">
        <EmptyState
          tone="error"
          eyebrow="Category Missing"
          title="Category not found."
          description="This category link does not exist anymore."
          action={
            <Link className="product-detail-back-link" to="/products">
              Back to products
            </Link>
          }
        />
      </section>
    );
  }

  if (state === "error") {
    return (
      <section className="category-products-page">
        <EmptyState
          tone="error"
          eyebrow="Category Error"
          title="Category jammed."
          description={
            error || "Failed to load category products. Please try again."
          }
          action={
            <button
              className="products-retry-button"
              type="button"
              onClick={handleRetry}
            >
              Retry
            </button>
          }
        />
      </section>
    );
  }

  const categoryTitle = category?.name || "Loading category";
  const categoryDescription =
    category?.description?.trim() ||
    "Products grouped under this catalog category.";

  return (
    <section
      className="category-products-page"
      aria-labelledby="category-products-title"
    >
      <nav className="product-detail-breadcrumbs" aria-label="Breadcrumb">
        <Link to="/dashboard">Home</Link>
        <span aria-hidden="true">/</span>
        <Link to="/products">Categories</Link>
        <span aria-hidden="true">/</span>
        <span aria-current="page">{categoryTitle}</span>
      </nav>

      <header className="category-products-hero">
        <span className="products-eyebrow">Category Shelf</span>
        <h1 className="products-title" id="category-products-title">
          {categoryTitle}
        </h1>
        <p className="products-copy">{categoryDescription}</p>
        {category?.slug && (
          <Link
            className="category-products-filter-link"
            to={`/products?category=${encodeURIComponent(category.slug)}`}
          >
            Open filtered catalog
          </Link>
        )}
      </header>

      <div className="products-status-line">
        <span>
          {state === "loading"
            ? "Loading category rows..."
            : `${meta.total} product${meta.total === 1 ? "" : "s"} in ${categoryTitle}`}
        </span>
      </div>

      {state === "loading" ? (
        <CategoryProductsSkeleton />
      ) : products.length === 0 ? (
        <EmptyState
          eyebrow="Category Shelf"
          title="No products here."
          description="This category exists, but it has no products attached yet."
          action={
            <Link className="products-retry-button" to="/products">
              Back to catalog
            </Link>
          }
        />
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

export default CategoryProducts;
