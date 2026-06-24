import { useEffect, useMemo, useState, type FormEvent } from "react";
import { Link, useSearchParams } from "react-router-dom";
import {
  AlertTriangle,
  Edit3,
  ImageOff,
  PackageSearch,
  Plus,
  Search,
  Trash2,
} from "lucide-react";

import ProductImage from "../../components/ProductImage";
import {
  deleteProduct,
  getProductErrorMessage,
  listProducts,
} from "../../services/productService";
import type {
  Product,
  ProductListMeta,
  ProductListParams,
} from "../../types/product";
import {
  formatRupiah,
  getProductImage,
  getStockState,
} from "../../utils/product";

const ADMIN_PRODUCT_LIMIT = 10;

const EMPTY_META: ProductListMeta = {
  page: 1,
  limit: ADMIN_PRODUCT_LIMIT,
  total: 0,
  total_pages: 0,
};

function getPositivePage(value: string | null): number {
  const parsed = Number(value);

  if (!Number.isInteger(parsed) || parsed < 1) {
    return 1;
  }

  return parsed;
}

function getProductStatus(product: Product): {
  label: string;
  className: string;
} {
  if (product.is_active === false) {
    return {
      label: "Inactive",
      className: "is-inactive",
    };
  }

  return {
    label: "Active",
    className: "is-active",
  };
}

function getCategoryName(product: Product): string {
  return product.category?.name || "Uncategorized";
}

function AdminProductsSkeleton() {
  return (
    <div className="admin-products-list" aria-label="Loading products">
      {Array.from({ length: 5 }, (_, index) => (
        <div className="admin-products-skeleton-row" key={index}>
          <div className="admin-products-skeleton-image" />
          <div className="admin-products-skeleton-copy">
            <div className="admin-products-skeleton-line short" />
            <div className="admin-products-skeleton-line" />
          </div>
          <div className="admin-products-skeleton-line tiny" />
        </div>
      ))}
    </div>
  );
}

