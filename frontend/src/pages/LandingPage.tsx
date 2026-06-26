import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import {
  ArrowRight,
  BadgeCheck,
  Grid2X2,
  Mail,
  MapPin,
  Phone,
  RotateCcw,
  ShieldCheck,
  Sparkles,
  Truck,
} from "lucide-react";

import ProductImage from "../components/ProductImage";
import { useAuth } from "../hooks/useAuth";
import { listProducts } from "../services/productService";
import type { Product } from "../types/product";
import { getDashboardPath } from "../utils/authRouting";
import { formatRupiah, getProductImage, getStockState } from "../utils/product";

const FEATURED_LIMIT = 8;

type FeaturedStatus = "loading" | "ready" | "error";

function LandingFeaturedSkeleton() {
  return (
    <div
      className="landing-featured-grid"
      aria-label="Loading featured products"
    >
      {Array.from({ length: FEATURED_LIMIT }, (_, index) => (
        <article className="landing-product-card is-skeleton" key={index}>
          <div className="landing-product-skeleton media" />
          <div className="landing-product-skeleton body">
            <span />
            <strong />
            <small />
          </div>
        </article>
      ))}
    </div>
  );
}

function LandingPage() {
  const [featuredProducts, setFeaturedProducts] = useState<Product[]>([]);
  const [featuredStatus, setFeaturedStatus] =
    useState<FeaturedStatus>("loading");
  const [featuredError, setFeaturedError] = useState("");
  const { user, isAuthenticated, isLoading: isAuthLoading } = useAuth();

  const dashboardPath = getDashboardPath(user);
  const primaryHeroPath = isAuthenticated ? dashboardPath : "/products";
  const primaryHeroLabel = isAuthenticated ? "Buka Dashboard" : "Mulai Belanja";
  const secondaryHeroPath = isAuthenticated ? "/products" : "/register";
  const secondaryHeroLabel = isAuthenticated
    ? "Lanjut Belanja"
    : "Daftar Sekarang";

  useEffect(() => {
    let isActive = true;

    async function loadFeaturedProducts() {
      setFeaturedStatus("loading");
      setFeaturedError("");

      try {
        const result = await listProducts({
          page: 1,
          limit: FEATURED_LIMIT,
          sort_by: "created_at",
          sort_order: "desc",
        });

        if (isActive) {
          setFeaturedProducts(result.data.slice(0, FEATURED_LIMIT));
          setFeaturedStatus("ready");
        }
      } catch {
        if (isActive) {
          setFeaturedProducts([]);
          setFeaturedStatus("error");
          setFeaturedError("Produk belum tersedia");
        }
      }
    }

    loadFeaturedProducts();

    return () => {
      isActive = false;
    };
  }, []);

  return (
    <main className="landing-page">
      <nav className="landing-navbar" aria-label="Landing navigation">
        <Link className="landing-brand" to="/">
          <span className="landing-brand-mark" aria-hidden="true">
            <Grid2X2 size={22} />
          </span>
          <span>Bjir E-Commerce</span>
        </Link>

        <div className="landing-nav-actions">
          {isAuthLoading ? (
            <>
              <span className="landing-nav-placeholder" aria-hidden="true" />
              <span
                className="landing-nav-placeholder is-short"
                aria-hidden="true"
              />
            </>
          ) : isAuthenticated ? (
            <>
              <Link className="landing-nav-link" to="/products">
                Products
              </Link>
              <Link className="landing-nav-cta" to={dashboardPath}>
                Dashboard
              </Link>
            </>
          ) : (
            <>
              <Link className="landing-nav-link" to="/login">
                Login
              </Link>
              <Link className="landing-nav-cta" to="/register">
                Daftar
              </Link>
            </>
          )}
        </div>
      </nav>

      <section
        className={`landing-hero${isAuthLoading ? " is-loading" : ""}`}
        aria-labelledby="landing-title"
      >
        <div className="landing-hero-copy">
          <span className="landing-eyebrow">
            <Sparkles size={15} aria-hidden="true" />
            Brutal storefront
          </span>

          <h1 id="landing-title">
            Commerce with elbows. Loud shelves, fast checkout.
          </h1>

          <p>
            Bjir E-Commerce brings the catalog, cart, checkout, and admin shelf
            into one sharp storefront built for products that refuse to look
            generic.
          </p>

          {isAuthLoading ? (
            <div
              className="landing-hero-actions is-loading"
              aria-label="Loading hero actions"
            >
              <span
                className="landing-hero-cta-placeholder is-wide"
                aria-hidden="true"
              />
              <span
                className="landing-hero-cta-placeholder"
                aria-hidden="true"
              />
            </div>
          ) : (
            <div className="landing-hero-actions">
              <Link className="landing-primary-button" to={primaryHeroPath}>
                {primaryHeroLabel}
                <ArrowRight size={17} aria-hidden="true" />
              </Link>
              <Link className="landing-secondary-button" to={secondaryHeroPath}>
                {secondaryHeroLabel}
              </Link>
            </div>
          )}
        </div>

        <aside className="landing-hero-panel" aria-label="Store highlights">
          <div className="landing-hero-ticket">
            <span>01</span>
            <strong>Catalog</strong>
            <small>Search, categories, product details.</small>
          </div>
          <div className="landing-hero-ticket accent">
            <span>02</span>
            <strong>Checkout</strong>
            <small>Cart flow, shipping, orders.</small>
          </div>
          <div className="landing-hero-ticket blue">
            <span>03</span>
            <strong>Admin</strong>
            <small>Products, categories, order control.</small>
          </div>
        </aside>
      </section>

      <section
        className="landing-section"
        aria-labelledby="landing-featured-title"
      >
        <div className="landing-section-heading">
          <span className="landing-eyebrow">Featured Products</span>
          <h2 id="landing-featured-title">Produk unggulan terbaru.</h2>
          <p>
            Delapan produk terbaru dari katalog, lengkap dengan gambar, harga,
            dan status stok.
          </p>
        </div>

        {featuredStatus === "loading" ? (
          <LandingFeaturedSkeleton />
        ) : featuredStatus === "error" ? (
          <div className="landing-featured-state" role="status">
            {featuredError}
          </div>
        ) : featuredProducts.length === 0 ? (
          <div className="landing-featured-state" role="status">
            Produk belum tersedia
          </div>
        ) : (
          <div className="landing-featured-grid">
            {featuredProducts.map((product) => {
              const stockState = getStockState(product.stock);

              return (
                <Link
                  className="landing-product-card"
                  to={`/products/${product.slug}`}
                  key={product.id}
                  aria-label={`Lihat produk ${product.name}`}
                >
                  <ProductImage
                    className="landing-product-image"
                    src={getProductImage(product)}
                    alt={product.name}
                    width={640}
                    height={480}
                    sizes="(max-width: 720px) 100vw, (max-width: 1180px) 50vw, 25vw"
                  />

                  <div className="landing-product-body">
                    <span>{product.category?.name || "Uncategorized"}</span>
                    <h3>{product.name}</h3>

                    <div className="landing-product-meta">
                      <strong>{formatRupiah(product.price)}</strong>
                      <span
                        className={`landing-product-stock ${stockState.className}`}
                      >
                        {stockState.label}
                      </span>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </section>

      <section className="landing-about" aria-labelledby="landing-about-title">
        <div className="landing-about-story">
          <span className="landing-eyebrow">Tentang Bjir</span>
          <h2 id="landing-about-title">Belanja cepat, aman, tanpa drama.</h2>
          <p>
            Bjir E-Commerce dibuat sebagai etalase digital yang rapi, tegas, dan
            langsung ke inti: pembeli bisa menemukan produk, melihat stok,
            checkout, lalu memantau pesanan tanpa alur yang bertele-tele.
          </p>
        </div>

        <div
          className="landing-about-grid"
          aria-label="Keunggulan Bjir E-Commerce"
        >
          <article>
            <Truck aria-hidden="true" />
            <h3>Pengiriman Cepat</h3>
            <p>
              Pesanan diproses dengan alur yang jelas supaya barang lebih cepat
              masuk antrean pengiriman.
            </p>
          </article>

          <article>
            <ShieldCheck aria-hidden="true" />
            <h3>Pembayaran Aman</h3>
            <p>
              Flow checkout dibuat terstruktur, dengan status pembayaran dan
              pesanan yang mudah dipantau.
            </p>
          </article>

          <article>
            <BadgeCheck aria-hidden="true" />
            <h3>Produk Original</h3>
            <p>
              Katalog menampilkan informasi produk, harga, gambar, dan stok
              secara transparan.
            </p>
          </article>

          <article>
            <RotateCcw aria-hidden="true" />
            <h3>Garansi Uang Kembali</h3>
            <p>
              Belanja lebih tenang dengan dukungan pengembalian dana sesuai
              kebijakan toko.
            </p>
          </article>
        </div>
      </section>

      <section
        className="landing-contact"
        aria-labelledby="landing-contact-title"
      >
        <div>
          <span className="landing-eyebrow">Contact</span>
          <h2 id="landing-contact-title">Need the shelf opened?</h2>
          <p>
            Reach the store team for order support, admin access, or catalog
            questions.
          </p>
        </div>

        <div className="landing-contact-list">
          <a href="mailto:support@bjir-commerce.test">
            <Mail size={18} aria-hidden="true" />
            support@bjir-commerce.test
          </a>
          <a href="tel:+62000000000">
            <Phone size={18} aria-hidden="true" />
            +62 000 0000 000
          </a>
          <span>
            <MapPin size={18} aria-hidden="true" />
            Indonesia storefront lab
          </span>
          <span>
            <Truck size={18} aria-hidden="true" />
            Fulfillment-ready catalog
          </span>
        </div>
      </section>

      <footer className="landing-footer">
        <strong>Bjir E-Commerce</strong>
        <span>Sharp catalog. Hard shadows. No generic shelf.</span>
      </footer>
    </main>
  );
}

export default LandingPage;
