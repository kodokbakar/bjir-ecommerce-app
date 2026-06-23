import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";

import { getImageUrl, getProductBySlug } from "../services/productService";
import { C } from "../styles/tokens";
import type { Product } from "../types/product";

type DetailState = "loading" | "ready" | "error" | "not-found";

type StockState = {
  label: string;
  className: string;
};

function formatRupiah(value: number): string {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);
}

function formatDate(value?: string): string {
  if (!value) {
    return "Unknown";
  }

  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "Unknown";
  }

  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

function getStockState(stock: number): StockState {
  if (stock <= 0) {
    return {
      label: "Out of Stock",
      className: "is-out",
    };
  }

  if (stock <= 5) {
    return {
      label: "Low Stock",
      className: "is-low",
    };
  }

  return {
    label: "In Stock",
    className: "is-in",
  };
}

function getProductImage(product: Product): string {
  const galleryImage = product.images?.[0]?.image_url;
  return product.image_url || galleryImage || "";
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

  return "Failed to load product detail. Please try again.";
}

function ProductDetailPlaceholder() {
  return (
    <span className="product-detail-placeholder" aria-hidden="true">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.4">
        <path d="M4 7h16v12H4z" />
        <path d="M8 7a4 4 0 0 1 8 0" />
        <path d="M8 13h8" />
      </svg>
    </span>
  );
}

function ProductDetailSkeleton() {
  return (
    <section className="product-detail-page" aria-label="Loading product detail">
      <div className="product-detail-skeleton breadcrumb" />
      <div className="product-detail-shell">
        <div className="product-detail-skeleton image" />
        <div className="product-detail-panel">
          <div className="product-detail-skeleton line tiny" />
          <div className="product-detail-skeleton line title" />
          <div className="product-detail-skeleton line price" />
          <div className="product-detail-skeleton line" />
          <div className="product-detail-skeleton line short" />
          <div className="product-detail-skeleton button" />
        </div>
      </div>
    </section>
  );
}

function ProductDetail() {
  const { slug } = useParams<{ slug: string }>();

  const [product, setProduct] = useState<Product | null>(null);
  const [state, setState] = useState<DetailState>("loading");
  const [error, setError] = useState<string | null>(null);
  const [imageFailed, setImageFailed] = useState(false);
  const [reloadKey, setReloadKey] = useState(0);

  const stockState = product ? getStockState(product.stock) : null;
  const imagePath = product ? getProductImage(product) : "";
  const imageUrl = imagePath ? getImageUrl(imagePath) : "";
  const categoryName = product?.category?.name || "Uncategorized";
  const categorySlug = product?.category?.slug || "";

  const descriptionBlocks = useMemo(() => {
    const description = product?.description?.trim();

    if (!description) {
      return ["No description has been written for this product yet."];
    }

    return description.split(/\n+/).filter(Boolean);
  }, [product?.description]);

  useEffect(() => {
    let isMounted = true;

    async function loadProduct() {
      if (!slug) {
        setProduct(null);
        setState("not-found");
        return;
      }

      setState("loading");
      setError(null);
      setImageFailed(false);

      try {
        const result = await getProductBySlug(slug);

        if (isMounted) {
          setProduct(result);
          setState("ready");
        }
      } catch (loadError) {
        if (!isMounted) {
          return;
        }

        setProduct(null);

        if (isNotFoundError(loadError)) {
          setState("not-found");
          return;
        }

        setError(getErrorMessage(loadError));
        setState("error");
      }
    }

    loadProduct();

    return () => {
      isMounted = false;
    };
  }, [slug, reloadKey]);

  function handleRetry() {
    setReloadKey((current) => current + 1);
  }

  if (state === "loading") {
    return <ProductDetailSkeleton />;
  }

  if (state === "not-found") {
    return (
      <section className="product-detail-page">
        <div className="products-panel-state" role="alert">
          <div>
            <h2>Product not found.</h2>
            <p>
              This product slug does not exist anymore, or the catalog link is outdated.
            </p>
            <Link className="product-detail-back-link" to="/products">
              Back to products
            </Link>
          </div>
        </div>
      </section>
    );
  }

  if (state === "error" || !product || !stockState) {
    return (
      <section className="product-detail-page">
        <div className="products-panel-state" role="alert">
          <div>
            <h2>Detail jammed.</h2>
            <p>{error || "Failed to load product detail. Please try again."}</p>
            <button className="products-retry-button" type="button" onClick={handleRetry}>
              Retry
            </button>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="product-detail-page" aria-labelledby="product-detail-title">
      <nav className="product-detail-breadcrumbs" aria-label="Breadcrumb">
        <Link to="/dashboard">Home</Link>
        <span aria-hidden="true">/</span>
        <Link to="/products">Products</Link>
        <span aria-hidden="true">/</span>
        {categorySlug ? (
          <Link to={`/products?category=${encodeURIComponent(categorySlug)}`}>
            {categoryName}
          </Link>
        ) : (
          <span>{categoryName}</span>
        )}
        <span aria-hidden="true">/</span>
        <span aria-current="page">{product.name}</span>
      </nav>

      <div className="product-detail-shell">
        <div className="product-detail-image-stage">
          <span className={`product-detail-stock ${stockState.className}`}>
            {stockState.label}
          </span>

          {imageUrl && !imageFailed ? (
            <img
              src={imageUrl}
              alt={product.name}
              onError={() => setImageFailed(true)}
            />
          ) : (
            <ProductDetailPlaceholder />
          )}
        </div>

        <article className="product-detail-panel">
          {categorySlug ? (
            <Link
              className="product-detail-category"
              to={`/products?category=${encodeURIComponent(categorySlug)}`}
            >
              {categoryName}
            </Link>
          ) : (
            <span className="product-detail-category">{categoryName}</span>
          )}

          <h1 className="product-detail-title" id="product-detail-title">
            {product.name}
          </h1>

          <p className="product-detail-price" style={{ color: C.primaryDark }}>
            {formatRupiah(product.price)}
          </p>

          <div className="product-detail-stock-row">
            <span className={`product-detail-stock-pill ${stockState.className}`}>
              {stockState.label}
            </span>
            <span className="product-detail-stock-copy">
              {product.stock > 0
                ? `${product.stock} unit${product.stock === 1 ? "" : "s"} available`
                : "This product is currently unavailable"}
            </span>
          </div>

          <div className="product-detail-description">
            <h2>Description</h2>
            {descriptionBlocks.map((block) => (
              <p key={block}>{block}</p>
            ))}
          </div>

          <button
            className="product-detail-cart-button"
            type="button"
            disabled={product.stock <= 0}
            style={{
              background: product.stock <= 0 ? C.textMuted : C.primary,
            }}
          >
            {product.stock <= 0 ? "Out of Stock" : "Add to Cart"}
          </button>

          <dl className="product-detail-meta">
            <div>
              <dt>Created</dt>
              <dd>{formatDate(product.created_at)}</dd>
            </div>
            <div>
              <dt>Updated</dt>
              <dd>{formatDate(product.updated_at)}</dd>
            </div>
          </dl>
        </article>
      </div>
    </section>
  );
}

export default ProductDetail;