function AdminProducts() {
  const [searchParams, setSearchParams] = useSearchParams();
  const page = getPositivePage(searchParams.get("page"));
  const search = searchParams.get("search")?.trim() ?? "";

  const [searchInput, setSearchInput] = useState(search);
  const [products, setProducts] = useState<Product[]>([]);
  const [meta, setMeta] = useState<ProductListMeta>(EMPTY_META);
  const [isLoading, setIsLoading] = useState(true);
  const [deletingProductID, setDeletingProductID] = useState<string | null>(
    null,
  );
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);

  const hasProducts = products.length > 0;

  const query = useMemo<ProductListParams>(
    () => ({
      page,
      limit: ADMIN_PRODUCT_LIMIT,
      search: search || undefined,
      sort_by: "created_at",
      sort_order: "desc",
    }),
    [page, search],
  );

  const pageSummary = useMemo(() => {
    if (meta.total === 0) {
      return "0 product";
    }

    const start = (meta.page - 1) * meta.limit + 1;
    const end = Math.min(meta.page * meta.limit, meta.total);

    return `${start}-${end} of ${meta.total} products`;
  }, [meta]);

  useEffect(() => {
    let isActive = true;

    async function loadProducts() {
      setIsLoading(true);
      setError(null);

      try {
        const result = await listProducts(query);

        if (isActive) {
          setProducts(result.data);
          setMeta(result.meta);
        }
      } catch (loadError) {
        if (isActive) {
          setProducts([]);
          setMeta({
            ...EMPTY_META,
            page,
          });
          setError(
            getProductErrorMessage(
              loadError,
              "Admin product list could not be loaded.",
            ),
          );
        }
      } finally {
        if (isActive) {
          setIsLoading(false);
        }
      }
    }

    loadProducts();

    return () => {
      isActive = false;
    };
  }, [page, query]);

  function updateParams(nextPage: number, nextSearch: string) {
    const params = new URLSearchParams();

    if (nextPage > 1) {
      params.set("page", String(nextPage));
    }

    if (nextSearch.trim()) {
      params.set("search", nextSearch.trim());
    }

    setSearchParams(params);
  }

  function handleSearchSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    updateParams(1, searchInput);
  }

  function handleClearSearch() {
    setSearchInput("");
    updateParams(1, "");
  }

  function handlePageChange(nextPage: number) {
    updateParams(Math.max(1, nextPage), search);
  }

  async function handleDeleteProduct(product: Product) {
    const confirmed = window.confirm(
      `Delete "${product.name}"? This action cannot be undone from this page.`,
    );

    if (!confirmed) {
      return;
    }

    setDeletingProductID(product.id);
    setError(null);
    setNotice(null);

    try {
      await deleteProduct(product.id);

      setProducts((currentProducts) =>
        currentProducts.filter((item) => item.id !== product.id),
      );
      setMeta((currentMeta) => ({
        ...currentMeta,
        total: Math.max(0, currentMeta.total - 1),
      }));
      setNotice(`${product.name} deleted.`);
    } catch (deleteError) {
      setError(
        getProductErrorMessage(
          deleteError,
          "Product could not be deleted. Try again.",
        ),
      );
    } finally {
      setDeletingProductID(null);
    }
  }

  return (
    <section className="admin-page" aria-labelledby="admin-products-title">
      <header className="admin-page-header">
        <span>Admin Products</span>
        <h1 id="admin-products-title">Products.</h1>
        <p>Search, review, edit, and remove products from the admin catalog.</p>
      </header>

      <div className="admin-products-toolbar">
        <form className="admin-products-search" onSubmit={handleSearchSubmit}>
          <label htmlFor="admin-product-search">Search products</label>

          <div className="admin-products-search-row">
            <Search className="h-5 w-5" aria-hidden="true" />
            <input
              key={search}
              id="admin-product-search"
              type="search"
              placeholder="Search by product name..."
              defaultValue={search}
              onChange={(event) => setSearchInput(event.target.value)}
            />
            {search && (
              <button type="button" onClick={handleClearSearch}>
                Clear
              </button>
            )}
            <button type="submit">Search</button>
          </div>
        </form>

        <Link className="admin-products-create-button" to="/admin/products/new">
          <Plus className="h-5 w-5" aria-hidden="true" />
          Tambah Produk
        </Link>
      </div>

      {error && (
        <div className="admin-products-notice is-error" role="alert">
          <AlertTriangle className="h-5 w-5" aria-hidden="true" />
          <span>{error}</span>
        </div>
      )}

      {notice && (
        <div className="admin-products-notice is-success" role="status">
          <PackageSearch className="h-5 w-5" aria-hidden="true" />
          <span>{notice}</span>
        </div>
      )}

      {isLoading ? (
        <AdminProductsSkeleton />
      ) : !hasProducts ? (
        <div className="admin-products-empty">
          <div>
            <ImageOff className="mx-auto mb-3 h-10 w-10" aria-hidden="true" />
            <h2>No products found.</h2>
            <p>
              {search
                ? "Try another search keyword or clear the current filter."
                : "Create your first product to start filling the catalog."}
            </p>
            <Link
              className="admin-products-create-button"
              to="/admin/products/new"
            >
              <Plus className="h-5 w-5" aria-hidden="true" />
              Tambah Produk
            </Link>
          </div>
        </div>
      ) : (
        <>
          <div className="admin-products-status-line">
            <span>{pageSummary}</span>
            <span>
              Page {meta.page} / {Math.max(meta.total_pages, 1)}
            </span>
          </div>

          <div className="admin-products-table" aria-label="Admin product list">
            <div className="admin-products-table-head" aria-hidden="true">
              <span>Product</span>
              <span>Category</span>
              <span>Price</span>
              <span>Stock</span>
              <span>Status</span>
              <span>Actions</span>
            </div>

            <div className="admin-products-list">
              {products.map((product) => {
                const imagePath = getProductImage(product);
                const stockState = getStockState(product.stock);
                const status = getProductStatus(product);
                const isDeleting = deletingProductID === product.id;

                return (
                  <article className="admin-products-row" key={product.id}>
                    <div className="admin-products-product-cell">
                      <div className="admin-products-thumb">
                        <ProductImage src={imagePath} alt={product.name} />
                      </div>

                      <div>
                        <strong>{product.name}</strong>
                        <small>{product.slug}</small>
                      </div>
                    </div>

                    <span className="admin-products-muted">
                      {getCategoryName(product)}
                    </span>

                    <strong className="admin-products-price">
                      {formatRupiah(product.price)}
                    </strong>

                    <span
                      className={`admin-products-stock ${stockState.className}`}
                    >
                      {product.stock} · {stockState.label}
                    </span>

                    <span
                      className={`admin-products-status ${status.className}`}
                    >
                      {status.label}
                    </span>

                    <div className="admin-products-actions">
                      <Link to={`/admin/products/edit?id=${encodeURIComponent(product.id)}`}>
                        <Edit3 className="h-4 w-4" aria-hidden="true" />
                        Edit
                      </Link>

                      <button
                        type="button"
                        disabled={isDeleting}
                        onClick={() => void handleDeleteProduct(product)}
                      >
                        <Trash2 className="h-4 w-4" aria-hidden="true" />
                        {isDeleting ? "Deleting..." : "Delete"}
                      </button>
                    </div>
                  </article>
                );
              })}
            </div>
          </div>

          <nav
            className="admin-products-pagination"
            aria-label="Product pagination"
          >
            <button
              className="pagination-button"
              type="button"
              disabled={meta.page <= 1}
              onClick={() => handlePageChange(meta.page - 1)}
            >
              Previous
            </button>

            <span>
              {meta.page} / {Math.max(meta.total_pages, 1)}
            </span>

            <button
              className="pagination-button"
              type="button"
              disabled={meta.page >= meta.total_pages}
              onClick={() => handlePageChange(meta.page + 1)}
            >
              Next
            </button>
          </nav>
        </>
      )}
    </section>
  );
}

export default AdminProducts;
